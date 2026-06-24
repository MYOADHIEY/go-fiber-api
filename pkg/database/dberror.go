package database

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/lib/pq"
)

// DBError is your clean domain error
type DBError struct {
	Code    string // internal code, e.g. "NOT_FOUND", "DUPLICATE", "INVALID_INPUT"
	Message string // human-readable
	Field   string // populated when pq gives us a column name
}

func (e *DBError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s (field: %s)", e.Message, e.Field)
	}
	return e.Message
}

// Map converts a raw DB error into a *DBError.
// Pass the field hints map to enrich "invalid input" errors when pq
// doesn't give us a column name — key = positional arg index (1-based),
// value = human field name.
//
// Example:
//
//	dberror.Map(err, map[int]string{1: "name", 2: "email", 3: "phone"})
func DBErrorMap(err error, args *[]interface{}, hints *map[int]string) (res error, code int) {
	if err == nil {
		return nil, 200
	}

	// --- sql.ErrNoRows ---
	if errors.Is(err, sql.ErrNoRows) {
		return &DBError{Code: "NOT_FOUND", Message: "data not found"}, http.StatusNotFound
	}

	// --- pq error ---
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		field := columnFromPQ(pqErr)

		switch pqErr.Code {

		// Class 23 — Integrity Constraint Violations
		case "23505": // unique_violation
			return &DBError{
				Code:    "DUPLICATE",
				Message: fmt.Sprintf("data already exists: %s", constraintHint(pqErr)),
				Field:   field,
			}, http.StatusConflict
		case "23503": // foreign_key_violation
			return &DBError{
				Code:    "INVALID_REFERENCE",
				Message: fmt.Sprintf("referenced data does not exist: %s", constraintHint(pqErr)),
				Field:   field,
			}, http.StatusUnprocessableEntity
		case "23502": // not_null_violation
			col := pqErr.Column
			if col == "" {
				col = field
			}
			return &DBError{
				Code:    "MISSING_REQUIRED",
				Message: "required field is missing",
				Field:   col,
			}, http.StatusUnprocessableEntity
		case "23514": // check_violation
			return &DBError{
				Code:    "CONSTRAINT_FAILED",
				Message: fmt.Sprintf("value failed validation: %s", constraintHint(pqErr)),
				Field:   field,
			}, http.StatusUnprocessableEntity

		// Class 22 — Data Exception (type mismatch, out of range, etc.)
		case "22P02": // invalid_text_representation  e.g. "invalid input syntax for type integer"
			return &DBError{
				Code:    "INVALID_INPUT",
				Message: "invalid value format",
				Field:   guessField(pqErr, args, hints),
			}, http.StatusBadRequest
		case "22003": // numeric_value_out_of_range
			return &DBError{
				Code:    "INVALID_INPUT",
				Message: "value is out of allowed range",
				Field:   guessField(pqErr, args, hints),
			}, http.StatusUnprocessableEntity
		case "22001": // string_data_right_truncation
			return &DBError{
				Code:    "INVALID_INPUT",
				Message: "value is too long for this field",
				Field:   guessField(pqErr, args, hints),
			}, http.StatusUnprocessableEntity

		// Class 42 — Syntax / Schema errors (you want to know during dev)
		case "42P01": // undefined_table
			return &DBError{Code: "SCHEMA_ERROR", Message: fmt.Sprintf("table not found: %s", pqErr.Table)}, http.StatusInternalServerError
		case "42703": // undefined_column
			return &DBError{Code: "SCHEMA_ERROR", Message: fmt.Sprintf("column not found: %s", pqErr.Column)}, http.StatusInternalServerError
		}

		// Fallback: still a pq error but unhandled code
		return &DBError{
			Code:    string(pqErr.Code),
			Message: pqErr.Message, // raw pg message, at least it's honest
			Field:   field,
		}, http.StatusUnprocessableEntity
	}

	// Passthrough for anything else
	return err, http.StatusUnprocessableEntity
}

// IsNotFound checks if a mapped error is a NOT_FOUND.
func IsNotFound(err error) bool {
	var dbErr *DBError
	return errors.As(err, &dbErr) && dbErr.Code == "NOT_FOUND"
}

// IsDuplicate checks for unique constraint violations.
func IsDuplicate(err error) bool {
	var dbErr *DBError
	return errors.As(err, &dbErr) && dbErr.Code == "DUPLICATE"
}

// --- helpers ---

func columnFromPQ(e *pq.Error) string {
	if e.Column != "" {
		return e.Column
	}
	return ""
}

func constraintHint(e *pq.Error) string {
	if e.Constraint != "" {
		return e.Constraint
	}
	return e.Table
}

// guessField tries to extract a column name from the pq detail/message
// string, then falls back to the positional hint map you passed in.
func guessField(e *pq.Error, args *[]interface{}, hints *map[int]string) string {
	// pq sometimes gives column directly (lucky path)
	if e.Column != "" {
		return e.Column
	}

	var hintList map[int]string
	if hints == nil {
		return "please provide hints"
	}
	if hints != nil {
		hintList = *hints
		// Extract the bad value from pq message
		// e.g. `invalid input syntax for type integer: "abc"`
		badValue := extractBadValue(e.Message)
		if badValue != "" && len(*args) > 0 {
			// Find which arg matches the bad value
			for i, arg := range *args {
				strVal := fmt.Sprintf("%v", arg)
				if strVal == badValue {
					if name, ok := hintList[i+1]; ok { // hints are 1-based
						return name
					}
				}
			}
		}
	}

	return "unknown"
}

// extractBadValue pulls the rejected value out of pq error message.
// Input:  `invalid input syntax for type integer: "abc"`
// Output: `abc`
func extractBadValue(msg string) string {
	// pq wraps the bad value in double quotes after the colon
	idx := strings.LastIndex(msg, ": \"")
	if idx == -1 {
		return ""
	}
	rest := msg[idx+3:]
	end := strings.Index(rest, "\"")
	if end == -1 {
		return ""
	}
	return rest[:end]
}

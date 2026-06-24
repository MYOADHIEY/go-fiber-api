package base

import (
	"database/sql"
	"kbaa-fiber-api/pkg/jwe"
	"kbaa-fiber-api/pkg/jwt"
	"kbaa-fiber-api/pkg/str"
	basevm "kbaa-fiber-api/usecase/base/viewmodels"

	"github.com/go-playground/validator/v10"
	jwtFiber "github.com/gofiber/jwt/v2"
)

var (
	defaultLimit    = 10
	maxLimit        = 50
	defaultSort     = "asc"
	sortWhitelist   = []string{"asc", "desc"}
	passwordLength  = 6
	defaultLastPage = 0

	// DefaultLocation ...
	DefaultLocation = "Asia/Jakarta"
	// DefaultTimezone ...
	DefaultTimezone = "+07:00"

	// OtpLifetime ...
	OtpLifetime = "8m"
	// MaxOtpSubmitRetry ...
	MaxOtpSubmitRetry = 5.0
)

type BaseUc struct {
	ReqID     string
	UserID    string
	EnvConfig map[string]string
	DB        *sql.DB
	TX        *sql.Tx
	Validate  *validator.Validate
	JweCred   jwe.Credential
	JwtCred   jwt.Credential
	JwtConfig jwtFiber.Config
}

func (uc BaseUc) SetPaginationParameter(page, limit int, orderBy, sort string, orderByWhiteLists, orderByStringWhiteLists []string) (int, int, int, string, string) {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 || limit > maxLimit {
		limit = defaultLimit
	}

	orderBy = uc.CheckWhiteList(orderBy, orderByWhiteLists)
	if str.Contains(orderByStringWhiteLists, orderBy) {
		orderBy = `LOWER(` + orderBy + `)`
	}

	if !str.Contains(sortWhitelist, sort) {
		sort = defaultSort
	}
	offset := (page - 1) * limit

	return offset, limit, page, orderBy, sort
}

func (uc BaseUc) CheckWhiteList(orderBy string, whiteLists []string) string {
	for _, whiteList := range whiteLists {
		if orderBy == whiteList {
			return orderBy
		}
	}

	return "def.updated_at"
}

func (uc BaseUc) SetPaginationResponse(page, limit, total int) (paginationResponse basevm.BasePaginationVM) {
	var lastPage int

	if total > 0 {
		lastPage = total / limit

		if total%limit != 0 {
			lastPage = lastPage + 1
		}
	} else {
		lastPage = defaultLastPage
	}

	paginationResponse = basevm.BasePaginationVM{
		CurrentPage: page,
		LastPage:    lastPage,
		Total:       total,
		PerPage:     limit,
	}

	return paginationResponse
}

func (uc BaseUc) TestUC(test string) {

}

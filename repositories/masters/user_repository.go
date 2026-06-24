package masters

import (
	"context"
	"database/sql"
	"kbaa-fiber-api/pkg/database"
	imodels "kbaa-fiber-api/repositories/masters/models"
	"kbaa-fiber-api/usecase/master/viewmodels"
	iviewmodels "kbaa-fiber-api/usecase/master/viewmodels"
	"net/http"
	"strconv"
)

type IUserRepository interface {
	FindAll(c context.Context, params imodels.UserParameters) ([]iviewmodels.UserVM, int, int, error)
	FindByID(c context.Context, params imodels.UserParameters) (iviewmodels.UserVM, int, error)
	Add(c context.Context, data *imodels.User) (string, int, error)
	Update(c context.Context, data *imodels.User) (string, int, error)
	Delete(c context.Context, data *imodels.User) (string, int, error)
}

type UserRepository struct {
	DB database.DBTX
}

func NewUserRepository(DB database.DBTX) IUserRepository {
	return &UserRepository{DB: DB}
}

func (repository UserRepository) scanRows(rows *sql.Rows) (res viewmodels.UserVM, err error) {
	err = rows.Scan(
		&res.ID, &res.Name, &res.Email, &res.Phone, &res.Address,
		&res.CreatedAt, &res.UpdatedAt, &res.DeletedAt,
		&res.CreatedBy, &res.UpdatedBy, &res.DeletedBy,
	)
	if err != nil {
		return
	}

	return
}

func (repository UserRepository) scanRow(row *sql.Row) (res viewmodels.UserVM, err error) {
	err = row.Scan(
		&res.ID, &res.Name, &res.Email, &res.Phone, &res.Address,
		&res.CreatedAt, &res.UpdatedAt, &res.DeletedAt,
		&res.CreatedBy, &res.UpdatedBy, &res.DeletedBy,
	)
	if err != nil {
		return
	}
	return
}

func (repository UserRepository) FindAll(c context.Context, params imodels.UserParameters) (res []iviewmodels.UserVM, count int, responseCode int, err error) {
	responseCode = 200
	args := []interface{}{
		params.Offset, params.Limit,
	}

	whereClause := ``
	counter := 3
	if params.Type != "" {
		whereClause += ` and def._type = $` + strconv.Itoa(counter)
		args = append(args, params.Type)
	}

	statement := imodels.UserSelectStatement + imodels.UserWhereSelectStatement + whereClause + ` ORDER BY ` +
		params.By + ` ` + params.Sort + ` OFFSET $1 limit $2 `

	rows, err := repository.DB.QueryContext(c, statement, args...)

	if err != nil {
		err, responseCode = database.DBErrorMap(err, nil, nil)
		return
	}
	defer rows.Close()
	for rows.Next() {
		temp, err := repository.scanRows(rows)
		if err != nil {
			return res, count, http.StatusBadRequest, err
		}

		res = append(res, temp)
	}

	err = rows.Err()
	if err != nil {

		return res, count, responseCode, err
	}
	statement = ` select count(*) from _user def ` +
		imodels.UserWhereSelectStatement
	err = repository.DB.QueryRowContext(c, statement).Scan(&count)
	if err != nil {

		return
	}

	return
}

func (repository UserRepository) FindByID(c context.Context, params imodels.UserParameters) (res iviewmodels.UserVM, resCode int, err error) {

	statement := imodels.UserSelectStatement + imodels.UserWhereSelectStatement +
		` and def.id = $1 `
	row := repository.DB.QueryRowContext(c, statement, params.ID)
	res, err = repository.scanRow(row)
	if err != nil {
		err, resCode = database.DBErrorMap(err, nil, nil)
		return
	}

	return
}

func (repository UserRepository) Add(c context.Context, data *imodels.User) (res string, code int, err error) {
	// hint example to help debug and help shows column caused the error
	code = 200
	utype := "ajsh"
	args := []interface{}{
		data.Name, data.Email, data.Phone, data.Address,
		data.CreatedAt, data.CreatedBy, utype,
	}
	hints := map[int]string{
		1: "_name",
		2: "_email",
		3: "_phone",
		4: "address",
		5: "created_at",
		6: "created_by",
		7: "_type",
	}

	statement := `insert into _user
	(_name, _email, _phone,address, created_at, created_by, _type)
	values ($1, $2, $3, $4, $5, $6, $7) returning id
	`
	err = repository.DB.QueryRowContext(c, statement,
		args...,
	).Scan(&res)
	if err != nil {
		err, code = database.DBErrorMap(err, &args, &hints)
		return
	}
	return
}

func (repository UserRepository) Update(c context.Context, data *imodels.User) (res string, resCode int, err error) {
	resCode = 200
	statement := ` update _user set
	_name = $1, _email = $2, _phone =$3 ,address = $4,
	updated_at = $5, updated_by = $6
	where id = $7 returning id 
	`
	err = repository.DB.QueryRowContext(c, statement,
		data.Name, data.Email, data.Phone, data.Address,
		data.UpdatedAt, data.UpdatedBy,
		data.ID,
	).Scan(&res)
	return
}

func (repository UserRepository) Delete(c context.Context, data *imodels.User) (res string, resCode int, err error) {
	resCode = 200
	statement := ` update _user set
	deleted_at = $1, deleted_by = $2
	where id = $3 returning id 
	`
	err = repository.DB.QueryRowContext(c, statement,
		data.DeletedAt, data.DeletedBy,
		data.ID,
	).Scan(&res)
	return
}

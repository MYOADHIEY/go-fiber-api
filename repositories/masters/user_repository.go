package masters

import (
	"context"
	"database/sql"
	"kbaa-fiber-api/pkg/database"
	imodels "kbaa-fiber-api/repositories/masters/models"
	"kbaa-fiber-api/usecase/master/viewmodels"
	iviewmodels "kbaa-fiber-api/usecase/master/viewmodels"
)

type IUserRepository interface {
	FindAll(c context.Context, params imodels.UserParameters) ([]iviewmodels.UserVM, int, error)
	FindByID(c context.Context, params imodels.UserParameters) (iviewmodels.UserVM, error)
	Add(c context.Context, data *imodels.User) (string, error)
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

func (repository UserRepository) FindAll(c context.Context, params imodels.UserParameters) (res []iviewmodels.UserVM, count int, err error) {

	statement := imodels.UserSelectStatement + imodels.UserWhereSelectStatement + ` ORDER BY ` +
		params.By + ` ` + params.Sort + ` OFFSET $1 limit $2 `

	rows, err := repository.DB.QueryContext(c, statement, params.Offset, params.Limit)

	if err != nil {

		return
	}
	defer rows.Close()
	for rows.Next() {
		temp, err := repository.scanRows(rows)
		if err != nil {

			return res, count, err
		}

		res = append(res, temp)
	}

	err = rows.Err()
	if err != nil {

		return res, count, err
	}
	statement = ` select count(*) from _user def ` +
		imodels.UserWhereSelectStatement
	err = repository.DB.QueryRowContext(c, statement).Scan(&count)
	if err != nil {

		return
	}

	return
}

func (repository UserRepository) FindByID(c context.Context, params imodels.UserParameters) (res iviewmodels.UserVM, err error) {

	statement := imodels.UserSelectStatement + imodels.UserWhereSelectStatement +
		` and def.id = $1 `
	row := repository.DB.QueryRowContext(c, statement, params.ID)
	res, err = repository.scanRow(row)
	if err != nil {
		return
	}

	return
}

func (repository UserRepository) Add(c context.Context, data *imodels.User) (res string, err error) {
	statement := `insert into _user
	(_name, _email, _phone,address, created_at, created_by)
	values ($1, $2, $3, $4, $5, $6) returning id
	`
	err = repository.DB.QueryRowContext(c, statement,
		data.Name, data.Email, data.Phone, data.Address,
		data.CreatedAt, data.CreatedBy,
	).Scan(&res)
	if err != nil {
		return
	}
	return
}

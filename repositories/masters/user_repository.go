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
	FindAll(c context.Context, params imodels.UserParameters) (res []iviewmodels.UserVM, count int, err error)
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
	)
	if err != nil {
		return
	}

	return
}

func (repository UserRepository) scanRow(row *sql.Row) (res viewmodels.UserVM, err error) {
	err = row.Scan(
		&res.ID, &res.Name, &res.Email, &res.Phone, &res.Address,
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

package useraccount

import (
	"context"
	"database/sql"
	"kbaa-fiber-api/pkg/database"
	imodels "kbaa-fiber-api/repositories/useraccount/models"
	iviewmodels "kbaa-fiber-api/usecase/useraccount/viewmodels"
)

type IAccountRepostiry interface {
	Add(c context.Context, data imodels.Registration) (res string, responseCode int, err error)
}

type RegistrationRepository struct {
	DB database.DBTX
}

func NewRegistrationRepository(DB database.DBTX) IAccountRepostiry {
	return &RegistrationRepository{DB: DB}
}

func (repository RegistrationRepository) scanRows(rows *sql.Rows) (res iviewmodels.RegistrationtVM, err error) {
	err = rows.Scan(
		&res.ID, &res.Email, &res.Password,
	)
	if err != nil {
		return
	}

	return
}

func (repository RegistrationRepository) scanRow(row *sql.Row) (res iviewmodels.RegistrationtVM, err error) {
	err = row.Scan(
		&res.ID, &res.Email, &res.Password,
	)
	if err != nil {
		return
	}
	return
}

func (repository RegistrationRepository) Add(c context.Context, data imodels.Registration) (res string, responseCode int, err error) {
	statement := `insert into _user
	( _email, _password)
	values ($1, $2) returning id
	`
	err = repository.DB.QueryRowContext(c, statement,
		data.Email,
		data.Password,
	).Scan(&res)
	if err != nil {
		err, responseCode = database.DBErrorMap(err, nil, nil)
		return
	}
	return
}

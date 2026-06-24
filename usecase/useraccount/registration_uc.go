package useraccount

import (
	"context"
	irepository "kbaa-fiber-api/repositories/useraccount"
	imodels "kbaa-fiber-api/repositories/useraccount/models"
	irequest "kbaa-fiber-api/server/requests/useraccount"
	baseUc "kbaa-fiber-api/usecase/base"
	iviewmodels "kbaa-fiber-api/usecase/useraccount/viewmodels"
)

type RegistrationUc struct {
	*baseUc.BaseUc
}

func (uc RegistrationUc) Register(c context.Context, data *irequest.RegistrationRequest) (res iviewmodels.RegistrationtVM, responseCode int, err error) {
	responseCode = 200
	repo := irepository.NewRegistrationRepository(uc.DB)
	inputModel := imodels.Registration{
		Email:    data.Email,
		Password: data.Password,
	}
	res.ID, responseCode, err = repo.Add(c, inputModel)
	if err != nil {
		return
	}
	return
}

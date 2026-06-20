package master

import (
	"context"
	"kbaa-fiber-api/pkg/functioncaller"
	"kbaa-fiber-api/pkg/logruslogger"
	irepository "kbaa-fiber-api/repositories/masters"
	imodles "kbaa-fiber-api/repositories/masters/models"
	baseUC "kbaa-fiber-api/usecase/base"
	baseviewmodels "kbaa-fiber-api/usecase/base/viewmodels"
	iviewmodels "kbaa-fiber-api/usecase/master/viewmodels"
)

type UserUC struct {
	*baseUC.BaseUc
}

func (uc UserUC) BuildBody(res *iviewmodels.UserVM) {

}

func (uc UserUC) FindAll(c context.Context, parameter imodles.UserParameters) (res []iviewmodels.UserVM, p baseviewmodels.BasePaginationVM, err error) {
	parameter.Offset, parameter.Limit, parameter.Page, parameter.By, parameter.Sort = uc.SetPaginationParameter(parameter.Page, parameter.Limit, parameter.By, parameter.Sort, imodles.UserOrderBy, imodles.UserByString)
	var count int

	repository := irepository.NewUserRepository(uc.DB)
	res, count, err = repository.FindAll(c, parameter)

	if err != nil {

		logruslogger.Log(logruslogger.WarnLevel, err.Error(), functioncaller.PrintFuncName(), "query", c.Value("requestid"))
		return res, p, err
	}

	p = uc.SetPaginationResponse(parameter.Page, parameter.Limit, count)

	for i := range res {
		uc.BuildBody(&res[i])
	}

	return
}

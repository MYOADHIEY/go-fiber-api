package useraccount

import (
	"context"
	baseUC "kbaa-fiber-api/usecase/base"
	utilsUc "kbaa-fiber-api/usecase/utils"
	utilsVM "kbaa-fiber-api/usecase/utils/viewmodels"
)

type AuthUC struct {
	*baseUC.BaseUc
}

func (uc AuthUC) BuildJWTToken(c context.Context) (res utilsVM.JWTVM, err error) {
	jwtUC := utilsUc.JWTUC{BaseUc: uc.BaseUc}
	payload := map[string]interface{}{
		"user_id": "1010",
	}
	err = jwtUC.GenerateJWTToken(c, payload, &res)
	if err != nil {
		return
	}
	return
}

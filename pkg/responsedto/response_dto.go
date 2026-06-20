package responsedto

import (
	basevm "kbaa-fiber-api/usecase/base/viewmodels"
)

func SuccessResponse(code int, data interface{}, meta interface{}) basevm.BaseResponseVM {
	return basevm.BaseResponseVM{
		StatusCode: code,
		Data:       data,
		Meta:       meta,
		Errors:     []string{},
	}
}

func ErrorRepsonse(code int, err ...string) basevm.BaseResponseVM {
	return basevm.BaseResponseVM{
		StatusCode: code,
		Data:       nil,
		Meta:       nil,
		Errors:     err,
	}
}

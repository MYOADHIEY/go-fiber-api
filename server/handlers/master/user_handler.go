package master

import (
	"context"
	"fmt"
	"kbaa-fiber-api/pkg/str"
	ihandler "kbaa-fiber-api/server/handlers/base"

	imodels "kbaa-fiber-api/repositories/masters/models"

	iusecase "kbaa-fiber-api/usecase/master"

	"github.com/gofiber/fiber/v2"
)

type UserHanlder struct {
	ihandler.BaseHandler
}

func (h *UserHanlder) Find(ctx *fiber.Ctx) error {
	c := ctx.Locals("ctx").(context.Context)
	parameters := imodels.UserParameters{
		ResponseCode: str.StringToInt(ctx.Query("response_code")),
		Search:       ctx.Query("search"),
		Page:         str.StringToInt(ctx.Query("page")),
		Limit:        str.StringToInt(ctx.Query("limit")),
	}

	fmt.Println("code", parameters.ResponseCode, ctx.Query("response_code"), parameters.Page, parameters.Limit)
	uc := iusecase.UserUC{BaseUc: h.BaseUC}
	res, meta, err := uc.FindAll(c, parameters)
	if err != nil {
		fmt.Println("error ", err)
	}
	return h.SendResponse(ctx, res, meta, err, parameters.ResponseCode)
}

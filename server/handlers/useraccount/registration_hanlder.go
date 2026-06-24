package useraccount

import (
	"context"
	"fmt"
	ihandler "kbaa-fiber-api/server/handlers/base"
	irequest "kbaa-fiber-api/server/requests/useraccount"
	iusecase "kbaa-fiber-api/usecase/useraccount"

	"github.com/gofiber/fiber/v2"
)

type RegistrationHandler struct {
	ihandler.BaseHandler
}

func (h *RegistrationHandler) Register(ctx *fiber.Ctx) error {
	c := ctx.Locals("ctx").(context.Context)
	data := new(irequest.RegistrationRequest)
	if err := ctx.BodyParser(data); err != nil {

		return h.SendResponse(ctx, nil, nil, err, 400)
	}

	if err := h.Validator.Struct(data); err != nil {

		return h.SendResponse(ctx, nil, nil, err, 400)
	}
	accounUC := iusecase.RegistrationUc{BaseUc: h.BaseUC}
	res, err := accounUC.Register(c, data)
	if err != nil {
		fmt.Println("err", err.Error())
		return h.SendResponse(ctx, nil, nil, err.Error(), 400)
	}
	return h.SendResponse(ctx, res, nil, nil, 200)
}

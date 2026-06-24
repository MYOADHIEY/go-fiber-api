package useraccount

import (
	"context"
	ihandler "kbaa-fiber-api/server/handlers/base"
	irequest "kbaa-fiber-api/server/requests/useraccount"
	iusecase "kbaa-fiber-api/usecase/useraccount"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type RegistrationHandler struct {
	ihandler.BaseHandler
}

func (h *RegistrationHandler) Register(ctx *fiber.Ctx) error {
	c := ctx.Locals("ctx").(context.Context)
	data := new(irequest.RegistrationRequest)
	if err := ctx.BodyParser(data); err != nil {

		return h.SendResponse(ctx, nil, nil, err, http.StatusBadRequest)
	}

	if err := h.Validator.Struct(data); err != nil {

		return h.SendResponse(ctx, nil, nil, err, http.StatusBadRequest)
	}
	accounUC := iusecase.RegistrationUc{BaseUc: h.BaseUC}
	res, responseCode, err := accounUC.Register(c, data)
	if err != nil {
		return h.SendResponse(ctx, nil, nil, err.Error(), responseCode)
	}
	return h.SendResponse(ctx, res, nil, err, responseCode)
}

package useraccount

import (
	"context"
	ihandler "kbaa-fiber-api/server/handlers/base"
	iusecase "kbaa-fiber-api/usecase/useraccount"

	"github.com/gofiber/fiber/v2"
)

type AuthHanlder struct {
	ihandler.BaseHandler
}

func (h *AuthHanlder) BuildJWTToken(ctx *fiber.Ctx) error {
	c := ctx.Locals("ctx").(context.Context)
	authUc := iusecase.AuthUC{BaseUc: h.BaseUC}
	res, err := authUc.BuildJWTToken(c)
	if err != nil {
		return h.SendResponse(ctx, nil, nil, err, 400)
	}
	return h.SendResponse(ctx, res, nil, nil, 200)
}

package master

import (
	"context"
	"kbaa-fiber-api/pkg/str"
	ihandler "kbaa-fiber-api/server/handlers/base"

	imodels "kbaa-fiber-api/repositories/masters/models"

	irequests "kbaa-fiber-api/server/requests/master"
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

	uc := iusecase.UserUC{BaseUc: h.BaseUC}
	res, meta, err := uc.FindAll(c, parameters)
	if err != nil {
		return h.SendResponse(ctx, nil, nil, err, 400)
	}
	return h.SendResponse(ctx, res, meta, err, 200)
}

func (h *UserHanlder) FindByID(ctx *fiber.Ctx) error {

	c := ctx.Locals("ctx").(context.Context)
	id := ctx.Params("id")
	uc := iusecase.UserUC{BaseUc: h.BaseUC}

	res, err := uc.FindById(c, imodels.UserParameters{ID: id})
	if err != nil {
		return h.SendResponse(ctx, nil, nil, err, 400)
	}
	return h.SendResponse(ctx, res, nil, err, 200)
}

func (h *UserHanlder) Add(ctx *fiber.Ctx) error {
	c := ctx.Locals("ctx").(context.Context)

	inputData := new(irequests.UserRequest)

	if err := ctx.BodyParser(inputData); err != nil {

		return h.SendResponse(ctx, nil, nil, err, 400)
	}

	if err := h.Validator.Struct(inputData); err != nil {

		return h.SendResponse(ctx, nil, nil, err, 400)
	}

	uc := iusecase.UserUC{BaseUc: h.BaseUC}
	res, err := uc.Add(c, inputData)
	if err != nil {
		return h.SendResponse(ctx, nil, nil, err, 400)
	}
	return h.SendResponse(ctx, nil, res, nil, 200)
}

func (h *UserHanlder) Update(ctx *fiber.Ctx) error {
	c := ctx.Locals("ctx").(context.Context)

	id := ctx.Params("id")
	if id == "" {
		return h.SendResponse(ctx, nil, nil, "err", 400)
	}
	inputData := new(irequests.UserRequest)

	if err := ctx.BodyParser(inputData); err != nil {

		return h.SendResponse(ctx, nil, nil, err, 400)
	}

	if err := h.Validator.Struct(inputData); err != nil {

		return h.SendResponse(ctx, nil, nil, err, 400)
	}
	inputData.ID = id
	uc := iusecase.UserUC{BaseUc: h.BaseUC}
	res, err := uc.Update(c, inputData)
	if err != nil {
		return h.SendResponse(ctx, nil, nil, err, 400)
	}
	return h.SendResponse(ctx, nil, res, nil, 200)
}

func (h *UserHanlder) Delete(ctx *fiber.Ctx) error {
	c := ctx.Locals("ctx").(context.Context)

	id := ctx.Params("id")
	if id == "" {
		return h.SendResponse(ctx, nil, nil, "err", 400)
	}
	inputData := new(irequests.UserRequest)

	if err := ctx.BodyParser(inputData); err != nil {

		return h.SendResponse(ctx, nil, nil, err, 400)
	}

	if err := h.Validator.Struct(inputData); err != nil {

		return h.SendResponse(ctx, nil, nil, err, 400)
	}
	inputData.ID = id
	uc := iusecase.UserUC{BaseUc: h.BaseUC}
	res, err := uc.Delete(c, inputData)
	if err != nil {
		return h.SendResponse(ctx, nil, nil, err, 400)
	}
	return h.SendResponse(ctx, nil, res, nil, 200)
}

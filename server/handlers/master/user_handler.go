package master

import (
	"context"
	"kbaa-fiber-api/pkg/str"
	ihandler "kbaa-fiber-api/server/handlers/base"
	"net/http"

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
		Type:         ctx.Query("_type"),
		By:           ctx.Query("_by"),
	}

	uc := iusecase.UserUC{BaseUc: h.BaseUC}
	res, meta, responseCode, err := uc.FindAll(c, parameters)
	if err != nil {
		return h.SendResponse(ctx, nil, nil, err, responseCode)
	}
	return h.SendResponse(ctx, res, meta, err, responseCode)
}

func (h *UserHanlder) FindByID(ctx *fiber.Ctx) error {

	c := ctx.Locals("ctx").(context.Context)
	id := ctx.Params("id")
	uc := iusecase.UserUC{BaseUc: h.BaseUC}

	res, responseCode, err := uc.FindById(c, imodels.UserParameters{ID: id})
	if err != nil {
		return h.SendResponse(ctx, nil, nil, err, responseCode)
	}
	return h.SendResponse(ctx, res, nil, err, responseCode)
}

func (h *UserHanlder) Add(ctx *fiber.Ctx) error {
	c := ctx.Locals("ctx").(context.Context)

	inputData := new(irequests.UserRequest)

	if err := ctx.BodyParser(inputData); err != nil {

		return h.SendResponse(ctx, nil, nil, err, http.StatusBadRequest)
	}

	if err := h.Validator.Struct(inputData); err != nil {

		return h.SendResponse(ctx, nil, nil, err, http.StatusBadRequest)
	}

	uc := iusecase.UserUC{BaseUc: h.BaseUC}
	res, responseCode, err := uc.Add(c, inputData)
	if err != nil {
		return h.SendResponse(ctx, nil, nil, err, responseCode)
	}
	return h.SendResponse(ctx, nil, res, err, responseCode)
}

func (h *UserHanlder) Update(ctx *fiber.Ctx) error {
	c := ctx.Locals("ctx").(context.Context)

	id := ctx.Params("id")
	if id == "" {
		return h.SendResponse(ctx, nil, nil, "err", http.StatusBadRequest)
	}
	inputData := new(irequests.UserRequest)

	if err := ctx.BodyParser(inputData); err != nil {

		return h.SendResponse(ctx, nil, nil, err, http.StatusBadRequest)
	}

	if err := h.Validator.Struct(inputData); err != nil {

		return h.SendResponse(ctx, nil, nil, err, http.StatusBadRequest)
	}
	inputData.ID = id
	uc := iusecase.UserUC{BaseUc: h.BaseUC}
	res, responseCode, err := uc.Update(c, inputData)
	if err != nil {
		return h.SendResponse(ctx, nil, nil, err, responseCode)
	}
	return h.SendResponse(ctx, nil, res, nil, responseCode)
}

func (h *UserHanlder) Delete(ctx *fiber.Ctx) error {
	c := ctx.Locals("ctx").(context.Context)

	id := ctx.Params("id")
	if id == "" {
		return h.SendResponse(ctx, nil, nil, "Please Provide ID", http.StatusBadRequest)
	}
	inputData := new(irequests.UserRequest)

	if err := ctx.BodyParser(inputData); err != nil {

		return h.SendResponse(ctx, nil, nil, err, http.StatusBadRequest)
	}

	if err := h.Validator.Struct(inputData); err != nil {

		return h.SendResponse(ctx, nil, nil, err, http.StatusBadRequest)
	}
	inputData.ID = id
	uc := iusecase.UserUC{BaseUc: h.BaseUC}
	res, responseCode, err := uc.Delete(c, inputData)
	if err != nil {
		return h.SendResponse(ctx, nil, nil, err, responseCode)
	}
	return h.SendResponse(ctx, nil, res, nil, responseCode)
}

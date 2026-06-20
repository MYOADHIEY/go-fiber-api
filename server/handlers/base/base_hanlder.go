package base

import (
	"database/sql"
	"fmt"

	"kbaa-fiber-api/pkg/responsedto"
	baseUC "kbaa-fiber-api/usecase/base"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtFiber "github.com/gofiber/jwt/v2"
)

type BaseHandler struct {
	FiberApp   *fiber.App
	Validator  *validator.Validate
	DB         *sql.DB
	Translator ut.Translator
	JwtConfig  jwtFiber.Config
	BaseUC     *baseUC.BaseUc
}

func (h BaseHandler) SendResponse(ctx *fiber.Ctx, data interface{}, meta interface{}, err interface{}, code int) error {

	if code >= 400 {

		return h.SendErrorResponse(ctx, code, h.errors(err)...)
	}
	return h.SendSuccessResponse(ctx, code, data, meta)
}

// send success response
func (h BaseHandler) SendSuccessResponse(ctx *fiber.Ctx, code int, data interface{}, meta interface{}) error {
	resposnse := responsedto.SuccessResponse(code, data, meta)

	return ctx.Status(code).JSON(resposnse)
}

// send success response
func (h BaseHandler) SendErrorResponse(ctx *fiber.Ctx, code int, err ...string) error {
	resposnse := responsedto.ErrorRepsonse(code, err...)

	return ctx.Status(code).JSON(resposnse)
}

func (h BaseHandler) errors(err interface{}) []string {

	var errMsgs []string

	switch e := err.(type) {
	case string:
		fmt.Println("type string")
		errMsgs = []string{e}
	case []string:
		fmt.Println("type string []")
		errMsgs = e
	case nil:
		fmt.Println("type nil")
		errMsgs = nil
	default:
		fmt.Println("type defalut")
		errMsgs = []string{fmt.Sprintf("%v", e)}
	}
	return errMsgs
}

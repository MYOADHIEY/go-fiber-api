package bootsrap

import (
	baseuc "kbaa-fiber-api/usecase/base"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AppBoot struct {
	App        *fiber.App
	BaseUC     baseuc.BaseUc
	Validator  *validator.Validate
	Translator ut.Translator
}

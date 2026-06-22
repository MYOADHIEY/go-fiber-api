package bootsrap

import (
	"net/http"

	masterRoute "kbaa-fiber-api/server/bootsrap/routers/master"
	useraccountRoute "kbaa-fiber-api/server/bootsrap/routers/useraccount"
	ihanlder "kbaa-fiber-api/server/handlers/base"

	"github.com/gofiber/fiber/v2"
)

func (appBoot AppBoot) RegisterRouters() {
	handler := ihanlder.BaseHandler{
		FiberApp:   appBoot.App,
		BaseUC:     &appBoot.BaseUC,
		Validator:  appBoot.Validator,
		Translator: appBoot.Translator,
	}

	appBoot.App.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).JSON("work")
	})

	apiV1 := appBoot.App.Group("/v1")

	masterRoutes := masterRoute.MasterRouters{RouteGroup: apiV1, Handler: handler}
	masterRoutes.RegisterRouters()

	userAccountRoutes := useraccountRoute.UserAccountRouters{RouteGroup: apiV1, Handler: handler}
	userAccountRoutes.RegisterRouters()

}

package master

import (
	"kbaa-fiber-api/pkg/str"
	ihandler "kbaa-fiber-api/server/handlers/base"
	"kbaa-fiber-api/server/middlewares"
	"time"

	ihandlers "kbaa-fiber-api/server/handlers/master"

	"github.com/gofiber/fiber/v2"
)

type UserRoute struct {
	RouteGroup fiber.Router
	Handler    ihandler.BaseHandler
}

func (route UserRoute) RegisterRouters() {
	handler := ihandlers.UserHanlder{BaseHandler: route.Handler}
	jwtMiddleware := middlewares.JWTMiddleware{BaseUc: handler.BaseUC}

	r := route.RouteGroup.Group("/user")
	r.Use(middlewares.SavingContextValue(time.Duration(str.StringToInt(route.Handler.BaseUC.EnvConfig["APP_TIMEOUT"])) * time.Second))
	r.Use(jwtMiddleware.VerifyBasic)
	r.Get("/", handler.Find)
	r.Get("/:id", handler.FindByID)
	r.Post("/", handler.Add)
}

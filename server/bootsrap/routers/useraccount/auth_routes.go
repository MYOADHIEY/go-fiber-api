package useraccount

import (
	"kbaa-fiber-api/pkg/str"
	ihandler "kbaa-fiber-api/server/handlers/base"
	"kbaa-fiber-api/server/middlewares"
	"time"

	ihandlers "kbaa-fiber-api/server/handlers/useraccount"

	"github.com/gofiber/fiber/v2"
)

type AuthRoute struct {
	RouteGroup fiber.Router
	Handler    ihandler.BaseHandler
}

func (route AuthRoute) RegisterRouters() {
	handler := ihandlers.AuthHanlder{BaseHandler: route.Handler}
	// jwtMiddleware := middlewares.JWTMiddleware{BaseUc: handler.BaseUC}

	r := route.RouteGroup.Group("/auth")
	r.Use(middlewares.SavingContextValue(time.Duration(str.StringToInt(route.Handler.BaseUC.EnvConfig["APP_TIMEOUT"])) * time.Second))
	r.Get("/", handler.BuildJWTToken)

}

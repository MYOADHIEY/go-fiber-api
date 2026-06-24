package useraccount

import (
	ihandler "kbaa-fiber-api/server/handlers/base"

	"github.com/gofiber/fiber/v2"
)

type UserAccountRouters struct {
	RouteGroup fiber.Router
	Handler    ihandler.BaseHandler
}

func (route UserAccountRouters) RegisterRouters() {
	apiPath := route.RouteGroup.Group("/api/user-account")
	authRoutes := AuthRoute{RouteGroup: apiPath, Handler: route.Handler}
	authRoutes.RegisterRouters()

	registrationRoutes := RegistrationRoute{RouteGroup: apiPath, Handler: route.Handler}
	registrationRoutes.RegisterRouters()
}

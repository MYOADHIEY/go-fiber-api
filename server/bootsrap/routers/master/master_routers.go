package master

import (
	ihandler "kbaa-fiber-api/server/handlers/base"

	"github.com/gofiber/fiber/v2"
)

type MasterRouters struct {
	RouteGroup fiber.Router
	Handler    ihandler.BaseHandler
}

func (route MasterRouters) RegisterRouters() {
	apiPath := route.RouteGroup.Group("/api/master")
	userRoutes := UserRoute{RouteGroup: apiPath, Handler: route.Handler}
	userRoutes.RegisterRouters()
}

package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nurhamsah1998/ppdb_be/controller"
	"github.com/nurhamsah1998/ppdb_be/service"
)

func RouteInit(app *fiber.App) {

	api := app.Group("/api")
	api.Get("/health", service.Health)
	controller.RegisterController(api.Group("/auth"))

}

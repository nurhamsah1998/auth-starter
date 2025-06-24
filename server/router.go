package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nurhamsah1998/ppdb_be/controller/auth"
	"github.com/nurhamsah1998/ppdb_be/service/health"
)

func RouteInit(app *fiber.App) {

	api := app.Group("/api")
	api.Get("/health", health.Health)
	auth.RegisterController(api.Group("/auth"))

}

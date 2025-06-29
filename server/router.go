package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nurhamsah1998/auth-starter/controller"
	"github.com/nurhamsah1998/auth-starter/service"
)

// / main router
func RouteInit(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/health", service.Health)
	/// inject controller ke RouteInit
	controller.AuthController(api.Group("/auth"))
	controller.ProfileController(api)
}

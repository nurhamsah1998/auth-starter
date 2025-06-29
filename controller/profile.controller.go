package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nurhamsah1998/auth-starter/internal/middleware"
	"github.com/nurhamsah1998/auth-starter/service"
)

// / profile controller.
func ProfileController(app fiber.Router) {
	service := service.ProfileHandler()
	/// inject service profile ke profile controller
	app.Get("/profile", middleware.Guard, service.MyProfile)
}

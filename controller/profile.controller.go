package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nurhamsah1998/auth-starter/internal/middleware"
	"github.com/nurhamsah1998/auth-starter/service/profile"
)

// / profile controller.
func ProfileController(app fiber.Router) {
	service := profile.ProfileHandler()
	/// inject service profile ke profile controller
	app.Get("/profile", middleware.Guard, service.MyProfile)
	app.Patch("/profile/:profile_id", middleware.Guard, service.UpdateProfile)
}

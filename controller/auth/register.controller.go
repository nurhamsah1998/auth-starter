package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nurhamsah1998/ppdb_be/service/auth"
)

func RegisterController(route fiber.Router) {
	service := auth.AuthHandler()
	route.Post("/register", service.Register)
}

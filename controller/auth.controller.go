package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nurhamsah1998/ppdb_be/service"
)

func RegisterController(route fiber.Router) {
	service := service.AuthHandler()
	route.Post("/register", service.Register)
	route.Post("/activation/:token_activation", service.Activation)
}

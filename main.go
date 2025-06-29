package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/nurhamsah1998/auth-starter/internal"
	"github.com/nurhamsah1998/auth-starter/server"
)

func main() {
	internal.DbGormInit()
	app := fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			var errorMessage string
			var errorStatus int
			if err.Error() != "" {
				errorMessage = err.Error()
				errorStatus = fiber.StatusBadRequest
			}
			return ctx.Status(errorStatus).JSON(fiber.Map{"error": true, "message": errorMessage})
		},
	})
	app.Use(recover.New())
	server.RouteInit(app)
	app.Listen(":3000")
}

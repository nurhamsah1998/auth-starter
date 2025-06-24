package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/nurhamsah1998/ppdb_be/server"
)

func main() {
	app := fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			log.Println(err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Internal server error."})
		},
	})
	app.Use(recover.New())
	server.RouteInit(app)
	app.Listen(":3000")
}

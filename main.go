package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/nurhamsah1998/auth-starter/internal"
	"github.com/nurhamsah1998/auth-starter/server"
)

func main() {
	/// konek ke database
	internal.DbGormInit()
	app := fiber.New(fiber.Config{
		/// handler untuk mengatasi semua jenis error (global)
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
	/// menghandle agar tidak terjadi server down ketika terjadi error secara fatal
	app.Use(recover.New())
	/// start server
	server.RouteInit(app)
	app.Listen(":3000")
}

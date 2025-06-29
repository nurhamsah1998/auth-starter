package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nurhamsah1998/auth-starter/service"
)

// / auth controller
func AuthController(route fiber.Router) {
	service := service.AuthHandler()
	/// inject service auth ke auth controller
	route.Post("/register", service.Register)
	route.Post("/login", service.Login)
	route.Post("/activation/:token_activation", service.Activation)
	/// forgot password, client mengirim email ke server untuk,
	/// meminta link reset password
	route.Post("/forgot-password", service.ForgotPassword)
	/// setelah mendapatkan link dari forgot password tadi,
	/// client bisa melakukan reset password pada end-point dibawah ini
	route.Post("/reset-password/:reset_pwd_token", service.ResetPassword)
}

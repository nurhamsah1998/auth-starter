package auth

import (
	"errors"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nurhamsah1998/auth-starter/internal"
	"github.com/nurhamsah1998/auth-starter/internal/model"
)

// / form validasi
type FormForgotPassword struct {
	Email string `json:"email" validate:"required,email"`
}

func (s *AuthService) ForgotPassword(c *fiber.Ctx) error {
	user := model.User{}
	bodyPayload := FormForgotPassword{}

	if err := c.BodyParser(&bodyPayload); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid body", "error": true})
	}
	/// validasi body payload yang dikirim oleh client (frontend)
	if err := internal.ClassValidate.Struct(bodyPayload); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = e.Tag()
		}
		return c.Status(422).JSON(fiber.Map{
			"message": errors,
			"error":   true,
		})
	}
	/// pencarian client/user berdasarkan email yang dikirim client/user
	res := internal.DB.First(&user, "email = ?", bodyPayload.Email)
	/// respon error ketika user/client tidak ditemukan
	if res.RowsAffected == 0 {
		return errors.New("user not found")
	}

	/// proses pembuatan akses token untuk reset password
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		/// token reset password akan expired/kadaluarsa dalam 168 jam kedepan,
		//  setelah berhasil mengirimkan email
		"exp": time.Now().Add(168 * time.Hour).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("RESET_PASSWORD_TOKEN"))) /// <--- secret key untuk reset password (mengambil dari file .env)
	///
	/// pada block dibawah ini bisa diimplementasikan dengan service pihak ketiga,
	/// untuk mengirimkan link reset password.
	//
	//
	// {masukan kodemu disini}
	//
	//
	/// contoh link yang akan diakses oleh user/client dari sisi frontend,
	/// ketika mendapatkan link reset password dari service pihak ketiga tadi : https://domainmu.com/auth/reset-password/{token-reset-password}.
	/// pada tampilan/UI kurang lebih bisa berisi form new password dan retype new password,
	/// ketika user/client click tombol reset password, maka akan hit api ke https://domainmu.com/api/auth/reset-password/{token-reset-password},
	/// untuk proses reset passwordnya
	///
	/// NOTE : Link diatas hanya contoh, bisa dicustom sendiri sesuai kebutuhan
	///
	return c.Status(200).JSON(fiber.Map{"message": "Successfully sending link to your email for reset password", "error": false, "data": tokenString})
}

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
	"golang.org/x/crypto/bcrypt"
)

type FormLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	user := model.User{}
	bodyPayload := FormLogin{}
	/// validasi format json
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
	/// proses pencarian data by email
	res := internal.DB.Preload("Profile").Find(&user, "email = ?", bodyPayload.Email)
	/// jika pencarian data by email di rables users tidak ditemukan
	if res.RowsAffected == 0 {
		return errors.New("invalid credential")
	}
	/// melakukan compare/perbandingan antara password yang dikirim client
	/// dan yang ada di database
	errPwd := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(bodyPayload.Password))
	/// error ketika password tidak sama
	if errPwd != nil {
		return errors.New("invalid credential")
	}
	/// error ketika client yang sudah daftar,
	/// tapi belum melakukan activasi mencoba untuk login.
	if user.Activation != "" {
		return errors.New("invalid credential")
	}
	/// proses pembuatan akses token untuk client
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		/// token JWT akan expired/kadaluarsa dalam 168 jam kedepan setelah berhasil login
		"exp": time.Now().Add(168 * time.Hour).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("ACCESS_TOKEN"))) /// <--- secret key untuk token login (mengambil dari file .env)
	data := fiber.Map{
		"token":   tokenString,
		"id":      user.ID,
		"email":   user.Email,
		"profile": user.Profile,
	}
	return c.Status(200).JSON(fiber.Map{"message": "Login successfully", "data": data})
}

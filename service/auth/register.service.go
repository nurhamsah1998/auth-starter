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

// / form validasi
type FormRegister struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required,numeric,min=10,max=15"`
	Password    string `json:"password" validate:"required,min=8,max=100"`
}

func (s *AuthService) Register(c *fiber.Ctx) error {
	user := model.User{}
	bodyPayload := FormRegister{}

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

	user.Email = bodyPayload.Email
	/// hashing password
	pwdHash, errPwdH := bcrypt.GenerateFromPassword([]byte(bodyPayload.Password), 10)
	if errPwdH != nil {
		return errors.New("failed hashing password")
	}
	user.Password = string(pwdHash)
	user.Profile.Name = bodyPayload.Name
	user.Profile.PhoneNumber = bodyPayload.PhoneNumber
	/// proses insert ke database
	res := internal.DB.Create(&user)
	/// jika terjadi error ketika proses insert
	/// contoh : error unique username
	if res.RowsAffected == 0 {
		return errors.New(res.Error.Error())
	}

	/// membuat token JWT untuk aktivasi akun
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":              user.ID,
		"email":           user.Email,
		"code_activation": internal.KeyGenerate(10),
		/// token aktivasi akan expired/kadaluarsa dalam 168 jam kedepan setelah berhasil register
		"exp": time.Now().Add(168 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("ACTIVATION_TOKEN"))) /// <--- secret key untuk token aktivasi (mengambil dari file .env)
	/// jika terjadi error ketika proses pembuatan token
	if err != nil {
		return errors.New("failed create token activation")
	}
	/// proses update token aktivasi di database
	internal.DB.Model(&user).Update("Activation", tokenString)

	return c.Status(201).JSON(fiber.Map{"message": "Register successfully. We send a unique code for activation", "data": bodyPayload})
}

package auth

import (
	"errors"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nurhamsah1998/auth-starter/internal"
	"github.com/nurhamsah1998/auth-starter/internal/model"
)

// / form validasi
type FormActivation struct {
	Activation string `json:"activation" validate:"required"`
}

func (s *AuthService) Activation(c *fiber.Ctx) error {
	user := model.User{}
	/// mengambil token dari param url
	activation := c.Params("token_activation")
	bodyPayload := FormActivation{}

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

	/// proses validasi user dan kode aktivasi
	token, err := jwt.Parse(activation, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("ACTIVATION_TOKEN")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return errors.New(err.Error())
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		id := claims["id"]
		code_activation := claims["code_activation"]
		res := internal.DB.First(&user, "id = ?", id)
		/// respon error ketika user/client tidak ditemukan
		if res.RowsAffected == 0 {
			return errors.New("user not found")
		}
		/// respon error ketika kode aktivasi yang dikirim client,
		/// tidak sama dengan kode aktivasi yang di database
		if code_activation != bodyPayload.Activation {
			return errors.New("invalid code activation")
		}
		/// respon error ketika client mencoba untuk aktivasi lagi setelah berhasil
		if user.Activation == "" {
			return errors.New("your account is already activated")

		}
	} else {
		return errors.New(err.Error())
	}
	/// proses update kolom "Activation" pada tabel users,
	/// set ke nil, menandakan client sudah melakukan aktivasi
	internal.DB.Model(&user).Update("Activation", nil)
	return c.Status(200).JSON(fiber.Map{"message": "Successfully activated account", "error": false})
}

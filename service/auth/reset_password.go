package auth

import (
	"errors"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nurhamsah1998/auth-starter/internal"
	"github.com/nurhamsah1998/auth-starter/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type FormResetPassword struct {
	NewPassword    string `json:"new_password" validate:"required,min=8,max=100"`
	ReTypePassword string `json:"retype_password" validate:"required,min=8,max=100"`
}

func (s *AuthService) ResetPassword(c *fiber.Ctx) error {
	user := model.User{}
	bodyPayload := FormResetPassword{}
	/// mengambil token dari param url
	resetPwdToken := c.Params("reset_pwd_token")
	/// validasi format json
	if err := c.BodyParser(&bodyPayload); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid body", "error": true})
	}

	/// proses validasi client
	token, err := jwt.Parse(resetPwdToken, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("RESET_PASSWORD_TOKEN")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return errors.New(err.Error())
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		email := claims["email"]
		res := internal.DB.First(&user, "email = ?", email)
		/// respon error ketika user/client tidak ditemukan
		if res.RowsAffected == 0 {
			return errors.New("cannot reset password")
		}
	} else {
		return errors.New(err.Error())
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
	/// respon error jika password baru dan retype password tidak sama
	if bodyPayload.NewPassword != bodyPayload.ReTypePassword {
		return errors.New("password not match")
	}

	/// hashing password baru
	pwdHash, errPwdH := bcrypt.GenerateFromPassword([]byte(bodyPayload.NewPassword), 10)
	if errPwdH != nil {
		return errors.New("failed hashing password")
	}
	/// proses update kolom "Password" pada tabel users berdasarkan,
	/// email yang dikirim user ketika forgot password
	internal.DB.Model(&user).Update("Password", string(pwdHash))
	return c.Status(200).JSON(fiber.Map{"message": "Reset password successfully", "error": false})
}

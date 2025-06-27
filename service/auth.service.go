package service

import (
	"errors"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nurhamsah1998/ppdb_be/internal"
	"github.com/nurhamsah1998/ppdb_be/internal/model"
	"github.com/nurhamsah1998/ppdb_be/pkg/utils"
)

type (
	AuthService  struct{}
	FormRegister struct {
		Name        string `json:"name" validate:"required"`
		Email       string `json:"email" validate:"required,email"`
		PhoneNumber string `json:"phone_number" validate:"required,numeric,min=10,max=15"`
		Password    string `json:"password" validate:"required,min=8,max=100"`
		SchoolName  string `json:"school_name" validate:"required"`
	}
	FormActivation struct {
		Activation string `json:"activation" validate:"required"`
	}
)

func AuthHandler() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Register(c *fiber.Ctx) error {
	user := model.User{}
	bodyPayload := FormRegister{}
	if err := c.BodyParser(&bodyPayload); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid body", "error": true})
	}

	if err := internal.ClassValidate.Struct(bodyPayload); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = e.Tag()
		}
		return c.Status(422).JSON(fiber.Map{
			"validation_errors": errors,
		})
	}

	user.Email = bodyPayload.Email
	user.Password = bodyPayload.Password
	user.Profile.Name = bodyPayload.Name
	user.Profile.SchoolName = bodyPayload.SchoolName
	user.Profile.PhoneNumber = bodyPayload.PhoneNumber
	res := internal.DB.Create(&user)
	if res.RowsAffected == 0 {
		return errors.New(res.Error.Error())
	}

	/// JWT SIGNING
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":              user.ID,
		"email":           user.Email,
		"code_activation": utils.KeyGenerate(10),
		"exp":             time.Now().Add(168 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("ACCESS_TOKEN")))
	///
	if err != nil {
		return errors.New("failed create token")
	}
	internal.DB.Model(&user).Update("Activation", tokenString)

	return c.Status(201).JSON(fiber.Map{"message": "Register successfully. We send a unique code for activation account", "data": bodyPayload})
}

func (s *AuthService) Activation(c *fiber.Ctx) error {
	user := model.User{}
	activation := c.Params("token_activation")
	bodyPayload := FormActivation{}
	if err := c.BodyParser(&bodyPayload); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid body", "error": true})
	}
	if err := internal.ClassValidate.Struct(bodyPayload); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = e.Tag()
		}
		return c.Status(422).JSON(fiber.Map{
			"validation_errors": errors,
		})
	}
	token, err := jwt.Parse(activation, func(token *jwt.Token) (interface{}, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("ACCESS_TOKEN")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return errors.New(err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		id := claims["id"]
		code_activation := claims["code_activation"]
		res := internal.DB.First(&user, "id = ?", id)
		if res.RowsAffected == 0 {
			return errors.New("user not found")
		}
		if code_activation != bodyPayload.Activation {
			return errors.New("invalid code activation")
		}
	} else {
		return errors.New(err.Error())
	}
	internal.DB.Model(&user).Update("Activation", nil)
	return c.Status(200).JSON(fiber.Map{"message": "Successfully activated account", "error": false})
}

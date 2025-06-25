package auth

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/nurhamsah1998/ppdb_be/internal"
	"github.com/nurhamsah1998/ppdb_be/internal/model"
)

type (
	AuthService struct{}
	Body        struct {
		Name        string `json:"name" validate:"required"`
		Email       string `json:"email" validate:"required,email"`
		PhoneNumber string `json:"phone_number" validate:"required,numeric,min=10,max=15"`
		Password    string `json:"password" validate:"required,min=8,max=100"`
		SchoolName  string `json:"school_name" validate:"required"`
	}
)

func AuthHandler() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Register(c *fiber.Ctx) error {
	user := model.User{}
	bodyPayload := Body{}
	if err := c.BodyParser(&bodyPayload); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid body", "error": false})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Fail send to db"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Register successfully", "data": &Body{}})
}

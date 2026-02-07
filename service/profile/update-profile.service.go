package profile

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/nurhamsah1998/auth-starter/internal"
	"github.com/nurhamsah1998/auth-starter/internal/middleware"
	"github.com/nurhamsah1998/auth-starter/internal/model"
)

type FormUpdateProfile struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	FullAddress string `json:"full_address"`
	PhoneNumber string `json:"phone_number" validate:"min=8,max=100"`
}

func (s *ProfileService) UpdateMyProfile(c *fiber.Ctx) error {
	userSession := c.Locals("user").(middleware.UserSession)
	profile := model.Profile{}
	bodyPayload := FormUpdateProfile{}
	internal.DB.First(&profile, "id = ?", userSession.ID)
	if err := c.BodyParser(&bodyPayload); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid body", "error": true})
	}

	// validation form value
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

	profile.PhoneNumber = bodyPayload.PhoneNumber
	profile.Name = bodyPayload.Name
	profile.FullAddress = bodyPayload.FullAddress

	res := internal.DB.Model(&model.Profile{}).Where("id = ?", userSession.ID).Updates(map[string]interface{}{
		"name":         profile.Name,
		"phone_number": profile.PhoneNumber,
		"full_address": profile.FullAddress,
	})
	/// jika terjadi error ketika proses update
	if res.RowsAffected == 0 {
		fmt.Println(res.Error.Error())
		return errors.New(res.Error.Error())
	}

	return c.Status(200).JSON(fiber.Map{"message": "Successfully updating profile", "error": false, "data": profile})
}

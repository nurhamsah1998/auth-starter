package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nurhamsah1998/auth-starter/internal/middleware"
)

type (
	ProfileService struct{}
)

// / service handler untuk menginject servis ke controller
func ProfileHandler() *ProfileService {
	return &ProfileService{}
}

func (s *ProfileService) MyProfile(c *fiber.Ctx) error {
	/// userCred kepanjangan dari userCredential, yang berisi data client,
	/// hasil dari decode akses token. berisi : id client dan email client,
	/// kurang lebih kalau di express JS seperti req.user.id / req.user.email
	userCred := c.Locals("user").(middleware.UserCredential)

	return c.Status(200).JSON(fiber.Map{"message": "OK", "error": false, "data": userCred})
}

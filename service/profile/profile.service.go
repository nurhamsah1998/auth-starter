package profile

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nurhamsah1998/auth-starter/internal"
	"github.com/nurhamsah1998/auth-starter/internal/middleware"
	"github.com/nurhamsah1998/auth-starter/internal/model"
)

func (s *ProfileService) MyProfile(c *fiber.Ctx) error {
	/// userSession berisi data client,
	/// hasil dari decode akses token. berisi : id client dan email client,
	/// kurang lebih kalau di express JS seperti req.user.id / req.user.email
	userSession := c.Locals("user").(middleware.UserSession)

	profile := model.Profile{}
	internal.DB.First(&profile, "id = ?", userSession.ID)

	return c.Status(200).JSON(fiber.Map{"message": "OK", "error": false, "data": profile})
}

package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserCredential struct {
	ID    float64
	Email string
}

// / guard (middleware) berfungsi untuk memvalidasi client ketika ingin mengakses resource,
// /  yang dilindungi (dibutuhkan token untuk mengakses resource tertentu)
func Guard(c *fiber.Ctx) error {
	/// get header dari client
	headerAuth := c.Get("Authorization")
	/// respon error jika client tidak menyertakan header "Authorization"
	if headerAuth == "" {
		return c.Status(401).JSON(fiber.Map{"error": true, "message": "header is not valid"})
	}
	stringToken := strings.Replace(headerAuth, "Bearer ", "", 1)
	/// decode token
	token, err := jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_TOKEN")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": true, "message": err.Error()})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userData := UserCredential{
			Email: claims["email"].(string),
			ID:    claims["id"].(float64),
		}
		/// set user/client data
		c.Locals("user", userData)
	} else {
		return c.Status(401).JSON(fiber.Map{"error": true, "message": err.Error()})
	}

	return c.Next()
}

package auth

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nurhamsah1998/auth-starter/internal"
	"github.com/nurhamsah1998/auth-starter/internal/model"
)

func (s *AuthService) RefreshToken(c *fiber.Ctx) error {
	user := model.User{}
	/// get header refresh token dari client
	headerAuth := c.Get("Authorization")
	/// respon error jika client tidak menyertakan header "Authorization"
	if headerAuth == "" {
		return c.Status(401).JSON(fiber.Map{"error": true, "message": "header is not valid"})
	}
	refreshToken := strings.Replace(headerAuth, "Refresh ", "", 1)
	/// pencarian user/client berdasarkan refresh token,
	/// yang tersimpan di database

	res := internal.DB.First(&user, "refresh_token = ?", refreshToken)
	/// respon error ketika user/client tidak ditemukan
	if res.RowsAffected == 0 {
		return c.Status(401).JSON(fiber.Map{"error": true, "message": "Unauthorized"})
	}
	/// proses validasi refresh token dan user/client
	token, err := jwt.Parse(user.RefreshToken, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("REFRESH_TOKEN")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		/// proses update refresh token diubah ke nil,
		//  jika ada kegagalan saat refresh token (token expired/kadaluarsa)
		internal.DB.Model(&user).Update("RefreshToken", nil)
		return c.Status(401).JSON(fiber.Map{"error": true, "message": err.Error()})
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		email := claims["email"]
		/// respon error ketika email user/client tidak sama dengan email yang berada di database
		if email != user.Email {
			/// proses update refresh token diubah ke nil,
			//  jika ada kegagalan saat refresh token
			internal.DB.Model(&user).Update("RefreshToken", nil)
			return c.Status(401).JSON(fiber.Map{"error": true, "message": "Unauthorized"})
		}
	} else {
		/// proses update refresh token diubah ke nil,
		//  jika ada kegagalan saat refresh token
		internal.DB.Model(&user).Update("RefreshToken", nil)
		return c.Status(401).JSON(fiber.Map{"error": true, "message": err.Error()})
	}
	/// proses pembuatan akses token baru untuk client
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		/// akses token JWT akan expired/kadaluarsa dalam 5 jam kedepan setelah berhasil login
		"exp": time.Now().Add(5 * time.Hour).Unix(),
	})
	accessTokenString, _ := accessToken.SignedString([]byte(os.Getenv("ACCESS_TOKEN"))) /// <--- secret key untuk token login (mengambil dari file .env)
	return c.Status(200).JSON(fiber.Map{"message": "Successfully create new access token", "data": fiber.Map{"refresh_token": accessTokenString}})
}

package middleware

import (
	"strings"

	"go-commerce/internal/handler/response"
	"go-commerce/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(jwtManager *jwt.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Unauthorized(c, "Authorization header required")
		}

		// Check Bearer format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return response.Unauthorized(c, "Invalid authorization format")
		}

		// Validate token
		claims, err := jwtManager.ValidateToken(tokenParts[1])
		if err != nil {
			return response.Unauthorized(c, "Invalid or expired token")
		}

		// Set user info in context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("is_admin", claims.IsAdmin)

		return c.Next()
	}
}

func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isAdmin, ok := c.Locals("is_admin").(bool)
		if !ok || !isAdmin {
			return response.Forbidden(c, "Admin access required")
		}
		return c.Next()
	}
}

func GetUserID(c *fiber.Ctx) uint64 {
	userID, _ := c.Locals("user_id").(uint64)
	return userID
}

func GetUserEmail(c *fiber.Ctx) string {
	email, _ := c.Locals("user_email").(string)
	return email
}

func IsAdmin(c *fiber.Ctx) bool {
	isAdmin, _ := c.Locals("is_admin").(bool)
	return isAdmin
}
package middleware

import (
	"strings"
	"time"

	"github.com/Marugo/birdlax/internal/shared/security"
	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ah := c.Get("Authorization")
		if ah == "" || !strings.HasPrefix(ah, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": "UNAUTHORIZED", "message": "missing bearer token"})
		}
		token := strings.TrimPrefix(ah, "Bearer ")
		claims, err := security.ParseAccess(token)
		if err != nil || claims.ExpiresAt.Time.Before(time.Now()) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": "UNAUTHORIZED", "message": "invalid or expired token"})
		}
		// inject context
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)
		c.Locals("emp_code", claims.Emp)
		return c.Next()
	}
}

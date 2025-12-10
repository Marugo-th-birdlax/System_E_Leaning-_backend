// internal/middlewares/roles.go
package middleware
        
import (
    "github.com/gofiber/fiber/v2"
    usermodels "github.com/Marugo/birdlax/internal/modules/user/models"
)

func RequireRoles(allowed ...usermodels.Role) fiber.Handler {
    allowedSet := map[usermodels.Role]struct{}{}
    for _, r := range allowed { allowedSet[r] = struct{}{} }

    return func(c *fiber.Ctx) error {
        // สมมติ auth middleware ยัด user role ลง c.Locals("role") เป็น string
        v, _ := c.Locals("role").(string)
        role := usermodels.Role(v)

        if _, ok := allowedSet[role]; !ok {
            return fiber.NewError(fiber.StatusForbidden, "forbidden")
        }
        return c.Next()
    }
}

// ตัวช่วย “ขั้นต่ำ”
func RequireAtLeast(min usermodels.Role) fiber.Handler {
    return func(c *fiber.Ctx) error {
        v, _ := c.Locals("role").(string)
        role := usermodels.Role(v)
        if !usermodels.IsAtLeast(role, min) {
            return fiber.NewError(fiber.StatusForbidden, "forbidden")
        }
        return c.Next()
    }
}

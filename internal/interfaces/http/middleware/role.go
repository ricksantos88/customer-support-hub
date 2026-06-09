package middleware

import (
	"github.com/gofiber/fiber/v3"
)

// RequireRole returns a middleware that allows access only if the authenticated agent
// has one of the specified roles. Must be applied after BearerAuthMiddleware.
func RequireRole(roles ...string) fiber.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(c fiber.Ctx) error {
		role, ok := c.Locals("agent_role").(string)
		if !ok || role == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "missing role context"})
		}
		if _, permitted := allowed[role]; !permitted {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "insufficient permissions"})
		}
		return c.Next()
	}
}

package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"

	"github.com/ricksantos88/customer-support-hub/internal/interfaces/http/middleware"
)

func TestRequireRole(t *testing.T) {
	app := fiber.New()

	app.Get("/admin-only",
		func(c fiber.Ctx) error {
			c.Locals("agent_role", "admin")
			return c.Next()
		},
		middleware.RequireRole("admin"),
		func(c fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		},
	)

	app.Get("/agent-forbidden",
		func(c fiber.Ctx) error {
			c.Locals("agent_role", "agent")
			return c.Next()
		},
		middleware.RequireRole("admin"),
		func(c fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		},
	)

	app.Get("/missing-role",
		middleware.RequireRole("admin"),
		func(c fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		},
	)

	t.Run("permitted role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("forbidden role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/agent-forbidden", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("missing role context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/missing-role", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}

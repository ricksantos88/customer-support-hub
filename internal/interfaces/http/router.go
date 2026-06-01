package http

import (
	"github.com/gofiber/fiber/v3"

	handlers "github.com/ricksantos88/customer-support-hub/internal/interfaces/http/handlers"
)

type RouterDependencies struct {
	AuthHandler    *handlers.AuthHandler
	AuthMiddleware fiber.Handler
}

func NewRouter(deps RouterDependencies) *fiber.App {
	app := fiber.New()

	app.Use(func(c fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		return c.Next()
	})

	app.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	auth := app.Group("/auth")
	auth.Post("/login", deps.AuthHandler.Login)
	auth.Post("/refresh", deps.AuthHandler.Refresh)
	auth.Use(deps.AuthMiddleware)
	auth.Post("/logout", deps.AuthHandler.Logout)

	return app
}

package http

import (
	"github.com/gofiber/fiber/v3"
)

func NewRouter() *fiber.App {
	app := fiber.New()

	app.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	return app
}

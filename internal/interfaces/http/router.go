package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/ricksantos88/swaggor"

	"github.com/ricksantos88/customer-support-hub/internal/interfaces/dto"
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

	setupSwagger(app)

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

func setupSwagger(app *fiber.App) {
	engine := swaggor.NewEngine("Customer Support Hub", "v1.0.0",
		swaggor.WithDescription("API for managing customer support conversations, agents, and sessions."),
		swaggor.WithServer("http://localhost:8080", "Local"),
		swaggor.WithSecurityScheme("bearer", swaggor.BearerJWT()),
	)

	engine.AddRoute("/auth/login", "POST", "Login", "Authenticate an agent with email and password and receive access and refresh tokens.",
		swaggor.WithTags("Auth"),
		swaggor.WithRequestBody("Agent credentials", true, dto.LoginRequest{}),
		swaggor.WithResponse(200, "Authenticated successfully", dto.AuthResponse{}),
		swaggor.WithResponse(400, "Invalid request body or missing fields", dto.ErrorResponse{}),
		swaggor.WithResponse(401, "Invalid credentials", dto.ErrorResponse{}),
	)

	engine.AddRoute("/auth/refresh", "POST", "Refresh Token", "Exchange a valid refresh token for a new access token.",
		swaggor.WithTags("Auth"),
		swaggor.WithRequestBody("Refresh token payload", true, dto.RefreshRequest{}),
		swaggor.WithResponse(200, "Token refreshed successfully", dto.AuthResponse{}),
		swaggor.WithResponse(400, "Missing or malformed refresh token", dto.ErrorResponse{}),
		swaggor.WithResponse(401, "Expired or invalid refresh token", dto.ErrorResponse{}),
	)

	engine.AddRoute("/auth/logout", "POST", "Logout", "Invalidate the current session. Requires a valid Bearer token.",
		swaggor.WithTags("Auth"),
		swaggor.WithSecurity("bearer"),
		swaggor.WithResponse(200, "Logged out successfully", dto.StatusResponse{}),
		swaggor.WithResponse(401, "Missing or invalid session", dto.ErrorResponse{}),
		swaggor.WithResponse(500, "Internal error while revoking session", dto.ErrorResponse{}),
	)

	app.All("/swaggor/*", adaptor.HTTPHandler(engine.Handler()))
}

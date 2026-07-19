package http

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/ricksantos88/swaggor"

	"github.com/ricksantos88/customer-support-hub/internal/interfaces/dto"
	handlers "github.com/ricksantos88/customer-support-hub/internal/interfaces/http/handlers"
	"github.com/ricksantos88/customer-support-hub/internal/interfaces/http/middleware"
)

type RouterDependencies struct {
	AuthHandler        *handlers.AuthHandler
	AdminHandler       *handlers.AdminHandler
	AuthMiddleware     fiber.Handler
	IPRateLimiter      fiber.Handler
	AgentRateLimiter   fiber.Handler
	CORSAllowedOrigins string
}

func NewRouter(deps RouterDependencies) *fiber.App {
	app := fiber.New()

	allowedOrigins := deps.CORSAllowedOrigins
	if strings.TrimSpace(allowedOrigins) == "" {
		allowedOrigins = "*"
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins:  strings.Split(allowedOrigins, ","),
		AllowMethods:  []string{fiber.MethodGet, fiber.MethodPost, fiber.MethodPut, fiber.MethodDelete, fiber.MethodOptions},
		AllowHeaders:  []string{"Authorization", "Content-Type", "Accept"},
		ExposeHeaders: []string{"Retry-After"},
	}))

	app.Use(func(c fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		return c.Next()
	})

	if deps.IPRateLimiter != nil {
		app.Use(deps.IPRateLimiter)
	}

	setupSwagger(app)

	app.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	auth := app.Group("/auth")
	auth.Post("/login", deps.AuthHandler.Login)
	auth.Post("/refresh", deps.AuthHandler.Refresh)
	auth.Use(deps.AuthMiddleware)
	if deps.AgentRateLimiter != nil {
		auth.Use(deps.AgentRateLimiter)
	}
	auth.Post("/logout", deps.AuthHandler.Logout)

	adminGroup := app.Group("/admin")
	adminGroup.Use(deps.AuthMiddleware)
	adminGroup.Use(RequireRole("admin"))
	if deps.AgentRateLimiter != nil {
		adminGroup.Use(deps.AgentRateLimiter)
	}
	adminGroup.Post("/agents", deps.AdminHandler.CreateAgent)
	adminGroup.Get("/agents", deps.AdminHandler.ListAgents)
	adminGroup.Put("/agents/:id", deps.AdminHandler.UpdateAgent)
	adminGroup.Delete("/agents/:id", deps.AdminHandler.DeleteAgent)
	adminGroup.Get("/sessions", deps.AdminHandler.ListActiveSessions)
	adminGroup.Delete("/sessions/:id", deps.AdminHandler.RevokeSession)
	adminGroup.Get("/status", deps.AdminHandler.GetSystemStatus)

	app.Get("/admin-panel", func(c fiber.Ctx) error {
		if !strings.HasSuffix(c.Path(), "/") {
			return c.Redirect().To("/admin-panel/")
		}
		return c.Next()
	})
	app.Get("/admin-panel/*", static.New("./web/admin"))

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

	engine.AddRoute("/admin/agents", "POST", "Create Agent", "Create a new agent. Requires a valid admin Bearer token.",
		swaggor.WithTags("Admin"),
		swaggor.WithSecurity("bearer"),
		swaggor.WithRequestBody("Agent info", true, dto.CreateAgentRequest{}),
		swaggor.WithResponse(201, "Agent created successfully", dto.AgentResponse{}),
		swaggor.WithResponse(400, "Invalid fields or parameters", dto.ErrorResponse{}),
		swaggor.WithResponse(401, "Unauthorized", dto.ErrorResponse{}),
		swaggor.WithResponse(403, "Forbidden", dto.ErrorResponse{}),
		swaggor.WithResponse(409, "Agent email already exists", dto.ErrorResponse{}),
	)

	engine.AddRoute("/admin/agents", "GET", "List Agents", "List registered agents with offset and limit pagination. Requires admin Bearer token.",
		swaggor.WithTags("Admin"),
		swaggor.WithSecurity("bearer"),
		swaggor.WithQueryParam("limit", "Default 10, max 100", false),
		swaggor.WithQueryParam("offset", "Default 0", false),
		swaggor.WithResponse(200, "Successfully retrieved agents list", []dto.AgentResponse{}),
		swaggor.WithResponse(401, "Unauthorized", dto.ErrorResponse{}),
		swaggor.WithResponse(403, "Forbidden", dto.ErrorResponse{}),
	)

	engine.AddRoute("/admin/agents/{id}", "PUT", "Update Agent", "Update an agent details. Requires admin Bearer token.",
		swaggor.WithTags("Admin"),
		swaggor.WithSecurity("bearer"),
		swaggor.WithPathParam("id", "Agent ID"),
		swaggor.WithRequestBody("Agent details to update", true, dto.UpdateAgentRequest{}),
		swaggor.WithResponse(200, "Agent updated successfully", dto.AgentResponse{}),
		swaggor.WithResponse(400, "Validation error or self-demotion attempt", dto.ErrorResponse{}),
		swaggor.WithResponse(401, "Unauthorized", dto.ErrorResponse{}),
		swaggor.WithResponse(403, "Forbidden", dto.ErrorResponse{}),
		swaggor.WithResponse(404, "Agent not found", dto.ErrorResponse{}),
	)

	engine.AddRoute("/admin/agents/{id}", "DELETE", "Delete Agent", "Inactivate and delete an agent, revoking all active sessions. Requires admin Bearer token.",
		swaggor.WithTags("Admin"),
		swaggor.WithSecurity("bearer"),
		swaggor.WithPathParam("id", "Agent ID"),
		swaggor.WithResponse(204, "Agent deleted successfully", nil),
		swaggor.WithResponse(400, "Cannot delete yourself", dto.ErrorResponse{}),
		swaggor.WithResponse(401, "Unauthorized", dto.ErrorResponse{}),
		swaggor.WithResponse(403, "Forbidden", dto.ErrorResponse{}),
		swaggor.WithResponse(404, "Agent not found", dto.ErrorResponse{}),
	)

	engine.AddRoute("/admin/sessions", "GET", "List Active Sessions", "List all active agent sessions in the system. Requires admin Bearer token.",
		swaggor.WithTags("Admin"),
		swaggor.WithSecurity("bearer"),
		swaggor.WithResponse(200, "Successfully retrieved sessions list", []dto.SessionResponse{}),
		swaggor.WithResponse(401, "Unauthorized", dto.ErrorResponse{}),
		swaggor.WithResponse(403, "Forbidden", dto.ErrorResponse{}),
	)

	engine.AddRoute("/admin/sessions/{id}", "DELETE", "Revoke Session", "Revoke an active session. Requires admin Bearer token.",
		swaggor.WithTags("Admin"),
		swaggor.WithSecurity("bearer"),
		swaggor.WithPathParam("id", "Session ID"),
		swaggor.WithResponse(204, "Session revoked successfully", nil),
		swaggor.WithResponse(401, "Unauthorized", dto.ErrorResponse{}),
		swaggor.WithResponse(403, "Forbidden", dto.ErrorResponse{}),
		swaggor.WithResponse(404, "Session not found", dto.ErrorResponse{}),
	)

	engine.AddRoute("/admin/status", "GET", "Get System Status", "Get status of external dependencies (DB, Redis, WhatsApp API). Requires admin Bearer token.",
		swaggor.WithTags("Admin"),
		swaggor.WithSecurity("bearer"),
		swaggor.WithResponse(200, "System status check", dto.SystemStatusResponse{}),
		swaggor.WithResponse(401, "Unauthorized", dto.ErrorResponse{}),
		swaggor.WithResponse(403, "Forbidden", dto.ErrorResponse{}),
	)

	app.All("/swaggor/*", adaptor.HTTPHandler(engine.Handler()))
}

// RequireRole returns a middleware that checks the agent_role local set by BearerAuthMiddleware.
func RequireRole(roles ...string) fiber.Handler {
	return middleware.RequireRole(roles...)
}

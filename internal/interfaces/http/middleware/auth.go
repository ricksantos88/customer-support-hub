package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/ricksantos88/customer-support-hub/internal/models"
	"github.com/ricksantos88/customer-support-hub/internal/repositories"
)

func ValidateJWT(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

type BearerAuthMiddleware struct {
	secret  string
	session repositories.SessionRepository
}

func NewBearerAuthMiddleware(secret string, session repositories.SessionRepository) fiber.Handler {
	return (&BearerAuthMiddleware{secret: secret, session: session}).Handle
}

func (m *BearerAuthMiddleware) Handle(c fiber.Ctx) error {
	tokenString, err := extractBearerToken(c.Get("Authorization"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	claims, err := ValidateJWT(tokenString, m.secret)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	if tokenType, _ := claims["token_type"].(string); tokenType != "access" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token type"})
	}

	agentID, err := claimUUID(claims, "agent_id")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid agent_id claim"})
	}

	sessionID, err := claimUUID(claims, "session_id")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid session_id claim"})
	}

	session, err := m.session.GetActiveByID(c.Context(), sessionID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "session is not active"})
	}

	if session.AgentID != agentID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "session does not belong to agent"})
	}

	c.Locals("agent_id", agentID.String())
	c.Locals("session_id", sessionID.String())
	c.Locals("session", session)
	return c.Next()
}

func extractBearerToken(header string) (string, error) {
	header = strings.TrimSpace(header)
	if header == "" {
		return "", fmt.Errorf("authorization header is required")
	}

	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", fmt.Errorf("authorization header must use Bearer token")
	}

	return strings.TrimSpace(parts[1]), nil
}

func claimUUID(claims jwt.MapClaims, key string) (uuid.UUID, error) {
	raw, ok := claims[key].(string)
	if !ok || strings.TrimSpace(raw) == "" {
		return uuid.Nil, fmt.Errorf("missing claim %s", key)
	}
	return uuid.Parse(raw)
}

func SessionFromContext(c fiber.Ctx) (*models.Session, bool) {
	session, ok := c.Locals("session").(*models.Session)
	return session, ok
}

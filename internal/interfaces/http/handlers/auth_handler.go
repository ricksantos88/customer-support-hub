package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/ricksantos88/customer-support-hub/internal/application/auth"
	"github.com/ricksantos88/customer-support-hub/internal/interfaces/dto"
)

type AuthHandler struct {
	service *auth.Service
}

func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req dto.LoginRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email and password are required"})
	}

	result, err := h.service.Login(context.Background(), req.Email, req.Password, c.Get("User-Agent"), c.IP())
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(dto.AuthResponse{
		TokenType:             auth.TokenTypeAccess,
		AccessToken:           result.AccessToken,
		RefreshToken:          result.RefreshToken,
		AccessTokenExpiresIn:  int64(time.Until(result.AccessTokenExpiresAt).Seconds()),
		RefreshTokenExpiresIn: int64(time.Until(result.RefreshTokenExpiresAt).Seconds()),
		SessionID:             result.SessionID.String(),
		AgentID:               result.AgentID.String(),
	})
}

func (h *AuthHandler) Refresh(c fiber.Ctx) error {
	var req dto.RefreshRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	req.RefreshToken = strings.TrimSpace(req.RefreshToken)
	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "refresh_token is required"})
	}

	result, err := h.service.Refresh(context.Background(), req.RefreshToken, c.Get("User-Agent"), c.IP())
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(dto.AuthResponse{
		TokenType:             auth.TokenTypeAccess,
		AccessToken:           result.AccessToken,
		RefreshToken:          result.RefreshToken,
		AccessTokenExpiresIn:  int64(time.Until(result.AccessTokenExpiresAt).Seconds()),
		RefreshTokenExpiresIn: int64(time.Until(result.RefreshTokenExpiresAt).Seconds()),
		SessionID:             result.SessionID.String(),
		AgentID:               result.AgentID.String(),
	})
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	sessionIDValue, ok := c.Locals("session_id").(string)
	if !ok || strings.TrimSpace(sessionIDValue) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing session context"})
	}

	sessionID, err := uuid.Parse(sessionIDValue)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": fmt.Sprintf("invalid session_id: %v", err)})
	}

	if err := h.service.Logout(context.Background(), sessionID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(dto.StatusResponse{Status: "logged_out"})
}

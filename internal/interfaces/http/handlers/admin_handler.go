package handlers

import (
	"encoding/json"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/ricksantos88/customer-support-hub/internal/application/admin"
	"github.com/ricksantos88/customer-support-hub/internal/interfaces/dto"
)

type AdminHandler struct {
	service *admin.Service
}

func NewAdminHandler(service *admin.Service) *AdminHandler {
	return &AdminHandler{service: service}
}

func (h *AdminHandler) CreateAgent(c fiber.Ctx) error {
	var req dto.CreateAgentRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "invalid request body"})
	}

	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "name, email, and password are required"})
	}

	agent, err := h.service.CreateAgent(c.Context(), req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		slog.Error("admin failed to create agent", "err", err)
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique") {
			return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{Error: "email already exists"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.AgentResponse{
		ID:         agent.ID.String(),
		Name:       agent.Name,
		Email:      agent.Email,
		Role:       agent.Role,
		CreatedAt:  agent.CreatedAt,
		LastActive: agent.LastActive,
	})
}

func (h *AdminHandler) ListAgents(c fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	agents, err := h.service.ListAgents(c.Context(), limit, offset)
	if err != nil {
		slog.Error("admin failed to list agents", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "failed to list agents"})
	}

	res := make([]dto.AgentResponse, len(agents))
	for i, agent := range agents {
		res[i] = dto.AgentResponse{
			ID:         agent.ID.String(),
			Name:       agent.Name,
			Email:      agent.Email,
			Role:       agent.Role,
			CreatedAt:  agent.CreatedAt,
			LastActive: agent.LastActive,
		}
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *AdminHandler) UpdateAgent(c fiber.Ctx) error {
	adminIDStr, ok := c.Locals("agent_id").(string)
	if !ok || adminIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{Error: "missing agent context"})
	}
	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{Error: "invalid admin session"})
	}

	targetAgentIDStr := c.Params("id")
	targetAgentID, err := uuid.Parse(targetAgentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "invalid agent id parameter"})
	}

	var req dto.UpdateAgentRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "invalid request body"})
	}

	agent, err := h.service.UpdateAgent(c.Context(), adminID, targetAgentID, req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		slog.Error("admin failed to update agent", "err", err)
		if strings.Contains(err.Error(), "demote") {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: err.Error()})
		}
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "agent not found"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(dto.AgentResponse{
		ID:         agent.ID.String(),
		Name:       agent.Name,
		Email:      agent.Email,
		Role:       agent.Role,
		CreatedAt:  agent.CreatedAt,
		LastActive: agent.LastActive,
	})
}

func (h *AdminHandler) DeleteAgent(c fiber.Ctx) error {
	adminIDStr, ok := c.Locals("agent_id").(string)
	if !ok || adminIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{Error: "missing agent context"})
	}
	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{Error: "invalid admin session"})
	}

	targetAgentIDStr := c.Params("id")
	targetAgentID, err := uuid.Parse(targetAgentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "invalid agent id parameter"})
	}

	err = h.service.DeleteAgent(c.Context(), adminID, targetAgentID)
	if err != nil {
		slog.Error("admin failed to delete agent", "err", err)
		if strings.Contains(err.Error(), "delete yourself") {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: err.Error()})
		}
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "agent not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AdminHandler) ListActiveSessions(c fiber.Ctx) error {
	sessions, err := h.service.ListActiveSessions(c.Context())
	if err != nil {
		slog.Error("admin failed to list active sessions", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "failed to list active sessions"})
	}

	res := make([]dto.SessionResponse, len(sessions))
	for i, sess := range sessions {
		res[i] = dto.SessionResponse{
			ID:        sess.ID.String(),
			AgentID:   sess.AgentID.String(),
			AgentName: sess.Agent.Name,
			IPAddress: sess.IPAddress,
			UserAgent: sess.UserAgent,
			CreatedAt: sess.CreatedAt,
			ExpiresAt: sess.ExpiresAt,
		}
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *AdminHandler) RevokeSession(c fiber.Ctx) error {
	sessionIDStr := c.Params("id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "invalid session id parameter"})
	}

	err = h.service.RevokeSession(c.Context(), sessionID)
	if err != nil {
		slog.Error("admin failed to revoke session", "err", err)
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "session not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "failed to revoke session"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AdminHandler) GetSystemStatus(c fiber.Ctx) error {
	status, err := h.service.GetSystemStatus(c.Context())
	if err != nil {
		slog.Error("admin failed to get system status", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "failed to check system status"})
	}

	return c.Status(fiber.StatusOK).JSON(dto.SystemStatusResponse{
		Database:            status.Database,
		Redis:               status.Redis,
		WhatsAppIntegration: status.WhatsAppIntegration,
	})
}

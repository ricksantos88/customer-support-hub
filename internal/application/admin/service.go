package admin

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ricksantos88/customer-support-hub/internal/models"
	"github.com/ricksantos88/customer-support-hub/internal/repositories"
)

type SessionCache interface {
	DeleteSession(ctx context.Context, id uuid.UUID) error
	DeleteRefreshTokenSessionID(ctx context.Context, refreshTokenHash string) error
	Ping(ctx context.Context) error
}

type Service struct {
	agents              repositories.AgentRepository
	sessions            repositories.SessionRepository
	cache               SessionCache
	dbPing              func(context.Context) error
	whatsAppConfigCheck func() bool
	now                 func() time.Time
}

func NewService(
	agents repositories.AgentRepository,
	sessions repositories.SessionRepository,
	cache SessionCache,
	dbPing func(context.Context) error,
	whatsAppConfigCheck func() bool,
) *Service {
	return &Service{
		agents:              agents,
		sessions:            sessions,
		cache:               cache,
		dbPing:              dbPing,
		whatsAppConfigCheck: whatsAppConfigCheck,
		now:                 func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) CreateAgent(ctx context.Context, name, email, password, role string) (*models.Agent, error) {
	agent := &models.Agent{
		Name:     strings.TrimSpace(name),
		Email:    strings.TrimSpace(email),
		Password: password,
		Role:     strings.ToLower(strings.TrimSpace(role)),
	}

	if err := s.agents.Create(ctx, agent); err != nil {
		return nil, err
	}

	return agent, nil
}

func (s *Service) ListAgents(ctx context.Context, limit, offset int) ([]models.Agent, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.agents.List(ctx, limit, offset)
}

func (s *Service) UpdateAgent(ctx context.Context, adminID, targetAgentID uuid.UUID, name, email, password, role string) (*models.Agent, error) {
	agent, err := s.agents.GetByID(ctx, targetAgentID)
	if err != nil {
		return nil, err
	}

	role = strings.ToLower(strings.TrimSpace(role))
	if adminID == targetAgentID && role != "" && role != models.AgentRoleAdmin {
		return nil, fmt.Errorf("cannot demote yourself")
	}

	if strings.TrimSpace(name) != "" {
		agent.Name = strings.TrimSpace(name)
	}
	if strings.TrimSpace(email) != "" {
		agent.Email = strings.TrimSpace(email)
	}
	if strings.TrimSpace(password) != "" {
		agent.Password = password
	}
	if role != "" {
		agent.Role = role
	}

	if err := s.agents.Update(ctx, agent); err != nil {
		return nil, err
	}

	return agent, nil
}

func (s *Service) DeleteAgent(ctx context.Context, adminID, targetAgentID uuid.UUID) error {
	if adminID == targetAgentID {
		return fmt.Errorf("cannot delete yourself")
	}

	// 1. Invalidate all active sessions for this agent in DB
	now := s.now()
	if err := s.sessions.RevokeAllByAgentID(ctx, targetAgentID, now); err != nil {
		return err
	}

	// 2. Invalidate all cached sessions of this agent in Redis
	activeSessions, err := s.sessions.ListActive(ctx)
	if err == nil {
		for _, sess := range activeSessions {
			if sess.AgentID == targetAgentID {
				_ = s.cache.DeleteSession(ctx, sess.ID)
				_ = s.cache.DeleteRefreshTokenSessionID(ctx, sess.RefreshTokenHash)
			}
		}
	}

	// 3. Delete Agent (soft-delete via GORM)
	return s.agents.Delete(ctx, targetAgentID)
}

func (s *Service) ListActiveSessions(ctx context.Context) ([]models.Session, error) {
	return s.sessions.ListActive(ctx)
}

func (s *Service) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	session, err := s.sessions.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}

	now := s.now()
	if err := s.sessions.Revoke(ctx, sessionID, now); err != nil {
		return err
	}

	_ = s.cache.DeleteSession(ctx, sessionID)
	_ = s.cache.DeleteRefreshTokenSessionID(ctx, session.RefreshTokenHash)

	return nil
}

type SystemStatus struct {
	Database            string `json:"database"`
	Redis               string `json:"redis"`
	WhatsAppIntegration string `json:"whatsapp_integration"`
}

func (s *Service) GetSystemStatus(ctx context.Context) (*SystemStatus, error) {
	dbStatus := "connected"
	if err := s.dbPing(ctx); err != nil {
		dbStatus = "disconnected"
	}

	redisStatus := "connected"
	if err := s.cache.Ping(ctx); err != nil {
		redisStatus = "disconnected"
	}

	waStatus := "not_configured"
	if s.whatsAppConfigCheck != nil && s.whatsAppConfigCheck() {
		waStatus = "configured"
	}

	return &SystemStatus{
		Database:            dbStatus,
		Redis:               redisStatus,
		WhatsAppIntegration: waStatus,
	}, nil
}

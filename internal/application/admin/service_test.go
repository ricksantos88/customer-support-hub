package admin_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ricksantos88/customer-support-hub/internal/application/admin"
	"github.com/ricksantos88/customer-support-hub/internal/models"
	"github.com/ricksantos88/customer-support-hub/internal/repositories/mocks"
)

type SessionCacheMock struct {
	mock.Mock
}

func (m *SessionCacheMock) DeleteSession(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *SessionCacheMock) DeleteRefreshTokenSessionID(ctx context.Context, refreshTokenHash string) error {
	args := m.Called(ctx, refreshTokenHash)
	return args.Error(0)
}

func (m *SessionCacheMock) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestAdminService_CreateAgent(t *testing.T) {
	agentRepo := new(mocks.AgentRepositoryMock)
	sessionRepo := new(mocks.SessionRepositoryMock)
	cacheMock := new(SessionCacheMock)

	svc := admin.NewService(agentRepo, sessionRepo, cacheMock, nil, nil)

	t.Run("success", func(t *testing.T) {
		agentRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Agent")).Return(nil).Once()

		res, err := svc.CreateAgent(context.Background(), "Alice", "alice@test.com", "password123", "admin")
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "Alice", res.Name)
		assert.Equal(t, "alice@test.com", res.Email)
		assert.Equal(t, "admin", res.Role)
		agentRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		agentRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Agent")).Return(errors.New("db error")).Once()

		res, err := svc.CreateAgent(context.Background(), "Alice", "alice@test.com", "password123", "admin")
		assert.Error(t, err)
		assert.Nil(t, res)
		agentRepo.AssertExpectations(t)
	})
}

func TestAdminService_ListAgents(t *testing.T) {
	agentRepo := new(mocks.AgentRepositoryMock)
	sessionRepo := new(mocks.SessionRepositoryMock)
	cacheMock := new(SessionCacheMock)

	svc := admin.NewService(agentRepo, sessionRepo, cacheMock, nil, nil)

	t.Run("success", func(t *testing.T) {
		expected := []models.Agent{
			{ID: uuid.New(), Name: "Alice"},
			{ID: uuid.New(), Name: "Bob"},
		}
		agentRepo.On("List", mock.Anything, 10, 0).Return(expected, nil).Once()

		res, err := svc.ListAgents(context.Background(), 10, 0)
		assert.NoError(t, err)
		assert.Len(t, res, 2)
		assert.Equal(t, "Alice", res[0].Name)
		agentRepo.AssertExpectations(t)
	})
}

func TestAdminService_UpdateAgent(t *testing.T) {
	agentRepo := new(mocks.AgentRepositoryMock)
	sessionRepo := new(mocks.SessionRepositoryMock)
	cacheMock := new(SessionCacheMock)

	svc := admin.NewService(agentRepo, sessionRepo, cacheMock, nil, nil)
	adminID := uuid.New()
	targetID := uuid.New()

	t.Run("success", func(t *testing.T) {
		existing := &models.Agent{ID: targetID, Name: "Alice", Email: "alice@test.com", Role: "agent"}
		agentRepo.On("GetByID", mock.Anything, targetID).Return(existing, nil).Once()
		agentRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Agent")).Return(nil).Once()

		res, err := svc.UpdateAgent(context.Background(), adminID, targetID, "Alice Updated", "", "", "admin")
		assert.NoError(t, err)
		assert.Equal(t, "Alice Updated", res.Name)
		assert.Equal(t, "admin", res.Role)
		agentRepo.AssertExpectations(t)
	})

	t.Run("self demotion check", func(t *testing.T) {
		existing := &models.Agent{ID: adminID, Name: "Admin", Email: "admin@test.com", Role: "admin"}
		agentRepo.On("GetByID", mock.Anything, adminID).Return(existing, nil).Once()

		res, err := svc.UpdateAgent(context.Background(), adminID, adminID, "", "", "", "agent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot demote yourself")
		assert.Nil(t, res)
		agentRepo.AssertExpectations(t)
	})
}

func TestAdminService_DeleteAgent(t *testing.T) {
	agentRepo := new(mocks.AgentRepositoryMock)
	sessionRepo := new(mocks.SessionRepositoryMock)
	cacheMock := new(SessionCacheMock)

	svc := admin.NewService(agentRepo, sessionRepo, cacheMock, nil, nil)
	adminID := uuid.New()
	targetID := uuid.New()

	t.Run("cannot delete self", func(t *testing.T) {
		err := svc.DeleteAgent(context.Background(), adminID, adminID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot delete yourself")
	})

	t.Run("success with session and cache revocation", func(t *testing.T) {
		sessionID := uuid.New()
		activeSessions := []models.Session{
			{ID: sessionID, AgentID: targetID, RefreshTokenHash: "hash123"},
			{ID: uuid.New(), AgentID: uuid.New(), RefreshTokenHash: "other"},
		}

		sessionRepo.On("RevokeAllByAgentID", mock.Anything, targetID, mock.Anything).Return(nil).Once()
		sessionRepo.On("ListActive", mock.Anything).Return(activeSessions, nil).Once()
		cacheMock.On("DeleteSession", mock.Anything, sessionID).Return(nil).Once()
		cacheMock.On("DeleteRefreshTokenSessionID", mock.Anything, "hash123").Return(nil).Once()
		agentRepo.On("Delete", mock.Anything, targetID).Return(nil).Once()

		err := svc.DeleteAgent(context.Background(), adminID, targetID)
		assert.NoError(t, err)

		sessionRepo.AssertExpectations(t)
		cacheMock.AssertExpectations(t)
		agentRepo.AssertExpectations(t)
	})
}

func TestAdminService_RevokeSession(t *testing.T) {
	agentRepo := new(mocks.AgentRepositoryMock)
	sessionRepo := new(mocks.SessionRepositoryMock)
	cacheMock := new(SessionCacheMock)

	svc := admin.NewService(agentRepo, sessionRepo, cacheMock, nil, nil)
	sessionID := uuid.New()

	t.Run("success", func(t *testing.T) {
		session := &models.Session{ID: sessionID, RefreshTokenHash: "hash123"}
		sessionRepo.On("GetByID", mock.Anything, sessionID).Return(session, nil).Once()
		sessionRepo.On("Revoke", mock.Anything, sessionID, mock.Anything).Return(nil).Once()
		cacheMock.On("DeleteSession", mock.Anything, sessionID).Return(nil).Once()
		cacheMock.On("DeleteRefreshTokenSessionID", mock.Anything, "hash123").Return(nil).Once()

		err := svc.RevokeSession(context.Background(), sessionID)
		assert.NoError(t, err)

		sessionRepo.AssertExpectations(t)
		cacheMock.AssertExpectations(t)
	})
}

func TestAdminService_GetSystemStatus(t *testing.T) {
	agentRepo := new(mocks.AgentRepositoryMock)
	sessionRepo := new(mocks.SessionRepositoryMock)
	cacheMock := new(SessionCacheMock)

	t.Run("all connected", func(t *testing.T) {
		dbPing := func(ctx context.Context) error { return nil }
		waCheck := func() bool { return true }
		cacheMock.On("Ping", mock.Anything).Return(nil).Once()

		svc := admin.NewService(agentRepo, sessionRepo, cacheMock, dbPing, waCheck)
		res, err := svc.GetSystemStatus(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, "connected", res.Database)
		assert.Equal(t, "connected", res.Redis)
		assert.Equal(t, "configured", res.WhatsAppIntegration)
		cacheMock.AssertExpectations(t)
	})

	t.Run("all failed", func(t *testing.T) {
		dbPing := func(ctx context.Context) error { return errors.New("failed") }
		waCheck := func() bool { return false }
		cacheMock.On("Ping", mock.Anything).Return(errors.New("failed")).Once()

		svc := admin.NewService(agentRepo, sessionRepo, cacheMock, dbPing, waCheck)
		res, err := svc.GetSystemStatus(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, "disconnected", res.Database)
		assert.Equal(t, "disconnected", res.Redis)
		assert.Equal(t, "not_configured", res.WhatsAppIntegration)
		cacheMock.AssertExpectations(t)
	})
}

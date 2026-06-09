package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/ricksantos88/customer-support-hub/internal/models"
	"github.com/ricksantos88/customer-support-hub/internal/repositories/mocks"
)

type sessionCacheStub struct {
	setSessionFn                     func(context.Context, *models.Session, time.Duration) error
	getSessionByIDFn                 func(context.Context, uuid.UUID) (*models.Session, error)
	deleteSessionFn                  func(context.Context, uuid.UUID) error
	setRefreshTokenSessionIDFn       func(context.Context, string, uuid.UUID, time.Duration) error
	getSessionIDByRefreshTokenHashFn func(context.Context, string) (uuid.UUID, error)
	deleteRefreshTokenSessionIDFn    func(context.Context, string) error
}

func (s *sessionCacheStub) SetSession(ctx context.Context, session *models.Session, ttl time.Duration) error {
	if s.setSessionFn != nil {
		return s.setSessionFn(ctx, session, ttl)
	}
	return nil
}

func (s *sessionCacheStub) GetSessionByID(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	if s.getSessionByIDFn != nil {
		return s.getSessionByIDFn(ctx, id)
	}
	return nil, nil
}

func (s *sessionCacheStub) DeleteSession(ctx context.Context, id uuid.UUID) error {
	if s.deleteSessionFn != nil {
		return s.deleteSessionFn(ctx, id)
	}
	return nil
}

func (s *sessionCacheStub) SetRefreshTokenSessionID(ctx context.Context, refreshTokenHash string, sessionID uuid.UUID, ttl time.Duration) error {
	if s.setRefreshTokenSessionIDFn != nil {
		return s.setRefreshTokenSessionIDFn(ctx, refreshTokenHash, sessionID, ttl)
	}
	return nil
}

func (s *sessionCacheStub) GetSessionIDByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (uuid.UUID, error) {
	if s.getSessionIDByRefreshTokenHashFn != nil {
		return s.getSessionIDByRefreshTokenHashFn(ctx, refreshTokenHash)
	}
	return uuid.Nil, errors.New("cache miss")
}

func (s *sessionCacheStub) DeleteRefreshTokenSessionID(ctx context.Context, refreshTokenHash string) error {
	if s.deleteRefreshTokenSessionIDFn != nil {
		return s.deleteRefreshTokenSessionIDFn(ctx, refreshTokenHash)
	}
	return nil
}

func TestService_Login(t *testing.T) {
	agentRepo := &mocks.AgentRepositoryMock{}
	sessionRepo := &mocks.SessionRepositoryMock{}
	cache := &sessionCacheStub{}

	agentID := uuid.New()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	require.NoError(t, err)
	agent := &models.Agent{ID: agentID, Email: "agent@example.com", PasswordHash: string(passwordHash)}
	agentRepo.On("GetByEmail", mock.Anything, "agent@example.com").Return(agent, nil)
	agentRepo.On("UpdateLastActive", mock.Anything, agentID, mock.AnythingOfType("time.Time")).Return(nil)

	sessionRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Session")).Return(nil).Run(func(args mock.Arguments) {
		session := args.Get(1).(*models.Session)
		session.ID = uuid.New()
	})

	service := NewService(agentRepo, sessionRepo, cache, "secret", "customer-support-hub", time.Minute, time.Hour, time.Hour)
	result, err := service.Login(context.Background(), "agent@example.com", "password", "agent/1.0", "127.0.0.1")
	require.NoError(t, err)
	require.NotEmpty(t, result.AccessToken)
	require.NotEmpty(t, result.RefreshToken)
	require.Equal(t, agentID, result.AgentID)
}

func TestService_Refresh(t *testing.T) {
	agentRepo := &mocks.AgentRepositoryMock{}
	sessionRepo := &mocks.SessionRepositoryMock{}
	cache := &sessionCacheStub{}

	agentID := uuid.New()
	sessionID := uuid.New()
	refreshToken := randomRefreshTokenForTest(t)
	refreshHash := hashToken(refreshToken)

	session := &models.Session{ID: sessionID, AgentID: agentID, RefreshTokenHash: refreshHash, ExpiresAt: time.Now().UTC().Add(time.Hour)}
	sessionRepo.On("GetByRefreshTokenHash", mock.Anything, refreshHash).Return(session, nil)
	sessionRepo.On("RotateRefreshToken", mock.Anything, sessionID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(nil)
	agentRepo.On("GetByID", mock.Anything, agentID).Return(&models.Agent{ID: agentID, Email: "agent@example.com"}, nil)
	agentRepo.On("UpdateLastActive", mock.Anything, agentID, mock.AnythingOfType("time.Time")).Return(nil)

	service := NewService(agentRepo, sessionRepo, cache, "secret", "customer-support-hub", time.Minute, time.Hour, time.Hour)
	result, err := service.Refresh(context.Background(), refreshToken, "agent/1.0", "127.0.0.1")
	require.NoError(t, err)
	require.NotEmpty(t, result.AccessToken)
	require.NotEmpty(t, result.RefreshToken)
	require.Equal(t, agentID, result.AgentID)
}

func TestService_Logout(t *testing.T) {
	agentRepo := &mocks.AgentRepositoryMock{}
	sessionRepo := &mocks.SessionRepositoryMock{}
	cache := &sessionCacheStub{}

	sessionID := uuid.New()
	sessionRepo.On("Revoke", mock.Anything, sessionID, mock.AnythingOfType("time.Time")).Return(nil)

	service := NewService(agentRepo, sessionRepo, cache, "secret", "customer-support-hub", time.Minute, time.Hour, time.Hour)
	require.NoError(t, service.Logout(context.Background(), sessionID))
}

func randomRefreshTokenForTest(t *testing.T) string {
	t.Helper()
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	require.NoError(t, err)
	return base64.RawURLEncoding.EncodeToString(bytes)
}

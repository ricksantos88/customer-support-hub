package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/ricksantos88/customer-support-hub/internal/models"
)

type SessionRepositoryMock struct {
	mock.Mock
}

func (m *SessionRepositoryMock) Create(ctx context.Context, session *models.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *SessionRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	args := m.Called(ctx, id)
	if session, ok := args.Get(0).(*models.Session); ok {
		return session, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SessionRepositoryMock) GetByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*models.Session, error) {
	args := m.Called(ctx, refreshTokenHash)
	if session, ok := args.Get(0).(*models.Session); ok {
		return session, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SessionRepositoryMock) GetActiveByID(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	args := m.Called(ctx, id)
	if session, ok := args.Get(0).(*models.Session); ok {
		return session, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SessionRepositoryMock) RotateRefreshToken(ctx context.Context, sessionID uuid.UUID, refreshTokenHash string, lastUsedAt, expiresAt time.Time) error {
	args := m.Called(ctx, sessionID, refreshTokenHash, lastUsedAt, expiresAt)
	return args.Error(0)
}

func (m *SessionRepositoryMock) Revoke(ctx context.Context, sessionID uuid.UUID, revokedAt time.Time) error {
	args := m.Called(ctx, sessionID, revokedAt)
	return args.Error(0)
}

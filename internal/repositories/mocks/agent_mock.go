package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/ricksantos88/customer-support-hub/internal/models"
)

type AgentRepositoryMock struct {
	mock.Mock
}

func (m *AgentRepositoryMock) Create(ctx context.Context, agent *models.Agent) error {
	args := m.Called(ctx, agent)
	return args.Error(0)
}

func (m *AgentRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*models.Agent, error) {
	args := m.Called(ctx, id)
	if agent, ok := args.Get(0).(*models.Agent); ok {
		return agent, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *AgentRepositoryMock) GetByEmail(ctx context.Context, email string) (*models.Agent, error) {
	args := m.Called(ctx, email)
	if agent, ok := args.Get(0).(*models.Agent); ok {
		return agent, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *AgentRepositoryMock) Update(ctx context.Context, agent *models.Agent) error {
	args := m.Called(ctx, agent)
	return args.Error(0)
}

func (m *AgentRepositoryMock) UpdateLastActive(ctx context.Context, agentID uuid.UUID, lastActive time.Time) error {
	args := m.Called(ctx, agentID, lastActive)
	return args.Error(0)
}

func (m *AgentRepositoryMock) List(ctx context.Context, limit, offset int) ([]models.Agent, error) {
	args := m.Called(ctx, limit, offset)
	if agents, ok := args.Get(0).([]models.Agent); ok {
		return agents, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *AgentRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

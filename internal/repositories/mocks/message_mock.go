package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/ricksantos88/customer-support-hub/internal/models"
)

type MessageRepositoryMock struct {
	mock.Mock
}

func (m *MessageRepositoryMock) Create(ctx context.Context, message *models.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MessageRepositoryMock) GetByConversation(ctx context.Context, conversationID uuid.UUID) ([]models.Message, error) {
	args := m.Called(ctx, conversationID)
	if messages, ok := args.Get(0).([]models.Message); ok {
		return messages, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MessageRepositoryMock) GetLastMessage(ctx context.Context, conversationID uuid.UUID) (*models.Message, error) {
	args := m.Called(ctx, conversationID)
	if message, ok := args.Get(0).(*models.Message); ok {
		return message, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MessageRepositoryMock) ListByDirection(ctx context.Context, conversationID uuid.UUID, direction string) ([]models.Message, error) {
	args := m.Called(ctx, conversationID, direction)
	if messages, ok := args.Get(0).([]models.Message); ok {
		return messages, args.Error(1)
	}
	return nil, args.Error(1)
}

package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/ricksantos88/customer-support-hub/internal/models"
)

type ConversationRepositoryMock struct {
	mock.Mock
}

func (m *ConversationRepositoryMock) Create(ctx context.Context, conversation *models.Conversation) error {
	args := m.Called(ctx, conversation)
	return args.Error(0)
}

func (m *ConversationRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error) {
	args := m.Called(ctx, id)
	if conversation, ok := args.Get(0).(*models.Conversation); ok {
		return conversation, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ConversationRepositoryMock) GetActiveByContact(ctx context.Context, contactID uuid.UUID) (*models.Conversation, error) {
	args := m.Called(ctx, contactID)
	if conversation, ok := args.Get(0).(*models.Conversation); ok {
		return conversation, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ConversationRepositoryMock) AssignAgent(ctx context.Context, conversationID, agentID uuid.UUID) error {
	args := m.Called(ctx, conversationID, agentID)
	return args.Error(0)
}

func (m *ConversationRepositoryMock) UpdateStatus(ctx context.Context, conversationID uuid.UUID, status string) error {
	args := m.Called(ctx, conversationID, status)
	return args.Error(0)
}

func (m *ConversationRepositoryMock) Close(ctx context.Context, conversationID uuid.UUID) error {
	args := m.Called(ctx, conversationID)
	return args.Error(0)
}

func (m *ConversationRepositoryMock) ListByAgent(ctx context.Context, agentID uuid.UUID) ([]models.Conversation, error) {
	args := m.Called(ctx, agentID)
	if conversations, ok := args.Get(0).([]models.Conversation); ok {
		return conversations, args.Error(1)
	}
	return nil, args.Error(1)
}

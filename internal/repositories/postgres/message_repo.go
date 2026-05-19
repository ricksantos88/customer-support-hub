package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ricksantos88/customer-support-hub/internal/models"
	"github.com/ricksantos88/customer-support-hub/internal/repositories"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) repositories.MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, message *models.Message) error {
	if err := r.db.WithContext(ctx).Create(message).Error; err != nil {
		return fmt.Errorf("create message: %w", err)
	}
	return nil
}

func (r *MessageRepository) GetByConversation(ctx context.Context, conversationID uuid.UUID) ([]models.Message, error) {
	messages := make([]models.Message, 0)
	if err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("get messages by conversation: %w", err)
	}
	return messages, nil
}

func (r *MessageRepository) GetLastMessage(ctx context.Context, conversationID uuid.UUID) (*models.Message, error) {
	var message models.Message
	if err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Order("created_at DESC").
		First(&message).Error; err != nil {
		return nil, fmt.Errorf("get last message: %w", err)
	}
	return &message, nil
}

func (r *MessageRepository) ListByDirection(ctx context.Context, conversationID uuid.UUID, direction string) ([]models.Message, error) {
	direction = strings.ToLower(strings.TrimSpace(direction))
	if direction != models.MessageDirectionInbound && direction != models.MessageDirectionOutbound {
		return nil, fmt.Errorf("invalid message direction: %s", direction)
	}

	messages := make([]models.Message, 0)
	if err := r.db.WithContext(ctx).
		Where("conversation_id = ? AND direction = ?", conversationID, direction).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("list messages by direction: %w", err)
	}
	return messages, nil
}

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

type ConversationRepository struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) repositories.ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) Create(ctx context.Context, conversation *models.Conversation) error {
	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(conversation).Error; err != nil {
			return fmt.Errorf("create conversation: %w", err)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (r *ConversationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error) {
	var conversation models.Conversation
	if err := r.db.WithContext(ctx).
		Preload("Contact").
		Preload("AssignedAgent").
		First(&conversation, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("get conversation by id: %w", err)
	}
	return &conversation, nil
}

func (r *ConversationRepository) GetActiveByContact(ctx context.Context, contactID uuid.UUID) (*models.Conversation, error) {
	var conversation models.Conversation
	if err := r.db.WithContext(ctx).
		Where("contact_id = ? AND status IN ?", contactID, []string{models.ConversationStatusOpen, models.ConversationStatusPending}).
		Order("created_at DESC").
		First(&conversation).Error; err != nil {
		return nil, fmt.Errorf("get active conversation by contact: %w", err)
	}
	return &conversation, nil
}

func (r *ConversationRepository) AssignAgent(ctx context.Context, conversationID, agentID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var agent models.Agent
		if err := tx.First(&agent, "id = ?", agentID).Error; err != nil {
			return fmt.Errorf("assign agent validate agent: %w", err)
		}

		result := tx.Model(&models.Conversation{}).
			Where("id = ?", conversationID).
			Update("assigned_agent_id", agentID)
		if result.Error != nil {
			return fmt.Errorf("assign agent: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("assign agent: conversation not found")
		}
		return nil
	})
}

func (r *ConversationRepository) UpdateStatus(ctx context.Context, conversationID uuid.UUID, status string) error {
	status = strings.ToLower(strings.TrimSpace(status))
	switch status {
	case models.ConversationStatusOpen, models.ConversationStatusPending, models.ConversationStatusClosed:
	default:
		return fmt.Errorf("invalid conversation status: %s", status)
	}

	result := r.db.WithContext(ctx).
		Model(&models.Conversation{}).
		Where("id = ?", conversationID).
		Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("update conversation status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("update conversation status: conversation not found")
	}
	return nil
}

func (r *ConversationRepository) Close(ctx context.Context, conversationID uuid.UUID) error {
	return r.UpdateStatus(ctx, conversationID, models.ConversationStatusClosed)
}

func (r *ConversationRepository) ListByAgent(ctx context.Context, agentID uuid.UUID) ([]models.Conversation, error) {
	conversations := make([]models.Conversation, 0)
	if err := r.db.WithContext(ctx).
		Where("assigned_agent_id = ?", agentID).
		Order("created_at DESC").
		Find(&conversations).Error; err != nil {
		return nil, fmt.Errorf("list conversations by agent: %w", err)
	}
	return conversations, nil
}

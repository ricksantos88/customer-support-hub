package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	MessageDirectionInbound  = "inbound"
	MessageDirectionOutbound = "outbound"
)

type Message struct {
	ID             uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ConversationID uuid.UUID    `gorm:"type:uuid;not null;index:idx_messages_conversation_created_at,priority:1"`
	Conversation   Conversation `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:ConversationID"`
	Content        string       `gorm:"type:text;not null"`
	Direction      string       `gorm:"size:20;not null"`
	SenderID       uuid.UUID    `gorm:"type:uuid;not null"`
	CreatedAt      time.Time    `gorm:"not null;autoCreateTime"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}

	m.Direction = strings.ToLower(strings.TrimSpace(m.Direction))
	if m.Direction != MessageDirectionInbound && m.Direction != MessageDirectionOutbound {
		return fmt.Errorf("invalid message direction: %s", m.Direction)
	}

	if strings.TrimSpace(m.Content) == "" {
		return fmt.Errorf("message content is required")
	}

	if m.SenderID == uuid.Nil {
		return fmt.Errorf("sender_id is required")
	}

	if m.ConversationID == uuid.Nil {
		return fmt.Errorf("conversation_id is required")
	}

	if err := m.validateSenderExists(tx); err != nil {
		return err
	}

	return nil
}

func (m *Message) validateSenderExists(tx *gorm.DB) error {
	switch m.Direction {
	case MessageDirectionInbound:
		var count int64
		if err := tx.Model(&Contact{}).Where("id = ?", m.SenderID).Count(&count).Error; err != nil {
			return fmt.Errorf("validate inbound sender: %w", err)
		}
		if count == 0 {
			return fmt.Errorf("inbound sender contact does not exist")
		}
	case MessageDirectionOutbound:
		var count int64
		if err := tx.Model(&Agent{}).Where("id = ?", m.SenderID).Count(&count).Error; err != nil {
			return fmt.Errorf("validate outbound sender: %w", err)
		}
		if count == 0 {
			return fmt.Errorf("outbound sender agent does not exist")
		}
	default:
		return fmt.Errorf("invalid message direction: %s", m.Direction)
	}

	var convoCount int64
	if err := tx.Model(&Conversation{}).Where("id = ?", m.ConversationID).Count(&convoCount).Error; err != nil {
		return fmt.Errorf("validate conversation exists: %w", err)
	}
	if convoCount == 0 {
		return fmt.Errorf("conversation does not exist")
	}

	return nil
}

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
	ID              uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ConversationID  uuid.UUID    `gorm:"type:uuid;not null;index:idx_messages_conversation_created_at,priority:1"`
	Conversation    Conversation `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:ConversationID"`
	Content         string       `gorm:"type:text;not null"`
	Direction       string       `gorm:"size:20;not null"`
	SenderContactID *uuid.UUID   `gorm:"type:uuid"`
	SenderContact   *Contact     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:SenderContactID"`
	SenderAgentID   *uuid.UUID   `gorm:"type:uuid"`
	SenderAgent     *Agent       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:SenderAgentID"`
	SenderID        uuid.UUID    `gorm:"-"`
	CreatedAt       time.Time    `gorm:"not null;autoCreateTime"`
}

func (m *Message) BeforeCreate(_ *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}

	m.Content = strings.TrimSpace(m.Content)
	m.Direction = strings.ToLower(strings.TrimSpace(m.Direction))
	if m.Direction != MessageDirectionInbound && m.Direction != MessageDirectionOutbound {
		return fmt.Errorf("invalid message direction: %s", m.Direction)
	}

	if m.Content == "" {
		return fmt.Errorf("message content is required")
	}

	if m.ConversationID == uuid.Nil {
		return fmt.Errorf("conversation_id is required")
	}

	if err := m.normalizeSenderFields(); err != nil {
		return err
	}

	return nil
}

func (m *Message) normalizeSenderFields() error {
	switch m.Direction {
	case MessageDirectionInbound:
		if m.SenderID != uuid.Nil && m.SenderContactID == nil {
			sender := m.SenderID
			m.SenderContactID = &sender
		}
		if m.SenderContactID == nil || *m.SenderContactID == uuid.Nil {
			return fmt.Errorf("sender_contact_id is required for inbound messages")
		}
		m.SenderAgentID = nil
		m.SenderID = *m.SenderContactID
	case MessageDirectionOutbound:
		if m.SenderID != uuid.Nil && m.SenderAgentID == nil {
			sender := m.SenderID
			m.SenderAgentID = &sender
		}
		if m.SenderAgentID == nil || *m.SenderAgentID == uuid.Nil {
			return fmt.Errorf("sender_agent_id is required for outbound messages")
		}
		m.SenderContactID = nil
		m.SenderID = *m.SenderAgentID
	default:
		return fmt.Errorf("invalid message direction: %s", m.Direction)
	}

	return nil
}

func (m *Message) AfterFind(_ *gorm.DB) error {
	switch m.Direction {
	case MessageDirectionInbound:
		if m.SenderContactID != nil {
			m.SenderID = *m.SenderContactID
		}
	case MessageDirectionOutbound:
		if m.SenderAgentID != nil {
			m.SenderID = *m.SenderAgentID
		}
	}

	return nil
}

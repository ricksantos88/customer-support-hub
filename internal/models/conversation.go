package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ConversationStatusOpen    = "open"
	ConversationStatusPending = "pending"
	ConversationStatusClosed  = "closed"
)

type Conversation struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ContactID       uuid.UUID  `gorm:"type:uuid;not null;index"`
	Contact         Contact    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:ContactID"`
	AssignedAgentID *uuid.UUID `gorm:"type:uuid;index"`
	AssignedAgent   *Agent     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:AssignedAgentID"`
	Status          string     `gorm:"size:20;not null;default:open;index:idx_conversations_assigned_agent_status,priority:2"`
	CreatedAt       time.Time  `gorm:"not null;autoCreateTime"`
	UpdatedAt       time.Time  `gorm:"not null;autoUpdateTime"`
	Messages        []Message  `gorm:"foreignKey:ConversationID"`
}

func (c *Conversation) BeforeCreate(_ *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}

	c.Status = strings.ToLower(strings.TrimSpace(c.Status))
	if c.Status == "" {
		c.Status = ConversationStatusOpen
	}

	switch c.Status {
	case ConversationStatusOpen, ConversationStatusPending, ConversationStatusClosed:
		return nil
	default:
		return fmt.Errorf("invalid conversation status: %s", c.Status)
	}
}

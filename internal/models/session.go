package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID               uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AgentID          uuid.UUID  `gorm:"type:uuid;not null;index:idx_sessions_agent_active,priority:1" json:"agent_id"`
	Agent            Agent      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:AgentID" json:"-"`
	RefreshTokenHash string     `gorm:"column:refresh_token_hash;size:128;not null;uniqueIndex" json:"-"`
	UserAgent        string     `gorm:"size:255" json:"user_agent,omitempty"`
	IPAddress        string     `gorm:"column:ip_address;size:45" json:"ip_address,omitempty"`
	CreatedAt        time.Time  `gorm:"not null;autoCreateTime" json:"created_at"`
	LastUsedAt       time.Time  `gorm:"not null;autoUpdateTime" json:"last_used_at"`
	ExpiresAt        time.Time  `gorm:"not null;index:idx_sessions_expires_at" json:"expires_at"`
	RevokedAt        *time.Time `gorm:"index" json:"revoked_at,omitempty"`
}

func (Session) TableName() string {
	return "auth_sessions"
}

func (s *Session) BeforeCreate(_ *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}

	if s.AgentID == uuid.Nil {
		return fmt.Errorf("agent_id is required")
	}

	if strings.TrimSpace(s.RefreshTokenHash) == "" {
		return fmt.Errorf("refresh_token_hash is required")
	}

	if s.ExpiresAt.IsZero() {
		return fmt.Errorf("expires_at is required")
	}

	if s.LastUsedAt.IsZero() {
		s.LastUsedAt = time.Now().UTC()
	}

	return nil
}

func (s *Session) IsActive(now time.Time) bool {
	if s.RevokedAt != nil {
		return false
	}
	return now.Before(s.ExpiresAt)
}

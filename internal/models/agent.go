package models

import (
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const bcryptCost = 10

const (
	AgentRoleAdmin = "admin"
	AgentRoleAgent = "agent"
)

type Agent struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name          string         `gorm:"size:255;not null"`
	Email         string         `gorm:"size:255;not null;uniqueIndex"`
	PasswordHash  string         `gorm:"column:password_hash;not null"`
	Password      string         `gorm:"-"`
	Role          string         `gorm:"size:20;not null;default:agent"`
	CreatedAt     time.Time      `gorm:"not null;autoCreateTime"`
	LastActive    time.Time      `gorm:"not null;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Conversations []Conversation `gorm:"foreignKey:AssignedAgentID"`
}

// BeforeSave runs for both create and update. We use it to:
// - normalize email (when present)
// - hash plaintext Password into PasswordHash (when provided)
//
// Note: the JWT signing secret should live in env/config (e.g. JWT_SECRET), not in the DB.
func (a *Agent) BeforeSave(_ *gorm.DB) error {
	if strings.TrimSpace(a.Email) != "" {
		parsed, err := mail.ParseAddress(strings.TrimSpace(a.Email))
		if err != nil || parsed.Address == "" {
			return fmt.Errorf("invalid email address")
		}
		a.Email = strings.ToLower(strings.TrimSpace(parsed.Address))
	}

	if strings.TrimSpace(a.Password) != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcryptCost)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}
		a.PasswordHash = string(hash)
		a.Password = ""
	}

	return nil
}

func (a *Agent) BeforeCreate(_ *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}

	if strings.TrimSpace(a.Name) == "" {
		return fmt.Errorf("agent name is required")
	}

	if strings.TrimSpace(a.Email) == "" {
		return fmt.Errorf("email is required")
	}

	// If no plaintext Password was provided, PasswordHash must already be set.
	if strings.TrimSpace(a.PasswordHash) == "" {
		return fmt.Errorf("password is required")
	}

	if a.Role == "" {
		a.Role = AgentRoleAgent
	}
	if a.Role != AgentRoleAdmin && a.Role != AgentRoleAgent {
		return fmt.Errorf("invalid role: must be 'admin' or 'agent'")
	}

	return nil
}

func (a *Agent) CheckPassword(password string) error {
	if strings.TrimSpace(a.PasswordHash) == "" {
		return fmt.Errorf("password hash is empty")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password)); err != nil {
		return fmt.Errorf("invalid password")
	}
	return nil
}

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

type Agent struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name          string         `gorm:"size:255;not null"`
	Email         string         `gorm:"size:255;not null;uniqueIndex"`
	JWTHash       string         `gorm:"column:jwt_hash;not null"`
	JWTSecret     string         `gorm:"-"`
	CreatedAt     time.Time      `gorm:"not null;autoCreateTime"`
	LastActive    time.Time      `gorm:"not null;autoUpdateTime"`
	Conversations []Conversation `gorm:"foreignKey:AssignedAgentID"`
}

func (a *Agent) BeforeCreate(_ *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}

	if strings.TrimSpace(a.Name) == "" {
		return fmt.Errorf("agent name is required")
	}

	parsed, err := mail.ParseAddress(strings.TrimSpace(a.Email))
	if err != nil || parsed.Address == "" {
		return fmt.Errorf("invalid email address")
	}
	a.Email = strings.ToLower(strings.TrimSpace(parsed.Address))

	if strings.TrimSpace(a.JWTSecret) == "" && strings.TrimSpace(a.JWTHash) == "" {
		return fmt.Errorf("jwt secret is required")
	}

	if strings.TrimSpace(a.JWTSecret) != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(a.JWTSecret), bcryptCost)
		if err != nil {
			return fmt.Errorf("hash jwt secret: %w", err)
		}
		a.JWTHash = string(hash)
		a.JWTSecret = ""
	}

	if strings.TrimSpace(a.JWTHash) == "" {
		return fmt.Errorf("jwt hash is required")
	}

	return nil
}

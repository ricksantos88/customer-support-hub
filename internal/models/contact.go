package models

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var phoneE164Regex = regexp.MustCompile(`^\+[0-9]{7,}$`)

type Contact struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Phone         string         `gorm:"size:20;not null;uniqueIndex:idx_contacts_phone_active_unique,where:deleted_at IS NULL"`
	Name          string         `gorm:"size:255;not null"`
	CreatedAt     time.Time      `gorm:"not null;autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"not null;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Conversations []Conversation `gorm:"foreignKey:ContactID"`
}

func (c *Contact) BeforeCreate(_ *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}

	c.Phone = strings.TrimSpace(c.Phone)
	if !phoneE164Regex.MatchString(c.Phone) {
		return fmt.Errorf("invalid phone format: must be E.164 with + and digits only, min 8 chars")
	}

	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("contact name is required")
	}

	return nil
}

func (c *Contact) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("UpdatedAt", time.Now().UTC())
	return nil
}

func (c *Contact) Validate(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if !phoneE164Regex.MatchString(strings.TrimSpace(c.Phone)) {
		return fmt.Errorf("invalid phone format: must be E.164 with + and digits only, min 8 chars")
	}

	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("contact name is required")
	}

	return nil
}

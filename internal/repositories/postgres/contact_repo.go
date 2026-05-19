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

type ContactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) repositories.ContactRepository {
	return &ContactRepository{db: db}
}

func (r *ContactRepository) Create(ctx context.Context, contact *models.Contact) error {
	if err := contact.Validate(ctx); err != nil {
		return fmt.Errorf("validate contact: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(contact).Error; err != nil {
		return fmt.Errorf("create contact: %w", err)
	}
	return nil
}

func (r *ContactRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Contact, error) {
	var contact models.Contact
	if err := r.db.WithContext(ctx).First(&contact, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("get contact by id: %w", err)
	}
	return &contact, nil
}

func (r *ContactRepository) GetByPhone(ctx context.Context, phone string) (*models.Contact, error) {
	var contact models.Contact
	if err := r.db.WithContext(ctx).First(&contact, "phone = ?", strings.TrimSpace(phone)).Error; err != nil {
		return nil, fmt.Errorf("get contact by phone: %w", err)
	}
	return &contact, nil
}

func (r *ContactRepository) Update(ctx context.Context, contact *models.Contact) error {
	if err := contact.Validate(ctx); err != nil {
		return fmt.Errorf("validate contact: %w", err)
	}

	if err := r.db.WithContext(ctx).Save(contact).Error; err != nil {
		return fmt.Errorf("update contact: %w", err)
	}
	return nil
}

func (r *ContactRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.Contact{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("soft delete contact: %w", err)
	}
	return nil
}

func (r *ContactRepository) List(ctx context.Context, limit, offset int) ([]models.Contact, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	contacts := make([]models.Contact, 0)
	query := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Offset(offset).Find(&contacts)
	if query.Error != nil {
		return nil, fmt.Errorf("list contacts: %w", query.Error)
	}
	return contacts, nil
}

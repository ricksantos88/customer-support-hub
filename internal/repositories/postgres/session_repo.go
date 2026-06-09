package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ricksantos88/customer-support-hub/internal/models"
	"github.com/ricksantos88/customer-support-hub/internal/repositories"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) repositories.SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

func (r *SessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	var session models.Session
	if err := r.db.WithContext(ctx).First(&session, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("get session by id: %w", err)
	}
	return &session, nil
}

func (r *SessionRepository) GetByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*models.Session, error) {
	var session models.Session
	if err := r.db.WithContext(ctx).
		Where("refresh_token_hash = ? AND revoked_at IS NULL AND expires_at > ?", refreshTokenHash, time.Now().UTC()).
		First(&session).Error; err != nil {
		return nil, fmt.Errorf("get session by refresh token hash: %w", err)
	}
	return &session, nil
}

func (r *SessionRepository) GetActiveByID(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	var session models.Session
	if err := r.db.WithContext(ctx).
		Where("id = ? AND revoked_at IS NULL AND expires_at > ?", id, time.Now().UTC()).
		First(&session).Error; err != nil {
		return nil, fmt.Errorf("get active session by id: %w", err)
	}
	return &session, nil
}

func (r *SessionRepository) RotateRefreshToken(ctx context.Context, sessionID uuid.UUID, refreshTokenHash string, lastUsedAt, expiresAt time.Time) error {
	result := r.db.WithContext(ctx).
		Model(&models.Session{}).
		Where("id = ? AND revoked_at IS NULL", sessionID).
		Updates(map[string]interface{}{
			"refresh_token_hash": refreshTokenHash,
			"last_used_at":       lastUsedAt.UTC(),
			"expires_at":         expiresAt.UTC(),
		})
	if result.Error != nil {
		return fmt.Errorf("rotate session refresh token: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("rotate session refresh token: session not found")
	}
	return nil
}

func (r *SessionRepository) Revoke(ctx context.Context, sessionID uuid.UUID, revokedAt time.Time) error {
	result := r.db.WithContext(ctx).
		Model(&models.Session{}).
		Where("id = ? AND revoked_at IS NULL", sessionID).
		Updates(map[string]interface{}{
			"revoked_at":   revokedAt.UTC(),
			"last_used_at": revokedAt.UTC(),
		})
	if result.Error != nil {
		return fmt.Errorf("revoke session: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("revoke session: session not found")
	}
	return nil
}

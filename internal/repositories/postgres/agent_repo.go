package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ricksantos88/customer-support-hub/internal/models"
	"github.com/ricksantos88/customer-support-hub/internal/repositories"
)

type AgentRepository struct {
	db *gorm.DB
}

func NewAgentRepository(db *gorm.DB) repositories.AgentRepository {
	return &AgentRepository{db: db}
}

func (r *AgentRepository) Create(ctx context.Context, agent *models.Agent) error {
	if err := r.db.WithContext(ctx).Create(agent).Error; err != nil {
		return fmt.Errorf("create agent: %w", err)
	}
	return nil
}

func (r *AgentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Agent, error) {
	var agent models.Agent
	if err := r.db.WithContext(ctx).First(&agent, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("get agent by id: %w", err)
	}
	return &agent, nil
}

func (r *AgentRepository) GetByEmail(ctx context.Context, email string) (*models.Agent, error) {
	var agent models.Agent
	if err := r.db.WithContext(ctx).First(&agent, "email = ?", strings.ToLower(strings.TrimSpace(email))).Error; err != nil {
		return nil, fmt.Errorf("get agent by email: %w", err)
	}
	return &agent, nil
}

func (r *AgentRepository) Update(ctx context.Context, agent *models.Agent) error {
	if err := r.db.WithContext(ctx).Save(agent).Error; err != nil {
		return fmt.Errorf("update agent: %w", err)
	}
	return nil
}

func (r *AgentRepository) UpdateLastActive(ctx context.Context, agentID uuid.UUID, lastActive time.Time) error {
	result := r.db.WithContext(ctx).
		Model(&models.Agent{}).
		Where("id = ?", agentID).
		Update("last_active", lastActive.UTC())
	if result.Error != nil {
		return fmt.Errorf("update agent last active: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("update agent last active: agent not found")
	}
	return nil
}

func (r *AgentRepository) List(ctx context.Context, limit, offset int) ([]models.Agent, error) {
	var agents []models.Agent
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at desc").Find(&agents).Error; err != nil {
		return nil, fmt.Errorf("list agents: %w", err)
	}
	return agents, nil
}

func (r *AgentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.Agent{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("delete agent: %w", err)
	}
	return nil
}

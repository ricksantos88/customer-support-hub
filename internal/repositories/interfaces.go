package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ricksantos88/customer-support-hub/internal/models"
)

type AgentRepository interface {
	Create(ctx context.Context, agent *models.Agent) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Agent, error)
	GetByEmail(ctx context.Context, email string) (*models.Agent, error)
	Update(ctx context.Context, agent *models.Agent) error
	UpdateLastActive(ctx context.Context, agentID uuid.UUID, lastActive time.Time) error
	List(ctx context.Context, limit, offset int) ([]models.Agent, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error)
	GetByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*models.Session, error)
	GetActiveByID(ctx context.Context, id uuid.UUID) (*models.Session, error)
	RotateRefreshToken(ctx context.Context, sessionID uuid.UUID, refreshTokenHash string, lastUsedAt, expiresAt time.Time) error
	Revoke(ctx context.Context, sessionID uuid.UUID, revokedAt time.Time) error
	ListActive(ctx context.Context) ([]models.Session, error)
	RevokeAllByAgentID(ctx context.Context, agentID uuid.UUID, revokedAt time.Time) error
}

type ContactRepository interface {
	Create(ctx context.Context, contact *models.Contact) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Contact, error)
	GetByPhone(ctx context.Context, phone string) (*models.Contact, error)
	Update(ctx context.Context, contact *models.Contact) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]models.Contact, error)
}

type ConversationRepository interface {
	Create(ctx context.Context, conversation *models.Conversation) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error)
	GetActiveByContact(ctx context.Context, contactID uuid.UUID) (*models.Conversation, error)
	AssignAgent(ctx context.Context, conversationID, agentID uuid.UUID) error
	UpdateStatus(ctx context.Context, conversationID uuid.UUID, status string) error
	Close(ctx context.Context, conversationID uuid.UUID) error
	ListByAgent(ctx context.Context, agentID uuid.UUID) ([]models.Conversation, error)
}

type MessageRepository interface {
	Create(ctx context.Context, message *models.Message) error
	GetByConversation(ctx context.Context, conversationID uuid.UUID) ([]models.Message, error)
	GetLastMessage(ctx context.Context, conversationID uuid.UUID) (*models.Message, error)
	ListByDirection(ctx context.Context, conversationID uuid.UUID, direction string) ([]models.Message, error)
}

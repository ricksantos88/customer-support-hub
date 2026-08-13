package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/ricksantos88/customer-support-hub/internal/models"
)

type SessionCache struct {
	client *redis.Client
}

func NewSessionCache(client *redis.Client) *SessionCache {
	return &SessionCache{client: client}
}

func (s *SessionCache) SetSession(ctx context.Context, session *models.Session, ttl time.Duration) error {
	if s == nil || s.client == nil {
		return fmt.Errorf("redis unavailable")
	}
	payload, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("marshal session cache payload: %w", err)
	}
	return s.client.Set(ctx, sessionKey(session.ID), payload, ttl).Err()
}

func (s *SessionCache) GetSessionByID(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	if s == nil || s.client == nil {
		return nil, fmt.Errorf("redis unavailable")
	}
	value, err := s.client.Get(ctx, sessionKey(id)).Bytes()
	if err != nil {
		return nil, err
	}
	var session models.Session
	if err := json.Unmarshal(value, &session); err != nil {
		return nil, fmt.Errorf("unmarshal session cache payload: %w", err)
	}
	return &session, nil
}

func (s *SessionCache) DeleteSession(ctx context.Context, id uuid.UUID) error {
	if s == nil || s.client == nil {
		return fmt.Errorf("redis unavailable")
	}
	return s.client.Del(ctx, sessionKey(id)).Err()
}

func (s *SessionCache) SetRefreshTokenSessionID(ctx context.Context, refreshTokenHash string, sessionID uuid.UUID, ttl time.Duration) error {
	if s == nil || s.client == nil {
		return fmt.Errorf("redis unavailable")
	}
	return s.client.Set(ctx, refreshTokenKey(refreshTokenHash), sessionID.String(), ttl).Err()
}

func (s *SessionCache) DeleteRefreshTokenSessionID(ctx context.Context, refreshTokenHash string) error {
	if s == nil || s.client == nil {
		return fmt.Errorf("redis unavailable")
	}
	return s.client.Del(ctx, refreshTokenKey(refreshTokenHash)).Err()
}

func (s *SessionCache) Ping(ctx context.Context) error {
	if s == nil || s.client == nil {
		return fmt.Errorf("redis unavailable")
	}
	return s.client.Ping(ctx).Err()
}

func (s *SessionCache) GetSessionIDByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (uuid.UUID, error) {
	if s == nil || s.client == nil {
		return uuid.Nil, fmt.Errorf("redis unavailable")
	}
	value, err := s.client.Get(ctx, refreshTokenKey(refreshTokenHash)).Result()
	if err != nil {
		return uuid.Nil, err
	}
	parsed, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse cached session id: %w", err)
	}
	return parsed, nil
}

func sessionKey(id uuid.UUID) string {
	return fmt.Sprintf("auth:session:%s", id.String())
}

func refreshTokenKey(hash string) string {
	return fmt.Sprintf("auth:refresh:%s", hash)
}

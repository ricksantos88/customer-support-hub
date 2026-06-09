package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/ricksantos88/customer-support-hub/internal/models"
	"github.com/ricksantos88/customer-support-hub/internal/repositories"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type SessionCache interface {
	SetSession(ctx context.Context, session *models.Session, ttl time.Duration) error
	GetSessionByID(ctx context.Context, id uuid.UUID) (*models.Session, error)
	DeleteSession(ctx context.Context, id uuid.UUID) error
	SetRefreshTokenSessionID(ctx context.Context, refreshTokenHash string, sessionID uuid.UUID, ttl time.Duration) error
	GetSessionIDByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (uuid.UUID, error)
	DeleteRefreshTokenSessionID(ctx context.Context, refreshTokenHash string) error
}

type Service struct {
	agents   repositories.AgentRepository
	sessions repositories.SessionRepository
	cache    SessionCache
	secret   string
	issuer   string

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	sessionTTL      time.Duration
	now             func() time.Time
}

type SessionResult struct {
	AgentID               uuid.UUID `json:"agent_id"`
	SessionID             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

func NewService(agents repositories.AgentRepository, sessions repositories.SessionRepository, cache SessionCache, secret, issuer string, accessTokenTTL, refreshTokenTTL, sessionTTL time.Duration) *Service {
	return &Service{
		agents:          agents,
		sessions:        sessions,
		cache:           cache,
		secret:          secret,
		issuer:          issuer,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		sessionTTL:      sessionTTL,
		now:             func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Login(ctx context.Context, email, password, userAgent, ipAddress string) (*SessionResult, error) {
	agent, err := s.agents.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("authenticate agent: %w", err)
	}

	if err := agent.CheckPassword(password); err != nil {
		return nil, err
	}

	now := s.now()
	refreshToken, refreshTokenHash, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	session := &models.Session{
		AgentID:          agent.ID,
		RefreshTokenHash: refreshTokenHash,
		UserAgent:        strings.TrimSpace(userAgent),
		IPAddress:        strings.TrimSpace(ipAddress),
		CreatedAt:        now,
		LastUsedAt:       now,
		ExpiresAt:        now.Add(s.refreshTokenTTL),
	}

	if err := s.sessions.Create(ctx, session); err != nil {
		return nil, err
	}

	if err := s.agents.UpdateLastActive(ctx, agent.ID, now); err != nil {
		return nil, err
	}

	if s.cache != nil {
		_ = s.cache.SetSession(ctx, session, s.sessionTTL)
		_ = s.cache.SetRefreshTokenSessionID(ctx, refreshTokenHash, session.ID, s.sessionTTL)
	}

	accessToken, err := s.generateAccessToken(agent, session.ID, now)
	if err != nil {
		return nil, err
	}

	return &SessionResult{
		AgentID:               agent.ID,
		SessionID:             session.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  now.Add(s.accessTokenTTL),
		RefreshTokenExpiresAt: now.Add(s.refreshTokenTTL),
	}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken, userAgent, ipAddress string) (*SessionResult, error) {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token is required")
	}

	refreshTokenHash := hashToken(refreshToken)
	now := s.now()

	var session *models.Session
	var err error

	if s.cache != nil {
		if sessionID, cacheErr := s.cache.GetSessionIDByRefreshTokenHash(ctx, refreshTokenHash); cacheErr == nil && sessionID != uuid.Nil {
			if cached, cacheErr := s.sessions.GetActiveByID(ctx, sessionID); cacheErr == nil && cached.RefreshTokenHash == refreshTokenHash {
				session = cached
			}
		}
	}

	if session == nil {
		session, err = s.sessions.GetByRefreshTokenHash(ctx, refreshTokenHash)
	}
	if err != nil {
		return nil, err
	}

	if !session.IsActive(now) {
		return nil, fmt.Errorf("session is not active")
	}

	refreshToken, newHash, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	newExpiresAt := now.Add(s.refreshTokenTTL)
	if err := s.sessions.RotateRefreshToken(ctx, session.ID, newHash, now, newExpiresAt); err != nil {
		return nil, err
	}

	if s.cache != nil {
		updatedSession := &models.Session{
			ID:               session.ID,
			AgentID:          session.AgentID,
			RefreshTokenHash: newHash,
			UserAgent:        session.UserAgent,
			IPAddress:        session.IPAddress,
			CreatedAt:        session.CreatedAt,
			LastUsedAt:       now,
			ExpiresAt:        newExpiresAt,
		}
		_ = s.cache.DeleteRefreshTokenSessionID(ctx, refreshTokenHash)
		_ = s.cache.SetSession(ctx, updatedSession, s.sessionTTL)
		_ = s.cache.SetRefreshTokenSessionID(ctx, newHash, session.ID, s.sessionTTL)
	}

	agent, err := s.agents.GetByID(ctx, session.AgentID)
	if err != nil {
		return nil, err
	}

	if err := s.agents.UpdateLastActive(ctx, agent.ID, now); err != nil {
		return nil, err
	}

	accessToken, err := s.generateAccessToken(agent, session.ID, now)
	if err != nil {
		return nil, err
	}

	return &SessionResult{
		AgentID:               agent.ID,
		SessionID:             session.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  now.Add(s.accessTokenTTL),
		RefreshTokenExpiresAt: newExpiresAt,
	}, nil
}

func (s *Service) Logout(ctx context.Context, sessionID uuid.UUID) error {
	now := s.now()
	if err := s.sessions.Revoke(ctx, sessionID, now); err != nil {
		return err
	}
	if s.cache != nil {
		_ = s.cache.DeleteSession(ctx, sessionID)
	}
	return nil
}

func (s *Service) ValidateAccessToken(tokenString string) (jwt.MapClaims, error) {
	return validateToken(tokenString, s.secret)
}

func (s *Service) generateAccessToken(agent *models.Agent, sessionID uuid.UUID, now time.Time) (string, error) {
	claims := jwt.MapClaims{
		"agent_id":   agent.ID.String(),
		"session_id": sessionID.String(),
		"token_type": TokenTypeAccess,
		"iss":        s.issuer,
		"sub":        agent.ID.String(),
		"iat":        now.Unix(),
		"nbf":        now.Unix(),
		"exp":        now.Add(s.accessTokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func validateToken(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func generateRefreshToken() (string, string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}
	refreshToken := base64.RawURLEncoding.EncodeToString(bytes)
	return refreshToken, hashToken(refreshToken), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

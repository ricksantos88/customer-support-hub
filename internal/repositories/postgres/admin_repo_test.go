package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/ricksantos88/customer-support-hub/internal/models"
	postgresrepo "github.com/ricksantos88/customer-support-hub/internal/repositories/postgres"
)

func TestAgentRepository_ListAndDelete(t *testing.T) {
	repo := postgresrepo.NewAgentRepository(sharedDB)
	ctx := context.Background()

	// 1. Create a unique agent
	agent := &models.Agent{
		Name:         "Admin Tester",
		Email:        "admin-tester@example.com",
		PasswordHash: "some-hash",
		Role:         "admin",
	}
	require.NoError(t, repo.Create(ctx, agent))

	// 2. List agents and check if present
	agents, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	require.NotEmpty(t, agents)

	found := false
	for _, a := range agents {
		if a.ID == agent.ID {
			found = true
			require.Equal(t, "Admin Tester", a.Name)
			require.Equal(t, "admin", a.Role)
		}
	}
	require.True(t, found)

	// 3. Delete (soft-delete) agent
	require.NoError(t, repo.Delete(ctx, agent.ID))

	// 4. Retrieve should fail with record not found
	_, err = repo.GetByID(ctx, agent.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestSessionRepository_ListActiveAndRevokeAll(t *testing.T) {
	agentRepo := postgresrepo.NewAgentRepository(sharedDB)
	sessionRepo := postgresrepo.NewSessionRepository(sharedDB)
	ctx := context.Background()

	// 1. Create agent
	agent := &models.Agent{
		Name:         "Sess Agent",
		Email:        "sess-agent@example.com",
		PasswordHash: "hash",
		Role:         "agent",
	}
	require.NoError(t, agentRepo.Create(ctx, agent))

	// 2. Create active session
	sess := &models.Session{
		AgentID:          agent.ID,
		RefreshTokenHash: "hash123",
		ExpiresAt:        time.Now().Add(1 * time.Hour),
	}
	require.NoError(t, sessionRepo.Create(ctx, sess))

	// 3. List active sessions and check
	sessions, err := sessionRepo.ListActive(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, sessions)

	found := false
	for _, s := range sessions {
		if s.ID == sess.ID {
			found = true
			require.Equal(t, agent.Name, s.Agent.Name)
		}
	}
	require.True(t, found)

	// 4. Revoke all by agent ID
	require.NoError(t, sessionRepo.RevokeAllByAgentID(ctx, agent.ID, time.Now()))

	// 5. Active sessions list should not contain this session
	sessions2, err := sessionRepo.ListActive(ctx)
	require.NoError(t, err)

	found2 := false
	for _, s := range sessions2 {
		if s.ID == sess.ID {
			found2 = true
		}
	}
	require.False(t, found2)
}

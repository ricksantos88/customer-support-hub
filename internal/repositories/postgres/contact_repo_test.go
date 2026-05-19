package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	pgcontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ricksantos88/customer-support-hub/internal/models"
	postgresrepo "github.com/ricksantos88/customer-support-hub/internal/repositories/postgres"
)

func TestCreate(t *testing.T) {
	db, terminate := setupTestDB(t)
	defer terminate()

	repo := postgresrepo.NewContactRepository(db)
	ctx := context.Background()

	contact := &models.Contact{Phone: "+5511999990001", Name: "Contato Teste"}
	require.NoError(t, repo.Create(ctx, contact))
	require.NotEqual(t, "", contact.ID.String())
}

func TestGetByPhone(t *testing.T) {
	db, terminate := setupTestDB(t)
	defer terminate()

	repo := postgresrepo.NewContactRepository(db)
	ctx := context.Background()

	contact := &models.Contact{Phone: "+5511999990002", Name: "Contato Consulta"}
	require.NoError(t, repo.Create(ctx, contact))

	found, err := repo.GetByPhone(ctx, "+5511999990002")
	require.NoError(t, err)
	require.Equal(t, contact.Phone, found.Phone)
	require.Equal(t, contact.Name, found.Name)
}

func TestSoftDelete(t *testing.T) {
	db, terminate := setupTestDB(t)
	defer terminate()

	repo := postgresrepo.NewContactRepository(db)
	ctx := context.Background()

	contact := &models.Contact{Phone: "+5511999990003", Name: "Contato Removido"}
	require.NoError(t, repo.Create(ctx, contact))
	require.NoError(t, repo.SoftDelete(ctx, contact.ID))

	_, err := repo.GetByPhone(ctx, "+5511999990003")
	require.Error(t, err)
	require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	t.Helper()

	ctx := context.Background()
	container, err := pgcontainer.Run(ctx, "postgres:15-alpine",
		pgcontainer.WithDatabase("customer_support_test"),
		pgcontainer.WithUsername("support"),
		pgcontainer.WithPassword("support123"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").WithStartupTimeout(2*time.Minute),
		),
	)
	require.NoError(t, err)

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(&models.Contact{}, &models.Agent{}, &models.Conversation{}, &models.Message{}))

	cleanup := func() {
		_ = container.Terminate(ctx)
	}

	return db, cleanup
}

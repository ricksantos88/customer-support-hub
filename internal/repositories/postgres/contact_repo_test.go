package postgres_test

import (
	"context"
	"errors"
	"os"
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

// sharedDB is the single database connection shared across all tests in this package.
var sharedDB *gorm.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := pgcontainer.Run(ctx, "postgres:15-alpine",
		pgcontainer.WithDatabase("customer_support_test"),
		pgcontainer.WithUsername("support"),
		pgcontainer.WithPassword("support123"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(2*time.Minute),
		),
	)
	if err != nil {
		panic("failed to start postgres container: " + err.Error())
	}
	defer container.Terminate(ctx) //nolint:errcheck

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic("failed to get connection string: " + err.Error())
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to open gorm connection: " + err.Error())
	}

	if err := db.AutoMigrate(
		&models.Contact{},
		&models.Agent{},
		&models.Conversation{},
		&models.Message{},
		&models.Session{},
	); err != nil {
		panic("failed to auto-migrate: " + err.Error())
	}

	sharedDB = db
	os.Exit(m.Run())
}

func TestContactRepository_Create(t *testing.T) {
	repo := postgresrepo.NewContactRepository(sharedDB)
	ctx := context.Background()

	contact := &models.Contact{Phone: "+5511999990001", Name: "Contato Teste"}
	require.NoError(t, repo.Create(ctx, contact))
	require.NotEqual(t, "", contact.ID.String())
}

func TestContactRepository_GetByPhone(t *testing.T) {
	repo := postgresrepo.NewContactRepository(sharedDB)
	ctx := context.Background()

	contact := &models.Contact{Phone: "+5511999990002", Name: "Contato Consulta"}
	require.NoError(t, repo.Create(ctx, contact))

	found, err := repo.GetByPhone(ctx, "+5511999990002")
	require.NoError(t, err)
	require.Equal(t, contact.Phone, found.Phone)
	require.Equal(t, contact.Name, found.Name)
}

func TestContactRepository_SoftDelete(t *testing.T) {
	repo := postgresrepo.NewContactRepository(sharedDB)
	ctx := context.Background()

	contact := &models.Contact{Phone: "+5511999990003", Name: "Contato Removido"}
	require.NoError(t, repo.Create(ctx, contact))
	require.NoError(t, repo.SoftDelete(ctx, contact.ID))

	_, err := repo.GetByPhone(ctx, "+5511999990003")
	require.Error(t, err)
	require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

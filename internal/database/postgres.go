package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ricksantos88/customer-support-hub/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}

	setPoolDefaults(sqlDB)

	if shouldAutoMigrate() {
		if err := db.AutoMigrate(
			&models.Contact{},
			&models.Agent{},
			&models.Conversation{},
			&models.Message{},
		); err != nil {
			return nil, fmt.Errorf("auto migrate postgres models: %w", err)
		}
	}

	return db, nil
}

func setPoolDefaults(sqlDB interface {
	SetMaxIdleConns(int)
	SetMaxOpenConns(int)
	SetConnMaxLifetime(time.Duration)
}) {
	maxOpen := getEnvInt("DB_MAX_OPEN_CONNS", 25)
	maxIdle := getEnvInt("DB_MAX_IDLE_CONNS", 25)
	maxLifetime := time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME_MINUTES", 5)) * time.Minute

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(maxLifetime)
}

func shouldAutoMigrate() bool {
	value := os.Getenv("DB_AUTO_MIGRATE")
	return value == "1" || value == "true" || value == "TRUE" || value == "yes" || value == "YES"
}

func getEnvInt(name string, fallback int) int {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

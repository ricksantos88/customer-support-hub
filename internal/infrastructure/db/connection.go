package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ricksantos88/customer-support-hub/internal/models"
)

func NewPostgresConnection(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&models.Session{}); err != nil {
		return nil, err
	}
	return db, nil
}

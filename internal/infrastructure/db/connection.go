package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

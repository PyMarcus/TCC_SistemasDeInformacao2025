package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgresConn: create a new postgres connection
func NewPostgresConn(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}

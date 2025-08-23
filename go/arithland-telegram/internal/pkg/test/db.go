package test

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to initiate db: %w", err)
	}

	return db, nil
}

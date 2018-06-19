package store

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type DBMigrator struct {
	db *gorm.DB
}

func (d *DBMigrator) Migrate() error {
	return d.db.AutoMigrate(Recipe{}).Error
}

func NewDBMigrator(db *gorm.DB) (*DBMigrator, error) {
	if db == nil {
		return nil, errors.New("store: db *gorm.DB is nil")
	}

	return &DBMigrator{db}, nil
}

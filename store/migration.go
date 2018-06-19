package store

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type GormMigrator struct {
	db *gorm.DB
}

func (d *GormMigrator) Migrate() error {
	return d.db.AutoMigrate(Recipe{}).Error
}

func NewGormMigrator(db *gorm.DB) (*GormMigrator, error) {
	if db == nil {
		return nil, errors.New("store: db *gorm.DB is nil")
	}

	return &GormMigrator{db}, nil
}

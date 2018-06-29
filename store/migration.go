package store

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// GormMigrator ...
type GormMigrator struct {
	db *gorm.DB
}

// Migrate ...
func (d *GormMigrator) Migrate() error {
	return errors.WithStack(d.db.AutoMigrate(Recipe{}).Error)
}

// NewGormMigrator ...
func NewGormMigrator(db *gorm.DB) (*GormMigrator, error) {
	if db == nil {
		return nil, errors.New("store: db *gorm.DB is nil")
	}

	return &GormMigrator{db}, nil
}

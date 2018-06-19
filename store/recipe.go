package store

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

//noinspection GoUnusedGlobalVariable
var (
	// ErrRecipeNotFound not found error
	ErrRecipeNotFound = errors.New("store: recipe not found")
)

type Recipe struct {
	gorm.Model

	Name        string
	Description string
}

type RecipeLister interface {
	All() ([]*Recipe, error)
}

type RecipeFinder interface {
	Find(ID uint) (*Recipe, error)
}

type RecipeSearcher interface {
	Search(query string) ([]*Recipe, error)
}

type RecipeInserter interface {
	Insert(r *Recipe) error
}

type RecipeUpdater interface {
	Update(r *Recipe) error
}

// Implementation using GORM database ORM
type GormRecipesStore struct {
	db *gorm.DB
}

func NewGormRecipesStore(db *gorm.DB) (*GormRecipesStore, error) {
	if db == nil {
		return nil, errors.New("store: db *gorm.DB is nil")
	}
	return &GormRecipesStore{db}, nil
}

func (d *GormRecipesStore) Update(r *Recipe) error {
	err := d.db.Model(Recipe{}).Update(r).Error
	if gorm.IsRecordNotFoundError(err) {
		return ErrRecipeNotFound
	}

	return err
}

func (d *GormRecipesStore) Insert(r *Recipe) error {
	return d.db.Model(Recipe{}).Create(r).Error
}

func (d *GormRecipesStore) Find(ID uint) (*Recipe, error) {
	var err error
	r := &Recipe{}

	err = d.db.Find(r, ID).Error
	if gorm.IsRecordNotFoundError(err) {
		return r, ErrRecipeNotFound
	}
	return r, err
}

// TODO: Pagination
func (d *GormRecipesStore) All() ([]*Recipe, error) {
	var recipes []*Recipe
	err := d.db.Find(recipes).Error
	return recipes, err
}

// TODO: Pagination
func (d *GormRecipesStore) Search(query string) ([]*Recipe, error) {
	var recipes []*Recipe
	err := d.db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", query)).Find(recipes).Error
	return recipes, err
}

package settlement

import (
	"gorm.io/gorm"
)

type IRepository interface {
	WithTransaction(f func(db *gorm.DB) error) error
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (s *Repository) WithTransaction(f func(db *gorm.DB) error) error {
	return s.db.Transaction(f)
}

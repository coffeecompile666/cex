package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) WithTransaction(fn func(tx *gorm.DB) error) error {
	return repo.db.Transaction(fn)
}

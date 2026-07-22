package repository

import (
	"icon_exchange/internal/order/model"
	"icon_exchange/internal/shared"

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

func (repo *Repository) GetByUserID(userID uint, offset int, limit int) ([]*model.Order, error) {
	var orders []*model.Order

	err := repo.db.Where("userID = ?", userID).Offset(offset).Limit(limit).Find(&orders).Error

	if err != nil {
		return nil, shared.ErrInternalServerError
	}

	return orders, nil
}

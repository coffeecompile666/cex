package model

import (
	"icon_exchange/internal/market/model"
	"icon_exchange/internal/shared"
	"math"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model

	UserID   uint `gorm:"type:bigint;not null;index:idx_order_user_id"`
	MarketID uint `gorm:"type:bigint;not null;index:idx_order_market_id"`
	Quantity uint `gorm:"type:bigint;not null"`
	Price    uint `gorm:"type:bigint;not null"`

	Market model.Market `gorm:"foreignkey:MarketID"`
}

func (o *Order) ValidateAmount() error {
	if o.Quantity == 0 {
		return shared.ErrInvalidOrderAmount
	}

	if o.Price > math.MaxUint/o.Quantity {
		return shared.ErrInvalidOrderAmount
	}

	return nil
}

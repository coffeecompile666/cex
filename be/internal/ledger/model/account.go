package model

import (
	"icon_exchange/internal/market/model"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model

	UserID          *uint `gorm:"uniqueIndex:uk_account"`
	IsSystemAccount bool
	MarketID        uint         `gorm:"not null;uniqueIndex:uk_account"`
	Market          model.Market `gorm:"foreignKey:MarketID"`
	Type            string       `gorm:"not null;uniqueIndex:uk_account"`
}

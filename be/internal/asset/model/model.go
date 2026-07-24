package model

import (
	"icon_exchange/internal/market/model"
	"icon_exchange/internal/shared"
	model2 "icon_exchange/internal/user/model"

	"gorm.io/gorm"
)

type Asset struct {
	gorm.Model

	UserID uint        `gorm:"not null;uniqueIndex:idx_asset_user_market;index:idx_asset_user_id"`
	User   model2.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	MarketID uint         `gorm:"not null;uniqueIndex:idx_asset_user_market"`
	Market   model.Market `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Amount       uint `gorm:"not null;default:0"`
	LockedAmount uint `gorm:"not null;default:0"`
}

func (a *Asset) GetAvailableAmount() uint {
	return a.Amount - a.LockedAmount
}

func (a *Asset) LockAmount(amount uint) error {
	if a.Amount < amount {
		return shared.ErrAmountNotSufficient
	}

	a.LockedAmount += amount
	return nil
}

func (a *Asset) UnLock(amount uint) error {
	if amount > a.LockedAmount {
		return shared.ErrUnlockAmountFailed
	}
	a.LockedAmount -= amount
	return nil
}

func (a *Asset) Debit(amount uint) error {
	if amount < a.GetAvailableAmount() {
		return shared.ErrAmountNotSufficient
	}
	a.Amount -= amount
	return nil
}

func (a *Asset) Credit(amount uint) {
	a.Amount += amount
}

package model

import (
	ordermodel "icon_exchange/internal/order/model"

	"gorm.io/gorm"
)

type Trade struct {
	gorm.Model

	TakerOrderID uint             `gorm:"index;not null"`
	TakerOrder   ordermodel.Order `gorm:"foreignKey:TakerOrderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	MakerOrderID uint             `gorm:"index;not null"`
	MakerOrder   ordermodel.Order `gorm:"foreignKey:MakerOrderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	Price    uint `gorm:"not null"`
	Quantity uint `gorm:"not null"`
}

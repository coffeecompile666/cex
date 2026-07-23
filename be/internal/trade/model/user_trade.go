package model

import (
	ordermodel "icon_exchange/internal/order/model"
	"icon_exchange/internal/shared"
	usermodel "icon_exchange/internal/user/model"

	"gorm.io/gorm"
)

type UserTrade struct {
	gorm.Model

	UserID uint           `gorm:"index;not null"`
	User   usermodel.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	TradeID uint  `gorm:"index;not null"`
	Trade   Trade `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	OrderID uint             `gorm:"index;not null"`
	Order   ordermodel.Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	Side shared.OrderSide `gorm:"type:varchar(10);not null"`
	Role shared.TradeRole `gorm:"type:varchar(10);not null"`

	Price    uint `gorm:"not null"`
	Quantity uint `gorm:"not null"`
}

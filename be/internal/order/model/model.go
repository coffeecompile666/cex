package model

import (
	"icon_exchange/internal/market/model"
	"icon_exchange/internal/shared"
	"math"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusOpen            OrderStatus = "OPEN"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusCancelled       OrderStatus = "CANCELLED"
	OrderStatusRejected        OrderStatus = "REJECTED"
	OrderStatusExpired         OrderStatus = "EXPIRED"
)

var allowedStatus = map[OrderStatus]map[OrderStatus]bool{
	OrderStatusOpen: {
		OrderStatusPartiallyFilled: true,
		OrderStatusFilled:          true,
		OrderStatusCancelled:       true,
		OrderStatusExpired:         true,
	},

	OrderStatusPartiallyFilled: {
		OrderStatusPartiallyFilled: true,
		OrderStatusFilled:          true,
		OrderStatusCancelled:       true,
		OrderStatusExpired:         true,
	},

	OrderStatusFilled:    {},
	OrderStatusCancelled: {},
	OrderStatusRejected:  {},
	OrderStatusExpired:   {},
}

type OrderType string

const (
	OrderTypeMarket OrderType = "MARKET"
	OrderTypeLimit  OrderType = "LIMIT"
)

type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

type Order struct {
	gorm.Model

	UserID            uint        `gorm:"type:bigint;not null;index:idx_order_user_id"`
	MarketID          uint        `gorm:"type:bigint;not null;index:idx_order_market_id"`
	Quantity          uint        `gorm:"type:bigint;not null"`
	RemainingQuantity uint        `gorm:"type:bigint;not null"`
	Price             uint        `gorm:"type:bigint;not null"`
	Status            OrderStatus `gorm:"type:varchar(32); not null"`
	Type              OrderType   `gorm:"type:varchar(32); not null"`
	Side              OrderSide   `gorm:"type:varchar(32); not null"`

	Market model.Market `gorm:"foreignkey:MarketID"`
}

func (o *Order) SetStatus(status OrderStatus) error {
	transitions, ok := allowedStatus[o.Status]
	if !ok {
		return shared.ErrOrderStatusNotAllowed
	}

	if transitions[status] != true {
		return shared.ErrOrderStatusNotAllowed
	}

	o.Status = status
	return nil
}

func (o *Order) ValidatePrice() error {
	if o.Quantity == 0 {
		return shared.ErrInvalidOrderAmount
	}

	if o.Price > math.MaxUint/o.Quantity {
		return shared.ErrInvalidOrderAmount
	}

	switch o.Type {
	case OrderTypeLimit:
		if o.Price <= 0 {
			return shared.ErrInvalidOrderAmount
		}
	case OrderTypeMarket:
		if o.Price != 0 {
			return shared.ErrOrderPriceForMarketTypeInvalid
		}
	}

	return nil
}

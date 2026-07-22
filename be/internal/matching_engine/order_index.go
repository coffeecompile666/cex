package matching_engine

import (
	template "icon_exchange/internal/order_book"
	"icon_exchange/internal/shared"
)

type IOrderIndex interface {
	Push(cmd Command) error
	Pop() Command
	GetLength() uint
}

type TradeType int

const (
	MarketTrade TradeType = iota
	LimitTrade  TradeType = iota
)

type CommandType int

const (
	NewOrder    CommandType = iota
	AmendOrder  CommandType = iota
	CancelOrder CommandType = iota
)

type Command struct {
	OrderID     uint
	Price       uint
	Quantity    uint
	Side        template.Side
	TradeType   TradeType
	CommandType CommandType
}

const MaxBufferSize = 100

type OrderIndex struct {
	oi chan Command
}

func NewOrderIndex() *OrderIndex {
	return &OrderIndex{
		oi: make(chan Command, MaxBufferSize),
	}
}

func (o *OrderIndex) Push(cmd Command) error {
	select {
	case o.oi <- cmd:
		return nil
	default:
		return shared.ErrQueueFull
	}
}

func (o *OrderIndex) Pop() Command {
	return <-o.oi
}

func (o *OrderIndex) GetLength() uint {
	return uint(len(o.oi))
}

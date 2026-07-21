package matching_engine

import (
	"fmt"
	template "icon_exchange/internal/order_book"
)

type IMatching interface {
	GetMarketID() uint
	Start()
}

type Matching struct {
	marketID   uint
	orderIndex IOrderIndex
	orderBook  template.IOrderBook
}

func NewMatching(marketID uint, orderIndex IOrderIndex, orderBook template.IOrderBook) *Matching {
	return &Matching{
		marketID:   marketID,
		orderIndex: orderIndex,
		orderBook:  orderBook,
	}
}

func (m *Matching) GetMarketID() uint {
	return m.marketID
}

func (m *Matching) Start() {
	go m.run()
}

func (m *Matching) Push(cmd Command) error {
	return m.orderIndex.Push(cmd)
}

func (m *Matching) run() {
	for {
		cmd := m.orderIndex.Pop()
		fmt.Println("Matching order ", cmd)
	}
}

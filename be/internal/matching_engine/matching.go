package matching_engine

import (
	"fmt"
	template "icon_exchange/internal/order_book"
	"icon_exchange/internal/shared"
)

type IMatching interface {
	GetMarketID() uint
	Start()
	Push(cmd Command) error
}

type Matching struct {
	marketID   uint
	orderIndex IOrderIndex
	orderBook  template.IOrderBook
}

type Matched struct {
	Quantity uint
	Price    uint
	OrderId  uint
}
type MatchingResult struct {
	OrderID uint
	Price   uint
	Matched []Matched
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
		m.match(cmd)
	}
}

func (m *Matching) match(cmd Command) {
	switch cmd.CommandType {
	case CommandTypeNew:
		isMarket := cmd.OrderType == shared.OrderTypeMarket
		isLimit := cmd.OrderType == shared.OrderTypeLimit
		isSell := cmd.OrderSide == shared.OrderSideSell
		isBuy := cmd.OrderSide == shared.OrderSideBuy

		var e error
		if isMarket && isSell {
			e = m.sellMarket(cmd)
			if e != nil {
				fmt.Println("Error matching order", e)
				// todo: send order match event
			} else {
				fmt.Println("Order matched", cmd.OrderID)
				// todo: send order match event
			}
		} else if isMarket && isBuy {
			result, e := m.buyMarket(cmd)
			if e != nil {
				fmt.Println("Error matching order", e)
				// todo: send order match event
			} else {
				fmt.Println("Order matched", result)
				// todo: send order match event
			}
		} else if isLimit && isSell {
			e = m.sellLimit(cmd)
			if e != nil {
				fmt.Println("Error matching order", e)
				// todo: send order match event
			} else {
				fmt.Println("Order matched", cmd.OrderID)
				// todo: send order match event
			}
		} else if isLimit && isBuy {
			e = m.buyLimit(cmd)
			if e != nil {
				fmt.Println("Error matching order", e)
				// todo: send order match event
			} else {
				fmt.Println("Order matched", cmd.OrderID)
				// todo: send order match event
			}
		}

	case CommandTypeCancel:
		e := m.orderBook.Remove(cmd.OrderID)
		if e != nil {
			_ = fmt.Errorf("error removing order %d: %s", cmd.OrderID, e)
		} else {
			fmt.Println("Order removed", cmd.OrderID)
			// todo: send order cancel event
		}
	case CommandTypeAmend:
		e := m.orderBook.Amend(template.AmendCommand{
			OrderID:  cmd.OrderID,
			Price:    cmd.Price,
			Quantity: cmd.Quantity,
		})
		if e != nil {
			_ = fmt.Errorf("error amending order %d: %s", cmd.OrderID, e)
		} else {
			fmt.Println("Order amended", cmd.OrderID)
			// todo: send order amend event
		}
	}
}

func (m *Matching) buyMarket(cmd Command) (*MatchingResult, error) {
	matchingResult := MatchingResult{
		OrderID: cmd.OrderID,
		Matched: []Matched{},
	}

	for cmd.Quantity > 0 {
		bestAsk := m.orderBook.GetBestAsk()

		if bestAsk == nil {
			return nil, shared.ErrOrderBookEmpty
		}

		for o := bestAsk.Orders.Front(); o != nil; o = o.Next() {

		}
	}

	return &matchingResult, nil
}

func (m *Matching) sellMarket(cmd Command) error {
	return nil
}

func (m *Matching) sellLimit(cmd Command) error {
	return nil
}

func (m *Matching) buyLimit(cmd Command) error {
	return nil
}

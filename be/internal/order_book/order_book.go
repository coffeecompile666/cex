package template

import (
	"container/list"
	"icon_exchange/internal/shared"

	"github.com/emirpasic/gods/trees/redblacktree"
)

type Side string

const (
	BUY  Side = "buy"
	SELL Side = "sell"
)

type Order struct {
	ID       uint
	Price    uint
	Quantity uint
	Side     Side
	Element  *list.Element
}

type SnapshotLevel struct {
	price  uint
	volume uint
}
type Snapshot struct {
	Bids []SnapshotLevel
	Asks []SnapshotLevel
}

type IOrderBook interface {
	Add(o *Order)
	Remove(orderID uint) error
	Amend(a AmendCommand) error
	GetBestAsk() *PriceLevel
	GetBestBid() *PriceLevel
	GetSnapshot() *Snapshot
}

type PriceLevel struct {
	Price  uint
	Volume uint
	Orders *list.List
}

type OrderEntry struct {
	order *Order
	level *PriceLevel
}

type OrderBook struct {
	bidBook *redblacktree.Tree
	askBook *redblacktree.Tree
	orders  map[uint]*OrderEntry
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		bidBook: redblacktree.NewWithIntComparator(),
		askBook: redblacktree.NewWithIntComparator(),
		orders:  make(map[uint]*OrderEntry),
	}
}

func (o *OrderBook) Add(order *Order) {
	var tree *redblacktree.Tree

	switch order.Side {
	case BUY:
		tree = o.bidBook
	case SELL:
		tree = o.askBook
	}

	if tree == nil {
		return
	}

	value, found := tree.Get(order.Price)

	var level *PriceLevel

	if found {
		level = value.(*PriceLevel)
	} else {
		level = &PriceLevel{
			Price:  order.Price,
			Volume: 0,
			Orders: list.New(),
		}
		tree.Put(order.Price, level)
	}

	level.Volume += order.Quantity

	order.Element = level.Orders.PushBack(order)

	o.orders[order.ID] = &OrderEntry{
		order: order,
		level: level,
	}
}

func (o *OrderBook) Remove(orderID uint) error {
	orderEntry, ok := o.orders[orderID]
	if !ok {
		// Todo: return error or just success
		return nil
	}

	level := orderEntry.level
	level.Orders.Remove(orderEntry.order.Element)
	level.Volume -= orderEntry.order.Quantity

	if level.Orders.Len() == 0 {
		if orderEntry.order.Side == BUY {
			o.bidBook.Remove(orderEntry.level.Price)
		} else {
			o.askBook.Remove(orderEntry.level.Price)
		}
	}

	delete(o.orders, orderID)
	return nil
}

type AmendCommand struct {
	OrderID  uint
	Price    uint
	Quantity uint
}

func (o *OrderBook) Amend(a AmendCommand) error {
	if a.Price < 1 || a.Quantity < 1 {
		return shared.ErrOrderInOrderBookInvalid
	}

	entry, ok := o.orders[a.OrderID]
	if !ok {
		return nil
	}

	e := o.Remove(a.OrderID)
	if e != nil {
		return e
	}

	entry.order.Quantity = a.Quantity
	entry.order.Price = a.Price
	o.Add(entry.order)

	return nil
}

func (o *OrderBook) GetBestAsk() *PriceLevel {
	node := o.askBook.Left()
	if node == nil {
		return nil
	}
	return node.Value.(*PriceLevel)
}

func (o *OrderBook) GetBestBid() *PriceLevel {
	node := o.bidBook.Right()
	if node == nil {
		return nil
	}
	return node.Value.(*PriceLevel)
}

func (o *OrderBook) GetSnapshot() *Snapshot {
	snapshot := &Snapshot{
		Bids: make([]SnapshotLevel, o.bidBook.Size()),
		Asks: make([]SnapshotLevel, o.askBook.Size()),
	}

	// bid: high -> low
	it := o.bidBook.Iterator()
	for it.End(); it.Prev(); {
		level := it.Value().(*PriceLevel)
		snapshot.Bids = append(snapshot.Bids, SnapshotLevel{price: level.Price, volume: level.Volume})
	}

	// ask: low -> high
	it = o.askBook.Iterator()
	for it.Begin(); it.Next(); {
		level := it.Value().(*PriceLevel)
		snapshot.Asks = append(snapshot.Asks, SnapshotLevel{price: level.Price, volume: level.Volume})
	}

	return snapshot
}

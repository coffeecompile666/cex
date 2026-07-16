package matching_engine

type IOrderBookManager interface {
	PlaceOrder()
	CancelOrder()
	GetOrderBook()
}

type OrderBookManager struct {
	OrderBook OrderBook
}

type OrderBook struct {
	orders []Order
}

type Order struct {
	name     string
	price    float64
	quantity int
}

type MatchEngine interface {
}

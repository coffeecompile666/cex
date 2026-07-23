package shared

type Response[T any] struct {
	Data T
}

type Empty struct{}

type OrderType string

const (
	OrderTypeLimit  OrderType = "LIMIT"
	OrderTypeMarket OrderType = "MARKET"
)

type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

type TradeRole string

const (
	TradeRoleTaker TradeRole = "TAKER"
	TradeRoleMaker TradeRole = "MAKER"
)

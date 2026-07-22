package matching_engine

import (
	template "icon_exchange/internal/order_book"
	"icon_exchange/internal/shared"
)

type MatchingEngine struct {
	markets map[uint]*MarketMatchingInstance
}

type MarketMatchingInstance struct {
	marketID uint
	matching *Matching
}

func NewMatchingEngine() *MatchingEngine {
	btc := uint(1)

	oi := NewOrderIndex()
	ob := template.NewOrderBook()

	btcMatching := NewMatching(btc, oi, ob)

	markets := make(map[uint]*MarketMatchingInstance)

	markets[btc] = &MarketMatchingInstance{
		marketID: btc,
		matching: btcMatching,
	}

	return &MatchingEngine{
		markets: markets,
	}
}

func (engine *MatchingEngine) Start() {
	for _, market := range engine.markets {
		market.matching.Start()
	}
}

func (engine *MatchingEngine) PushOrder(marketID uint, cmd Command) error {
	market, ok := engine.markets[marketID]
	if !ok {
		return shared.ErrMarketNotFound
	}

	if cmd.Price < 1 || cmd.Quantity < 1 {
		return shared.ErrOrderInOrderBookInvalid
	}

	return market.matching.Push(cmd)
}

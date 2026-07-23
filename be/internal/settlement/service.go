package settlement

import "context"

type MatchResult struct {
	TakerOrderID uint `json:"taker_order_id"`
	MakerOrderID uint `json:"maker_order_id"`
	Price        uint `json:"price"`
	Quantity     uint `json:"quantity"`
	TakerUserID  uint `json:"taker_user_id"`
	MakerUserID  uint `json:"maker_user_id"`
}

// ISettlementService định nghĩa các hành động xử lý sau khi khớp lệnh thành công.
type ISettlementService interface {
	ProcessSettlement(ctx context.Context, result MatchResult) error
}

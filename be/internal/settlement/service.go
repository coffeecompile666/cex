package settlement

import (
	service4 "icon_exchange/internal/asset/service"
	"icon_exchange/internal/ledger/service"
	orderService "icon_exchange/internal/order/service"
	"icon_exchange/internal/shared"
	service2 "icon_exchange/internal/trade/service"

	"gorm.io/gorm"
)

type MatchResult struct {
	TakerOrderID uint
	MakerOrderID uint
	Price        uint
	Quantity     uint
	TakerRole    shared.OrderSide
}

// ISettlementService định nghĩa các hành động xử lý sau khi khớp lệnh thành công.
type ISettlementService interface {
	ProcessSettlement(result MatchResult) error
}

type Service struct {
	repo          IRepository
	orderService  orderService.IOrderService
	tradeService  service2.ITradeService
	assetService  service4.IAssetService
	ledgerService service.ILedgerService
}

func NewSettlementService(repo IRepository, os orderService.IOrderService, ts service2.ITradeService, as service4.IAssetService, ls service.ILedgerService) *Service {
	return &Service{
		repo:          repo,
		orderService:  os,
		tradeService:  ts,
		assetService:  as,
		ledgerService: ls,
	}
}

func (s Service) ProcessSettlement(r MatchResult) error {
	err := s.repo.WithTransaction(func(tx *gorm.DB) error {
		var buyerUserID uint
		var sellerUserID uint
		var buyerOrderID uint
		var sellerOrderID uint

		if r.TakerRole == shared.OrderSideBuy {
			buyerUserID = r.TakerOrderID
			sellerUserID = r.MakerOrderID
			buyerOrderID = r.TakerOrderID
			sellerOrderID = r.MakerOrderID
		} else {
			buyerUserID = r.MakerOrderID
			sellerUserID = r.TakerOrderID
			buyerOrderID = r.MakerOrderID
			sellerOrderID = r.TakerOrderID
		}

		// -> write ledger asset moving
		bo, err := s.orderService.GetByIDForUpdate(tx, buyerOrderID)
		if err != nil {
			return err
		}
		so, err := s.orderService.GetByIDForUpdate(tx, sellerOrderID)
		if err != nil {
			return err
		}

		// -> moving asset between maker and taker
		// -> update order
		// -> write trade
		return nil

	})

	if err != nil {
		return err
	}

	return nil
}

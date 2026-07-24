package service

import (
	service2 "icon_exchange/internal/asset/service"
	"icon_exchange/internal/market/service"
	"icon_exchange/internal/matching_engine"
	"icon_exchange/internal/order/model"
	"icon_exchange/internal/order/repository"
	"icon_exchange/internal/shared"
	service3 "icon_exchange/internal/user/service"

	"gorm.io/gorm"
)

type IOrderService interface {
	CreateOrder(cmd CreateOrderCommand) (uint, error)
	CancelOrder(cmd CancelOrderCommand) (uint, error)
	AmendOrder(cmd AmendOrderCommand) (uint, error)
	GetByIDForUpdate(tx *gorm.DB, id uint) (*model.Order, error)
}
type Service struct {
	repo           *repository.Repository
	marketService  service.IMarketService
	assetService   service2.IAssetService
	userService    service3.IUserService
	matchingEngine *matching_engine.MatchingEngine
}

func NewOrderService(repo *repository.Repository, marketService service.IMarketService,
	assetService service2.IAssetService, userService service3.IUserService,
	matchingEngine *matching_engine.MatchingEngine) *Service {
	return &Service{
		repo:           repo,
		marketService:  marketService,
		assetService:   assetService,
		userService:    userService,
		matchingEngine: matchingEngine,
	}
}

type CreateOrderCommand struct {
	UserID        uint
	MarketID      uint
	FloatQuantity float64
	FloatPrice    float64
	Type          shared.OrderType
	Side          shared.OrderSide
}

type CancelOrderCommand struct {
	UserID  uint
	OrderID uint
}

type AmendOrderCommand struct {
	UserID        uint
	OrderID       uint
	FloatPrice    float64
	FloatQuantity float64
}

func (s Service) CreateOrder(cmd CreateOrderCommand) (uint, error) {
	var orderID uint

	err := s.repo.WithTransaction(func(tx *gorm.DB) error {
		// -> validate
		market, err := s.marketService.GetMarketByID(cmd.MarketID)
		if err != nil {
			return err
		}

		quantity, err := market.ToSmallestUnit(cmd.FloatQuantity)
		if err != nil {
			return err
		}

		baseCurrency, err := s.marketService.GetBaseCurrency()
		if err != nil {
			return err
		}

		price, err := baseCurrency.ToSmallestUnit(cmd.FloatPrice)
		if err != nil {
			return err
		}

		err = shared.CheckOverFlowUintWithMulOperator(quantity, price)
		if err != nil {
			return err
		}
		totalPrice := quantity * price

		// -> business validate
		baseAsset, err := s.assetService.GetAssetByMarketID(tx, cmd.UserID, baseCurrency.ID)
		if baseAsset.GetAvailableAmount() < totalPrice {
			return shared.ErrBalanceNotSufficient
		}

		err = baseAsset.LockAmount(totalPrice)
		if err != nil {
			return err
		}

		// -> create order
		order, err := s.repo.Create(tx, &model.Order{
			UserID:            cmd.UserID,
			MarketID:          cmd.MarketID,
			Quantity:          quantity,
			RemainingQuantity: quantity,
			Price:             price,
			TotalPrice:        totalPrice,
			Status:            model.OrderStatusOpen,
			Type:              cmd.Type,
			Side:              cmd.Side,
		})

		if err != nil {
			return err
		}

		orderID = order.ID

		// -> push to matching engine
		err = s.matchingEngine.PushOrder(order.MarketID, matching_engine.Command{
			OrderID:     order.ID,
			Price:       order.Price,
			Quantity:    order.Quantity,
			OrderSide:   order.Side,
			OrderType:   order.Type,
			CommandType: matching_engine.CommandTypeNew,
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return orderID, nil
}

func (s Service) CancelOrder(cmd CancelOrderCommand) error {
	err := s.repo.WithTransaction(func(tx *gorm.DB) error {
		// -> business validate
		order, err := s.repo.GetByID(cmd.UserID, cmd.OrderID)
		if err != nil {
			return err
		}

		err = order.IsAvailableToCancel()
		if err != nil {
			return err
		}
		// -> push command to matching engine
		err = s.matchingEngine.PushOrder(order.ID, matching_engine.Command{
			OrderID:     0,
			CommandType: matching_engine.CommandTypeCancel,
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (s Service) AmendOrder(cmd AmendOrderCommand) error {
	err := s.repo.WithTransaction(func(tx *gorm.DB) error {
		// -> validate
		order, err := s.repo.GetByID(cmd.UserID, cmd.OrderID)
		if err != nil {
			return err
		}

		err = order.IsAvailableToAmend()
		if err != nil {
			return err
		}

		market, err := s.marketService.GetMarketByID(order.MarketID)

		quantity, err := market.ToSmallestUnit(cmd.FloatQuantity)
		if err != nil {
			return err
		}

		baseCurrency, err := s.marketService.GetBaseCurrency()
		if err != nil {
			return err
		}

		price, err := baseCurrency.ToSmallestUnit(cmd.FloatPrice)
		if err != nil {
			return err
		}

		err = shared.CheckOverFlowUintWithMulOperator(quantity, price)
		if err != nil {
			return err
		}

		// -> push command to matching engine
		err = s.matchingEngine.PushOrder(order.ID, matching_engine.Command{
			OrderID:     order.ID,
			Price:       price,
			Quantity:    quantity,
			CommandType: matching_engine.CommandTypeAmend,
		})
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetByIDForUpdate(tx *gorm.DB, id uint) (*model.Order, error) {
	return s.repo.GetByIDForUpdate(tx, id)
}

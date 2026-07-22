package service

import (
	"icon_exchange/internal/order/model"
	"icon_exchange/internal/order/repository"
)

type IOrderService interface {
	CreateOrder(order *model.Order) (*model.Order, error)
	CancelOrder(order *model.Order) (*model.Order, error)
	AmendOrder(order *model.Order) (*model.Order, error)
}
type Service struct {
	repo *repository.Repository
}

func NewOrderService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

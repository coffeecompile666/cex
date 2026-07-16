package service

import (
	"icon_exchange/internal/order/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewOrderService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

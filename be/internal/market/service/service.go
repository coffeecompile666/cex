package service

import (
	"icon_exchange/internal/market/model"
	"icon_exchange/internal/market/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewMarketService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetMarketBySymbol(symbol string) (*model.Market, error) {
	return s.repo.GetBySymbol(symbol)
}

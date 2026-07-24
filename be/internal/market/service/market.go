package service

import (
	"icon_exchange/internal/market/model"
	"icon_exchange/internal/market/repository"
)

type IMarketService interface {
	GetMarketByID(id uint) (*model.Market, error)
	GetBaseCurrency() (*model.Market, error)
	GetAll() ([]*model.Market, error)
}

type MarketService struct {
	repo *repository.Repository
}

func NewMarketService(repo *repository.Repository) *MarketService {
	return &MarketService{repo: repo}
}

func (m MarketService) GetMarketByID(id uint) (*model.Market, error) {
	return m.repo.GetMarketByID(id)
}

func (m MarketService) GetBaseCurrency() (*model.Market, error) {
	return m.repo.GetBaseCurrency()
}

func (m MarketService) GetAll() ([]*model.Market, error) {
	return m.repo.GetAll()
}

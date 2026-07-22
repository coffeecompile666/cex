package service

import (
	"icon_exchange/internal/asset/model"
	"icon_exchange/internal/asset/repository"

	"gorm.io/gorm"
)

type IAssetService interface {
	GetAssetByMarketID(tx *gorm.DB, userID uint, assetID uint) (model.Asset, error)
}

type Service struct {
	repo *repository.Repository
}

func (s *Service) GetAssetByMarketID(tx *gorm.DB, userID uint, assetID uint) (*model.Asset, error) {
	return s.repo.GetByID(tx, userID, assetID)
}

func NewAssetService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

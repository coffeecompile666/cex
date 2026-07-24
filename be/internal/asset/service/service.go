package service

import (
	"icon_exchange/internal/asset/model"
	"icon_exchange/internal/asset/repository"
	"icon_exchange/internal/market/service"

	"gorm.io/gorm"
)

type IAssetService interface {
	GetAssetByMarketID(tx *gorm.DB, userID uint, assetID uint) (model.Asset, error)
	GetByIDForUpdate(tx *gorm.DB, userID uint, assetID uint) (model.Asset, error)
	SaveAsset(tx *gorm.DB, asset model.Asset) error
	InitUserAssets(tx *gorm.DB, userID uint) error
	LockAmount(tx *gorm.DB, userID uint, assetID uint, amount uint) error
	UnlockAmount(tx *gorm.DB, userID uint, assetID uint, amount uint) error
	Debit(tx *gorm.DB, userID uint, amount uint) error
	Credit(tx *gorm.DB, amount uint) error
}

type Service struct {
	repo          *repository.Repository
	marketService service.IMarketService
}

func NewAssetService(repo *repository.Repository, marketService service.IMarketService) *Service {
	return &Service{repo: repo, marketService: marketService}
}

func (s Service) GetAssetByMarketID(tx *gorm.DB, userID uint, assetID uint) (*model.Asset, error) {
	return s.repo.GetByID(tx, userID, assetID)
}

func (s Service) GetByIDForUpdate(tx *gorm.DB, userID uint, assetID uint) (*model.Asset, error) {
	return s.repo.GetForUpdate(tx, userID, assetID)
}

func (s Service) SaveAsset(tx *gorm.DB, asset model.Asset) error {
	return s.repo.Update(tx, &asset)
}

func (s Service) InitUserAssets(tx *gorm.DB, userID uint) error {
	markets, err := s.marketService.GetAll()
	if err != nil {
		return err
	}

	for _, market := range markets {
		err = s.repo.Save(tx, &model.Asset{
			UserID:       userID,
			MarketID:     market.ID,
			Amount:       0,
			LockedAmount: 0,
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func (s Service) LockAmount(tx *gorm.DB, userID uint, assetID uint, amount uint) error {
	asset, err := s.repo.GetForUpdate(tx, userID, assetID)
	if err != nil {
		return err
	}
	err = asset.LockAmount(amount)
	if err != nil {
		return err
	}
	return s.repo.Save(tx, asset)
}

func (s Service) UnlockAmount(tx *gorm.DB, userID uint, assetID uint, amount uint) error {
	asset, err := s.repo.GetForUpdate(tx, userID, assetID)
	if err != nil {
		return err
	}
	err = asset.UnLock(amount)
	if err != nil {
		return err
	}
	return s.repo.Save(tx, asset)
}

func (s Service) Debit(tx *gorm.DB, userID uint, assetID uint, amount uint) error {
	// -> update asset
	asset, err := s.repo.GetForUpdate(tx, userID, assetID)
	if err != nil {
		return err
	}
	err = asset.Debit(amount)
	if err != nil {
		return err
	}

	err = s.repo.Save(tx, asset)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) Credit(tx *gorm.DB, userID uint, assetID uint, amount uint) error {
	asset, err := s.repo.GetForUpdate(tx, userID, assetID)
	if err != nil {
		return err
	}
	asset.Credit(amount)
	return s.repo.Save(tx, asset)
}

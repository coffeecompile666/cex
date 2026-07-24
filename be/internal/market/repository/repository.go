package repository

import (
	"errors"
	"icon_exchange/internal/market/model"
	"icon_exchange/internal/shared"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewMarketRepo(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) WithTransaction(fn func(tx *gorm.DB) error) error {
	return repo.db.Transaction(fn)
}

func (repo *Repository) GetAll() ([]*model.Market, error) {
	var markets []*model.Market
	if err := repo.db.Find(&markets).Error; err != nil {
		return nil, shared.ErrInternalServerError
	}
	return markets, nil
}

func (repo *Repository) GetByID(id uint) (*model.Market, error) {
	var market model.Market
	if err := repo.db.First(&market, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNotFound
		}
		return nil, shared.ErrInternalServerError
	}
	return &market, nil
}

func (repo *Repository) GetBySymbol(symbol string) (*model.Market, error) {
	var market model.Market
	if err := repo.db.Where("symbol = ?", symbol).First(&market).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNotFound
		}
		return nil, shared.ErrInternalServerError
	}
	return &market, nil
}

func (repo *Repository) GetMarkets() ([]model.Market, error) {
	var markets []model.Market
	if err := repo.db.Where("IsBaseCurrency = ?", false).Find(&markets).Error; err != nil {
		return nil, shared.ErrInternalServerError
	}
	return markets, nil
}

func (repo *Repository) GetMarketByID(id uint) (*model.Market, error) {
	var market model.Market
	if err := repo.db.Where("IsBaseCurrency", false).First(&market, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		}
	}
	return &market, nil
}

func (repo *Repository) GetBaseCurrencyByID(id uint) (*model.Market, error) {
	var market model.Market
	if err := repo.db.Where("IsBaseCurrency", true).First(&market, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		}
	}
	return &market, nil
}

func (repo *Repository) GetBaseCurrency() (*model.Market, error) {
	var market model.Market
	if err := repo.db.Where("IsBaseCurrency", true).First(&market).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		}
	}
	return &market, nil
}

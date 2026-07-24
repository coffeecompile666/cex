package repository

import (
	"errors"
	model2 "icon_exchange/internal/asset/model"
	"icon_exchange/internal/shared"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) WithTransaction(fn func(tx *gorm.DB) error) error {
	return repo.db.Transaction(fn)
}

func (repo *Repository) GetByID(tx *gorm.DB, userID uint, marketID uint) (*model2.Asset, error) {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = repo.db
	}

	asset := &model2.Asset{}

	e := db.Where("userID = ? AND marketID = ?", userID, marketID).Find(asset).Error
	if e != nil {
		return nil, shared.ErrAssetNotFound
	}

	return asset, nil
}

func (repo *Repository) GetByUserID(tx *gorm.DB, userID uint) ([]*model2.Asset, error) {
	var assets []*model2.Asset

	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = repo.db
	}

	e := db.Where("userID = ?", userID).Find(&assets).Error
	if e != nil {
		return nil, shared.ErrAssetNotFound
	}

	return assets, nil
}

func (repo *Repository) GetForUpdate(tx *gorm.DB, userID uint, assetID uint) (*model2.Asset, error) {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = repo.db
	}
	asset := &model2.Asset{}
	err := db.Where("userID = ? AND assetID = ?", userID, assetID).Find(asset).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrAssetNotFound
		}
		return nil, err
	}
	return asset, nil
}

func (repo *Repository) Create(tx *gorm.DB, asset *model2.Asset) error {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = repo.db
	}

	err := db.Create(asset).Error
	if err != nil {
		return shared.ErrInternalServerError
	}
	return nil
}

func (repo *Repository) Update(tx *gorm.DB, asset *model2.Asset) error {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = repo.db
	}

	err := db.Save(asset).Error
	if err != nil {
		return shared.ErrInternalServerError
	}
	return nil
}

func (repo *Repository) Save(tx *gorm.DB, asset *model2.Asset) error {
	var db *gorm.DB
	if tx != nil {
		db = tx

	} else {
		db = repo.db
	}
	err := db.Save(asset).Error
	if err != nil {
		return shared.ErrInternalServerError
	}
	return nil
}

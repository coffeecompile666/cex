package repository

import (
	"errors"
	"icon_exchange/internal/ledger/model"
	"icon_exchange/internal/shared"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewLedgerRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) WithTransaction(fn func(tx *gorm.DB) error) error {
	return repo.db.Transaction(fn)
}

// Write save all ledger entries to db
func (repo *Repository) Write(journalEntry *model.JournalEntry, tx *gorm.DB) error {
	db := repo.db
	if tx != nil {
		db = tx
	}

	ledgerEntries := journalEntry.LedgerEntries
	return db.Create(ledgerEntries).Error
}

func (repo *Repository) GetAccountCondition(tx *gorm.DB, condition *model.Account) (*model.Account, error) {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = repo.db
	}
	var account model.Account
	err := db.Where(&condition).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNotFound
		}
		return nil, shared.ErrInternalServerError
	}
	return &account, nil
}

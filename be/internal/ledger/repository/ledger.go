package repository

import (
	"icon_exchange/internal/ledger/model"

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

package service

import (
	"icon_exchange/internal/ledger/model"
	"icon_exchange/internal/ledger/repository"

	"gorm.io/gorm"
)

type Service struct {
	repo *repository.Repository
}

func NewLedgerService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// Write log a transaction to ledger
func (s *Service) Write(journalCommand []model.JournalCommand, tx *gorm.DB) error {
	journalEntry, err := model.NewJournalEntry(journalCommand)
	if err != nil {
		return err
	}

	return s.repo.Write(journalEntry, tx)
}

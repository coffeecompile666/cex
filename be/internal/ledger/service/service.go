package service

import (
	"icon_exchange/internal/ledger/model"
	"icon_exchange/internal/ledger/repository"

	"gorm.io/gorm"
)

type ILedgerService interface {
	Write(tx *gorm.DB, journalCommand []model.JournalCommand) error
	GetAccount(tx *gorm.DB, conditions *model.Account) (*model.Account, error)
}

type Service struct {
	repo *repository.Repository
}

func NewLedgerService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// Write log a transaction to ledger
func (s *Service) Write(tx *gorm.DB, journalCommand []model.JournalCommand) error {
	journalEntry, err := model.NewJournalEntry(journalCommand)
	if err != nil {
		return err
	}

	return s.repo.Write(journalEntry, tx)
}

func (s *Service) GetAccountConditions(tx *gorm.DB, conditions *model.Account) (*model.Account, error) {
	return s.repo.GetAccountCondition(tx, conditions)
}

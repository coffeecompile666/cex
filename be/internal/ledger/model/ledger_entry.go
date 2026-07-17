package model

import (
	"icon_exchange/internal/shared"

	"gorm.io/gorm"
)

type LedgerSide string

const (
	SideDebit  LedgerSide = "DEBIT"
	SideCredit LedgerSide = "CREDIT"
)

type LedgerEntry struct {
	gorm.Model
	Amount    uint       `gorm:"not null"`
	Side      LedgerSide `gorm:"not null"`
	AccountID uint       `gorm:"not null;uniqueIndex:uk_account"`
	Account   Account    `gorm:"foreignKey:AccountID"`
}

type JournalCommand struct {
	Amount    uint
	Side      LedgerSide
	AccountID uint
}

type JournalEntry struct {
	LedgerEntries []LedgerEntry
}

func NewJournalEntry(command []JournalCommand) (*JournalEntry, error) {
	var debit, credit uint
	debit = 0
	credit = 0

	for _, command := range command {
		switch command.Side {
		case SideDebit:
			debit += command.Amount
		case SideCredit:
			credit += command.Amount
		}
	}

	if debit != credit {
		return nil, shared.ErrJournalEntryInvalid
	}

	ledgerEntries := make([]LedgerEntry, len(command))
	for i, command := range command {
		ledgerEntries[i] = LedgerEntry{
			Amount:    command.Amount,
			Side:      command.Side,
			AccountID: command.AccountID,
		}
	}

	return &JournalEntry{
		LedgerEntries: ledgerEntries,
	}, nil
}

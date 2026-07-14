package module_wallet

import (
	"log"

	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateWalletForUser creates a new empty wallet for a user.
// This is an internal API meant to be called by other modules (e.g. user module).
func (s *Service) CreateWalletForUser(userID uint, tx *gorm.DB) error {
	wallet := &Wallet{
		UserID:  userID,
		Balance: 0.0,
	}
	err := s.repo.Create(wallet, tx)
	if err != nil {
		log.Printf("Failed to create wallet for user %d: %v", userID, err)
		return err
	}
	return nil
}

// GetBalance returns the user's wallet
func (s *Service) GetBalance(userID uint) (*Wallet, error) {
	return s.repo.GetByUserID(userID)
}

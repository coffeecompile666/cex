package module_user

import (
	"log"

	"gorm.io/gorm"
)

// WalletServiceInterface defines the contract that module_user needs from module_wallet.
// This is optional for high velocity (you could just inject *module_wallet.Service directly),
// but using an interface is a good practice to keep modules slightly decoupled.
type WalletServiceInterface interface {
	CreateWalletForUser(userID uint, tx *gorm.DB) error
}

// MailerServiceInterface defines the contract that module_user needs from module_mailer.
type MailerServiceInterface interface {
	SendWelcomeEmail(toEmail string, username string) error
}

type Service struct {
	repo          *Repository
	walletService WalletServiceInterface
	mailerService MailerServiceInterface
}

func NewService(repo *Repository, walletService WalletServiceInterface, mailerService MailerServiceInterface) *Service {
	return &Service{
		repo:          repo,
		walletService: walletService,
		mailerService: mailerService,
	}
}

func (s *Service) CreateUser(username, email string) (*User, error) {
	user := &User{
		Username: username,
		Email:    email,
	}

	// Execute everything within a transaction
	err := s.repo.WithTransaction(func(tx *gorm.DB) error {
		// 1. Create the user in the database
		if err := s.repo.Create(user, tx); err != nil {
			log.Printf("Failed to create user: %v", err)
			return err
		}

		// 2. Inter-module communication: Ask Wallet module to create a wallet
		if err := s.walletService.CreateWalletForUser(user.ID, tx); err != nil {
			log.Printf("Warning: User %d created, but wallet creation failed: %v", user.ID, err)
			return err // Return error to rollback transaction
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 3. Inter-module communication: Ask Mailer module to send welcome email
	// This can be done asynchronously in a goroutine to not block the API response
	go func() {
		err := s.mailerService.SendWelcomeEmail(user.Email, user.Username)
		if err != nil {
			log.Printf("Warning: Failed to send welcome email: %v", err)
		}
	}()

	return user, nil
}

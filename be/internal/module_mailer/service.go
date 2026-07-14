package module_mailer

import (
	"log"
)

type Service struct {
	// Here you would inject SMTP configuration, SendGrid client, AWS SES, etc.
}

func NewService() *Service {
	return &Service{}
}

// SendWelcomeEmail sends a welcome email to the newly registered user.
// In a real application, this might use an external API and should ideally be asynchronous.
func (s *Service) SendWelcomeEmail(toEmail string, username string) error {
	// Dummy implementation
	log.Printf("[MAILER] 📧 Sending welcome email to %s (%s)...", username, toEmail)
	log.Printf("[MAILER] ✅ Email sent successfully to %s", toEmail)

	return nil
}

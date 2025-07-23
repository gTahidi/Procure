package services

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
)

// EmailService defines the interface for email-related operations
type EmailService interface {
	SendPasswordResetEmail(to string, resetLink string) error
}

// SMTPEmailService implements EmailService using SMTP
type SMTPEmailService struct {
	host     string
	port     string
	username string
	password string
	from     string
}

// NewSMTPEmailService creates a new SMTPEmailService
func NewSMTPEmailService() (*SMTPEmailService, error) {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	if host == "" || port == "" || username == "" || password == "" || from == "" {
		return nil, fmt.Errorf("missing SMTP configuration")
	}

	return &SMTPEmailService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}, nil
}

// SendPasswordResetEmail sends a password reset email
func (s *SMTPEmailService) SendPasswordResetEmail(to string, resetLink string) error {
	// Email template
	const emailTemplate = `
Subject: Password Reset Request

Hello,

You have requested to reset your password. Please click the link below to reset your password:

{{.ResetLink}}

This link will expire in 1 hour.

If you did not request a password reset, please ignore this email.

Best regards,
The Procurement Team
`

	// Parse the template
	tmpl, err := template.New("resetEmail").Parse(emailTemplate)
	if err != nil {
		return err
	}

	// Prepare the data for the template
	data := struct {
		ResetLink string
	}{
		ResetLink: resetLink,
	}

	// Execute the template
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	// Set up authentication
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	// Send the email
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	return smtp.SendMail(addr, auth, s.from, []string{to}, body.Bytes())
}

// MockEmailService is a mock implementation of EmailService for testing or development
type MockEmailService struct {
	LogEmails bool
}

// NewMockEmailService creates a new MockEmailService
func NewMockEmailService() *MockEmailService {
	return &MockEmailService{
		LogEmails: true,
	}
}

// SendPasswordResetEmail logs the email instead of sending it
func (s *MockEmailService) SendPasswordResetEmail(to string, resetLink string) error {
	if s.LogEmails {
		fmt.Printf("Mock email sent to %s with reset link: %s\n", to, resetLink)
	}
	return nil
}

// GetEmailService returns the appropriate email service based on environment
func GetEmailService() (EmailService, error) {
	// Always use mock service for now
	return NewMockEmailService(), nil
}
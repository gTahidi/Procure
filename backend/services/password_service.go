package services

import (
	"errors"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// PasswordService defines the interface for password-related operations
type PasswordService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
	ValidatePasswordStrength(password string) error
}

// BCryptPasswordService implements PasswordService using bcrypt
type BCryptPasswordService struct {
	minLength      int
	requireUpper   bool
	requireLower   bool
	requireNumber  bool
	requireSpecial bool
	cost          int
}

// NewBCryptPasswordService creates a new BCryptPasswordService
func NewBCryptPasswordService() *BCryptPasswordService {
	return &BCryptPasswordService{
		minLength:      8,
		requireUpper:   true,
		requireLower:   true,
		requireNumber:  true,
		requireSpecial: false, // Optional for better UX
		cost:          12,     // Higher cost means more secure but slower
	}
}

// HashPassword hashes a password using bcrypt
func (s *BCryptPasswordService) HashPassword(password string) (string, error) {
	if err := s.ValidatePasswordStrength(password); err != nil {
		return "", err
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPassword checks if the provided password matches the hashed password
func (s *BCryptPasswordService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// ValidatePasswordStrength checks if the password meets the strength requirements
func (s *BCryptPasswordService) ValidatePasswordStrength(password string) error {
	if len(password) < s.minLength {
		return errors.New("password must be at least 8 characters long")
	}

	if s.requireUpper && !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	if s.requireLower && !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	if s.requireNumber && !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return errors.New("password must contain at least one number")
	}

	if s.requireSpecial && !regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password) {
		return errors.New("password must contain at least one special character")
	}

	// Check for common passwords (this is a very basic check)
	commonPasswords := []string{"password", "123456", "qwerty", "admin"}
	passwordLower := strings.ToLower(password)
	for _, common := range commonPasswords {
		if passwordLower == common {
			return errors.New("password is too common")
		}
	}

	return nil
}
package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"unicode"

	"golang.org/x/crypto/bcrypt"

	"github.com/mnuddindev/jutsu-api/internal/domain"
)

// PasswordService handles password operations
type PasswordService struct {
	cost int
}

// NewPasswordService creates a new password service
func NewPasswordService(cost int) *PasswordService {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &PasswordService{
		cost: cost,
	}
}

// Hash generates a bcrypt hash of the password
func (s *PasswordService) Hash(password string) (string, error) {
	if len(password) > 72 {
		return "", domain.ErrPasswordTooLong
	}

	if len(password) < 8 {
		return "", domain.ErrPasswordTooShort
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(bytes), nil
}

// Verify checks if the password matches the hash
func (s *PasswordService) ComparePass(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// IsValid checks if the password matches the hash and returns boolean
func (s *PasswordService) IsValid(password, hash string) bool {
	return s.ComparePass(password, hash) == nil
}

// ValidatePasswordStrength validates password strength
func (s *PasswordService) ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return domain.ErrPasswordTooShort
	}

	if len(password) > 72 {
		return domain.ErrPasswordTooLong
	}

	var (
		hasUpper  = false
		hasLower  = false
		hasNumber = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber {
		return domain.ErrWeakPassword
	}

	return nil
}

// GenerateRandomPassword generates a random secure password
func (s *PasswordService) GenerateRandomPassword(length int) (string, error) {
	if length < 8 {
		length = 8
	}
	if length > 72 {
		length = 72
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random password: %w", err)
	}

	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}

	return string(bytes), nil
}

// GenerateRandomToken generates a secure random token
func (s *PasswordService) GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

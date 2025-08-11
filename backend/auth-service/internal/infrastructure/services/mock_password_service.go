package services

import (
	"golang.org/x/crypto/bcrypt"
)

// MockPasswordService implements PasswordService for testing
type MockPasswordService struct{}

// NewMockPasswordService creates a new mock password service
func NewMockPasswordService() *MockPasswordService {
	return &MockPasswordService{}
}

// HashPassword hashes a password using bcrypt
func (m *MockPasswordService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPassword verifies a password against its hash
func (m *MockPasswordService) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

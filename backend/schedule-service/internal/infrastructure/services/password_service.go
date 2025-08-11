package services

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

// PasswordService interface for password operations
type PasswordService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) bool
	GenerateRandomPassword(length int) string
}

// PasswordServiceImpl implements the PasswordService interface
type PasswordServiceImpl struct {
	cost int
}

// NewPasswordService creates a new password service
func NewPasswordService(cost int) PasswordService {
	return &PasswordServiceImpl{
		cost: cost,
	}
}

// HashPassword hashes a password using bcrypt
func (p *PasswordServiceImpl) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// VerifyPassword verifies a password against its hash
func (p *PasswordServiceImpl) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// GenerateRandomPassword generates a random password of specified length
func (p *PasswordServiceImpl) GenerateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"

	password := make([]byte, length)
	for i := range password {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password[i] = charset[randomIndex.Int64()]
	}

	return string(password)
}

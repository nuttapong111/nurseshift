package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
)

// PasswordResetService defines the interface for password reset operations
type PasswordResetService interface {
	GenerateResetToken() (string, error)
	ValidateResetToken(token string) bool
	StoreResetToken(userID uuid.UUID, token string) error
	GetUserIDByToken(token string) (uuid.UUID, error)
	ClearResetToken(token string) error
}

// InMemoryPasswordResetService implements PasswordResetService using in-memory storage
// In production, you should use Redis or database for this
type InMemoryPasswordResetService struct {
	tokens map[string]resetTokenData
}

type resetTokenData struct {
	UserID    uuid.UUID
	ExpiresAt time.Time
}

// NewInMemoryPasswordResetService creates a new in-memory password reset service
func NewInMemoryPasswordResetService() PasswordResetService {
	return &InMemoryPasswordResetService{
		tokens: make(map[string]resetTokenData),
	}
}

// GenerateResetToken generates a random 6-digit token
func (s *InMemoryPasswordResetService) GenerateResetToken() (string, error) {
	// Generate 6 random digits
	token := ""
	for i := 0; i < 6; i++ {
		randomNum, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		token += fmt.Sprintf("%d", randomNum.Int64())
	}
	return token, nil
}

// ValidateResetToken checks if a reset token is valid and not expired
func (s *InMemoryPasswordResetService) ValidateResetToken(token string) bool {
	data, exists := s.tokens[token]
	if !exists {
		return false
	}
	
	// Check if token is expired
	if time.Now().After(data.ExpiresAt) {
		// Remove expired token
		delete(s.tokens, token)
		return false
	}
	
	return true
}

// StoreResetToken stores a reset token with expiration (15 minutes)
func (s *InMemoryPasswordResetService) StoreResetToken(userID uuid.UUID, token string) error {
	s.tokens[token] = resetTokenData{
		UserID:    userID,
		ExpiresAt: time.Now().Add(15 * time.Minute), // Token expires in 15 minutes
	}
	return nil
}

// GetUserIDByToken retrieves the user ID associated with a reset token
func (s *InMemoryPasswordResetService) GetUserIDByToken(token string) (uuid.UUID, error) {
	data, exists := s.tokens[token]
	if !exists {
		return uuid.Nil, fmt.Errorf("reset token not found")
	}
	
	// Check if token is expired
	if time.Now().After(data.ExpiresAt) {
		// Remove expired token
		delete(s.tokens, token)
		return uuid.Nil, fmt.Errorf("reset token expired")
	}
	
	return data.UserID, nil
}

// ClearResetToken removes a reset token after use
func (s *InMemoryPasswordResetService) ClearResetToken(token string) error {
	delete(s.tokens, token)
	return nil
}

// CleanupExpiredTokens removes all expired tokens
func (s *InMemoryPasswordResetService) CleanupExpiredTokens() {
	now := time.Now()
	for token, data := range s.tokens {
		if now.After(data.ExpiresAt) {
			delete(s.tokens, token)
		}
	}
}

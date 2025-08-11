package services

import (
	"fmt"
	"time"

	"nurseshift/auth-service/internal/domain/entities"
	"nurseshift/auth-service/internal/domain/usecases"

	"github.com/google/uuid"
)

// MockJWTService implements JWTService for testing
type MockJWTService struct {
	accessTokens  map[string]*usecases.JWTClaims
	refreshTokens map[string]*usecases.JWTClaims
}

// NewMockJWTService creates a new mock JWT service
func NewMockJWTService() *MockJWTService {
	return &MockJWTService{
		accessTokens:  make(map[string]*usecases.JWTClaims),
		refreshTokens: make(map[string]*usecases.JWTClaims),
	}
}

// GenerateAccessToken generates a mock access token
func (s *MockJWTService) GenerateAccessToken(userID uuid.UUID, role entities.UserRole) (string, error) {
	token := fmt.Sprintf("mock-access-token-%s", userID.String())
	claims := &usecases.JWTClaims{
		UserID:   userID,
		Role:     role,
		Type:     "access",
		IssuedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	s.accessTokens[token] = claims
	return token, nil
}

// GenerateRefreshToken generates a mock refresh token
func (s *MockJWTService) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	token := fmt.Sprintf("mock-refresh-token-%s", userID.String())
	claims := &usecases.JWTClaims{
		UserID:   userID,
		Role:     entities.UserRole("user"), // Default role
		Type:     "refresh",
		IssuedAt: time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	s.refreshTokens[token] = claims
	return token, nil
}

// ValidateToken validates a mock token
func (s *MockJWTService) ValidateToken(token string) (*usecases.JWTClaims, error) {
	// Check access tokens first
	if claims, exists := s.accessTokens[token]; exists {
		if time.Now().After(claims.ExpiresAt) {
			return nil, fmt.Errorf("token expired")
		}
		return claims, nil
	}

	// Check refresh tokens
	if claims, exists := s.refreshTokens[token]; exists {
		if time.Now().After(claims.ExpiresAt) {
			return nil, fmt.Errorf("token expired")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// HashToken creates a hash of the token
func (s *MockJWTService) HashToken(token string) string {
	return fmt.Sprintf("hash-%s", token)
}

// Clear clears all mock data
func (s *MockJWTService) Clear() {
	s.accessTokens = make(map[string]*usecases.JWTClaims)
	s.refreshTokens = make(map[string]*usecases.JWTClaims)
}

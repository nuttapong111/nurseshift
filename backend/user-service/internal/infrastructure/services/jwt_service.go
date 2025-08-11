package services

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// UserRole represents user roles
type UserRole string

const (
	RoleAdmin        UserRole = "admin"
	RoleManager      UserRole = "manager"
	RoleNurse        UserRole = "nurse"
	RoleCoordinator  UserRole = "coordinator"
	RoleSupervisor   UserRole = "supervisor"
	RoleAdministrator UserRole = "administrator"
)

// JWTClaims represents JWT claims
type JWTClaims struct {
	UserID         uuid.UUID `json:"userId"`
	Role           UserRole  `json:"role"`
	OrganizationID uuid.UUID `json:"organizationId"`
	Type           string    `json:"type"`
	IssuedAt       time.Time `json:"issuedAt"`
	ExpiresAt      time.Time `json:"expiresAt"`
}

// JWTService interface for JWT operations
type JWTService interface {
	GenerateAccessToken(userID uuid.UUID, role UserRole, organizationID uuid.UUID) (string, error)
	GenerateRefreshToken(userID uuid.UUID) (string, error)
	ValidateAccessToken(tokenString string) (*JWTClaims, error)
	ValidateRefreshToken(tokenString string) (*JWTClaims, error)
	HashToken(token string) string
}

// JWTServiceImpl implements the JWTService interface
type JWTServiceImpl struct {
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewJWTService creates a new JWT service
func NewJWTService(secretKey string, accessTokenTTL, refreshTokenTTL time.Duration) JWTService {
	return &JWTServiceImpl{
		secretKey:       secretKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

// Claims represents JWT claims with custom fields
type Claims struct {
	UserID         uuid.UUID `json:"userId"`
	Role           UserRole  `json:"role"`
	OrganizationID uuid.UUID `json:"organizationId"`
	Type           string    `json:"type"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates an access token for a user
func (j *JWTServiceImpl) GenerateAccessToken(userID uuid.UUID, role UserRole, organizationID uuid.UUID) (string, error) {
	now := time.Now()
	expiresAt := now.Add(j.accessTokenTTL)

	claims := &Claims{
		UserID:         userID,
		Role:           role,
		OrganizationID: organizationID,
		Type:           "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "nurseshift-auth",
			Audience:  []string{"nurseshift-api"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a refresh token for a user
func (j *JWTServiceImpl) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	now := time.Now()
	expiresAt := now.Add(j.refreshTokenTTL)

	claims := &Claims{
		UserID: userID,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "nurseshift-auth",
			Audience:  []string{"nurseshift-refresh"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

// ValidateAccessToken validates an access token and returns claims
func (j *JWTServiceImpl) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.Type != "access" {
		return nil, fmt.Errorf("token is not an access token")
	}

	// Check if token is expired
	if time.Now().After(claims.RegisteredClaims.ExpiresAt.Time) {
		return nil, fmt.Errorf("token is expired")
	}

	return &JWTClaims{
		UserID:         claims.UserID,
		Role:           claims.Role,
		OrganizationID: claims.OrganizationID,
		Type:           claims.Type,
		IssuedAt:       claims.RegisteredClaims.IssuedAt.Time,
		ExpiresAt:      claims.RegisteredClaims.ExpiresAt.Time,
	}, nil
}

// ValidateRefreshToken validates a refresh token and returns claims
func (j *JWTServiceImpl) ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.Type != "refresh" {
		return nil, fmt.Errorf("token is not a refresh token")
	}

	// Check if token is expired
	if time.Now().After(claims.RegisteredClaims.ExpiresAt.Time) {
		return nil, fmt.Errorf("token is expired")
	}

	return &JWTClaims{
		UserID:    claims.UserID,
		Type:      claims.Type,
		IssuedAt:  claims.RegisteredClaims.IssuedAt.Time,
		ExpiresAt: claims.RegisteredClaims.ExpiresAt.Time,
	}, nil
}

// HashToken creates a hash of the token for storage
func (j *JWTServiceImpl) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}

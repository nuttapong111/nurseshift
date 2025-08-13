package usecases

import (
	"context"
	"fmt"
	"time"

	"nurseshift/auth-service/internal/domain/entities"
	"nurseshift/auth-service/internal/domain/repositories"

	"github.com/google/uuid"
)

// JWTService defines the interface for JWT operations
type JWTService interface {
	GenerateAccessToken(userID uuid.UUID, role entities.UserRole) (string, error)
	GenerateRefreshToken(userID uuid.UUID) (string, error)
	ValidateToken(token string) (*JWTClaims, error)
	HashToken(token string) string
}

// PasswordService defines the interface for password operations
type PasswordService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) bool
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID    uuid.UUID
	Role      entities.UserRole
	Type      string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// AuthUseCase defines the interface for authentication use cases
type AuthUseCase interface {
	Login(ctx context.Context, email, password, ipAddress, userAgent string) (*LoginResponse, error)
	Register(ctx context.Context, req *RegisterRequest) (*User, error)
	RefreshToken(ctx context.Context, tokenHash string) (*LoginResponse, error)
	Logout(ctx context.Context, sessionID uuid.UUID) error
	LogoutAllSessions(ctx context.Context, userID uuid.UUID) error
	ValidateSession(ctx context.Context, tokenHash string) (*User, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

// AuthUseCaseImpl implements AuthUseCase
type AuthUseCaseImpl struct {
	userRepo        repositories.UserRepository
	jwtService      JWTService
	passwordService PasswordService
	sessionTimeout  time.Duration
}

// NewAuthUseCase creates a new auth use case
func NewAuthUseCase(
	userRepo repositories.UserRepository,
	jwtService JWTService,
	passwordService PasswordService,
	sessionTimeout time.Duration,
) AuthUseCase {
	return &AuthUseCaseImpl{
		userRepo:        userRepo,
		jwtService:      jwtService,
		passwordService: passwordService,
		sessionTimeout:  sessionTimeout,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	IPAddress string `json:"-"`
	UserAgent string `json:"-"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=6"`
	FirstName string  `json:"firstName" validate:"required"`
	LastName  string  `json:"lastName" validate:"required"`
	Phone     *string `json:"phone,omitempty"`
	Position  *string `json:"position,omitempty"`
}

// LoginResponse represents a successful login response
type LoginResponse struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	User         *User     `json:"user"`
}

// User represents a user response (without sensitive data)
type User struct {
	ID             uuid.UUID           `json:"id"`
	Email          string              `json:"email"`
	FirstName      string              `json:"firstName"`
	LastName       string              `json:"lastName"`
	Phone          *string             `json:"phone,omitempty"`
	Role           entities.UserRole   `json:"role"`
	Status         entities.UserStatus `json:"status"`
	Position       *string             `json:"position,omitempty"`
	DaysRemaining  int                 `json:"remainingDays"`
	PackageType    string              `json:"packageType"`
	MaxDepartments int                 `json:"maxDepartments"`
	AvatarURL      *string             `json:"avatarUrl,omitempty"`
	CreatedAt      time.Time           `json:"createdAt"`
	UpdatedAt      time.Time           `json:"updatedAt"`
}

// Login authenticates a user and creates a session
func (uc *AuthUseCaseImpl) Login(ctx context.Context, email, password, ipAddress, userAgent string) (*LoginResponse, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, fmt.Errorf("account is not active")
	}

	// Verify password
	if !uc.passwordService.VerifyPassword(user.PasswordHash, password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check subscription status
	if user.SubscriptionExpiresAt != nil && user.SubscriptionExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("subscription has expired")
	}

	// Generate JWT tokens
	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token")
	}

	refreshToken, err := uc.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token")
	}

	// Create session
	tokenHash := uc.jwtService.HashToken(refreshToken)
	expiresAt := time.Now().Add(uc.sessionTimeout)

	session := &entities.UserSession{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		IPAddress: &ipAddress,
		UserAgent: &userAgent,
		ExpiresAt: expiresAt,
	}

	err = uc.userRepo.CreateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session")
	}

	// Update last login
	err = uc.userRepo.UpdateLastLogin(ctx, user.ID)
	if err != nil {
		// Log error but don't fail the login
		fmt.Printf("Warning: failed to update last login for user %s: %v\n", user.ID, err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         uc.mapUserToResponse(user),
	}, nil
}

// Register creates a new user account
func (uc *AuthUseCaseImpl) Register(ctx context.Context, req *RegisterRequest) (*User, error) {
	// Check if email already exists
	exists, err := uc.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence")
	}
	if exists {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := uc.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password")
	}

	// Create user with default values
	now := time.Now()
	trialExpiresAt := now.AddDate(0, 0, 30) // 30 days trial period

	user := &entities.User{
		ID:                    uuid.New(),
		Email:                 req.Email,
		PasswordHash:          hashedPassword,
		FirstName:             req.FirstName,
		LastName:              req.LastName,
		Phone:                 req.Phone,
		Role:                  entities.RoleUser, // Default to user role
		Status:                entities.StatusActive,
		Position:              req.Position,
		DaysRemaining:         30,              // Default trial period
		SubscriptionExpiresAt: &trialExpiresAt, // Set trial expiration date
		PackageType:           "trial",
		MaxDepartments:        2,
		Settings:              "{}",
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user")
	}

	return uc.mapUserToResponse(user), nil
}

// RefreshToken generates new tokens using a refresh token
func (uc *AuthUseCaseImpl) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	tokenHash := uc.jwtService.HashToken(refreshToken)

	// Get session
	session, err := uc.userRepo.GetSessionByToken(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, fmt.Errorf("account is not active")
	}

	// Check subscription status
	if user.SubscriptionExpiresAt != nil && user.SubscriptionExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("subscription has expired")
	}

	// Generate new tokens
	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token")
	}

	newRefreshToken, err := uc.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token")
	}

	// Update session with new token
	newTokenHash := uc.jwtService.HashToken(newRefreshToken)
	session.TokenHash = newTokenHash
	session.ExpiresAt = time.Now().Add(uc.sessionTimeout)

	// For simplicity, we'll create a new session instead of updating
	// In a real implementation, you might want to update the existing session
	newSession := &entities.UserSession{
		ID:        uuid.New(),
		UserID:    session.UserID,
		TokenHash: newTokenHash,
		IPAddress: session.IPAddress,
		UserAgent: session.UserAgent,
		ExpiresAt: time.Now().Add(uc.sessionTimeout),
	}

	err = uc.userRepo.CreateSession(ctx, newSession)
	if err != nil {
		return nil, fmt.Errorf("failed to create new session")
	}

	// Revoke old session
	err = uc.userRepo.RevokeSession(ctx, session.ID)
	if err != nil {
		// Log error but don't fail the refresh
		fmt.Printf("Warning: failed to revoke old session %s: %v\n", session.ID, err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    session.ExpiresAt,
		User:         uc.mapUserToResponse(user),
	}, nil
}

// Logout revokes a specific session
func (uc *AuthUseCaseImpl) Logout(ctx context.Context, sessionID uuid.UUID) error {
	return uc.userRepo.RevokeSession(ctx, sessionID)
}

// LogoutAllSessions revokes all sessions for a user
func (uc *AuthUseCaseImpl) LogoutAllSessions(ctx context.Context, userID uuid.UUID) error {
	return uc.userRepo.RevokeAllUserSessions(ctx, userID)
}

// ValidateSession validates a session token and returns the user
func (uc *AuthUseCaseImpl) ValidateSession(ctx context.Context, refreshToken string) (*User, error) {
	tokenHash := uc.jwtService.HashToken(refreshToken)

	session, err := uc.userRepo.GetSessionByToken(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("invalid session")
	}

	user, err := uc.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if !user.IsActive() {
		return nil, fmt.Errorf("account is not active")
	}

	return uc.mapUserToResponse(user), nil
}

// ChangePassword changes a user's password
func (uc *AuthUseCaseImpl) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Verify old password
	if !uc.passwordService.VerifyPassword(user.PasswordHash, oldPassword) {
		return fmt.Errorf("invalid old password")
	}

	// Hash new password
	hashedPassword, err := uc.passwordService.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password")
	}

	// Update password in database
	user.PasswordHash = hashedPassword
	user.UpdatedAt = time.Now()

	err = uc.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update password")
	}

	// Revoke all existing sessions to force re-login
	err = uc.userRepo.RevokeAllUserSessions(ctx, userID)
	if err != nil {
		// Log error but don't fail the password change
		fmt.Printf("Warning: failed to revoke all sessions for user %s: %v\n", userID, err)
	}

	return nil
}

// UpdatePassword changes a user's password
func (uc *AuthUseCaseImpl) UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Hash new password
	hashedPassword, err := uc.passwordService.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password")
	}

	// Update password in database
	user.PasswordHash = hashedPassword
	user.UpdatedAt = time.Now()

	err = uc.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update password")
	}

	// Revoke all existing sessions to force re-login
	err = uc.userRepo.RevokeAllUserSessions(ctx, userID)
	if err != nil {
		// Log error but don't fail the password change
		fmt.Printf("Warning: failed to revoke all sessions for user %s: %v\n", userID, err)
	}

	return nil
}

// GetUserByEmail retrieves a user by email
func (uc *AuthUseCaseImpl) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return uc.mapUserToResponse(user), nil
}

// mapUserToResponse converts a domain user entity to a response user
func (uc *AuthUseCaseImpl) mapUserToResponse(user *entities.User) *User {
	return &User{
		ID:             user.ID,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Phone:          user.Phone,
		Role:           user.Role,
		Status:         user.Status,
		Position:       user.Position,
		DaysRemaining:  user.DaysRemaining,
		PackageType:    user.PackageType,
		MaxDepartments: user.MaxDepartments,
		AvatarURL:      user.AvatarURL,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

package usecases

import (
	"context"
	"fmt"
	"nurseshift/user-service/internal/domain/entities"
	"nurseshift/user-service/internal/domain/repositories"
)

// UserUseCase defines the interface for user business logic operations
type UserUseCase interface {
	// Basic user operations
	GetProfile(ctx context.Context, userID string) (*entities.User, error)
	UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*entities.User, error)
	UploadAvatar(ctx context.Context, userID string, avatarURL string) error

	// User management operations (admin only)
	GetUsers(ctx context.Context, requesterID string, req repositories.GetUsersRequest) (*repositories.GetUsersResponse, error)
	SearchUsers(ctx context.Context, requesterID string, req repositories.SearchUsersRequest) (*repositories.SearchUsersResponse, error)
	GetUserStats(ctx context.Context, requesterID string) (*repositories.UserStatsResponse, error)

	// Email verification methods
	SendVerificationEmail(ctx context.Context, email string) error
	VerifyEmail(ctx context.Context, token string) error
	IsEmailVerified(ctx context.Context, email string) (bool, error)
}

// UserUseCaseImpl implements UserUseCase
type UserUseCaseImpl struct {
	userRepo repositories.UserRepository
}

// NewUserUseCase creates a new instance of UserUseCase
func NewUserUseCase(userRepo repositories.UserRepository) UserUseCase {
	return &UserUseCaseImpl{
		userRepo: userRepo,
	}
}

// GetProfile retrieves the profile of the authenticated user
func (uc *UserUseCaseImpl) GetProfile(ctx context.Context, userID string) (*entities.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Check email verification status
	emailVerified, err := uc.userRepo.IsEmailVerified(ctx, user.Email)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to check email verification status: %v\n", err)
		emailVerified = false
	}

	// Set email verification status
	user.EmailVerified = emailVerified

	return user, nil
}

// UpdateProfile updates the profile of the authenticated user
func (uc *UserUseCaseImpl) UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*entities.User, error) {
	// Get current user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}
	if req.Position != "" {
		user.Position = &req.Position
	}

	// Save to database
	err = uc.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	return user, nil
}

// UploadAvatar updates the user's avatar URL
func (uc *UserUseCaseImpl) UploadAvatar(ctx context.Context, userID string, avatarURL string) error {
	// Get current user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Update avatar URL
	user.AvatarURL = &avatarURL

	// Save to database
	err = uc.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update avatar: %w", err)
	}

	return nil
}

// GetUsers retrieves a list of users (admin only)
func (uc *UserUseCaseImpl) GetUsers(ctx context.Context, requesterID string, req repositories.GetUsersRequest) (*repositories.GetUsersResponse, error) {
	// Check if requester is admin
	requester, err := uc.userRepo.GetByID(ctx, requesterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get requester: %w", err)
	}

	if !requester.IsAdmin() {
		return nil, fmt.Errorf("insufficient permissions")
	}

	// Get users from repository
	response, err := uc.userRepo.GetUsers(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return response, nil
}

// SearchUsers searches for users based on query (admin only)
func (uc *UserUseCaseImpl) SearchUsers(ctx context.Context, requesterID string, req repositories.SearchUsersRequest) (*repositories.SearchUsersResponse, error) {
	// Check if requester is admin
	requester, err := uc.userRepo.GetByID(ctx, requesterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get requester: %w", err)
	}

	if !requester.IsAdmin() {
		return nil, fmt.Errorf("insufficient permissions")
	}

	// Search users from repository
	response, err := uc.userRepo.SearchUsers(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return response, nil
}

// GetUserStats retrieves user statistics (admin only)
func (uc *UserUseCaseImpl) GetUserStats(ctx context.Context, requesterID string) (*repositories.UserStatsResponse, error) {
	// Check if requester is admin
	requester, err := uc.userRepo.GetByID(ctx, requesterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get requester: %w", err)
	}

	if !requester.IsAdmin() {
		return nil, fmt.Errorf("insufficient permissions")
	}

	// Get stats from repository
	stats, err := uc.userRepo.GetUserStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	return stats, nil
}

// Email verification methods
func (uc *UserUseCaseImpl) SendVerificationEmail(ctx context.Context, email string) error {
	// Check if user exists
	_, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Check if email is already verified
	verified, err := uc.userRepo.IsEmailVerified(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to check email verification status: %w", err)
	}

	if verified {
		return fmt.Errorf("email is already verified")
	}

	// Send verification email
	err = uc.userRepo.SendVerificationEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

func (uc *UserUseCaseImpl) VerifyEmail(ctx context.Context, token string) error {
	// Verify email with token
	err := uc.userRepo.VerifyEmail(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	return nil
}

func (uc *UserUseCaseImpl) IsEmailVerified(ctx context.Context, email string) (bool, error) {
	// Check email verification status
	verified, err := uc.userRepo.IsEmailVerified(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to check email verification status: %w", err)
	}

	return verified, nil
}

// UpdateProfileRequest represents the request to update a user profile
type UpdateProfileRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Position  string `json:"position"`
}

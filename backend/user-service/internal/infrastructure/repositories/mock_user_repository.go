package repositories

import (
	"context"
	"time"

	"nurseshift/user-service/internal/domain/entities"
	repo "nurseshift/user-service/internal/domain/repositories"

	"github.com/google/uuid"
)

// MockUserRepository implements UserRepository for testing
type MockUserRepository struct{}

// NewMockUserRepository creates a new mock user repository
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{}
}

// GetByID retrieves a user by ID
func (r *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	// Return mock user data
	user := &entities.User{
		ID:                    id,
		Email:                 "admin@nurseshift.com",
		PasswordHash:          "hashed_password",
		FirstName:            "ผู้ดูแล",
		LastName:             "ระบบ",
		Phone:                stringPtr("0812345678"),
		Role:                 entities.RoleAdmin,
		Status:               entities.StatusActive,
		Position:             stringPtr("System Administrator"),
		DaysRemaining:        365,
		SubscriptionExpiresAt: timePtr(time.Now().AddDate(1, 0, 0)),
		PackageType:          entities.PackageEnterprise,
		MaxDepartments:       10,
		AvatarURL:            stringPtr("https://example.com/avatar.jpg"),
		Settings:             stringPtr("{}"),
		LastLoginAt:          timePtr(time.Now()),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	return user, nil
}

// GetByEmail retrieves a user by email
func (r *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	// Return mock user data
	user := &entities.User{
		ID:                    uuid.New(),
		Email:                 email,
		PasswordHash:          "hashed_password",
		FirstName:            "ผู้ดูแล",
		LastName:             "ระบบ",
		Phone:                stringPtr("0812345678"),
		Role:                 entities.RoleAdmin,
		Status:               entities.StatusActive,
		Position:             stringPtr("System Administrator"),
		DaysRemaining:        365,
		SubscriptionExpiresAt: timePtr(time.Now().AddDate(1, 0, 0)),
		PackageType:          entities.PackageEnterprise,
		MaxDepartments:       10,
		AvatarURL:            stringPtr("https://example.com/avatar.jpg"),
		Settings:             stringPtr("{}"),
		LastLoginAt:          timePtr(time.Now()),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	return user, nil
}

// Create creates a new user
func (r *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	// Mock implementation - just return success
	return nil
}

// Update updates an existing user
func (r *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	// Mock implementation - just return success
	return nil
}

// UpdateLastLogin updates the user's last login timestamp
func (r *MockUserRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	// Mock implementation - just return success
	return nil
}

// GetUsers returns paginated list of users
func (r *MockUserRepository) GetUsers(ctx context.Context, req *repo.GetUsersRequest) (*repo.GetUsersResponse, error) {
	// Return mock users data
	users := []*entities.User{
		{
			ID:                    uuid.New(),
			Email:                 "admin@nurseshift.com",
			FirstName:            "ผู้ดูแล",
			LastName:             "ระบบ",
			Phone:                stringPtr("0812345678"),
			Role:                 entities.RoleAdmin,
			Status:               entities.StatusActive,
			Position:             stringPtr("System Administrator"),
			DaysRemaining:        365,
			SubscriptionExpiresAt: timePtr(time.Now().AddDate(1, 0, 0)),
			PackageType:          entities.PackageEnterprise,
			MaxDepartments:       10,
			AvatarURL:            stringPtr("https://example.com/avatar.jpg"),
			Settings:             stringPtr("{}"),
			LastLoginAt:          timePtr(time.Now()),
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		},
		{
			ID:                    uuid.New(),
			Email:                 "user@nurseshift.com",
			FirstName:            "ผู้ใช้งาน",
			LastName:             "ทั่วไป",
			Phone:                stringPtr("0898765432"),
			Role:                 entities.RoleUser,
			Status:               entities.StatusActive,
			Position:             stringPtr("พยาบาล"),
			DaysRemaining:        30,
			SubscriptionExpiresAt: timePtr(time.Now().AddDate(0, 1, 0)),
			PackageType:          entities.PackageStandard,
			MaxDepartments:       2,
			AvatarURL:            nil,
			Settings:             stringPtr("{}"),
			LastLoginAt:          timePtr(time.Now().Add(-24 * time.Hour)),
			CreatedAt:            time.Now().AddDate(0, 0, -7),
			UpdatedAt:            time.Now().AddDate(0, 0, -1),
		},
	}

	totalPages := (len(users) + req.Limit - 1) / req.Limit

	return &repo.GetUsersResponse{
		Users:      users,
		Total:      len(users),
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// SearchUsers searches users by query
func (r *MockUserRepository) SearchUsers(ctx context.Context, req *repo.SearchUsersRequest) (*repo.SearchUsersResponse, error) {
	// Return mock search results
	users := []*entities.User{
		{
			ID:                    uuid.New(),
			Email:                 "admin@nurseshift.com",
			FirstName:            "ผู้ดูแล",
			LastName:             "ระบบ",
			Phone:                stringPtr("0812345678"),
			Role:                 entities.RoleAdmin,
			Status:               entities.StatusActive,
			Position:             stringPtr("System Administrator"),
			DaysRemaining:        365,
			SubscriptionExpiresAt: timePtr(time.Now().AddDate(1, 0, 0)),
			PackageType:          entities.PackageEnterprise,
			MaxDepartments:       10,
			AvatarURL:            stringPtr("https://example.com/avatar.jpg"),
			Settings:             stringPtr("{}"),
			LastLoginAt:          timePtr(time.Now()),
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		},
	}

	totalPages := (len(users) + req.Limit - 1) / req.Limit

	return &repo.SearchUsersResponse{
		Users:      users,
		Total:      len(users),
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetUserStats returns user statistics
func (r *MockUserRepository) GetUserStats(ctx context.Context) (*repo.UserStatsResponse, error) {
	// Return mock stats
	return &repo.UserStatsResponse{
		TotalUsers:     25,
		ActiveUsers:    22,
		InactiveUsers:  3,
		AdminCount:     2,
		UserCount:      23,
	}, nil
}

// EmailExists checks if an email already exists
func (r *MockUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	// Mock implementation - return false for testing
	return false, nil
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

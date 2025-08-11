package usecases

import (
	"context"
	"time"

	"nurseshift/user-service/internal/domain/entities"
	"nurseshift/user-service/internal/domain/repositories"

	"github.com/google/uuid"
)

// UserUseCase interface for user business logic
type UserUseCase interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*entities.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *UpdateProfileRequest) (*entities.User, error)
	GetUsers(ctx context.Context, organizationID uuid.UUID, req *GetUsersRequest) (*GetUsersResponse, error)
	GetUser(ctx context.Context, userID, requesterID uuid.UUID) (*entities.User, error)
	SearchUsers(ctx context.Context, organizationID uuid.UUID, req *SearchUsersRequest) (*SearchUsersResponse, error)
	GetUserStats(ctx context.Context, organizationID uuid.UUID) (*UserStatsResponse, error)
	UploadAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) error
}

// Request/Response types
type UpdateProfileRequest struct {
	FirstName   *string    `json:"firstName"`
	LastName    *string    `json:"lastName"`
	Phone       *string    `json:"phone"`
	Position    *string    `json:"position"`
	DateOfBirth *time.Time `json:"dateOfBirth"`
}

type GetUsersRequest struct {
	Role         *entities.UserRole   `json:"role"`
	Status       *entities.UserStatus `json:"status"`
	DepartmentID *uuid.UUID           `json:"departmentId"`
	Page         int                  `json:"page"`
	Limit        int                  `json:"limit"`
}

type GetUsersResponse struct {
	Users      []*entities.User `json:"users"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"totalPages"`
}

type SearchUsersRequest struct {
	Query  string               `json:"query"`
	Role   *entities.UserRole   `json:"role"`
	Status *entities.UserStatus `json:"status"`
	Page   int                  `json:"page"`
	Limit  int                  `json:"limit"`
}

type SearchUsersResponse struct {
	Users      []*entities.User `json:"users"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"totalPages"`
}

type UserStatsResponse struct {
	TotalUsers     int `json:"totalUsers"`
	ActiveUsers    int `json:"activeUsers"`
	InactiveUsers  int `json:"inactiveUsers"`
	NurseCount     int `json:"nurseCount"`
	AssistantCount int `json:"assistantCount"`
	ManagerCount   int `json:"managerCount"`
}

// UserUseCaseImpl implements UserUseCase
type UserUseCaseImpl struct {
	userRepo repositories.UserRepository
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(userRepo repositories.UserRepository) UserUseCase {
	return &UserUseCaseImpl{
		userRepo: userRepo,
	}
}

// GetProfile returns user profile
func (uc *UserUseCaseImpl) GetProfile(ctx context.Context, userID uuid.UUID) (*entities.User, error) {
	// Mock implementation
	user := &entities.User{
		ID:        userID,
		FirstName: "สมชาย",
		LastName:  "ใจดี",
		Email:     "test@example.com",
		Role:      entities.RoleNurse,
		Status:    entities.StatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return user, nil
}

// UpdateProfile updates user profile
func (uc *UserUseCaseImpl) UpdateProfile(ctx context.Context, userID uuid.UUID, req *UpdateProfileRequest) (*entities.User, error) {
	// Mock implementation
	user := &entities.User{
		ID:        userID,
		FirstName: getStringValue(req.FirstName, "สมชาย"),
		LastName:  getStringValue(req.LastName, "ใจดี"),
		Phone:     req.Phone,
		Position:  req.Position,
		UpdatedAt: time.Now(),
	}
	return user, nil
}

// GetUsers returns paginated list of users
func (uc *UserUseCaseImpl) GetUsers(ctx context.Context, organizationID uuid.UUID, req *GetUsersRequest) (*GetUsersResponse, error) {
	// Mock implementation
	users := []*entities.User{
		{
			ID:        uuid.New(),
			FirstName: "สมชาย",
			LastName:  "ใจดี",
			Email:     "nurse1@example.com",
			Role:      entities.RoleNurse,
			Status:    entities.StatusActive,
		},
		{
			ID:        uuid.New(),
			FirstName: "สมหญิง",
			LastName:  "รักดี",
			Email:     "nurse2@example.com",
			Role:      entities.RoleNurse,
			Status:    entities.StatusActive,
		},
	}

	totalPages := (len(users) + req.Limit - 1) / req.Limit

	return &GetUsersResponse{
		Users:      users,
		Total:      len(users),
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetUser returns specific user details
func (uc *UserUseCaseImpl) GetUser(ctx context.Context, userID, requesterID uuid.UUID) (*entities.User, error) {
	// Mock implementation
	user := &entities.User{
		ID:        userID,
		FirstName: "สมชาย",
		LastName:  "ใจดี",
		Email:     "test@example.com",
		Role:      entities.RoleNurse,
		Status:    entities.StatusActive,
	}
	return user, nil
}

// SearchUsers searches users by query
func (uc *UserUseCaseImpl) SearchUsers(ctx context.Context, organizationID uuid.UUID, req *SearchUsersRequest) (*SearchUsersResponse, error) {
	// Mock implementation
	users := []*entities.User{}

	totalPages := (len(users) + req.Limit - 1) / req.Limit

	return &SearchUsersResponse{
		Users:      users,
		Total:      len(users),
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetUserStats returns user statistics
func (uc *UserUseCaseImpl) GetUserStats(ctx context.Context, organizationID uuid.UUID) (*UserStatsResponse, error) {
	// Mock implementation
	return &UserStatsResponse{
		TotalUsers:     25,
		ActiveUsers:    22,
		InactiveUsers:  3,
		NurseCount:     15,
		AssistantCount: 8,
		ManagerCount:   2,
	}, nil
}

// UploadAvatar uploads user avatar
func (uc *UserUseCaseImpl) UploadAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) error {
	// Mock implementation - in real scenario, this would update the database
	return nil
}

// Helper functions
func getStringValue(ptr *string, defaultValue string) string {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}

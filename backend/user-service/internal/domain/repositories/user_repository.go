package repositories

import (
	"context"
	"nurseshift/user-service/internal/domain/entities"
)

// UserRepository defines the interface for user data access operations
type UserRepository interface {
	// Basic CRUD operations
	GetByID(ctx context.Context, id string) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
	UpdateLastLogin(ctx context.Context, userID string) error
	EmailExists(ctx context.Context, email string) (bool, error)

	// User management operations (admin only)
	GetUsers(ctx context.Context, req *GetUsersRequest) (*GetUsersResponse, error)
	SearchUsers(ctx context.Context, req *SearchUsersRequest) (*SearchUsersResponse, error)
	GetUserStats(ctx context.Context) (*UserStatsResponse, error)

	// Email verification operations
	SendVerificationEmail(ctx context.Context, email string) error
	VerifyEmail(ctx context.Context, token string) error
	IsEmailVerified(ctx context.Context, email string) (bool, error)
}

// GetUsersRequest represents request for getting users
type GetUsersRequest struct {
	Role   *entities.UserRole   `json:"role"`
	Status *entities.UserStatus `json:"status"`
	Page   int                  `json:"page"`
	Limit  int                  `json:"limit"`
}

// GetUsersResponse represents response for getting users
type GetUsersResponse struct {
	Users      []*entities.User `json:"users"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"totalPages"`
}

// SearchUsersRequest represents request for searching users
type SearchUsersRequest struct {
	Query  string               `json:"query"`
	Role   *entities.UserRole   `json:"role"`
	Status *entities.UserStatus `json:"status"`
	Page   int                  `json:"page"`
	Limit  int                  `json:"limit"`
}

// SearchUsersResponse represents response for searching users
type SearchUsersResponse struct {
	Users      []*entities.User `json:"users"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"totalPages"`
}

// UserStatsResponse represents user statistics
type UserStatsResponse struct {
	TotalUsers    int `json:"totalUsers"`
	ActiveUsers   int `json:"activeUsers"`
	InactiveUsers int `json:"inactiveUsers"`
	AdminCount    int `json:"adminCount"`
	UserCount     int `json:"userCount"`
}

// Email verification DTOs
type SendVerificationEmailRequest struct {
	Email string `json:"email"`
}

type SendVerificationEmailResponse struct {
	Message string `json:"message"`
}

type VerifyEmailRequest struct {
	Token string `json:"token"`
}

type VerifyEmailResponse struct {
	Message  string `json:"message"`
	Verified bool   `json:"verified"`
}

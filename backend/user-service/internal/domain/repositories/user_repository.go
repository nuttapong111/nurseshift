package repositories

import (
	"context"
	"nurseshift/user-service/internal/domain/entities"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// User CRUD operations
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByEmployeeID(ctx context.Context, employeeID string, organizationID uuid.UUID) (*entities.User, error)
	GetByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entities.User, error)
	GetByDepartment(ctx context.Context, departmentID uuid.UUID, limit, offset int) ([]*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
	UpdateProfile(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
	UpdateStatus(ctx context.Context, userID uuid.UUID, status entities.UserStatus) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Search and filtering
	Search(ctx context.Context, organizationID uuid.UUID, query string, role *entities.UserRole, status *entities.UserStatus, limit, offset int) ([]*entities.User, error)
	Count(ctx context.Context, organizationID uuid.UUID, role *entities.UserRole, status *entities.UserStatus) (int, error)

	// Organization operations
	GetOrganizationByID(ctx context.Context, id uuid.UUID) (*entities.Organization, error)
	GetOrganizationByUser(ctx context.Context, userID uuid.UUID) (*entities.Organization, error)

	// Department operations
	GetDepartmentByID(ctx context.Context, id uuid.UUID) (*entities.Department, error)
	GetDepartmentsByOrganization(ctx context.Context, organizationID uuid.UUID) ([]*entities.Department, error)
	GetUserDepartment(ctx context.Context, userID uuid.UUID) (*entities.Department, error)

	// Validation
	EmailExists(ctx context.Context, email string, excludeUserID *uuid.UUID) (bool, error)
	EmployeeIDExists(ctx context.Context, employeeID string, organizationID uuid.UUID, excludeUserID *uuid.UUID) (bool, error)
}



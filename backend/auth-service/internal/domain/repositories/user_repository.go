package repositories

import (
	"context"
	"nurseshift/auth-service/internal/domain/entities"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// User operations
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByEmployeeID(ctx context.Context, employeeID string) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error

	// Organization operations
	GetOrganizationByID(ctx context.Context, id uuid.UUID) (*entities.Organization, error)
	GetOrganizationByUser(ctx context.Context, userID uuid.UUID) (*entities.Organization, error)

	// Session operations
	CreateSession(ctx context.Context, session *entities.UserSession) error
	GetSessionByToken(ctx context.Context, tokenHash string) (*entities.UserSession, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error
	CleanExpiredSessions(ctx context.Context) error

	// User validation
	EmailExists(ctx context.Context, email string) (bool, error)
	EmployeeIDExists(ctx context.Context, employeeID string, organizationID uuid.UUID) (bool, error)
}


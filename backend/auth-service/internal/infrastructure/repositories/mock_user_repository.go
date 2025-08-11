package repositories

import (
	"context"
	"fmt"
	"time"

	"nurseshift/auth-service/internal/domain/entities"

	"github.com/google/uuid"
)

// MockUserRepository implements UserRepository for testing
type MockUserRepository struct {
	users    map[uuid.UUID]*entities.User
	sessions map[uuid.UUID]*entities.UserSession
	emails   map[string]bool
}

// NewMockUserRepository creates a new mock user repository
func NewMockUserRepository() *MockUserRepository {
	repo := &MockUserRepository{
		users:    make(map[uuid.UUID]*entities.User),
		sessions: make(map[uuid.UUID]*entities.UserSession),
		emails:   make(map[string]bool),
	}

	// Add mock users for testing
	adminID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")

	// Mock admin user
	adminUser := &entities.User{
		ID:             adminID,
		Email:          "admin@nurseshift.com",
		PasswordHash:   "$2a$10$OQJnRxKT4dwQD1blpI2lze9/1NPb4XVl.V8Hle6a3p1p7CsIC7I4m", // password: admin123
		FirstName:      "ผู้ดูแล",
		LastName:       "ระบบ",
		Role:           entities.RoleAdmin,
		Status:         entities.StatusActive,
		Position:       nil,
		DaysRemaining:  90,
		PackageType:    "enterprise",
		MaxDepartments: 20,
		AvatarURL:      nil,
		Settings:       "{}",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Mock regular user
	regularUser := &entities.User{
		ID:             userID,
		Email:          "user@nurseshift.com",
		PasswordHash:   "$2a$10$5G5QXW39SmBuPa9UF0Mft.rAXKMaLi0VeBuvesgjwFslILoE7ej.C", // password: user123
		FirstName:      "พยาบาล",
		LastName:       "ทดสอบ",
		Role:           entities.RoleUser,
		Status:         entities.StatusActive,
		Position:       stringPtr("หัวหน้าพยาบาล"),
		DaysRemaining:  85,
		PackageType:    "trial",
		MaxDepartments: 2,
		AvatarURL:      nil,
		Settings:       "{}",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	repo.users[adminID] = adminUser
	repo.users[userID] = regularUser
	repo.emails["admin@nurseshift.com"] = true
	repo.emails["user@nurseshift.com"] = true

	return repo
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

// GetByID retrieves a user by ID
func (r *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// GetByEmail retrieves a user by email
func (r *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// Create creates a new user
func (r *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	r.users[user.ID] = user
	r.emails[user.Email] = true
	return nil
}

// Update updates an existing user
func (r *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	if _, exists := r.users[user.ID]; !exists {
		return fmt.Errorf("user not found")
	}
	r.users[user.ID] = user
	return nil
}

// Delete deletes a user
func (r *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, exists := r.users[id]; !exists {
		return fmt.Errorf("user not found")
	}
	delete(r.users, id)
	return nil
}

// EmailExists checks if an email already exists
func (r *MockUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	return r.emails[email], nil
}

// UpdateLastLogin updates the last login timestamp
func (r *MockUserRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	user, exists := r.users[userID]
	if !exists {
		return fmt.Errorf("user not found")
	}
	now := time.Now()
	user.LastLoginAt = &now
	return nil
}

// CreateSession creates a new user session
func (r *MockUserRepository) CreateSession(ctx context.Context, session *entities.UserSession) error {
	r.sessions[session.ID] = session
	return nil
}

// GetSessionByToken retrieves a session by token hash
func (r *MockUserRepository) GetSessionByToken(ctx context.Context, tokenHash string) (*entities.UserSession, error) {
	for _, session := range r.sessions {
		if session.TokenHash == tokenHash && session.IsValid() {
			return session, nil
		}
	}
	return nil, fmt.Errorf("session not found")
}

// RevokeSession revokes a specific session
func (r *MockUserRepository) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	session, exists := r.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found")
	}
	now := time.Now()
	session.RevokedAt = &now
	return nil
}

// RevokeAllUserSessions revokes all sessions for a user
func (r *MockUserRepository) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	for _, session := range r.sessions {
		if session.UserID == userID {
			now := time.Now()
			session.RevokedAt = &now
		}
	}
	return nil
}

// CleanExpiredSessions removes expired sessions
func (r *MockUserRepository) CleanExpiredSessions(ctx context.Context) error {
	now := time.Now()
	for sessionID, session := range r.sessions {
		if session.ExpiresAt.Before(now) {
			delete(r.sessions, sessionID)
		}
	}
	return nil
}

// AddMockUser adds a mock user for testing
func (r *MockUserRepository) AddMockUser(user *entities.User) {
	r.users[user.ID] = user
	r.emails[user.Email] = true
}

// Clear clears all mock data
func (r *MockUserRepository) Clear() {
	r.users = make(map[uuid.UUID]*entities.User)
	r.sessions = make(map[uuid.UUID]*entities.UserSession)
	r.emails = make(map[string]bool)
}

// EmployeeIDExists checks if an employee ID already exists
func (r *MockUserRepository) EmployeeIDExists(ctx context.Context, employeeID string, organizationID uuid.UUID) (bool, error) {
	// For mock implementation, always return false
	return false, nil
}

// GetByEmployeeID retrieves a user by employee ID
func (r *MockUserRepository) GetByEmployeeID(ctx context.Context, employeeID string) (*entities.User, error) {
	// For mock implementation, return nil
	return nil, fmt.Errorf("user not found")
}

// GetOrganizationByID retrieves an organization by ID
func (r *MockUserRepository) GetOrganizationByID(ctx context.Context, id uuid.UUID) (*entities.Organization, error) {
	// For mock implementation, return nil
	return nil, fmt.Errorf("organization not found")
}

// GetOrganizationByUser retrieves an organization by user ID
func (r *MockUserRepository) GetOrganizationByUser(ctx context.Context, userID uuid.UUID) (*entities.Organization, error) {
	// For mock implementation, return nil
	return nil, fmt.Errorf("organization not found")
}

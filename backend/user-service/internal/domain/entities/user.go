package entities

import (
	"time"

	"github.com/google/uuid"
)

// UserRole represents user roles in the system
type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

// UserStatus represents user status
type UserStatus string

const (
	StatusActive    UserStatus = "active"
	StatusInactive  UserStatus = "inactive"
	StatusPending   UserStatus = "pending"
	StatusSuspended UserStatus = "suspended"
)

// PackageType represents package types
type PackageType string

const (
	PackageStandard   PackageType = "standard"
	PackageEnterprise PackageType = "enterprise"
	PackageTrial      PackageType = "trial"
)

// User represents the user entity matching the database schema
type User struct {
	ID                         uuid.UUID   `json:"id" db:"id"`
	Email                      string      `json:"email" db:"email"`
	PasswordHash               string      `json:"-" db:"password_hash"`
	FirstName                  string      `json:"firstName" db:"first_name"`
	LastName                   string      `json:"lastName" db:"last_name"`
	Phone                      *string     `json:"phone" db:"phone"`
	Role                       UserRole    `json:"role" db:"role"`
	Status                     UserStatus  `json:"status" db:"status"`
	Position                   *string     `json:"position" db:"position"`
	DaysRemaining              int         `json:"remainingDays" db:"days_remaining"`
	SubscriptionExpiresAt      *time.Time  `json:"subscriptionExpiresAt" db:"subscription_expires_at"`
	PackageType                PackageType `json:"packageType" db:"package_type"`
	MaxDepartments             int         `json:"maxDepartments" db:"max_departments"`
	AvatarURL                  *string     `json:"avatarUrl" db:"avatar_url"`
	Settings                   *string     `json:"settings" db:"settings"`
	LastLoginAt                *time.Time  `json:"lastLoginAt" db:"last_login_at"`
	EmailVerified              bool        `json:"emailVerified" db:"email_verified"`
	EmailVerificationToken     *string     `json:"emailVerificationToken,omitempty" db:"email_verification_token"`
	EmailVerificationExpiresAt *time.Time  `json:"emailVerificationExpiresAt,omitempty" db:"email_verification_expires_at"`
	CreatedAt                  time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt                  time.Time   `json:"updatedAt" db:"updated_at"`
}

// Department represents a department entity
type Department struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"userId" db:"user_id"`
	Name          string    `json:"name" db:"name"`
	Description   *string   `json:"description" db:"description"`
	MaxNurses     int       `json:"maxNurses" db:"max_nurses"`
	MaxAssistants int       `json:"maxAssistants" db:"max_assistants"`
	Settings      *string   `json:"settings" db:"settings"`
	IsActive      bool      `json:"isActive" db:"is_active"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

// UserProfile represents extended user profile information
type UserProfile struct {
	User        *User       `json:"user"`
	Department  *Department `json:"department,omitempty"`
	Permissions []string    `json:"permissions"`
}

// Methods for User entity

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// IsActive checks if the user is active
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// HasRole checks if the user has a specific role
func (u *User) HasRole(role UserRole) bool {
	return u.Role == role
}

// IsAdmin checks if the user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// CanManage checks if the user can manage other users
func (u *User) CanManage() bool {
	return u.Role == RoleAdmin
}

// GetRoleDisplayName returns the display name for the user's role
func (u *User) GetRoleDisplayName() string {
	switch u.Role {
	case RoleAdmin:
		return "ผู้ดูแลระบบ"
	case RoleUser:
		return "ผู้ใช้งาน"
	default:
		return "ไม่ระบุ"
	}
}

// GetPermissions returns user permissions based on role
func (u *User) GetPermissions() []string {
	permissions := []string{"profile:read", "profile:update"}

	switch u.Role {
	case RoleAdmin:
		permissions = append(permissions, []string{
			"users:create", "users:read", "users:update", "users:delete",
			"departments:create", "departments:read", "departments:update", "departments:delete",
			"schedules:create", "schedules:read", "schedules:update", "schedules:delete",
			"notifications:create", "notifications:read", "notifications:update", "notifications:delete",
			"organization:read", "organization:update",
		}...)
	case RoleUser:
		permissions = append(permissions, []string{
			"departments:read", "departments:update",
			"schedules:create", "schedules:read", "schedules:update",
			"notifications:read",
		}...)
	}

	return permissions
}

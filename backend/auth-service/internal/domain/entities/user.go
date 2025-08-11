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

// User represents the user entity
type User struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	Email                 string     `json:"email" db:"email"`
	PasswordHash          string     `json:"-" db:"password_hash"`
	FirstName             string     `json:"firstName" db:"first_name"`
	LastName              string     `json:"lastName" db:"last_name"`
	Phone                 *string    `json:"phone" db:"phone"`
	Role                  UserRole   `json:"role" db:"role"`
	Status                UserStatus `json:"status" db:"status"`
	Position              *string    `json:"position" db:"position"`
	DaysRemaining         int        `json:"remainingDays" db:"days_remaining"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at" db:"subscription_expires_at"`
	PackageType           string     `json:"packageType" db:"package_type"`
	MaxDepartments        int        `json:"maxDepartments" db:"max_departments"`
	AvatarURL             *string    `json:"avatarUrl" db:"avatar_url"`
	Settings              string     `json:"settings" db:"settings"`
	LastLoginAt           *time.Time `json:"lastLoginAt" db:"last_login_at"`
	CreatedAt             time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt             time.Time  `json:"updatedAt" db:"updated_at"`
}

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

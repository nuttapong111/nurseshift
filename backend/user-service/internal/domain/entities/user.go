package entities

import (
	"time"

	"github.com/google/uuid"
)

// UserRole represents user roles in the system
type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleManager   UserRole = "manager"
	RoleNurse     UserRole = "nurse"
	RoleAssistant UserRole = "assistant"
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
	ID             uuid.UUID  `json:"id" db:"id"`
	OrganizationID uuid.UUID  `json:"organization_id" db:"organization_id"`
	EmployeeID     string     `json:"employee_id" db:"employee_id"`
	Email          string     `json:"email" db:"email"`
	PasswordHash   string     `json:"-" db:"password_hash"`
	FirstName      string     `json:"first_name" db:"first_name"`
	LastName       string     `json:"last_name" db:"last_name"`
	Phone          *string    `json:"phone" db:"phone"`
	Role           UserRole   `json:"role" db:"role"`
	Status         UserStatus `json:"status" db:"status"`
	Position       *string    `json:"position" db:"position"`
	DateJoined     time.Time  `json:"date_joined" db:"date_joined"`
	DateOfBirth    *time.Time `json:"date_of_birth" db:"date_of_birth"`
	AvatarURL      *string    `json:"avatar_url" db:"avatar_url"`
	LastLoginAt    *time.Time `json:"last_login_at" db:"last_login_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// Organization represents the organization entity
type Organization struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	Name                  string     `json:"name" db:"name"`
	Description           *string    `json:"description" db:"description"`
	Email                 *string    `json:"email" db:"email"`
	Phone                 *string    `json:"phone" db:"phone"`
	Address               *string    `json:"address" db:"address"`
	Website               *string    `json:"website" db:"website"`
	LicenseNumber         *string    `json:"license_number" db:"license_number"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at" db:"subscription_expires_at"`
	PackageType           string     `json:"package_type" db:"package_type"`
	MaxUsers              int        `json:"max_users" db:"max_users"`
	MaxDepartments        int        `json:"max_departments" db:"max_departments"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`
}

// Department represents a department entity
type Department struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	OrganizationID uuid.UUID  `json:"organization_id" db:"organization_id"`
	Name           string     `json:"name" db:"name"`
	Description    *string    `json:"description" db:"description"`
	HeadUserID     *uuid.UUID `json:"head_user_id" db:"head_user_id"`
	MaxNurses      int        `json:"max_nurses" db:"max_nurses"`
	MaxAssistants  int        `json:"max_assistants" db:"max_assistants"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// UserProfile represents extended user profile information
type UserProfile struct {
	User         *User         `json:"user"`
	Organization *Organization `json:"organization"`
	Department   *Department   `json:"department,omitempty"`
	Permissions  []string      `json:"permissions"`
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

// IsManager checks if the user is a manager
func (u *User) IsManager() bool {
	return u.Role == RoleManager
}

// CanManage checks if the user can manage other users
func (u *User) CanManage() bool {
	return u.Role == RoleAdmin || u.Role == RoleManager
}

// GetRoleDisplayName returns the display name for the user's role
func (u *User) GetRoleDisplayName() string {
	switch u.Role {
	case RoleAdmin:
		return "ผู้ดูแลระบบ"
	case RoleManager:
		return "หัวหน้าแผนก"
	case RoleNurse:
		return "พยาบาลวิชาชีพ"
	case RoleAssistant:
		return "ผู้ช่วยพยาบาล"
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
	case RoleManager:
		permissions = append(permissions, []string{
			"users:read", "users:update",
			"departments:read", "departments:update",
			"schedules:create", "schedules:read", "schedules:update",
			"notifications:read",
		}...)
	case RoleNurse, RoleAssistant:
		permissions = append(permissions, []string{
			"schedules:read",
			"notifications:read",
		}...)
	}

	return permissions
}

// Methods for Organization entity

// IsSubscriptionActive checks if the organization's subscription is active
func (o *Organization) IsSubscriptionActive() bool {
	if o.SubscriptionExpiresAt == nil {
		return false
	}
	return time.Now().Before(*o.SubscriptionExpiresAt)
}

// GetPackageDisplayName returns the display name for the package type
func (o *Organization) GetPackageDisplayName() string {
	switch o.PackageType {
	case "trial":
		return "แพ็คเกจทดลองใช้"
	case "standard":
		return "แพ็คเกจมาตรฐาน"
	case "enterprise":
		return "แพ็คเกจระดับองค์กร"
	default:
		return "ไม่ระบุ"
	}
}



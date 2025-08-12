package entities

import (
	"time"

	"github.com/google/uuid"
)

// Department represents a department entity
type Department struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	Name          string     `json:"name" db:"name"`
	Description   *string    `json:"description" db:"description"`
	HeadUserID    *uuid.UUID `json:"head_user_id" db:"head_user_id"`
	MaxNurses     int        `json:"max_nurses" db:"max_nurses"`
	MaxAssistants int        `json:"max_assistants" db:"max_assistants"`
	Settings      *string    `json:"settings" db:"settings"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	CreatedBy     *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// DepartmentStaff represents staff in a department (from department_staff table)
type DepartmentStaff struct {
	ID           uuid.UUID `json:"id" db:"id"`
	DepartmentID uuid.UUID `json:"department_id" db:"department_id"`
	Name         string    `json:"name" db:"name"`
	Position     string    `json:"position" db:"position"`
	Phone        *string   `json:"phone" db:"phone"`
	Email        *string   `json:"email" db:"email"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// DepartmentUser represents the relationship between users and departments
type DepartmentUser struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	DepartmentID   uuid.UUID  `json:"department_id" db:"department_id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	DepartmentRole string     `json:"department_role" db:"department_role"`
	AssignedAt     time.Time  `json:"assigned_at" db:"assigned_at"`
	AssignedBy     *uuid.UUID `json:"assigned_by" db:"assigned_by"`
}

// User basic info for department context
type User struct {
	ID         uuid.UUID `json:"id"`
	EmployeeID string    `json:"employee_id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Phone      *string   `json:"phone"`
	Role       string    `json:"role"`
	Status     string    `json:"status"`
	Position   *string   `json:"position"`
}

// DepartmentWithStats represents department with statistics
type DepartmentWithStats struct {
	*Department
	TotalEmployees  int `json:"total_employees"`
	ActiveEmployees int `json:"active_employees"`
	NurseCount      int `json:"nurse_count"`
	AssistantCount  int `json:"assistant_count"`
}

// DepartmentDetail represents detailed department information
type DepartmentDetail struct {
	*Department
	Head  *User              `json:"head,omitempty"`
	Staff []*DepartmentStaff `json:"staff"`
	Stats *DepartmentStats   `json:"stats"`
}

// DepartmentStats represents department statistics
type DepartmentStats struct {
	TotalEmployees  int     `json:"total_employees"`
	ActiveEmployees int     `json:"active_employees"`
	NurseCount      int     `json:"nurse_count"`
	AssistantCount  int     `json:"assistant_count"`
	UtilizationRate float64 `json:"utilization_rate"`
}

// CreateDepartmentRequest represents the request structure for creating a department
type CreateDepartmentRequest struct {
	Name          string     `json:"name" validate:"required"`
	Description   *string    `json:"description"`
	HeadUserID    *uuid.UUID `json:"head_user_id"`
	MaxNurses     int        `json:"max_nurses" validate:"required,min=1"`
	MaxAssistants int        `json:"max_assistants" validate:"required,min=1"`
	Settings      *string    `json:"settings"`
}

// UpdateDepartmentRequest represents the request structure for updating a department
type UpdateDepartmentRequest struct {
	Name          *string    `json:"name"`
	Description   *string    `json:"description"`
	HeadUserID    *uuid.UUID `json:"head_user_id"`
	MaxNurses     *int       `json:"max_nurses"`
	MaxAssistants *int       `json:"max_assistants"`
	Settings      *string    `json:"settings"`
	IsActive      *bool      `json:"is_active"`
}

// Methods

// GetDisplayName returns formatted department name
func (d *Department) GetDisplayName() string {
	if d.Description != nil && *d.Description != "" {
		return d.Name + " - " + *d.Description
	}
	return d.Name
}

// HasHead checks if department has a head assigned
func (d *Department) HasHead() bool {
	return d.HeadUserID != nil
}

// GetMaxStaff returns total maximum staff
func (d *Department) GetMaxStaff() int {
	return d.MaxNurses + d.MaxAssistants
}

// GetUtilizationRate calculates staff utilization rate
func (stats *DepartmentStats) GetUtilizationRate() float64 {
	return stats.UtilizationRate
}

// IsCurrentlyActive checks if staff is currently active
func (s *DepartmentStaff) IsCurrentlyActive() bool {
	return s.IsActive
}

// GetFullName returns staff's full name
func (s *DepartmentStaff) GetFullName() string {
	return s.Name
}

// IsNurse checks if staff is a nurse
func (s *DepartmentStaff) IsNurse() bool {
	return s.Position == "nurse" || s.Position == "พยาบาล"
}

// IsAssistant checks if staff is an assistant
func (s *DepartmentStaff) IsAssistant() bool {
	return s.Position == "assistant" || s.Position == "ผู้ช่วยพยาบาล"
}

// GetFullName returns user's full name
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// IsNurse checks if user is a nurse
func (u *User) IsNurse() bool {
	return u.Role == "nurse"
}

// IsAssistant checks if user is an assistant
func (u *User) IsAssistant() bool {
	return u.Role == "assistant"
}

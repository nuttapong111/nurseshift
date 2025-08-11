package entities

import (
	"time"

	"github.com/google/uuid"
)

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

// Employee represents an employee in a department
type Employee struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	DepartmentID uuid.UUID  `json:"department_id" db:"department_id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	Position     string     `json:"position" db:"position"`
	StartDate    time.Time  `json:"start_date" db:"start_date"`
	EndDate      *time.Time `json:"end_date" db:"end_date"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
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

// EmployeeWithUser represents employee with user information
type EmployeeWithUser struct {
	*Employee
	User *User `json:"user"`
}

// DepartmentDetail represents detailed department information
type DepartmentDetail struct {
	*Department
	Head      *User               `json:"head,omitempty"`
	Employees []*EmployeeWithUser `json:"employees"`
	Stats     *DepartmentStats    `json:"stats"`
}

// DepartmentStats represents department statistics
type DepartmentStats struct {
	TotalEmployees  int     `json:"total_employees"`
	ActiveEmployees int     `json:"active_employees"`
	NurseCount      int     `json:"nurse_count"`
	AssistantCount  int     `json:"assistant_count"`
	UtilizationRate float64 `json:"utilization_rate"`
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

// IsCurrentlyActive checks if employee is currently active
func (e *Employee) IsCurrentlyActive() bool {
	return e.IsActive && (e.EndDate == nil || e.EndDate.After(time.Now()))
}

// GetDuration returns employment duration
func (e *Employee) GetDuration() time.Duration {
	endDate := time.Now()
	if e.EndDate != nil {
		endDate = *e.EndDate
	}
	return endDate.Sub(e.StartDate)
}

// GetFullName returns employee's full name
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

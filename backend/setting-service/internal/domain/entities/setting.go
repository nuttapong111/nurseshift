package entities

import (
	"time"

	"github.com/google/uuid"
)

// Shift represents a shift definition in a department
type Shift struct {
	ID                 uuid.UUID
	DepartmentID       uuid.UUID
	Name               string
	Type               string
	StartTime          string // HH:MM
	EndTime            string // HH:MM
	DurationHours      float64
	RequiredNurses     int
	RequiredAssistants int
	Color              string
	IsActive           bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// WorkingDay represents enabled days of week for a department
type WorkingDay struct {
	ID           uuid.UUID
	DepartmentID uuid.UUID
	DayOfWeek    int // 0=Sunday .. 6=Saturday
	IsWorkingDay bool
	CreatedAt    time.Time
}

// Holiday represents department holidays
type Holiday struct {
	ID           uuid.UUID
	DepartmentID uuid.UUID
	Name         string
	StartDate    time.Time
	EndDate      time.Time
	IsRecurring  bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// SettingsAggregate groups department settings for transport
type SettingsAggregate struct {
	WorkingDays []WorkingDay `json:"workingDays"`
	Shifts      []Shift      `json:"shifts"`
	Holidays    []Holiday    `json:"holidays"`
}

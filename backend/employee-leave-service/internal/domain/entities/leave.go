package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Error constants
var (
	ErrInvalidDateRange = errors.New("end date must be after start date")
	ErrLeaveNotFound    = errors.New("leave request not found")
	ErrInvalidStatus    = errors.New("invalid leave status")
)

// LeaveType represents the type of leave
type LeaveType string

const (
	LeaveTypeSick      LeaveType = "sick"
	LeaveTypePersonal  LeaveType = "personal"
	LeaveTypeVacation  LeaveType = "vacation"
	LeaveTypeEmergency LeaveType = "emergency"
	LeaveTypeMaternity LeaveType = "maternity"
)

// LeaveStatus represents the status of leave request
type LeaveStatus string

const (
	LeaveStatusPending   LeaveStatus = "pending"
	LeaveStatusApproved  LeaveStatus = "approved"
	LeaveStatusRejected  LeaveStatus = "rejected"
	LeaveStatusCancelled LeaveStatus = "cancelled"
)

// LeaveRequest represents a leave request entity
type LeaveRequest struct {
	ID              uuid.UUID   `json:"id"`
	StaffID         uuid.UUID   `json:"userId"`
	DepartmentID    uuid.UUID   `json:"departmentId"`
	LeaveType       LeaveType   `json:"leaveType"`
	StartDate       time.Time   `json:"startDate"`
	EndDate         time.Time   `json:"endDate"`
	Reason          *string     `json:"reason,omitempty"`
	Status          LeaveStatus `json:"status"`
	ApprovedBy      *uuid.UUID  `json:"approvedBy,omitempty"`
	ApprovedAt      *time.Time  `json:"approvedAt,omitempty"`
	RejectionReason *string     `json:"rejectionReason,omitempty"`
	Attachments     *string     `json:"attachments,omitempty"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
}

// LeaveRequestCreate represents the data needed to create a leave request
type LeaveRequestCreate struct {
	StaffID      uuid.UUID `json:"userId" validate:"required"`
	DepartmentID uuid.UUID `json:"departmentId" validate:"required"`
	LeaveType    LeaveType `json:"leaveType" validate:"required"`
	StartDate    time.Time `json:"startDate" validate:"required"`
	EndDate      time.Time `json:"endDate" validate:"required"`
	Reason       *string   `json:"reason,omitempty"`
}

// LeaveRequestUpdate represents the data needed to update a leave request
type LeaveRequestUpdate struct {
	LeaveType *LeaveType   `json:"leaveType,omitempty"`
	StartDate *time.Time   `json:"startDate,omitempty"`
	EndDate   *time.Time   `json:"endDate,omitempty"`
	Reason    *string      `json:"reason,omitempty"`
	Status    *LeaveStatus `json:"status,omitempty"`
}

// LeaveRequestFilter represents filters for querying leave requests
type LeaveRequestFilter struct {
	Month        *string      `json:"month,omitempty"` // YYYY-MM format
	StaffID      *uuid.UUID   `json:"userId,omitempty"`
	DepartmentID *uuid.UUID   `json:"departmentId,omitempty"`
	Status       *LeaveStatus `json:"status,omitempty"`
	StartDate    *time.Time   `json:"startDate,omitempty"`
	EndDate      *time.Time   `json:"endDate,omitempty"`
}

// LeaveRequestWithDetails represents a leave request with additional details
type LeaveRequestWithDetails struct {
	LeaveRequest
	UserName       string  `json:"userName"`
	DepartmentName string  `json:"departmentName"`
	ApproverName   *string `json:"approverName,omitempty"`
}

package repositories

import (
	"context"
	"nurseshift/employee-leave-service/internal/domain/entities"

	"github.com/google/uuid"
)

// LeaveRepository defines the interface for leave request data access
type LeaveRepository interface {
	// Create creates a new leave request
	Create(ctx context.Context, leave entities.LeaveRequest) (uuid.UUID, error)

	// GetByID retrieves a leave request by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.LeaveRequest, error)

	// GetByIDWithDetails retrieves a leave request by ID with additional details
	GetByIDWithDetails(ctx context.Context, id uuid.UUID) (*entities.LeaveRequestWithDetails, error)

	// GetByFilter retrieves leave requests based on filters
	GetByFilter(ctx context.Context, filter entities.LeaveRequestFilter) ([]entities.LeaveRequestWithDetails, error)

	// GetByUser retrieves all leave requests for a specific user
	GetByUser(ctx context.Context, userID uuid.UUID) ([]entities.LeaveRequestWithDetails, error)

	// GetByDepartment retrieves all leave requests for a specific department
	GetByDepartment(ctx context.Context, departmentID uuid.UUID) ([]entities.LeaveRequestWithDetails, error)

	// Update updates an existing leave request
	Update(ctx context.Context, id uuid.UUID, update entities.LeaveRequestUpdate) error

	// UpdateStatus updates the status of a leave request
	UpdateStatus(ctx context.Context, id uuid.UUID, status entities.LeaveStatus, approverID *uuid.UUID, rejectionReason *string) error

	// Delete deletes a leave request
	Delete(ctx context.Context, id uuid.UUID) error

	// ToggleActive toggles the active status (for soft delete/restore)
	ToggleActive(ctx context.Context, id uuid.UUID) error
}

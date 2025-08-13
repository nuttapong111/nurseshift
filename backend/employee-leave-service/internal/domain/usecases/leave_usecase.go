package usecases

import (
	"context"
	"time"

	"nurseshift/employee-leave-service/internal/domain/entities"
	"nurseshift/employee-leave-service/internal/domain/repositories"

	"github.com/google/uuid"
)

// LeaveUseCase defines the interface for leave request business logic
type LeaveUseCase interface {
	// CreateLeave creates a new leave request
	CreateLeave(ctx context.Context, req entities.LeaveRequestCreate) (uuid.UUID, error)

	// GetLeaveByID retrieves a leave request by ID
	GetLeaveByID(ctx context.Context, id uuid.UUID) (*entities.LeaveRequestWithDetails, error)

	// GetLeavesByFilter retrieves leave requests based on filters
	GetLeavesByFilter(ctx context.Context, filter entities.LeaveRequestFilter) ([]entities.LeaveRequestWithDetails, error)

	// GetLeavesByUser retrieves all leave requests for a specific user
	GetLeavesByUser(ctx context.Context, userID uuid.UUID) ([]entities.LeaveRequestWithDetails, error)

	// GetLeavesByDepartment retrieves all leave requests for a specific department
	GetLeavesByDepartment(ctx context.Context, departmentID uuid.UUID) ([]entities.LeaveRequestWithDetails, error)

	// UpdateLeave updates an existing leave request
	UpdateLeave(ctx context.Context, id uuid.UUID, update entities.LeaveRequestUpdate) error

	// ApproveLeave approves a leave request
	ApproveLeave(ctx context.Context, id uuid.UUID, approverID uuid.UUID) error

	// RejectLeave rejects a leave request
	RejectLeave(ctx context.Context, id uuid.UUID, approverID uuid.UUID, reason string) error

	// DeleteLeave deletes a leave request
	DeleteLeave(ctx context.Context, id uuid.UUID) error

	// ToggleLeaveStatus toggles the leave request status
	ToggleLeaveStatus(ctx context.Context, id uuid.UUID) error
}

// LeaveUseCaseImpl implements LeaveUseCase
type LeaveUseCaseImpl struct {
	repo repositories.LeaveRepository
}

// NewLeaveUseCase creates a new leave use case
func NewLeaveUseCase(repo repositories.LeaveRepository) LeaveUseCase {
	return &LeaveUseCaseImpl{
		repo: repo,
	}
}

// CreateLeave creates a new leave request
func (uc *LeaveUseCaseImpl) CreateLeave(ctx context.Context, req entities.LeaveRequestCreate) (uuid.UUID, error) {
	// Validate dates
	if req.EndDate.Before(req.StartDate) {
		return uuid.Nil, entities.ErrInvalidDateRange
	}

	// Create leave request
	leave := entities.LeaveRequest{
		StaffID:      req.StaffID,
		DepartmentID: req.DepartmentID,
		LeaveType:    req.LeaveType,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Reason:       req.Reason,
		Status:       entities.LeaveStatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return uc.repo.Create(ctx, leave)
}

// GetLeaveByID retrieves a leave request by ID
func (uc *LeaveUseCaseImpl) GetLeaveByID(ctx context.Context, id uuid.UUID) (*entities.LeaveRequestWithDetails, error) {
	return uc.repo.GetByIDWithDetails(ctx, id)
}

// GetLeavesByFilter retrieves leave requests based on filters
func (uc *LeaveUseCaseImpl) GetLeavesByFilter(ctx context.Context, filter entities.LeaveRequestFilter) ([]entities.LeaveRequestWithDetails, error) {
	return uc.repo.GetByFilter(ctx, filter)
}

// GetLeavesByUser retrieves all leave requests for a specific user
func (uc *LeaveUseCaseImpl) GetLeavesByUser(ctx context.Context, userID uuid.UUID) ([]entities.LeaveRequestWithDetails, error) {
	return uc.repo.GetByUser(ctx, userID)
}

// GetLeavesByDepartment retrieves all leave requests for a specific department
func (uc *LeaveUseCaseImpl) GetLeavesByDepartment(ctx context.Context, departmentID uuid.UUID) ([]entities.LeaveRequestWithDetails, error) {
	return uc.repo.GetByDepartment(ctx, departmentID)
}

// UpdateLeave updates an existing leave request
func (uc *LeaveUseCaseImpl) UpdateLeave(ctx context.Context, id uuid.UUID, update entities.LeaveRequestUpdate) error {
	// Validate dates if both are provided
	if update.StartDate != nil && update.EndDate != nil {
		if update.EndDate.Before(*update.StartDate) {
			return entities.ErrInvalidDateRange
		}
	}

	return uc.repo.Update(ctx, id, update)
}

// ApproveLeave approves a leave request
func (uc *LeaveUseCaseImpl) ApproveLeave(ctx context.Context, id uuid.UUID, approverID uuid.UUID) error {
	return uc.repo.UpdateStatus(ctx, id, entities.LeaveStatusApproved, &approverID, nil)
}

// RejectLeave rejects a leave request
func (uc *LeaveUseCaseImpl) RejectLeave(ctx context.Context, id uuid.UUID, approverID uuid.UUID, reason string) error {
	return uc.repo.UpdateStatus(ctx, id, entities.LeaveStatusRejected, &approverID, &reason)
}

// DeleteLeave deletes a leave request
func (uc *LeaveUseCaseImpl) DeleteLeave(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}

// ToggleLeaveStatus toggles the leave request status
func (uc *LeaveUseCaseImpl) ToggleLeaveStatus(ctx context.Context, id uuid.UUID) error {
	return uc.repo.ToggleActive(ctx, id)
}

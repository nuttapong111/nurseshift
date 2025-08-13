package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"nurseshift/employee-leave-service/internal/domain/entities"

	"github.com/google/uuid"
)

// PostgresLeaveRepository implements LeaveRepository for PostgreSQL
type PostgresLeaveRepository struct {
	db     *sql.DB
	schema string
}

// NewPostgresLeaveRepository creates a new PostgreSQL leave repository
func NewPostgresLeaveRepository(db *sql.DB, schema string) *PostgresLeaveRepository {
	return &PostgresLeaveRepository{
		db:     db,
		schema: schema,
	}
}

// Create creates a new leave request
func (r *PostgresLeaveRepository) Create(ctx context.Context, leave entities.LeaveRequest) (uuid.UUID, error) {
	query := fmt.Sprintf(`
        INSERT INTO %s.leave_requests (
            staff_id, department_id, leave_type, start_date, end_date, 
            reason, status, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id
    `, r.schema)

	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, query,
		leave.StaffID, leave.DepartmentID, leave.LeaveType, leave.StartDate, leave.EndDate,
		leave.Reason, leave.Status, leave.CreatedAt, leave.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create leave request: %w", err)
	}

	return id, nil
}

// GetByID retrieves a leave request by ID
func (r *PostgresLeaveRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.LeaveRequest, error) {
	query := fmt.Sprintf(`
        SELECT id, staff_id, department_id, leave_type, start_date, end_date,
			   reason, status, approved_by, approved_at, rejection_reason,
			   attachments, created_at, updated_at
		FROM %s.leave_requests
		WHERE id = $1
	`, r.schema)

	var leave entities.LeaveRequest
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&leave.ID, &leave.StaffID, &leave.DepartmentID, &leave.LeaveType,
		&leave.StartDate, &leave.EndDate, &leave.Reason, &leave.Status,
		&leave.ApprovedBy, &leave.ApprovedAt, &leave.RejectionReason,
		&leave.Attachments, &leave.CreatedAt, &leave.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get leave request: %w", err)
	}

	return &leave, nil
}

// GetByIDWithDetails retrieves a leave request by ID with additional details
func (r *PostgresLeaveRepository) GetByIDWithDetails(ctx context.Context, id uuid.UUID) (*entities.LeaveRequestWithDetails, error) {
	query := fmt.Sprintf(`
		SELECT 
            lr.id, lr.staff_id, lr.department_id, lr.leave_type, lr.start_date, lr.end_date,
			lr.reason, lr.status, lr.approved_by, lr.approved_at, lr.rejection_reason,
			lr.attachments, lr.created_at, lr.updated_at,
            ds.name as user_name,
			d.name as department_name,
            CONCAT(approver.first_name, ' ', approver.last_name) as approver_name
		FROM %s.leave_requests lr
        JOIN %s.department_staff ds ON lr.staff_id = ds.id
		JOIN %s.departments d ON lr.department_id = d.id
		LEFT JOIN %s.users approver ON lr.approved_by = approver.id
		WHERE lr.id = $1
    `, r.schema, r.schema, r.schema, r.schema)

	var leave entities.LeaveRequestWithDetails
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&leave.ID, &leave.StaffID, &leave.DepartmentID, &leave.LeaveType,
		&leave.StartDate, &leave.EndDate, &leave.Reason, &leave.Status,
		&leave.ApprovedBy, &leave.ApprovedAt, &leave.RejectionReason,
		&leave.Attachments, &leave.CreatedAt, &leave.UpdatedAt,
		&leave.UserName, &leave.DepartmentName, &leave.ApproverName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get leave request with details: %w", err)
	}

	return &leave, nil
}

// GetByFilter retrieves leave requests based on filters
func (r *PostgresLeaveRepository) GetByFilter(ctx context.Context, filter entities.LeaveRequestFilter) ([]entities.LeaveRequestWithDetails, error) {
	query := fmt.Sprintf(`
		SELECT 
            lr.id, lr.staff_id, lr.department_id, lr.leave_type, lr.start_date, lr.end_date,
			lr.reason, lr.status, lr.approved_by, lr.approved_at, lr.rejection_reason,
			lr.attachments, lr.created_at, lr.updated_at,
            ds.name as user_name,
			d.name as department_name,
			CONCAT(approver.first_name, ' ', approver.last_name) as approver_name
		FROM %s.leave_requests lr
        JOIN %s.department_staff ds ON lr.staff_id = ds.id
		JOIN %s.departments d ON lr.department_id = d.id
		LEFT JOIN %s.users approver ON lr.approved_by = approver.id
    `, r.schema, r.schema, r.schema, r.schema)

	var whereClauses []string
	var args []interface{}
	argIndex := 1

	if filter.Month != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("TO_CHAR(lr.start_date, 'YYYY-MM') = $%d", argIndex))
		args = append(args, *filter.Month)
		argIndex++
	}

	if filter.StaffID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("lr.staff_id = $%d", argIndex))
		args = append(args, *filter.StaffID)
		argIndex++
	}

	if filter.DepartmentID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("lr.department_id = $%d", argIndex))
		args = append(args, *filter.DepartmentID)
		argIndex++
	}

	if filter.Status != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("lr.status = $%d", argIndex))
		args = append(args, *filter.Status)
		argIndex++
	}

	if filter.StartDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("lr.start_date >= $%d", argIndex))
		args = append(args, *filter.StartDate)
		argIndex++
	}

	if filter.EndDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("lr.end_date <= $%d", argIndex))
		args = append(args, *filter.EndDate)
		argIndex++
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query += " ORDER BY lr.created_at DESC"

	// Debug log
	// NOTE: avoid logging sensitive data; this logs only schema and query shape
	fmt.Printf("[LeaveRepo] schema=%s executing query: %s\n", r.schema, query)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query leave requests: %w", err)
	}
	defer rows.Close()

	var leaves []entities.LeaveRequestWithDetails
	for rows.Next() {
		var leave entities.LeaveRequestWithDetails
		err := rows.Scan(
			&leave.ID, &leave.StaffID, &leave.DepartmentID, &leave.LeaveType,
			&leave.StartDate, &leave.EndDate, &leave.Reason, &leave.Status,
			&leave.ApprovedBy, &leave.ApprovedAt, &leave.RejectionReason,
			&leave.Attachments, &leave.CreatedAt, &leave.UpdatedAt,
			&leave.UserName, &leave.DepartmentName, &leave.ApproverName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan leave request: %w", err)
		}
		leaves = append(leaves, leave)
	}

	return leaves, nil
}

// GetByUser retrieves all leave requests for a specific staff (kept name for compatibility)
func (r *PostgresLeaveRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]entities.LeaveRequestWithDetails, error) {
	filter := entities.LeaveRequestFilter{StaffID: &userID}
	return r.GetByFilter(ctx, filter)
}

// GetByDepartment retrieves all leave requests for a specific department
func (r *PostgresLeaveRepository) GetByDepartment(ctx context.Context, departmentID uuid.UUID) ([]entities.LeaveRequestWithDetails, error) {
	filter := entities.LeaveRequestFilter{DepartmentID: &departmentID}
	return r.GetByFilter(ctx, filter)
}

// Update updates an existing leave request
func (r *PostgresLeaveRepository) Update(ctx context.Context, id uuid.UUID, update entities.LeaveRequestUpdate) error {
	setClauses := []string{"updated_at = $1"}
	args := []interface{}{time.Now()}
	argIndex := 2

	if update.LeaveType != nil {
		setClauses = append(setClauses, fmt.Sprintf("leave_type = $%d", argIndex))
		args = append(args, *update.LeaveType)
		argIndex++
	}

	if update.StartDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("start_date = $%d", argIndex))
		args = append(args, *update.StartDate)
		argIndex++
	}

	if update.EndDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("end_date = $%d", argIndex))
		args = append(args, *update.EndDate)
		argIndex++
	}

	if update.Reason != nil {
		setClauses = append(setClauses, fmt.Sprintf("reason = $%d", argIndex))
		args = append(args, *update.Reason)
		argIndex++
	}

	if update.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *update.Status)
		argIndex++
	}

	query := fmt.Sprintf("UPDATE %s.leave_requests SET %s WHERE id = $%d", r.schema, strings.Join(setClauses, ", "), argIndex)
	args = append(args, id)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update leave request: %w", err)
	}

	return nil
}

// UpdateStatus updates the status of a leave request
func (r *PostgresLeaveRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entities.LeaveStatus, approverID *uuid.UUID, rejectionReason *string) error {
	query := fmt.Sprintf(`
		UPDATE %s.leave_requests SET
			status = $2, approved_by = $3, approved_at = $4, 
			rejection_reason = $5, updated_at = $6
		WHERE id = $1
	`, r.schema)

	var approvedAt *time.Time
	if approverID != nil {
		now := time.Now()
		approvedAt = &now
	}

	_, err := r.db.ExecContext(ctx, query, id, status, approverID, approvedAt, rejectionReason, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update leave request status: %w", err)
	}

	return nil
}

// Delete deletes a leave request
func (r *PostgresLeaveRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s.leave_requests WHERE id = $1`, r.schema)

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete leave request: %w", err)
	}

	return nil
}

// ToggleActive toggles the active status (for soft delete/restore)
func (r *PostgresLeaveRepository) ToggleActive(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.leave_requests SET
			status = CASE 
				WHEN status = 'pending'::%s.leave_status THEN 'cancelled'::%s.leave_status 
				ELSE 'pending'::%s.leave_status 
			END,
			updated_at = $2
		WHERE id = $1
	`, r.schema, r.schema, r.schema, r.schema)

	_, err := r.db.ExecContext(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to toggle leave request status: %w", err)
	}

	return nil
}

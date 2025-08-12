package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"nurseshift/department-service/internal/domain/entities"

	"github.com/google/uuid"
)

// DepartmentRepository interface for department operations
type DepartmentRepository interface {
	Create(ctx context.Context, dept *entities.Department) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Department, error)
	GetAll(ctx context.Context, userID uuid.UUID) ([]*entities.Department, error)
	Update(ctx context.Context, dept *entities.Department) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetWithStats(ctx context.Context, userID uuid.UUID) ([]*entities.DepartmentWithStats, error)
	GetStaff(ctx context.Context, departmentID uuid.UUID) ([]*entities.DepartmentStaff, error)
	CreateStaff(ctx context.Context, staff *entities.DepartmentStaff) error
	DeleteStaff(ctx context.Context, staffID, departmentID uuid.UUID) error
}

// PostgresDepartmentRepository implements DepartmentRepository
type PostgresDepartmentRepository struct {
	db     *sql.DB
	schema string
}

// NewPostgresDepartmentRepository creates a new department repository
func NewPostgresDepartmentRepository(db *sql.DB, schema string) DepartmentRepository {
	return &PostgresDepartmentRepository{
		db:     db,
		schema: schema,
	}
}

// Create creates a new department
func (r *PostgresDepartmentRepository) Create(ctx context.Context, dept *entities.Department) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.departments (
			id, name, description, head_user_id, max_nurses, max_assistants, 
			settings, is_active, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, r.schema)

	_, err := r.db.ExecContext(ctx, query,
		dept.ID, dept.Name, dept.Description, dept.HeadUserID,
		dept.MaxNurses, dept.MaxAssistants, dept.Settings, dept.IsActive,
		dept.CreatedBy, dept.CreatedAt, dept.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create department: %w", err)
	}

	return nil
}

// GetByID retrieves a department by ID
func (r *PostgresDepartmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Department, error) {
	query := fmt.Sprintf(`
		SELECT id, name, description, head_user_id, max_nurses, max_assistants,
			   settings, is_active, created_by, created_at, updated_at
		FROM %s.departments
		WHERE id = $1 AND is_active = true
	`, r.schema)

	var dept entities.Department
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&dept.ID, &dept.Name, &dept.Description, &dept.HeadUserID,
		&dept.MaxNurses, &dept.MaxAssistants, &dept.Settings, &dept.IsActive,
		&dept.CreatedBy, &dept.CreatedAt, &dept.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("department not found")
		}
		return nil, fmt.Errorf("failed to get department: %w", err)
	}

	return &dept, nil
}

// GetAll retrieves all departments for a specific user
func (r *PostgresDepartmentRepository) GetAll(ctx context.Context, userID uuid.UUID) ([]*entities.Department, error) {
	query := fmt.Sprintf(`
		SELECT id, name, description, head_user_id, max_nurses, max_assistants,
			   settings, is_active, created_by, created_at, updated_at
		FROM %s.departments
		WHERE created_by = $1 AND is_active = true
		ORDER BY created_at DESC
	`, r.schema)

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query departments: %w", err)
	}
	defer rows.Close()

	var departments []*entities.Department
	for rows.Next() {
		var dept entities.Department
		err := rows.Scan(
			&dept.ID, &dept.Name, &dept.Description, &dept.HeadUserID,
			&dept.MaxNurses, &dept.MaxAssistants, &dept.Settings, &dept.IsActive,
			&dept.CreatedBy, &dept.CreatedAt, &dept.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan department: %w", err)
		}
		departments = append(departments, &dept)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating departments: %w", err)
	}

	return departments, nil
}

// Update updates a department
func (r *PostgresDepartmentRepository) Update(ctx context.Context, dept *entities.Department) error {
	query := fmt.Sprintf(`
		UPDATE %s.departments
		SET name = $1, description = $2, head_user_id = $3, max_nurses = $4,
			max_assistants = $5, settings = $6, is_active = $7, updated_at = $8
		WHERE id = $9
	`, r.schema)

	result, err := r.db.ExecContext(ctx, query,
		dept.Name, dept.Description, dept.HeadUserID, dept.MaxNurses,
		dept.MaxAssistants, dept.Settings, dept.IsActive, dept.UpdatedAt, dept.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update department: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("department not found or no changes made")
	}

	return nil
}

// Delete performs a soft delete on a department
func (r *PostgresDepartmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.departments
		SET is_active = false, updated_at = $1
		WHERE id = $2
	`, r.schema)

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete department: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("department not found")
	}

	return nil
}

// GetWithStats retrieves departments with statistics
func (r *PostgresDepartmentRepository) GetWithStats(ctx context.Context, userID uuid.UUID) ([]*entities.DepartmentWithStats, error) {
	query := fmt.Sprintf(`
		SELECT 
			d.id, d.name, d.description, d.head_user_id, d.max_nurses, d.max_assistants,
			d.settings, d.is_active, d.created_by, d.created_at, d.updated_at,
			COALESCE(staff_stats.total_employees, 0) as total_employees,
			COALESCE(staff_stats.nurse_count, 0) as nurse_count,
			COALESCE(staff_stats.assistant_count, 0) as assistant_count
		FROM %s.departments d
		LEFT JOIN (
			SELECT 
				ds.department_id,
				COUNT(*) as total_employees,
				COUNT(CASE WHEN ds.position = 'nurse' THEN 1 END) as nurse_count,
				COUNT(CASE WHEN ds.position = 'assistant' THEN 1 END) as assistant_count
			FROM %s.department_staff ds
			WHERE ds.is_active = true
			GROUP BY ds.department_id
		) staff_stats ON d.id = staff_stats.department_id
		WHERE d.created_by = $1 AND d.is_active = true
		ORDER BY d.created_at DESC
	`, r.schema, r.schema)

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query departments with stats: %w", err)
	}
	defer rows.Close()

	var departments []*entities.DepartmentWithStats
	for rows.Next() {
		var dept entities.Department
		var totalEmployees, nurseCount, assistantCount int

		err := rows.Scan(
			&dept.ID, &dept.Name, &dept.Description, &dept.HeadUserID,
			&dept.MaxNurses, &dept.MaxAssistants, &dept.Settings, &dept.IsActive,
			&dept.CreatedBy, &dept.CreatedAt, &dept.UpdatedAt,
			&totalEmployees, &nurseCount, &assistantCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan department with stats: %w", err)
		}

		deptWithStats := &entities.DepartmentWithStats{
			Department:      &dept,
			TotalEmployees:  totalEmployees,
			ActiveEmployees: totalEmployees, // Assuming all are active for now
			NurseCount:      nurseCount,
			AssistantCount:  assistantCount,
		}

		departments = append(departments, deptWithStats)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating departments with stats: %w", err)
	}

	return departments, nil
}

// GetStaff retrieves staff for a specific department
func (r *PostgresDepartmentRepository) GetStaff(ctx context.Context, departmentID uuid.UUID) ([]*entities.DepartmentStaff, error) {
	query := fmt.Sprintf(`
		SELECT 
			ds.id, ds.department_id, ds.name, ds.position, ds.phone, ds.email,
			ds.is_active, ds.created_at, ds.updated_at
		FROM %s.department_staff ds
		WHERE ds.department_id = $1 AND ds.is_active = true
		ORDER BY ds.name
	`, r.schema)

	rows, err := r.db.QueryContext(ctx, query, departmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query department staff: %w", err)
	}
	defer rows.Close()

	var staff []*entities.DepartmentStaff
	for rows.Next() {
		var s entities.DepartmentStaff
		err := rows.Scan(
			&s.ID, &s.DepartmentID, &s.Name, &s.Position, &s.Phone, &s.Email,
			&s.IsActive, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan staff: %w", err)
		}
		staff = append(staff, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating staff: %w", err)
	}

	return staff, nil
}

// CreateStaff creates a new staff member in a department
func (r *PostgresDepartmentRepository) CreateStaff(ctx context.Context, staff *entities.DepartmentStaff) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.department_staff (
			id, department_id, name, position, phone, email,
			is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, r.schema)

	_, err := r.db.ExecContext(ctx, query,
		staff.ID, staff.DepartmentID, staff.Name, staff.Position,
		staff.Phone, staff.Email, staff.IsActive, staff.CreatedAt, staff.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create department staff: %w", err)
	}

	return nil
}

// DeleteStaff deletes a staff member from a department
func (r *PostgresDepartmentRepository) DeleteStaff(ctx context.Context, staffID, departmentID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.department_staff
		WHERE id = $1 AND department_id = $2
	`, r.schema)

	result, err := r.db.ExecContext(ctx, query, staffID, departmentID)
	if err != nil {
		return fmt.Errorf("failed to delete department staff: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("staff member not found")
	}

	return nil
}

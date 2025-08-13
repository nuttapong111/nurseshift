package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domain "nurseshift/setting-service/internal/domain/entities"
	repo "nurseshift/setting-service/internal/domain/repositories"

	"github.com/google/uuid"
)

type PostgresSettingRepository struct {
	db     *sql.DB
	schema string
}

func NewPostgresSettingRepository(db *sql.DB, schema string) repo.SettingRepository {
	return &PostgresSettingRepository{db: db, schema: schema}
}

func (r *PostgresSettingRepository) GetWorkingDays(ctx context.Context, departmentID uuid.UUID) ([]domain.WorkingDay, error) {
	query := fmt.Sprintf(`SELECT id, department_id, day_of_week, is_working_day, created_at FROM %s.working_days WHERE department_id=$1 ORDER BY day_of_week`, r.schema)
	rows, err := r.db.QueryContext(ctx, query, departmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []domain.WorkingDay
	for rows.Next() {
		var d domain.WorkingDay
		if err := rows.Scan(&d.ID, &d.DepartmentID, &d.DayOfWeek, &d.IsWorkingDay, &d.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, nil
}

func (r *PostgresSettingRepository) UpsertWorkingDays(ctx context.Context, departmentID uuid.UUID, days []domain.WorkingDay) error {
	// Simple strategy: delete then insert
	delQ := fmt.Sprintf(`DELETE FROM %s.working_days WHERE department_id=$1`, r.schema)
	if _, err := r.db.ExecContext(ctx, delQ, departmentID); err != nil {
		return err
	}
	insQ := fmt.Sprintf(`INSERT INTO %s.working_days (id, department_id, day_of_week, is_working_day, created_at) VALUES ($1,$2,$3,$4,$5)`, r.schema)
	for _, d := range days {
		if d.ID == uuid.Nil {
			d.ID = uuid.New()
		}
		if d.CreatedAt.IsZero() {
			d.CreatedAt = time.Now()
		}
		if _, err := r.db.ExecContext(ctx, insQ, d.ID, departmentID, d.DayOfWeek, d.IsWorkingDay, d.CreatedAt); err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresSettingRepository) GetShifts(ctx context.Context, departmentID uuid.UUID) ([]domain.Shift, error) {
	query := fmt.Sprintf(`SELECT id, department_id, name, type, start_time, end_time, duration_hours, required_nurses, required_assistants, color, is_active, created_at, updated_at FROM %s.shifts WHERE department_id=$1 ORDER BY start_time`, r.schema)
	rows, err := r.db.QueryContext(ctx, query, departmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []domain.Shift
	for rows.Next() {
		var s domain.Shift
		var start, end time.Time
		if err := rows.Scan(&s.ID, &s.DepartmentID, &s.Name, &s.Type, &start, &end, &s.DurationHours, &s.RequiredNurses, &s.RequiredAssistants, &s.Color, &s.IsActive, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		s.StartTime = start.Format("15:04")
		s.EndTime = end.Format("15:04")
		result = append(result, s)
	}
	return result, nil
}

func (r *PostgresSettingRepository) CreateShift(ctx context.Context, shift domain.Shift) (uuid.UUID, error) {
	if shift.ID == uuid.Nil {
		shift.ID = uuid.New()
	}
	query := fmt.Sprintf(`INSERT INTO %s.shifts (id, department_id, name, type, start_time, end_time, duration_hours, required_nurses, required_assistants, color, is_active, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11, NOW(), NOW())`, r.schema)
	// Parse HH:MM to time
	start, _ := time.Parse("15:04", shift.StartTime)
	end, _ := time.Parse("15:04", shift.EndTime)
	_, err := r.db.ExecContext(ctx, query, shift.ID, shift.DepartmentID, shift.Name, shift.Type, start, end, shift.DurationHours, shift.RequiredNurses, shift.RequiredAssistants, shift.Color, shift.IsActive)
	if err != nil {
		return uuid.Nil, err
	}
	return shift.ID, nil
}

func (r *PostgresSettingRepository) UpdateShift(ctx context.Context, shift domain.Shift) error {
	query := fmt.Sprintf(`UPDATE %s.shifts SET name=$2, type=$3, start_time=$4, end_time=$5, required_nurses=$6, required_assistants=$7, color=$8, updated_at=NOW() WHERE id=$1`, r.schema)
	start, _ := time.Parse("15:04", shift.StartTime)
	end, _ := time.Parse("15:04", shift.EndTime)
	_, err := r.db.ExecContext(ctx, query, shift.ID, shift.Name, shift.Type, start, end, shift.RequiredNurses, shift.RequiredAssistants, shift.Color)
	return err
}

func (r *PostgresSettingRepository) UpdateShiftStatus(ctx context.Context, shiftID uuid.UUID, isActive bool) error {
	query := fmt.Sprintf(`UPDATE %s.shifts SET is_active=$2, updated_at=NOW() WHERE id=$1`, r.schema)
	_, err := r.db.ExecContext(ctx, query, shiftID, isActive)
	return err
}

func (r *PostgresSettingRepository) DeleteShift(ctx context.Context, shiftID uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s.shifts WHERE id=$1`, r.schema)
	_, err := r.db.ExecContext(ctx, query, shiftID)
	return err
}

func (r *PostgresSettingRepository) GetHolidays(ctx context.Context, departmentID uuid.UUID) ([]domain.Holiday, error) {
	query := fmt.Sprintf(`SELECT id, department_id, name, start_date, end_date, is_recurring, created_at, updated_at FROM %s.holidays WHERE department_id=$1 ORDER BY start_date`, r.schema)
	rows, err := r.db.QueryContext(ctx, query, departmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []domain.Holiday
	for rows.Next() {
		var h domain.Holiday
		if err := rows.Scan(&h.ID, &h.DepartmentID, &h.Name, &h.StartDate, &h.EndDate, &h.IsRecurring, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, nil
}

func (r *PostgresSettingRepository) CreateHoliday(ctx context.Context, holiday domain.Holiday) (uuid.UUID, error) {
	if holiday.ID == uuid.Nil {
		holiday.ID = uuid.New()
	}
	query := fmt.Sprintf(`INSERT INTO %s.holidays (id, department_id, name, start_date, end_date, is_recurring, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,NOW(),NOW())`, r.schema)
	_, err := r.db.ExecContext(ctx, query, holiday.ID, holiday.DepartmentID, holiday.Name, holiday.StartDate, holiday.EndDate, holiday.IsRecurring)
	if err != nil {
		return uuid.Nil, err
	}
	return holiday.ID, nil
}

func (r *PostgresSettingRepository) DeleteHoliday(ctx context.Context, holidayID uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s.holidays WHERE id=$1`, r.schema)
	_, err := r.db.ExecContext(ctx, query, holidayID)
	return err
}

func (r *PostgresSettingRepository) UpdateHoliday(ctx context.Context, holiday domain.Holiday) error {
	query := fmt.Sprintf(`UPDATE %s.holidays SET name=$2, start_date=$3, end_date=$4, is_recurring=$5, updated_at=NOW() WHERE id=$1`, r.schema)
	_, err := r.db.ExecContext(ctx, query, holiday.ID, holiday.Name, holiday.StartDate, holiday.EndDate, holiday.IsRecurring)
	return err
}

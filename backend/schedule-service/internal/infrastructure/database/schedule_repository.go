package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type ScheduleRecord struct {
	ID           string
	DepartmentID string
	UserID       string
	StaffID      string
	ShiftID      string
	ScheduleDate string // YYYY-MM-DD
	Status       string
	Notes        sql.NullString
}

// ScheduleWithRole joins schedule with department_users to know role in department
type ScheduleWithRole struct {
	ScheduleRecord
	DepartmentRole string
	UserFirstName  string
	UserLastName   string
}

type ScheduleRepository struct {
	conn   *Connection
	schema string
}

func NewScheduleRepository(conn *Connection) *ScheduleRepository {
	schema := "nurse_shift"
	if conn.Config != nil {
		// best-effort: parse from DSN not implemented; default schema
	}
	return &ScheduleRepository{conn: conn, schema: schema}
}

func (r *ScheduleRepository) table() string {
	return fmt.Sprintf("%s.schedules", r.schema)
}

func (r *ScheduleRepository) List(ctx context.Context, departmentID, month string) ([]ScheduleRecord, error) {
	base := fmt.Sprintf("SELECT id, department_id, user_id, shift_id, to_char(schedule_date,'YYYY-MM-DD'), status, notes FROM %s WHERE 1=1", r.table())
	var args []any
	idx := 1
	if departmentID != "" {
		base += fmt.Sprintf(" AND department_id = $%d", idx)
		args = append(args, departmentID)
		idx++
	}
	if month != "" { // YYYY-MM
		base += fmt.Sprintf(" AND to_char(schedule_date,'YYYY-MM') = $%d", idx)
		args = append(args, month)
		idx++
	}
	base += " ORDER BY schedule_date ASC"

	rows, err := r.conn.DB.QueryContext(ctx, base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ScheduleRecord
	for rows.Next() {
		var rec ScheduleRecord
		if err := rows.Scan(&rec.ID, &rec.DepartmentID, &rec.UserID, &rec.ShiftID, &rec.ScheduleDate, &rec.Status, &rec.Notes); err != nil {
			return nil, err
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}

func (r *ScheduleRepository) Create(ctx context.Context, rec *ScheduleRecord) error {
	q := fmt.Sprintf("INSERT INTO %s (id, department_id, user_id, shift_id, schedule_date, status, notes, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,COALESCE($6,'assigned'),$7,NOW(),NOW())", r.table())
	_, err := r.conn.DB.ExecContext(ctx, q, rec.ID, rec.DepartmentID, rec.UserID, rec.ShiftID, rec.ScheduleDate, rec.Status, rec.Notes)
	return err
}

func (r *ScheduleRepository) Update(ctx context.Context, id string, status *string, notes *string, shiftID *string) error {
	set := "updated_at = NOW()"
	var args []any
	idx := 1
	if status != nil {
		set += fmt.Sprintf(", status = $%d", idx)
		args = append(args, *status)
		idx++
	}
	if notes != nil {
		set += fmt.Sprintf(", notes = $%d", idx)
		args = append(args, *notes)
		idx++
	}
	if shiftID != nil {
		set += fmt.Sprintf(", shift_id = $%d", idx)
		args = append(args, *shiftID)
		idx++
	}
	q := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", r.table(), set, idx)
	args = append(args, id)
	_, err := r.conn.DB.ExecContext(ctx, q, args...)
	return err
}

// EnsureStaffSchedulingSchema adds staff_id column and relaxes user_id NOT NULL for staff-based scheduling
func (r *ScheduleRepository) EnsureStaffSchedulingSchema(ctx context.Context) error {
	// Add staff_id column
	q1 := fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS staff_id UUID REFERENCES %s.department_staff(id)", r.table(), r.schema)
	if _, err := r.conn.DB.ExecContext(ctx, q1); err != nil {
		return err
	}
	// Allow NULL user_id
	q2 := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN user_id DROP NOT NULL", r.table())
	if _, err := r.conn.DB.ExecContext(ctx, q2); err != nil { /* ignore if already nullable */
	}
	// Unique index to prevent duplicate (staff_id, date, shift)
	q3 := fmt.Sprintf("CREATE UNIQUE INDEX IF NOT EXISTS idx_schedules_staff_date_shift ON %s (staff_id, schedule_date, shift_id)", r.table())
	if _, err := r.conn.DB.ExecContext(ctx, q3); err != nil {
		return err
	}
	return nil
}

func (r *ScheduleRepository) Delete(ctx context.Context, id string) error {
	q := fmt.Sprintf("DELETE FROM %s WHERE id=$1", r.table())
	_, err := r.conn.DB.ExecContext(ctx, q, id)
	return err
}

// DeleteByDepartmentAndMonth deletes all schedules for a department in a given YYYY-MM month
func (r *ScheduleRepository) DeleteByDepartmentAndMonth(ctx context.Context, departmentID string, month string) error {
	q := fmt.Sprintf("DELETE FROM %s WHERE department_id=$1 AND to_char(schedule_date,'YYYY-MM')=$2", r.table())
	_, err := r.conn.DB.ExecContext(ctx, q, departmentID, month)
	return err
}

type ShiftRecord struct {
	ID            string
	DepartmentID  string
	Name          string
	Type          string
	StartTime     string
	EndTime       string
	RequiredNurse int
	RequiredAsst  int
	Color         string
}

func (r *ScheduleRepository) ListShifts(ctx context.Context, departmentID string) ([]ShiftRecord, error) {
	q := fmt.Sprintf("SELECT id, department_id, name, type, to_char(start_time,'HH24:MI'), to_char(end_time,'HH24:MI'), required_nurses, required_assistants, color FROM %s.shifts WHERE is_active = true", r.schema)
	args := []any{}
	if departmentID != "" {
		q += " AND department_id = $1"
		args = append(args, departmentID)
	}
	q += " ORDER BY name"
	rows, err := r.conn.DB.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ShiftRecord
	for rows.Next() {
		var rec ShiftRecord
		if err := rows.Scan(&rec.ID, &rec.DepartmentID, &rec.Name, &rec.Type, &rec.StartTime, &rec.EndTime, &rec.RequiredNurse, &rec.RequiredAsst, &rec.Color); err != nil {
			return nil, err
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}

// Working day config (0=Sunday .. 6=Saturday)
func (r *ScheduleRepository) ListWorkingDays(ctx context.Context, departmentID string) (map[int]bool, error) {
	q := fmt.Sprintf("SELECT day_of_week, is_working_day FROM %s.working_days WHERE department_id=$1", r.schema)
	rows, err := r.conn.DB.QueryContext(ctx, q, departmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Strategy:
	// - หากยังไม่เคยตั้งค่าเลย (ไม่มีเรคคอร์ด) ให้ค่าเริ่มต้นเป็น "ทำงานทุกวัน"
	// - หากมีการตั้งค่าบางวันแล้ว ให้ถือว่าเป็นโหมดกำหนดเอง: วันใดที่ไม่ได้ระบุ ให้เป็น "ไม่ทำงาน" โดยปริยาย
	//   เพื่อหลีกเลี่ยงเคสที่ระบุเฉพาะ จ-ศ แล้ว ส-อา ไม่มีเรคคอร์ดแต่ถูกนับเป็นวันทำงาน
	tmp := map[int]bool{}
	hasAnyRow := false
	for rows.Next() {
		var d int
		var w bool
		if err := rows.Scan(&d, &w); err != nil {
			return nil, err
		}
		tmp[d] = w
		hasAnyRow = true
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	res := map[int]bool{}
	if !hasAnyRow {
		// Try to fetch from setting-service as a fallback
		serviceURL := os.Getenv("SETTING_SERVICE_URL")
		if serviceURL == "" {
			serviceURL = "http://localhost:8085"
		}
		reqURL := fmt.Sprintf("%s/api/v1/settings?departmentId=%s", serviceURL, url.QueryEscape(departmentID))
		httpClient := &http.Client{Timeout: 2 * time.Second}
		if resp, err := httpClient.Get(reqURL); err == nil && resp != nil && resp.StatusCode < 300 {
			defer resp.Body.Close()
			var payload struct {
				Data struct {
					WorkingDays []struct {
						DayOfWeek    int  `json:"dayOfWeek"`
						IsWorkingDay bool `json:"isWorkingDay"`
					} `json:"workingDays"`
				} `json:"data"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&payload); err == nil && len(payload.Data.WorkingDays) > 0 {
				// custom mode → unspecified days = not working by default
				res = map[int]bool{0: false, 1: false, 2: false, 3: false, 4: false, 5: false, 6: false}
				for _, w := range payload.Data.WorkingDays {
					if w.DayOfWeek >= 0 && w.DayOfWeek <= 6 {
						res[w.DayOfWeek] = w.IsWorkingDay
					}
				}
				return res, nil
			}
		}
		// fallback ultimate: working everyday
		res = map[int]bool{0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true}
	} else {
		// custom mode → unspecified days = not working by default
		res = map[int]bool{0: false, 1: false, 2: false, 3: false, 4: false, 5: false, 6: false}
		for d, w := range tmp {
			res[d] = w
		}
	}
	return res, nil
}

type Holiday struct {
	Start string
	End   string
}

// ListHolidaysForMonth returns holidays overlapping the given YYYY-MM month
func (r *ScheduleRepository) ListHolidaysForMonth(ctx context.Context, departmentID string, month string) ([]Holiday, error) {
	// select where range overlaps month
	q := fmt.Sprintf(`
        SELECT to_char(start_date,'YYYY-MM-DD'), to_char(end_date,'YYYY-MM-DD')
        FROM %s.holidays 
        WHERE department_id = $1
          AND (
            to_char(start_date,'YYYY-MM') = $2 
            OR to_char(end_date,'YYYY-MM') = $2
            OR (
              start_date <= to_date($2 || '-01','YYYY-MM-DD') 
              AND end_date >= to_date($2 || '-28','YYYY-MM-DD')
            )
          )
    `, r.schema)
	rows, err := r.conn.DB.QueryContext(ctx, q, departmentID, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Holiday
	for rows.Next() {
		var h Holiday
		if err := rows.Scan(&h.Start, &h.End); err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	return out, rows.Err()
}

type LeaveRange struct {
	StaffID string
	Start   string
	End     string
}

// ListLeavesForMonth returns approved leave ranges overlapping the month
func (r *ScheduleRepository) ListLeavesForMonth(ctx context.Context, departmentID string, month string) ([]LeaveRange, error) {
	q := fmt.Sprintf(`
        SELECT staff_id, to_char(start_date,'YYYY-MM-DD'), to_char(end_date,'YYYY-MM-DD')
        FROM %s.leave_requests
        WHERE department_id = $1
          AND status <> 'cancelled'
          AND (
            to_char(start_date,'YYYY-MM') = $2 
            OR to_char(end_date,'YYYY-MM') = $2
            OR (
              start_date <= to_date($2 || '-01','YYYY-MM-DD') 
              AND end_date >= to_date($2 || '-28','YYYY-MM-DD')
            )
          )
    `, r.schema)
	rows, err := r.conn.DB.QueryContext(ctx, q, departmentID, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []LeaveRange
	for rows.Next() {
		var rec LeaveRange
		if err := rows.Scan(&rec.StaffID, &rec.Start, &rec.End); err != nil {
			return nil, err
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}

// GetPriorityValue returns numeric setting from scheduling_priorities.config->>'value'
func (r *ScheduleRepository) GetPriorityValue(ctx context.Context, departmentID string, priorityName string) (sql.NullInt64, error) {
	q := fmt.Sprintf(`
        SELECT NULLIF(TRIM((config->>'value')), '')::bigint
        FROM %s.scheduling_priorities
        WHERE department_id = $1 AND name = $2 AND is_active = true
        ORDER BY priority_order ASC
        LIMIT 1
    `, r.schema)
	var out sql.NullInt64
	err := r.conn.DB.QueryRowContext(ctx, q, departmentID, priorityName).Scan(&out)
	if err == sql.ErrNoRows {
		return sql.NullInt64{}, nil
	}
	return out, err
}

type Assignment struct {
	ID           string
	DepartmentID string
	UserID       string
	StaffID      string
	ShiftID      string
	ScheduleDate string
	Status       string
	Notes        sql.NullString
}

func (r *ScheduleRepository) BulkInsertAssignments(ctx context.Context, items []Assignment) error {
	if len(items) == 0 {
		return nil
	}
	q := fmt.Sprintf("INSERT INTO %s (id, department_id, user_id, shift_id, schedule_date, status, notes, created_at, updated_at) VALUES ", r.table())
	args := []any{}
	for i, a := range items {
		if i > 0 {
			q += ","
		}
		base := i*7 + 1
		q += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,COALESCE($%d,'assigned'),$%d,NOW(),NOW())", base, base+1, base+2, base+3, base+4, base+5, base+6)
		args = append(args, a.ID, a.DepartmentID, a.UserID, a.ShiftID, a.ScheduleDate, a.Status, a.Notes)
	}
	q += " ON CONFLICT (user_id, schedule_date, shift_id) DO NOTHING"
	_, err := r.conn.DB.ExecContext(ctx, q, args...)
	return err
}

// BulkInsertAssignmentsStaff inserts using staff_id column
func (r *ScheduleRepository) BulkInsertAssignmentsStaff(ctx context.Context, items []Assignment) error {
	if len(items) == 0 {
		return nil
	}
	q := fmt.Sprintf("INSERT INTO %s (id, department_id, staff_id, shift_id, schedule_date, status, notes, created_at, updated_at) VALUES ", r.table())
	args := []any{}
	for i, a := range items {
		if i > 0 {
			q += ","
		}
		base := i*7 + 1
		q += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,COALESCE($%d,'assigned'),$%d,NOW(),NOW())", base, base+1, base+2, base+3, base+4, base+5, base+6)
		args = append(args, a.ID, a.DepartmentID, a.StaffID, a.ShiftID, a.ScheduleDate, a.Status, a.Notes)
	}
	q += " ON CONFLICT (staff_id, schedule_date, shift_id) DO NOTHING"
	_, err := r.conn.DB.ExecContext(ctx, q, args...)
	return err
}

// ScheduleWithStaff joins schedules with department_staff to get name/position
type ScheduleWithStaff struct {
	ID           string
	DepartmentID string
	StaffID      string
	ShiftID      string
	ScheduleDate string
	Status       string
	Notes        sql.NullString
	StaffName    string
	StaffRole    string
}

func (r *ScheduleRepository) ListWithStaff(ctx context.Context, departmentID, month string) ([]ScheduleWithStaff, error) {
	base := fmt.Sprintf(`
        SELECT s.id, s.department_id, s.staff_id, s.shift_id, to_char(s.schedule_date,'YYYY-MM-DD'), s.status, s.notes,
               ds.name, ds.position
        FROM %s s
        LEFT JOIN %s.department_staff ds ON ds.id = s.staff_id
        WHERE 1=1`, r.table(), r.schema)
	var args []any
	idx := 1
	if departmentID != "" {
		base += fmt.Sprintf(" AND s.department_id = $%d", idx)
		args = append(args, departmentID)
		idx++
	}
	if month != "" {
		base += fmt.Sprintf(" AND to_char(s.schedule_date,'YYYY-MM') = $%d", idx)
		args = append(args, month)
		idx++
	}
	base += " ORDER BY s.schedule_date ASC"

	rows, err := r.conn.DB.QueryContext(ctx, base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ScheduleWithStaff
	for rows.Next() {
		var rec ScheduleWithStaff
		if err := rows.Scan(&rec.ID, &rec.DepartmentID, &rec.StaffID, &rec.ShiftID, &rec.ScheduleDate, &rec.Status, &rec.Notes, &rec.StaffName, &rec.StaffRole); err != nil {
			return nil, err
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}

// DeleteAssignmentByStaffAndShift deletes an assignment for a specific staff member, shift, and date
func (r *ScheduleRepository) DeleteAssignmentByStaffAndShift(ctx context.Context, staffID, shiftID, date string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s 
		WHERE staff_id = $1 AND shift_id = $2 AND schedule_date = $3
	`, r.table())

	_, err := r.conn.DB.ExecContext(ctx, query, staffID, shiftID, date)
	return err
}

// List schedules with role for aggregation
func (r *ScheduleRepository) ListWithRole(ctx context.Context, departmentID, month string) ([]ScheduleWithRole, error) {
	base := fmt.Sprintf(`
        SELECT s.id, s.department_id, s.user_id, s.shift_id, to_char(s.schedule_date,'YYYY-MM-DD'), s.status, s.notes,
               du.department_role, u.first_name, u.last_name
        FROM %s s
        JOIN %s.department_users du
          ON du.department_id = s.department_id AND du.user_id = s.user_id
        JOIN %s.users u ON u.id = s.user_id
        WHERE 1=1`, r.table(), r.schema, r.schema)
	var args []any
	idx := 1
	if departmentID != "" {
		base += fmt.Sprintf(" AND s.department_id = $%d", idx)
		args = append(args, departmentID)
		idx++
	}
	if month != "" {
		base += fmt.Sprintf(" AND to_char(s.schedule_date,'YYYY-MM') = $%d", idx)
		args = append(args, month)
		idx++
	}
	base += " ORDER BY s.schedule_date ASC"

	rows, err := r.conn.DB.QueryContext(ctx, base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ScheduleWithRole
	for rows.Next() {
		var rec ScheduleWithRole
		if err := rows.Scan(&rec.ID, &rec.DepartmentID, &rec.UserID, &rec.ShiftID, &rec.ScheduleDate, &rec.Status, &rec.Notes, &rec.DepartmentRole, &rec.UserFirstName, &rec.UserLastName); err != nil {
			return nil, err
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}

type DepartmentUser struct {
	UserID         string
	DepartmentRole string
}

func (r *ScheduleRepository) ListDepartmentUsers(ctx context.Context, departmentID string) ([]DepartmentUser, error) {
	q := fmt.Sprintf("SELECT user_id, department_role FROM %s.department_users WHERE department_id = $1", r.schema)
	rows, err := r.conn.DB.QueryContext(ctx, q, departmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DepartmentUser
	for rows.Next() {
		var du DepartmentUser
		if err := rows.Scan(&du.UserID, &du.DepartmentRole); err != nil {
			return nil, err
		}
		out = append(out, du)
	}
	return out, rows.Err()
}

type DepartmentStaff struct {
	ID           string
	DepartmentID string
	Name         string
	Position     string
}

func (r *ScheduleRepository) ListDepartmentStaff(ctx context.Context, departmentID string) ([]DepartmentStaff, error) {
	q := fmt.Sprintf("SELECT id, department_id, name, position FROM %s.department_staff WHERE department_id = $1 AND is_active = true", r.schema)
	rows, err := r.conn.DB.QueryContext(ctx, q, departmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DepartmentStaff
	for rows.Next() {
		var s DepartmentStaff
		if err := rows.Scan(&s.ID, &s.DepartmentID, &s.Name, &s.Position); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// AssignmentInterval represents an assignment on a date with concrete time window
type AssignmentInterval struct {
	StaffID   string
	ShiftID   string
	StartTime string
	EndTime   string
}

// ListAssignmentsWithShiftForDate joins schedules with shifts to get time intervals for a specific date
func (r *ScheduleRepository) ListAssignmentsWithShiftForDate(ctx context.Context, departmentID, date string) ([]AssignmentInterval, error) {
	q := fmt.Sprintf(`
		SELECT s.staff_id, s.shift_id, to_char(sh.start_time,'HH24:MI'), to_char(sh.end_time,'HH24:MI')
		FROM %s s
		JOIN %s.shifts sh ON sh.id = s.shift_id
		WHERE s.department_id = $1
		  AND to_char(s.schedule_date,'YYYY-MM-DD') = $2
		  AND s.staff_id IS NOT NULL
	`, r.table(), r.schema)
	rows, err := r.conn.DB.QueryContext(ctx, q, departmentID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AssignmentInterval
	for rows.Next() {
		var rec AssignmentInterval
		if err := rows.Scan(&rec.StaffID, &rec.ShiftID, &rec.StartTime, &rec.EndTime); err != nil {
			return nil, err
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}

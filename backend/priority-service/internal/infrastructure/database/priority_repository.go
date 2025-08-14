package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// PriorityRecord represents a row in scheduling_priorities
type PriorityRecord struct {
	ID            string
	DepartmentID  string
	Name          string
	Description   sql.NullString
	PriorityOrder int
	Config        sql.NullString
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// PriorityRepository provides CRUD for scheduling_priorities
type PriorityRepository struct {
	conn   *Connection
	schema string
}

func NewPriorityRepository(conn *Connection) *PriorityRepository {
	schema := conn.Config.Database.Schema
	if schema == "" {
		schema = "public"
	}
	return &PriorityRepository{conn: conn, schema: schema}
}

func (r *PriorityRepository) table() string {
	return fmt.Sprintf("%s.scheduling_priorities", r.schema)
}

func (r *PriorityRepository) GetByID(ctx context.Context, id string) (*PriorityRecord, error) {
	q := fmt.Sprintf("SELECT id, department_id, name, description, priority_order, config, is_active, created_at, updated_at FROM %s WHERE id = $1", r.table())
	rec := &PriorityRecord{}
	err := r.conn.DB.QueryRowContext(ctx, q, id).Scan(
		&rec.ID, &rec.DepartmentID, &rec.Name, &rec.Description, &rec.PriorityOrder, &rec.Config, &rec.IsActive, &rec.CreatedAt, &rec.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (r *PriorityRepository) ListByDepartment(ctx context.Context, departmentID string) ([]PriorityRecord, error) {
	q := fmt.Sprintf("SELECT id, department_id, name, description, priority_order, config, is_active, created_at, updated_at FROM %s WHERE department_id = $1 ORDER BY priority_order ASC", r.table())
	rows, err := r.conn.DB.QueryContext(ctx, q, departmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PriorityRecord
	for rows.Next() {
		var rec PriorityRecord
		if err := rows.Scan(&rec.ID, &rec.DepartmentID, &rec.Name, &rec.Description, &rec.PriorityOrder, &rec.Config, &rec.IsActive, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, rec)
	}
	return items, rows.Err()
}

func (r *PriorityRepository) Insert(ctx context.Context, rec *PriorityRecord) error {
	q := fmt.Sprintf("INSERT INTO %s (department_id, name, description, priority_order, config, is_active, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,NOW(),NOW()) RETURNING id", r.table())
	return r.conn.DB.QueryRowContext(ctx, q, rec.DepartmentID, rec.Name, rec.Description, rec.PriorityOrder, rec.Config, rec.IsActive).Scan(&rec.ID)
}

// UpdateSetting updates the JSONB config value for a priority
func (r *PriorityRepository) UpdateSetting(ctx context.Context, id string, value int) error {
	// Use jsonb_set for better compatibility across Postgres versions
	q := fmt.Sprintf("UPDATE %s SET config = jsonb_set(COALESCE(config,'{}'::jsonb), '{value}', to_jsonb($1::int), true), updated_at = NOW() WHERE id=$2", r.table())
	_, err := r.conn.DB.ExecContext(ctx, q, value, id)
	return err
}

func (r *PriorityRepository) UpdateActive(ctx context.Context, id string, isActive bool) error {
	q := fmt.Sprintf("UPDATE %s SET is_active = $1, updated_at = NOW() WHERE id = $2", r.table())
	_, err := r.conn.DB.ExecContext(ctx, q, isActive, id)
	return err
}

// UpdateOrderAndReorder moves one record to a new position and compacts others sequentially starting from 1
func (r *PriorityRepository) UpdateOrderAndReorder(ctx context.Context, id string, newOrder int) error {
	tx, err := r.conn.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// Fetch department for the target and current list
	var departmentID string
	var currentOrder int
	qGet := fmt.Sprintf("SELECT department_id, priority_order FROM %s WHERE id=$1", r.table())
	if err := tx.QueryRowContext(ctx, qGet, id).Scan(&departmentID, &currentOrder); err != nil {
		return err
	}

	qList := fmt.Sprintf("SELECT id FROM %s WHERE department_id=$1 ORDER BY priority_order ASC", r.table())
	rows, err := tx.QueryContext(ctx, qList, departmentID)
	if err != nil {
		return err
	}
	var ids []string
	for rows.Next() {
		var pid string
		if err := rows.Scan(&pid); err != nil {
			rows.Close()
			return err
		}
		ids = append(ids, pid)
	}
	rows.Close()

	if newOrder < 1 {
		newOrder = 1
	}
	if newOrder > len(ids) {
		newOrder = len(ids)
	}

	// Build new ordering by moving id to new index
	var reordered []string
	for _, pid := range ids {
		if pid != id {
			reordered = append(reordered, pid)
		}
	}
	idx := newOrder - 1
	if idx < 0 {
		idx = 0
	}
	if idx > len(reordered) {
		idx = len(reordered)
	}
	reordered = append(reordered[:idx], append([]string{id}, reordered[idx:]...)...)

	// Persist sequential levels
	qUpd := fmt.Sprintf("UPDATE %s SET priority_order = $1, updated_at = NOW() WHERE id = $2", r.table())
	for i, pid := range reordered {
		if _, err := tx.ExecContext(ctx, qUpd, i+1, pid); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// FindOrCreateDefaults ensures the department has default priorities; returns list after ensuring
func (r *PriorityRepository) FindOrCreateDefaults(ctx context.Context, departmentID string) ([]PriorityRecord, error) {
	items, err := r.ListByDepartment(ctx, departmentID)
	if err != nil {
		return nil, err
	}
	if len(items) > 0 {
		return items, nil
	}

	// Seed defaults (names align with frontend expectations)
	defaults := []PriorityRecord{
		{DepartmentID: departmentID, Name: "วันที่ขอหยุด", Description: sql.NullString{String: "ระบบจะหลีกเลี่ยงการจัดเวรในวันที่พนักงานขอหยุด", Valid: true}, PriorityOrder: 1, IsActive: true},
		{DepartmentID: departmentID, Name: "จำนวนเวรเท่ากันในแต่ละประเภท", Description: sql.NullString{String: "กระจายจำนวนเวรแต่ละประเภท (เช้า/บ่าย/ดึก) ให้แต่ละคนได้เท่าๆ กัน", Valid: true}, PriorityOrder: 2, IsActive: true},
		{DepartmentID: departmentID, Name: "จำนวนเวรดึกติดต่อกัน", Description: sql.NullString{String: "จำกัดจำนวนเวรดึกที่พนักงานคนหนึ่งทำติดกันไม่เกิน X วัน", Valid: true}, PriorityOrder: 3, IsActive: true},
		{DepartmentID: departmentID, Name: "จำนวนเวรติดต่อกัน", Description: sql.NullString{String: "จำกัดจำนวนเวรทุกประเภทที่พนักงานคนหนึ่งทำติดกันไม่เกิน X วัน", Valid: true}, PriorityOrder: 4, IsActive: true},
		{DepartmentID: departmentID, Name: "จำนวนชั่วโมงทำงานสูงสุดติดต่อกันโดยไม่พัก", Description: sql.NullString{String: "จำกัดจำนวนชั่วโมงการทำงานต่อเนื่องโดยไม่พักไม่เกิน X ชั่วโมง", Valid: true}, PriorityOrder: 5, IsActive: true},
		{DepartmentID: departmentID, Name: "จำนวนชั่วโมงการทำงานทั้งหมด", Description: sql.NullString{String: "จำกัดจำนวนชั่วโมงทำงานรวมในช่วงเวลาที่กำหนดไม่เกิน X ชั่วโมง", Valid: true}, PriorityOrder: 6, IsActive: true},
	}

	for i := range defaults {
		if err := r.Insert(ctx, &defaults[i]); err != nil {
			return nil, err
		}
	}
	return r.ListByDepartment(ctx, departmentID)
}

// GetMeta returns department_id and priority_level for an id
func (r *PriorityRepository) GetMeta(ctx context.Context, id string) (departmentID string, order int, err error) {
	q := fmt.Sprintf("SELECT department_id, priority_order FROM %s WHERE id=$1", r.table())
	err = r.conn.DB.QueryRowContext(ctx, q, id).Scan(&departmentID, &order)
	return
}

// SwapOrder swaps order of two ids in same department
func (r *PriorityRepository) SwapOrder(ctx context.Context, id1, id2 string) error {
	tx, err := r.conn.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var dept1, dept2 string
	var order1, order2 int
	qGet := fmt.Sprintf("SELECT department_id, priority_order FROM %s WHERE id=$1", r.table())
	if err := tx.QueryRowContext(ctx, qGet, id1).Scan(&dept1, &order1); err != nil {
		return err
	}
	if err := tx.QueryRowContext(ctx, qGet, id2).Scan(&dept2, &order2); err != nil {
		return err
	}
	if dept1 != dept2 {
		return fmt.Errorf("priorities in different departments")
	}

	// 3-step swap using department-scoped temporary value (lower than any existing)
	var minOrder int
	qMin := fmt.Sprintf("SELECT COALESCE(MIN(priority_order),0) FROM %s WHERE department_id=$1", r.table())
	if err := tx.QueryRowContext(ctx, qMin, dept1).Scan(&minOrder); err != nil {
		return err
	}
	temp := minOrder - 1

	qTmp := fmt.Sprintf("UPDATE %s SET priority_order = $1, updated_at = NOW() WHERE id=$2", r.table())
	if _, err := tx.ExecContext(ctx, qTmp, temp, id1); err != nil {
		return err
	}

	qSet := fmt.Sprintf("UPDATE %s SET priority_order = $1, updated_at = NOW() WHERE id=$2", r.table())
	if _, err := tx.ExecContext(ctx, qSet, order1, id2); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, qSet, order2, id1); err != nil {
		return err
	}

	return tx.Commit()
}

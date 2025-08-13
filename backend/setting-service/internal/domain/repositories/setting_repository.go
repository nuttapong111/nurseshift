package repositories

import (
	"context"
	"nurseshift/setting-service/internal/domain/entities"

	"github.com/google/uuid"
)

// SettingRepository defines methods to manage department settings
type SettingRepository interface {
	GetWorkingDays(ctx context.Context, departmentID uuid.UUID) ([]entities.WorkingDay, error)
	UpsertWorkingDays(ctx context.Context, departmentID uuid.UUID, days []entities.WorkingDay) error

	GetShifts(ctx context.Context, departmentID uuid.UUID) ([]entities.Shift, error)
	CreateShift(ctx context.Context, shift entities.Shift) (uuid.UUID, error)
	UpdateShift(ctx context.Context, shift entities.Shift) error
	UpdateShiftStatus(ctx context.Context, shiftID uuid.UUID, isActive bool) error
	DeleteShift(ctx context.Context, shiftID uuid.UUID) error

	GetHolidays(ctx context.Context, departmentID uuid.UUID) ([]entities.Holiday, error)
	CreateHoliday(ctx context.Context, holiday entities.Holiday) (uuid.UUID, error)
	DeleteHoliday(ctx context.Context, holidayID uuid.UUID) error
	UpdateHoliday(ctx context.Context, holiday entities.Holiday) error
}

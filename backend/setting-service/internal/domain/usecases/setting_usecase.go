package usecases

import (
	"context"
	"nurseshift/setting-service/internal/domain/entities"
	"nurseshift/setting-service/internal/domain/repositories"

	"github.com/google/uuid"
)

type SettingUseCase interface {
	GetSettings(ctx context.Context, departmentID uuid.UUID) (*entities.SettingsAggregate, error)
	UpdateWorkingDays(ctx context.Context, departmentID uuid.UUID, days []entities.WorkingDay) error
	CreateShift(ctx context.Context, shift entities.Shift) (uuid.UUID, error)
	UpdateShift(ctx context.Context, shift entities.Shift) error
	ToggleShift(ctx context.Context, shiftID uuid.UUID, isActive bool) error
	DeleteShift(ctx context.Context, shiftID uuid.UUID) error
	CreateHoliday(ctx context.Context, holiday entities.Holiday) (uuid.UUID, error)
	DeleteHoliday(ctx context.Context, holidayID uuid.UUID) error
	UpdateHoliday(ctx context.Context, holiday entities.Holiday) error
}

type SettingUseCaseImpl struct {
	repo repositories.SettingRepository
}

func NewSettingUseCase(repo repositories.SettingRepository) SettingUseCase {
	return &SettingUseCaseImpl{repo: repo}
}

func (uc *SettingUseCaseImpl) GetSettings(ctx context.Context, departmentID uuid.UUID) (*entities.SettingsAggregate, error) {
	working, err := uc.repo.GetWorkingDays(ctx, departmentID)
	if err != nil {
		return nil, err
	}
	shifts, err := uc.repo.GetShifts(ctx, departmentID)
	if err != nil {
		return nil, err
	}
	holidays, err := uc.repo.GetHolidays(ctx, departmentID)
	if err != nil {
		return nil, err
	}
	return &entities.SettingsAggregate{WorkingDays: working, Shifts: shifts, Holidays: holidays}, nil
}

func (uc *SettingUseCaseImpl) UpdateWorkingDays(ctx context.Context, departmentID uuid.UUID, days []entities.WorkingDay) error {
	return uc.repo.UpsertWorkingDays(ctx, departmentID, days)
}

func (uc *SettingUseCaseImpl) CreateShift(ctx context.Context, shift entities.Shift) (uuid.UUID, error) {
	return uc.repo.CreateShift(ctx, shift)
}

func (uc *SettingUseCaseImpl) UpdateShift(ctx context.Context, shift entities.Shift) error {
	return uc.repo.UpdateShift(ctx, shift)
}

func (uc *SettingUseCaseImpl) ToggleShift(ctx context.Context, shiftID uuid.UUID, isActive bool) error {
	return uc.repo.UpdateShiftStatus(ctx, shiftID, isActive)
}

func (uc *SettingUseCaseImpl) DeleteShift(ctx context.Context, shiftID uuid.UUID) error {
	return uc.repo.DeleteShift(ctx, shiftID)
}

func (uc *SettingUseCaseImpl) CreateHoliday(ctx context.Context, holiday entities.Holiday) (uuid.UUID, error) {
	return uc.repo.CreateHoliday(ctx, holiday)
}

func (uc *SettingUseCaseImpl) DeleteHoliday(ctx context.Context, holidayID uuid.UUID) error {
	return uc.repo.DeleteHoliday(ctx, holidayID)
}

func (uc *SettingUseCaseImpl) UpdateHoliday(ctx context.Context, holiday entities.Holiday) error {
	return uc.repo.UpdateHoliday(ctx, holiday)
}

package handlers

import (
	"context"
	"time"

	ent "nurseshift/setting-service/internal/domain/entities"
	usecase "nurseshift/setting-service/internal/domain/usecases"

	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// SettingHandler handles setting-related HTTP requests
type SettingHandler struct{ uc usecase.SettingUseCase }

// NewSettingHandler creates a new setting handler
func NewSettingHandler(uc usecase.SettingUseCase) *SettingHandler { return &SettingHandler{uc: uc} }

// GetSettings returns department settings
func (h *SettingHandler) GetSettings(c *fiber.Ctx) error {
	departmentIDStr := c.Query("departmentId")
	if departmentIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "ต้องระบุ departmentId"})
	}
	departmentID, err := uuid.Parse(departmentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "departmentId ไม่ถูกต้อง"})
	}

	settings, err := h.uc.GetSettings(context.Background(), departmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "ดึงการตั้งค่าสำเร็จ", "data": settings})
}

// UpdateSettings updates department working days and optionally shifts/holidays
func (h *SettingHandler) UpdateSettings(c *fiber.Ctx) error {
	departmentIDStr := c.Query("departmentId")
	if departmentIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "ต้องระบุ departmentId"})
	}
	departmentID, err := uuid.Parse(departmentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "departmentId ไม่ถูกต้อง"})
	}

	var req struct {
		WorkingDays []struct {
			ID      string `json:"id"`
			DayID   string `json:"idOrName"`
			Enabled bool   `json:"enabled"`
		} `json:"workingDays"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "ข้อมูลไม่ถูกต้อง"})
	}

	// Map weekdays (frontend uses monday..sunday). 0=Sunday..6=Saturday per schema
	weekdayMap := map[string]int{"sunday": 0, "monday": 1, "tuesday": 2, "wednesday": 3, "thursday": 4, "friday": 5, "saturday": 6}
	var days []ent.WorkingDay
	for _, d := range req.WorkingDays {
		dayKey := d.DayID
		if _, ok := weekdayMap[dayKey]; !ok {
			continue
		}
		days = append(days, ent.WorkingDay{DayOfWeek: weekdayMap[dayKey], IsWorkingDay: d.Enabled})
	}
	if err := h.uc.UpdateWorkingDays(context.Background(), departmentID, days); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "อัปเดตการตั้งค่าสำเร็จ"})
}

// Health returns service health status
func (h *SettingHandler) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok", "service": "setting-service", "timestamp": time.Now()})
}

// CreateShift creates a new shift
func (h *SettingHandler) CreateShift(c *fiber.Ctx) error {
	var req struct {
		DepartmentID   string `json:"departmentId"`
		Name           string `json:"name"`
		Type           string `json:"type"`
		StartTime      string `json:"startTime"`
		EndTime        string `json:"endTime"`
		NurseCount     int    `json:"nurseCount"`
		AssistantCount int    `json:"assistantCount"`
		Color          string `json:"color"`
		IsActive       bool   `json:"isActive"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "ข้อมูลไม่ถูกต้อง"})
	}
	deptID, err := uuid.Parse(req.DepartmentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "departmentId ไม่ถูกต้อง"})
	}
	// Accept free-text type, normalize trimming only
	st := strings.TrimSpace(req.Type)

	id, err := h.uc.CreateShift(context.Background(), ent.Shift{
		DepartmentID:       deptID,
		Name:               req.Name,
		Type:               st,
		StartTime:          req.StartTime,
		EndTime:            req.EndTime,
		DurationHours:      0,
		RequiredNurses:     req.NurseCount,
		RequiredAssistants: req.AssistantCount,
		Color:              req.Color,
		IsActive:           true,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "id": id})
}

// ToggleShift toggles shift active status
func (h *SettingHandler) ToggleShift(c *fiber.Ctx) error {
	shiftID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "shift id ไม่ถูกต้อง"})
	}
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "ข้อมูลไม่ถูกต้อง"})
	}
	if err := h.uc.ToggleShift(context.Background(), shiftID, req.Enabled); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

// DeleteShift deletes a shift
func (h *SettingHandler) DeleteShift(c *fiber.Ctx) error {
	shiftID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "shift id ไม่ถูกต้อง"})
	}
	if err := h.uc.DeleteShift(context.Background(), shiftID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

// UpdateShift updates an existing shift
func (h *SettingHandler) UpdateShift(c *fiber.Ctx) error {
	shiftID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "shift id ไม่ถูกต้อง"})
	}
	var req struct {
		Name           string `json:"name"`
		Type           string `json:"type"`
		StartTime      string `json:"startTime"`
		EndTime        string `json:"endTime"`
		NurseCount     int    `json:"nurseCount"`
		AssistantCount int    `json:"assistantCount"`
		Color          string `json:"color"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "ข้อมูลไม่ถูกต้อง"})
	}
	if err := h.uc.UpdateShift(context.Background(), ent.Shift{ID: shiftID, Name: req.Name, Type: strings.TrimSpace(req.Type), StartTime: req.StartTime, EndTime: req.EndTime, RequiredNurses: req.NurseCount, RequiredAssistants: req.AssistantCount, Color: req.Color}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

// CreateHoliday creates a new holiday
func (h *SettingHandler) CreateHoliday(c *fiber.Ctx) error {
	var req struct {
		DepartmentID string `json:"departmentId"`
		Name         string `json:"name"`
		StartDate    string `json:"startDate"`
		EndDate      string `json:"endDate"`
		IsRecurring  bool   `json:"isRecurring"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "ข้อมูลไม่ถูกต้อง"})
	}
	deptID, err := uuid.Parse(req.DepartmentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "departmentId ไม่ถูกต้อง"})
	}
	start, err1 := time.Parse("2006-01-02", req.StartDate)
	end, err2 := time.Parse("2006-01-02", req.EndDate)
	if err1 != nil || err2 != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "รูปแบบวันที่ต้องเป็น YYYY-MM-DD"})
	}
	id, err := h.uc.CreateHoliday(context.Background(), ent.Holiday{DepartmentID: deptID, Name: req.Name, StartDate: start, EndDate: end, IsRecurring: req.IsRecurring})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "id": id})
}

// DeleteHoliday deletes a holiday
func (h *SettingHandler) DeleteHoliday(c *fiber.Ctx) error {
	holidayID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "holiday id ไม่ถูกต้อง"})
	}
	if err := h.uc.DeleteHoliday(context.Background(), holidayID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

// UpdateHoliday updates an existing holiday
func (h *SettingHandler) UpdateHoliday(c *fiber.Ctx) error {
	holidayID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "holiday id ไม่ถูกต้อง"})
	}
	var req struct {
		Name        string `json:"name"`
		StartDate   string `json:"startDate"`
		EndDate     string `json:"endDate"`
		IsRecurring bool   `json:"isRecurring"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "ข้อมูลไม่ถูกต้อง"})
	}
	start, err1 := time.Parse("2006-01-02", req.StartDate)
	end, err2 := time.Parse("2006-01-02", req.EndDate)
	if err1 != nil || err2 != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "รูปแบบวันที่ต้องเป็น YYYY-MM-DD"})
	}

	// Proper update without recreating
	if err := h.uc.UpdateHoliday(context.Background(), ent.Holiday{ID: holidayID, Name: req.Name, StartDate: start, EndDate: end, IsRecurring: req.IsRecurring}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

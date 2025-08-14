package handlers

import (
	"context"
	"encoding/json"
	"time"

	"nurseshift/priority-service/internal/infrastructure/database"

	"github.com/gofiber/fiber/v2"
)

// PriorityHandler handles priority-related HTTP requests (DB-backed)
type PriorityHandler struct {
	repo *database.PriorityRepository
}

// NewPriorityHandler creates a new priority handler
func NewPriorityHandler(repo *database.PriorityRepository) *PriorityHandler {
	return &PriorityHandler{repo: repo}
}

// GetPriorities returns all priorities for a department
// Frontend จะส่ง departmentId มาใน query; ถ้าไม่ส่งให้บังคับ
func (h *PriorityHandler) GetPriorities(c *fiber.Ctx) error {
	departmentID := c.Query("departmentId")
	if departmentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "กรุณาระบุ departmentId",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()
	items, err := h.repo.FindOrCreateDefaults(ctx, departmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	// map to frontend shape
	result := make([]fiber.Map, 0, len(items))
	activeCount := 0
	for _, it := range items {
		if it.IsActive {
			activeCount++
		}
		// Derive UI setting metadata by name
		hasSettings := true
		settingType := ""
		settingUnit := ""
		settingLabel := ""
		var defaultVal int
		switch it.Name {
		case "วันที่ขอหยุด":
			hasSettings = false
		case "จำนวนเวรเท่ากันในแต่ละประเภท":
			settingType = "maxShiftTypeDifference"
			settingUnit = "เวร"
			settingLabel = "ความแตกต่างจำนวนเวรแต่ละประเภทสูงสุด"
			defaultVal = 2
		case "จำนวนเวรดึกติดต่อกัน":
			settingType = "maxConsecutiveNightShifts"
			settingUnit = "วัน"
			settingLabel = "จำนวนเวรดึกติดกันสูงสุด"
			defaultVal = 2
		case "จำนวนเวรติดต่อกัน":
			settingType = "maxConsecutiveShifts"
			settingUnit = "วัน"
			settingLabel = "จำนวนเวรติดต่อกันสูงสุด"
			defaultVal = 4
		case "จำนวนชั่วโมงทำงานสูงสุดติดต่อกันโดยไม่พัก":
			settingType = "maxConsecutiveWorkHours"
			settingUnit = "ชั่วโมง"
			settingLabel = "จำนวนชั่วโมงทำงานติดต่อกันสูงสุด"
			defaultVal = 48
		case "จำนวนชั่วโมงการทำงานทั้งหมด":
			// ตีความเป็นความแตกต่างชั่วโมงทำงานรวมระหว่างบุคคล
			settingType = "maxTotalWorkHoursDifference"
			settingUnit = "ชั่วโมง"
			settingLabel = "ความแตกต่างชั่วโมงการทำงานรวมสูงสุดระหว่างบุคคล"
			defaultVal = 16
		}

		// Extract setting value from config JSON if exists
		var settingValue *int
		if it.Config.Valid {
			var obj map[string]any
			if err := json.Unmarshal([]byte(it.Config.String), &obj); err == nil {
				if v, ok := obj["value"].(float64); ok {
					vv := int(v)
					settingValue = &vv
				}
			}
		}
		if settingValue == nil && hasSettings {
			v := defaultVal
			settingValue = &v
		}

		m := fiber.Map{
			"id":          it.ID,
			"name":        it.Name,
			"description": it.Description.String,
			"order":       it.PriorityOrder,
			"isActive":    it.IsActive,
			"hasSettings": hasSettings,
			"createdAt":   it.CreatedAt,
			"updatedAt":   it.UpdatedAt,
		}
		if hasSettings {
			m["settingType"] = settingType
			m["settingUnit"] = settingUnit
			m["settingLabel"] = settingLabel
			m["settingValue"] = settingValue
		}
		result = append(result, m)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลความสำคัญสำเร็จ",
		"data": fiber.Map{
			"priorities":  result,
			"total":       len(result),
			"activeCount": activeCount,
		},
	})
}

// UpdatePriority updates order or active status
func (h *PriorityHandler) UpdatePriority(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Order    *int  `json:"order"`
		IsActive *bool `json:"isActive"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "ข้อมูลที่ส่งมาไม่ถูกต้อง"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	if req.IsActive != nil {
		if err := h.repo.UpdateActive(ctx, id, *req.IsActive); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}
	}
	if req.Order != nil {
		if err := h.repo.UpdateOrderAndReorder(ctx, id, *req.Order); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}
	}

	rec, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "อัปเดตความสำคัญสำเร็จ",
		"data": fiber.Map{
			"id":          rec.ID,
			"name":        rec.Name,
			"description": rec.Description.String,
			"order":       rec.PriorityOrder,
			"isActive":    rec.IsActive,
			"hasSettings": false,
			"createdAt":   rec.CreatedAt,
			"updatedAt":   rec.UpdatedAt,
		},
	})
}

// UpdatePrioritySetting placeholder (ยังไม่เก็บค่า settings ใน schema นี้)
func (h *PriorityHandler) UpdatePrioritySetting(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		SettingValue int `json:"settingValue"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "ข้อมูลที่ส่งมาไม่ถูกต้อง"})
	}
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()
	if err := h.repo.UpdateSetting(ctx, id, req.SettingValue); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	rec, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{
		"id": rec.ID, "name": rec.Name, "description": rec.Description.String, "order": rec.PriorityOrder, "isActive": rec.IsActive,
	}})
}

// SwapPriorityOrder swaps order of two ids
func (h *PriorityHandler) SwapPriorityOrder(c *fiber.Ctx) error {
	var req struct {
		PriorityID1 string `json:"priorityId1"`
		PriorityID2 string `json:"priorityId2"`
	}
	if err := c.BodyParser(&req); err != nil || req.PriorityID1 == "" || req.PriorityID2 == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "ข้อมูลที่ส่งมาไม่ถูกต้อง"})
	}
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()
	if err := h.repo.SwapOrder(ctx, req.PriorityID1, req.PriorityID2); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "สลับลำดับความสำคัญสำเร็จ"})
}

// Health returns service health status
func (h *PriorityHandler) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "ok",
		"service":   "priority-service",
		"timestamp": time.Now(),
	})
}

package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// SettingHandler handles setting-related HTTP requests
type SettingHandler struct{}

// NewSettingHandler creates a new setting handler
func NewSettingHandler() *SettingHandler {
	return &SettingHandler{}
}

// GetSettings returns department settings
func (h *SettingHandler) GetSettings(c *fiber.Ctx) error {
	// Mock settings data
	settings := fiber.Map{
		"departmentSettings": fiber.Map{
			"workingDays": []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"},
			"shifts": []fiber.Map{
				{
					"id":        uuid.New().String(),
					"name":      "เวรเช้า",
					"startTime": "07:00",
					"endTime":   "15:00",
					"color":     "bg-blue-100",
					"enabled":   true,
				},
				{
					"id":        uuid.New().String(),
					"name":      "เวรบ่าย",
					"startTime": "15:00",
					"endTime":   "23:00",
					"color":     "bg-green-100",
					"enabled":   true,
				},
				{
					"id":        uuid.New().String(),
					"name":      "เวรดึก",
					"startTime": "23:00",
					"endTime":   "07:00",
					"color":     "bg-purple-100",
					"enabled":   true,
				},
			},
			"holidays": []fiber.Map{
				{
					"id":        uuid.New().String(),
					"name":      "วันปีใหม่",
					"startDate": "2025-01-01",
					"endDate":   "2025-01-01",
					"enabled":   true,
				},
				{
					"id":        uuid.New().String(),
					"name":      "วันสงกรานต์",
					"startDate": "2025-04-13",
					"endDate":   "2025-04-15",
					"enabled":   true,
				},
			},
		},
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงการตั้งค่าสำเร็จ",
		"data":    settings,
	})
}

// UpdateSettings updates department settings
func (h *SettingHandler) UpdateSettings(c *fiber.Ctx) error {
	var req fiber.Map
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลไม่ถูกต้อง",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "อัปเดตการตั้งค่าสำเร็จ",
		"data":    req,
	})
}

// Health returns service health status
func (h *SettingHandler) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "ok",
		"service":   "setting-service",
		"timestamp": time.Now(),
	})
}

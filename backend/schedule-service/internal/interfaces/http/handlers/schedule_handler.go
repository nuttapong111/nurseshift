package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ScheduleHandler handles schedule-related HTTP requests
type ScheduleHandler struct{}

// NewScheduleHandler creates a new schedule handler
func NewScheduleHandler() *ScheduleHandler {
	return &ScheduleHandler{}
}

// GetSchedules returns schedules for authenticated user's departments
func (h *ScheduleHandler) GetSchedules(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	departmentId := c.Query("departmentId")
	month := c.Query("month")

	// Mock schedule data filtered by user's departments
	var schedules []fiber.Map

	if userID == "user-1" {
		schedules = []fiber.Map{
			{
				"id":             uuid.New().String(),
				"userId":         userID,
				"departmentId":   "dept-1",
				"departmentName": "แผนกผู้ป่วยใน",
				"date":           "2024-03-15",
				"shifts": []fiber.Map{
					{
						"id":                 uuid.New().String(),
						"name":               "เวรเช้า",
						"startTime":          "07:00",
						"endTime":            "15:00",
						"nurses":             []string{"NUR001", "NUR002"},
						"assistants":         []string{"AST001"},
						"requiredNurses":     2,
						"requiredAssistants": 1,
					},
					{
						"id":                 uuid.New().String(),
						"name":               "เวรบ่าย",
						"startTime":          "15:00",
						"endTime":            "23:00",
						"nurses":             []string{"NUR003"},
						"assistants":         []string{"AST002"},
						"requiredNurses":     1,
						"requiredAssistants": 2,
					},
				},
			},
		}

		// Filter by department if specified
		if departmentId != "" {
			var filtered []fiber.Map
			for _, schedule := range schedules {
				if schedule["departmentId"] == departmentId {
					filtered = append(filtered, schedule)
				}
			}
			schedules = filtered
		}

		// Filter by month if specified
		if month != "" {
			var filtered []fiber.Map
			for _, schedule := range schedules {
				if scheduleDate, ok := schedule["date"].(string); ok {
					if scheduleDate[:7] == month { // YYYY-MM format
						filtered = append(filtered, schedule)
					}
				}
			}
			schedules = filtered
		}
	} else {
		// Return empty array for other users
		schedules = []fiber.Map{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลตารางเวรสำเร็จ",
		"data":    schedules,
	})
}

// CreateSchedule creates a new schedule
func (h *ScheduleHandler) CreateSchedule(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req struct {
		DepartmentID string `json:"departmentId" validate:"required"`
		Date         string `json:"date" validate:"required"`
		Shifts       []struct {
			Name               string   `json:"name"`
			StartTime          string   `json:"startTime"`
			EndTime            string   `json:"endTime"`
			RequiredNurses     int      `json:"requiredNurses"`
			RequiredAssistants int      `json:"requiredAssistants"`
			Nurses             []string `json:"nurses"`
			Assistants         []string `json:"assistants"`
		} `json:"shifts"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Mock response
	schedule := fiber.Map{
		"id":           uuid.New().String(),
		"userId":       userID,
		"departmentId": req.DepartmentID,
		"date":         req.Date,
		"shifts":       req.Shifts,
		"createdAt":    time.Now(),
		"status":       "created",
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "สร้างตารางเวรสำเร็จ",
		"data":    schedule,
	})
}

// GetScheduleStats returns schedule statistics for user's departments
func (h *ScheduleHandler) GetScheduleStats(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	// Mock statistics data filtered by user
	var stats fiber.Map

	if userID == "user-1" {
		stats = fiber.Map{
			"totalSchedules":     15,
			"thisMonthSchedules": 8,
			"totalNurses":        24,
			"totalAssistants":    18,
			"totalShifts":        156,
			"totalDepartments":   2,
			"departmentStats": []fiber.Map{
				{
					"departmentId":   "dept-1",
					"departmentName": "แผนกผู้ป่วยใน",
					"totalSchedules": 10,
					"totalShifts":    90,
				},
				{
					"departmentId":   "dept-2",
					"departmentName": "แผนก ICU",
					"totalSchedules": 5,
					"totalShifts":    66,
				},
			},
		}
	} else {
		stats = fiber.Map{
			"totalSchedules":     0,
			"thisMonthSchedules": 0,
			"totalNurses":        0,
			"totalAssistants":    0,
			"totalShifts":        0,
			"totalDepartments":   0,
			"departmentStats":    []fiber.Map{},
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงสถิติตารางเวรสำเร็จ",
		"data":    stats,
	})
}

// GetSchedule returns specific schedule details
func (h *ScheduleHandler) GetSchedule(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	scheduleID := c.Params("id")

	// Mock validation - check if schedule belongs to user
	// In real implementation, query database with userID filter

	schedule := fiber.Map{
		"id":           scheduleID,
		"userId":       userID,
		"departmentId": "dept-1",
		"date":         "2024-03-15",
		"shifts": []fiber.Map{
			{
				"id":         uuid.New().String(),
				"name":       "เวรเช้า",
				"startTime":  "07:00",
				"endTime":    "15:00",
				"nurses":     []string{"NUR001", "NUR002"},
				"assistants": []string{"AST001"},
			},
		},
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลตารางเวรสำเร็จ",
		"data":    schedule,
	})
}

// UpdateSchedule updates schedule information
func (h *ScheduleHandler) UpdateSchedule(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	scheduleID := c.Params("id")

	var req fiber.Map
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Mock validation - check if schedule belongs to user
	// In real implementation, query database with userID filter

	updatedSchedule := fiber.Map{
		"id":        scheduleID,
		"userId":    userID,
		"updatedAt": time.Now(),
	}

	// Merge request data
	for key, value := range req {
		updatedSchedule[key] = value
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "อัปเดตตารางเวรสำเร็จ",
		"data":    updatedSchedule,
	})
}

// DeleteSchedule deletes a schedule
func (h *ScheduleHandler) DeleteSchedule(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	scheduleID := c.Params("id")

	// Mock validation - check if schedule belongs to user
	// In real implementation, query database with userID filter and soft delete

	_ = userID     // Use userID for validation
	_ = scheduleID // Use scheduleID for deletion

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ลบตารางเวรสำเร็จ",
	})
}

// Health returns service health status
func (h *ScheduleHandler) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "ok",
		"service":   "schedule-service",
		"timestamp": time.Now(),
	})
}

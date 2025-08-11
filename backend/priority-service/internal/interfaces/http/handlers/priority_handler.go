package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// PriorityHandler handles priority-related HTTP requests
type PriorityHandler struct{}

// NewPriorityHandler creates a new priority handler
func NewPriorityHandler() *PriorityHandler {
	return &PriorityHandler{}
}

// Priority represents priority data structure
type Priority struct {
	ID           string                 `json:"id"`
	UserID       string                 `json:"userId"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Order        int                    `json:"order"`
	IsActive     bool                   `json:"isActive"`
	HasSettings  bool                   `json:"hasSettings"`
	SettingType  *string                `json:"settingType"`
	SettingValue *int                   `json:"settingValue"`
	SettingUnit  *string                `json:"settingUnit"`
	SettingLabel *string                `json:"settingLabel"`
	Settings     map[string]interface{} `json:"settings"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}

// Mock data
var mockPriorities = []Priority{
	{
		ID:           "1",
		UserID:       "user-1",
		Name:         "วันที่ขอหยุด",
		Description:  "ระบบจะหลีกเลี่ยงการจัดเวรในวันที่พนักงานขอหยุด",
		Order:        1,
		IsActive:     true,
		HasSettings:  false,
		Settings:     make(map[string]interface{}),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	},
	{
		ID:           "2",
		UserID:       "user-1",
		Name:         "จำนวนเวรในแต่ละประเภทเท่ากัน",
		Description:  "กระจายจำนวนเวรแต่ละประเภท (เช้า/บ่าย/ดึก) ให้แต่ละคนได้เท่าๆ กัน",
		Order:        2,
		IsActive:     true,
		HasSettings:  true,
		SettingType:  stringPtr("maxShiftTypeDifference"),
		SettingValue: intPtr(2),
		SettingUnit:  stringPtr("เวร"),
		SettingLabel: stringPtr("ความแตกต่างจำนวนเวรแต่ละประเภทสูงสุด"),
		Settings: map[string]interface{}{
			"min":  0,
			"max":  5,
			"step": 1,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:           "3",
		UserID:       "user-1",
		Name:         "เวรดึกติดกัน",
		Description:  "จำกัดจำนวนเวรดึกที่พนักงานคนหนึ่งทำติดกันไม่เกิน X วัน",
		Order:        3,
		IsActive:     true,
		HasSettings:  true,
		SettingType:  stringPtr("maxConsecutiveNightShifts"),
		SettingValue: intPtr(2),
		SettingUnit:  stringPtr("วัน"),
		SettingLabel: stringPtr("จำนวนเวรดึกติดกันสูงสุด"),
		Settings: map[string]interface{}{
			"min":  1,
			"max":  5,
			"step": 1,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:           "4",
		UserID:       "user-1",
		Name:         "เวรติดต่อกัน",
		Description:  "จำกัดจำนวนเวรทุกประเภทที่พนักงานคนหนึ่งทำติดกันไม่เกิน X วัน",
		Order:        4,
		IsActive:     true,
		HasSettings:  true,
		SettingType:  stringPtr("maxConsecutiveShifts"),
		SettingValue: intPtr(4),
		SettingUnit:  stringPtr("วัน"),
		SettingLabel: stringPtr("จำนวนเวรติดต่อกันสูงสุด"),
		Settings: map[string]interface{}{
			"min":  1,
			"max":  10,
			"step": 1,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:           "5",
		UserID:       "user-1",
		Name:         "จำนวนชั่วโมงการทำงานติดต่อกันสูงสุด",
		Description:  "จำกัดจำนวนชั่วโมงการทำงานต่อเนื่องของพนักงานคนหนึ่งไม่เกิน X ชั่วโมง",
		Order:        5,
		IsActive:     false,
		HasSettings:  true,
		SettingType:  stringPtr("maxConsecutiveWorkHours"),
		SettingValue: intPtr(48),
		SettingUnit:  stringPtr("ชั่วโมง"),
		SettingLabel: stringPtr("จำนวนชั่วโมงทำงานติดต่อกันสูงสุด"),
		Settings: map[string]interface{}{
			"min":  12,
			"max":  72,
			"step": 6,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

// GetPriorities returns all priorities for a user
func (h *PriorityHandler) GetPriorities(c *fiber.Ctx) error {
	userID := c.Query("userId", "user-1") // Default to user-1 for demo

	var userPriorities []Priority
	for _, priority := range mockPriorities {
		if priority.UserID == userID {
			userPriorities = append(userPriorities, priority)
		}
	}

	// Sort by order
	for i := 0; i < len(userPriorities)-1; i++ {
		for j := i + 1; j < len(userPriorities); j++ {
			if userPriorities[i].Order > userPriorities[j].Order {
				userPriorities[i], userPriorities[j] = userPriorities[j], userPriorities[i]
			}
		}
	}

	// Count active priorities
	activeCount := 0
	for _, priority := range userPriorities {
		if priority.IsActive {
			activeCount++
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลความสำคัญสำเร็จ",
		"data": fiber.Map{
			"priorities":   userPriorities,
			"total":        len(userPriorities),
			"activeCount":  activeCount,
		},
	})
}

// UpdatePriority updates priority information (including order and active status)
func (h *PriorityHandler) UpdatePriority(c *fiber.Ctx) error {
	priorityID := c.Params("id")

	var req struct {
		Order    *int  `json:"order"`
		IsActive *bool `json:"isActive"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Find and update priority
	for i, priority := range mockPriorities {
		if priority.ID == priorityID {
			if req.Order != nil {
				mockPriorities[i].Order = *req.Order
			}
			if req.IsActive != nil {
				mockPriorities[i].IsActive = *req.IsActive
			}
			mockPriorities[i].UpdatedAt = time.Now()

			// If order is updated, we might need to reorder other priorities
			if req.Order != nil {
				// Simple reorder logic - in real implementation, this would be more sophisticated
				h.reorderPriorities(priority.UserID)
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":  "success",
				"message": "อัปเดตความสำคัญสำเร็จ",
				"data":    mockPriorities[i],
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status":  "error",
		"message": "ไม่พบข้อมูลความสำคัญ",
	})
}

// UpdatePrioritySetting updates priority setting value
func (h *PriorityHandler) UpdatePrioritySetting(c *fiber.Ctx) error {
	priorityID := c.Params("id")

	var req struct {
		SettingValue int `json:"settingValue" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Find and update priority setting
	for i, priority := range mockPriorities {
		if priority.ID == priorityID {
			if !priority.HasSettings {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  "error",
					"message": "ความสำคัญนี้ไม่สามารถตั้งค่าได้",
				})
			}

			// Validate setting value range
			if minVal, ok := priority.Settings["min"].(int); ok {
				if req.SettingValue < minVal {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"status":  "error",
						"message": "ค่าที่ส่งมาต่ำกว่าค่าต่ำสุดที่อนุญาต",
					})
				}
			}

			if maxVal, ok := priority.Settings["max"].(int); ok {
				if req.SettingValue > maxVal {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"status":  "error",
						"message": "ค่าที่ส่งมาสูงกว่าค่าสูงสุดที่อนุญาต",
					})
				}
			}

			mockPriorities[i].SettingValue = &req.SettingValue
			mockPriorities[i].UpdatedAt = time.Now()

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":  "success",
				"message": "อัปเดตค่าตั้งความสำคัญสำเร็จ",
				"data":    mockPriorities[i],
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status":  "error",
		"message": "ไม่พบข้อมูลความสำคัญ",
	})
}

// SwapPriorityOrder swaps the order of two priorities
func (h *PriorityHandler) SwapPriorityOrder(c *fiber.Ctx) error {
	var req struct {
		PriorityID1 string `json:"priorityId1" validate:"required"`
		PriorityID2 string `json:"priorityId2" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	var priority1, priority2 *Priority
	var index1, index2 int

	// Find both priorities
	for i, priority := range mockPriorities {
		if priority.ID == req.PriorityID1 {
			priority1 = &priority
			index1 = i
		}
		if priority.ID == req.PriorityID2 {
			priority2 = &priority
			index2 = i
		}
	}

	if priority1 == nil || priority2 == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่พบข้อมูลความสำคัญที่ต้องการสลับ",
		})
	}

	// Swap orders
	mockPriorities[index1].Order, mockPriorities[index2].Order = mockPriorities[index2].Order, mockPriorities[index1].Order
	mockPriorities[index1].UpdatedAt = time.Now()
	mockPriorities[index2].UpdatedAt = time.Now()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "สลับลำดับความสำคัญสำเร็จ",
		"data": fiber.Map{
			"priority1": mockPriorities[index1],
			"priority2": mockPriorities[index2],
		},
	})
}

// Health returns service health status
func (h *PriorityHandler) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "ok",
		"service":   "priority-service",
		"timestamp": time.Now(),
	})
}

// Helper function to reorder priorities for a user
func (h *PriorityHandler) reorderPriorities(userID string) {
	var userPriorities []int
	
	for i, priority := range mockPriorities {
		if priority.UserID == userID {
			userPriorities = append(userPriorities, i)
		}
	}

	// Sort by order
	for i := 0; i < len(userPriorities)-1; i++ {
		for j := i + 1; j < len(userPriorities); j++ {
			if mockPriorities[userPriorities[i]].Order > mockPriorities[userPriorities[j]].Order {
				userPriorities[i], userPriorities[j] = userPriorities[j], userPriorities[i]
			}
		}
	}

	// Reassign sequential orders
	for i, priorityIndex := range userPriorities {
		mockPriorities[priorityIndex].Order = i + 1
		mockPriorities[priorityIndex].UpdatedAt = time.Now()
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// EmployeeLeaveHandler handles employee leave-related HTTP requests
type EmployeeLeaveHandler struct{}

// NewEmployeeLeaveHandler creates a new employee leave handler
func NewEmployeeLeaveHandler() *EmployeeLeaveHandler {
	return &EmployeeLeaveHandler{}
}

// EmployeeLeave represents employee leave data structure
type EmployeeLeave struct {
	ID             string    `json:"id"`
	EmployeeID     string    `json:"employeeId"`
	EmployeeName   string    `json:"employeeName"`
	DepartmentID   string    `json:"departmentId"`
	DepartmentName string    `json:"departmentName"`
	Date           string    `json:"date"`
	Reason         string    `json:"reason"`
	IsActive       bool      `json:"isActive"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// Mock data
var mockEmployeeLeaves = []EmployeeLeave{
	{
		ID:             "1",
		EmployeeID:     "1",
		EmployeeName:   "สมหญิง ใจดี",
		DepartmentID:   "emergency",
		DepartmentName: "แผนกฉุกเฉิน",
		Date:           "2024-03-15",
		Reason:         "ลาป่วย",
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
	{
		ID:             "2",
		EmployeeID:     "2",
		EmployeeName:   "สมชาย ขยัน",
		DepartmentID:   "emergency",
		DepartmentName: "แผนกฉุกเฉิน",
		Date:           "2024-03-20",
		Reason:         "ลากิจส่วนตัว",
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
	{
		ID:             "3",
		EmployeeID:     "4",
		EmployeeName:   "สมใส เก่งกาจ",
		DepartmentID:   "internal",
		DepartmentName: "แผนกอายุรกรรม",
		Date:           "2024-03-25",
		Reason:         "ลาพักร้อน",
		IsActive:       false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
	{
		ID:             "4",
		EmployeeID:     "7",
		EmployeeName:   "มณี ใสใส",
		DepartmentID:   "surgery",
		DepartmentName: "แผนกศัลยกรรม",
		Date:           "2024-03-18",
		Reason:         "ลาคลอด",
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	},
}

// GetEmployeeLeaves returns all employee leaves with filters
func (h *EmployeeLeaveHandler) GetEmployeeLeaves(c *fiber.Ctx) error {
	month := c.Query("month")
	employeeId := c.Query("employeeId")
	departmentId := c.Query("departmentId")

	filteredLeaves := mockEmployeeLeaves

	// Filter by month
	if month != "" {
		var filtered []EmployeeLeave
		for _, leave := range filteredLeaves {
			if leave.Date[:7] == month { // YYYY-MM format
				filtered = append(filtered, leave)
			}
		}
		filteredLeaves = filtered
	}

	// Filter by employee
	if employeeId != "" {
		var filtered []EmployeeLeave
		for _, leave := range filteredLeaves {
			if leave.EmployeeID == employeeId {
				filtered = append(filtered, leave)
			}
		}
		filteredLeaves = filtered
	}

	// Filter by department
	if departmentId != "" {
		var filtered []EmployeeLeave
		for _, leave := range filteredLeaves {
			if leave.DepartmentID == departmentId {
				filtered = append(filtered, leave)
			}
		}
		filteredLeaves = filtered
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลวันหยุดพนักงานสำเร็จ",
		"data":    filteredLeaves,
	})
}

// CreateEmployeeLeave creates a new employee leave
func (h *EmployeeLeaveHandler) CreateEmployeeLeave(c *fiber.Ctx) error {
	var req struct {
		EmployeeID     string `json:"employeeId" validate:"required"`
		EmployeeName   string `json:"employeeName" validate:"required"`
		DepartmentID   string `json:"departmentId" validate:"required"`
		DepartmentName string `json:"departmentName" validate:"required"`
		Date           string `json:"date" validate:"required"`
		Reason         string `json:"reason"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Create new employee leave
	newLeave := EmployeeLeave{
		ID:             uuid.New().String(),
		EmployeeID:     req.EmployeeID,
		EmployeeName:   req.EmployeeName,
		DepartmentID:   req.DepartmentID,
		DepartmentName: req.DepartmentName,
		Date:           req.Date,
		Reason:         req.Reason,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Add to mock data
	mockEmployeeLeaves = append(mockEmployeeLeaves, newLeave)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "สร้างวันหยุดพนักงานสำเร็จ",
		"data":    newLeave,
	})
}

// GetLeavesByDepartment returns leaves for specific department
func (h *EmployeeLeaveHandler) GetLeavesByDepartment(c *fiber.Ctx) error {
	departmentId := c.Params("departmentId")
	
	var filtered []EmployeeLeave
	for _, leave := range mockEmployeeLeaves {
		if leave.DepartmentID == departmentId {
			filtered = append(filtered, leave)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลวันหยุดพนักงานตามแผนกสำเร็จ",
		"data":    filtered,
	})
}

// GetLeavesByEmployee returns leaves for specific employee
func (h *EmployeeLeaveHandler) GetLeavesByEmployee(c *fiber.Ctx) error {
	employeeId := c.Params("employeeId")
	
	var filtered []EmployeeLeave
	for _, leave := range mockEmployeeLeaves {
		if leave.EmployeeID == employeeId {
			filtered = append(filtered, leave)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลวันหยุดของพนักงานสำเร็จ",
		"data":    filtered,
	})
}

// UpdateEmployeeLeave updates employee leave information
func (h *EmployeeLeaveHandler) UpdateEmployeeLeave(c *fiber.Ctx) error {
	leaveId := c.Params("id")
	
	var req struct {
		Date   *string `json:"date"`
		Reason *string `json:"reason"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Find and update leave
	for i, leave := range mockEmployeeLeaves {
		if leave.ID == leaveId {
			if req.Date != nil {
				mockEmployeeLeaves[i].Date = *req.Date
			}
			if req.Reason != nil {
				mockEmployeeLeaves[i].Reason = *req.Reason
			}
			mockEmployeeLeaves[i].UpdatedAt = time.Now()
			
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":  "success",
				"message": "อัปเดตวันหยุดพนักงานสำเร็จ",
				"data":    mockEmployeeLeaves[i],
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status":  "error",
		"message": "ไม่พบข้อมูลวันหยุดพนักงาน",
	})
}

// DeleteEmployeeLeave deletes employee leave
func (h *EmployeeLeaveHandler) DeleteEmployeeLeave(c *fiber.Ctx) error {
	leaveId := c.Params("id")
	
	// Find and remove leave
	for i, leave := range mockEmployeeLeaves {
		if leave.ID == leaveId {
			mockEmployeeLeaves = append(mockEmployeeLeaves[:i], mockEmployeeLeaves[i+1:]...)
			
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":  "success",
				"message": "ลบวันหยุดพนักงานสำเร็จ",
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status":  "error",
		"message": "ไม่พบข้อมูลวันหยุดพนักงาน",
	})
}

// ToggleEmployeeLeave toggles active status of employee leave
func (h *EmployeeLeaveHandler) ToggleEmployeeLeave(c *fiber.Ctx) error {
	leaveId := c.Params("id")
	
	// Find and toggle leave status
	for i, leave := range mockEmployeeLeaves {
		if leave.ID == leaveId {
			mockEmployeeLeaves[i].IsActive = !mockEmployeeLeaves[i].IsActive
			mockEmployeeLeaves[i].UpdatedAt = time.Now()
			
			status := "เปิดใช้งาน"
			if !mockEmployeeLeaves[i].IsActive {
				status = "ปิดใช้งาน"
			}
			
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":  "success",
				"message": status + "วันหยุดพนักงานสำเร็จ",
				"data":    mockEmployeeLeaves[i],
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status":  "error",
		"message": "ไม่พบข้อมูลวันหยุดพนักงาน",
	})
}

// Health returns service health status
func (h *EmployeeLeaveHandler) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "ok",
		"service":   "employee-leave-service",
		"timestamp": time.Now(),
	})
}

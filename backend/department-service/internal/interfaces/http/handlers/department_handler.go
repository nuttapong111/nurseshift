package handlers

import (
	"strconv"
	"time"

	"nurseshift/department-service/internal/domain/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// DepartmentHandler handles department-related HTTP requests
type DepartmentHandler struct {
	// departmentUseCase usecases.DepartmentUseCase
}

// NewDepartmentHandler creates a new department handler
func NewDepartmentHandler() *DepartmentHandler {
	return &DepartmentHandler{}
}

// GetDepartments returns list of departments for the authenticated user
func (h *DepartmentHandler) GetDepartments(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	// Mock data filtered by userID - only show departments created by this user
	var departments []*entities.DepartmentWithStats

	// Sample departments for different users
	if userID == "user-1" {
		departments = []*entities.DepartmentWithStats{
			{
				Department: &entities.Department{
					ID:             uuid.New(),
					OrganizationID: uuid.MustParse(userID), // Use userID as organizationID for demo
					Name:           "แผนกผู้ป่วยใน",
					Description:    stringPtr("แผนกดูแลผู้ป่วยที่ต้องพักรักษาตัวในโรงพยาบาล"),
					MaxNurses:      15,
					MaxAssistants:  8,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				},
				TotalEmployees:  18,
				ActiveEmployees: 16,
				NurseCount:      12,
				AssistantCount:  6,
			},
			{
				Department: &entities.Department{
					ID:             uuid.New(),
					OrganizationID: uuid.MustParse(userID),
					Name:           "แผนก ICU",
					Description:    stringPtr("แผนกผู้ป่วยหนักที่ต้องการการดูแลเข้มข้น"),
					MaxNurses:      12,
					MaxAssistants:  6,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				},
				TotalEmployees:  15,
				ActiveEmployees: 14,
				NurseCount:      10,
				AssistantCount:  5,
			},
		}
	} else {
		// Return empty array for other users or different departments
		departments = []*entities.DepartmentWithStats{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลแผนกสำเร็จ",
		"data":    departments,
	})
}

// GetDepartment returns specific department details
func (h *DepartmentHandler) GetDepartment(c *fiber.Ctx) error {
	departmentIDStr := c.Params("id")
	departmentID, err := uuid.Parse(departmentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "รหัสแผนกไม่ถูกต้อง",
		})
	}

	// Mock data for department detail
	department := &entities.DepartmentDetail{
		Department: &entities.Department{
			ID:             departmentID,
			OrganizationID: c.Locals("organizationID").(uuid.UUID),
			Name:           "แผนกผู้ป่วยใน",
			Description:    stringPtr("แผนกดูแลผู้ป่วยที่ต้องพักรักษาตัวในโรงพยาบาล"),
			MaxNurses:      15,
			MaxAssistants:  8,
			CreatedAt:      time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:      time.Now(),
		},
		Head: &entities.User{
			ID:         uuid.New(),
			EmployeeID: "MGR001",
			FirstName:  "ดร.สมชาย",
			LastName:   "ใจดี",
			Email:      "manager1@thephyathai.com",
			Role:       "manager",
			Status:     "active",
			Position:   stringPtr("หัวหน้าแผนกผู้ป่วยใน"),
		},
		Employees: []*entities.EmployeeWithUser{
			{
				Employee: &entities.Employee{
					ID:           uuid.New(),
					DepartmentID: departmentID,
					UserID:       uuid.New(),
					Position:     "พยาบาลวิชาชีพ",
					StartDate:    time.Now().Add(-180 * 24 * time.Hour),
					IsActive:     true,
					CreatedAt:    time.Now().Add(-180 * 24 * time.Hour),
					UpdatedAt:    time.Now(),
				},
				User: &entities.User{
					ID:         uuid.New(),
					EmployeeID: "NUR001",
					FirstName:  "สุภา",
					LastName:   "จิตรดี",
					Email:      "nurse1@thephyathai.com",
					Role:       "nurse",
					Status:     "active",
					Position:   stringPtr("พยาบาลวิชาชีพ"),
				},
			},
			{
				Employee: &entities.Employee{
					ID:           uuid.New(),
					DepartmentID: departmentID,
					UserID:       uuid.New(),
					Position:     "ผู้ช่วยพยาบาล",
					StartDate:    time.Now().Add(-120 * 24 * time.Hour),
					IsActive:     true,
					CreatedAt:    time.Now().Add(-120 * 24 * time.Hour),
					UpdatedAt:    time.Now(),
				},
				User: &entities.User{
					ID:         uuid.New(),
					EmployeeID: "AST001",
					FirstName:  "นิรันดร์",
					LastName:   "ช่วยดี",
					Email:      "assistant1@thephyathai.com",
					Role:       "assistant",
					Status:     "active",
					Position:   stringPtr("ผู้ช่วยพยาบาล"),
				},
			},
		},
		Stats: &entities.DepartmentStats{
			TotalEmployees:  18,
			ActiveEmployees: 16,
			NurseCount:      12,
			AssistantCount:  6,
			UtilizationRate: 88.9,
		},
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลแผนกสำเร็จ",
		"data":    department,
	})
}

// CreateDepartment creates a new department
func (h *DepartmentHandler) CreateDepartment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req struct {
		Name          string     `json:"name" validate:"required"`
		Description   *string    `json:"description"`
		HeadUserID    *uuid.UUID `json:"headUserId"`
		MaxNurses     int        `json:"maxNurses" validate:"required,min=1"`
		MaxAssistants int        `json:"maxAssistants" validate:"required,min=1"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Mock response - use userID as organizationID for demo
	userUUID, _ := uuid.Parse(userID)
	department := &entities.Department{
		ID:             uuid.New(),
		OrganizationID: userUUID,
		Name:           req.Name,
		Description:    req.Description,
		HeadUserID:     req.HeadUserID,
		MaxNurses:      req.MaxNurses,
		MaxAssistants:  req.MaxAssistants,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "สร้างแผนกสำเร็จ",
		"data":    department,
	})
}

// UpdateDepartment updates department information
func (h *DepartmentHandler) UpdateDepartment(c *fiber.Ctx) error {
	departmentIDStr := c.Params("id")
	departmentID, err := uuid.Parse(departmentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "รหัสแผนกไม่ถูกต้อง",
		})
	}

	var req struct {
		Name          *string    `json:"name"`
		Description   *string    `json:"description"`
		HeadUserID    *uuid.UUID `json:"headUserId"`
		MaxNurses     *int       `json:"maxNurses"`
		MaxAssistants *int       `json:"maxAssistants"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Mock response - in real implementation, update in database
	department := &entities.Department{
		ID:             departmentID,
		OrganizationID: c.Locals("organizationID").(uuid.UUID),
		Name:           getStringValue(req.Name, "แผนกที่อัปเดต"),
		Description:    req.Description,
		HeadUserID:     req.HeadUserID,
		MaxNurses:      getIntValue(req.MaxNurses, 15),
		MaxAssistants:  getIntValue(req.MaxAssistants, 8),
		CreatedAt:      time.Now().Add(-30 * 24 * time.Hour),
		UpdatedAt:      time.Now(),
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "อัปเดตแผนกสำเร็จ",
		"data":    department,
	})
}

// DeleteDepartment deletes a department
func (h *DepartmentHandler) DeleteDepartment(c *fiber.Ctx) error {
	departmentIDStr := c.Params("id")
	_, err := uuid.Parse(departmentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "รหัสแผนกไม่ถูกต้อง",
		})
	}

	// Mock deletion - in real implementation, soft delete in database
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ลบแผนกสำเร็จ",
	})
}

// GetDepartmentEmployees returns employees in a department
func (h *DepartmentHandler) GetDepartmentEmployees(c *fiber.Ctx) error {
	departmentIDStr := c.Params("id")
	departmentID, err := uuid.Parse(departmentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "รหัสแผนกไม่ถูกต้อง",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	// Mock employees data
	employees := []*entities.EmployeeWithUser{
		{
			Employee: &entities.Employee{
				ID:           uuid.New(),
				DepartmentID: departmentID,
				UserID:       uuid.New(),
				Position:     "พยาบาลวิชาชีพ",
				StartDate:    time.Now().Add(-180 * 24 * time.Hour),
				IsActive:     true,
				CreatedAt:    time.Now().Add(-180 * 24 * time.Hour),
				UpdatedAt:    time.Now(),
			},
			User: &entities.User{
				ID:         uuid.New(),
				EmployeeID: "NUR001",
				FirstName:  "สุภา",
				LastName:   "จิตรดี",
				Email:      "nurse1@thephyathai.com",
				Phone:      stringPtr("081-456-7890"),
				Role:       "nurse",
				Status:     "active",
				Position:   stringPtr("พยาบาลวิชาชีพ"),
			},
		},
		{
			Employee: &entities.Employee{
				ID:           uuid.New(),
				DepartmentID: departmentID,
				UserID:       uuid.New(),
				Position:     "ผู้ช่วยพยาบาล",
				StartDate:    time.Now().Add(-120 * 24 * time.Hour),
				IsActive:     true,
				CreatedAt:    time.Now().Add(-120 * 24 * time.Hour),
				UpdatedAt:    time.Now(),
			},
			User: &entities.User{
				ID:         uuid.New(),
				EmployeeID: "AST001",
				FirstName:  "นิรันดร์",
				LastName:   "ช่วยดี",
				Email:      "assistant1@thephyathai.com",
				Phone:      stringPtr("081-789-0123"),
				Role:       "assistant",
				Status:     "active",
				Position:   stringPtr("ผู้ช่วยพยาบาล"),
			},
		},
	}

	totalPages := (len(employees) + limit - 1) / limit

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลพนักงานสำเร็จ",
		"data": fiber.Map{
			"employees":  employees,
			"total":      len(employees),
			"page":       page,
			"limit":      limit,
			"totalPages": totalPages,
		},
	})
}

// GetDepartmentStats returns department statistics
func (h *DepartmentHandler) GetDepartmentStats(c *fiber.Ctx) error {
	organizationID := c.Locals("organizationID").(uuid.UUID)
	_ = organizationID

	stats := fiber.Map{
		"totalDepartments": 3,
		"totalEmployees":   55,
		"departmentStats": []fiber.Map{
			{
				"departmentName":  "แผนกผู้ป่วยใน",
				"totalEmployees":  18,
				"activeEmployees": 16,
				"nurseCount":      12,
				"assistantCount":  6,
				"utilizationRate": 88.9,
			},
			{
				"departmentName":  "แผนก ICU",
				"totalEmployees":  15,
				"activeEmployees": 14,
				"nurseCount":      10,
				"assistantCount":  5,
				"utilizationRate": 93.3,
			},
			{
				"departmentName":  "แผนกฉุกเฉิน",
				"totalEmployees":  22,
				"activeEmployees": 20,
				"nurseCount":      15,
				"assistantCount":  7,
				"utilizationRate": 90.9,
			},
		},
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงสถิติแผนกสำเร็จ",
		"data":    stats,
	})
}

// Health returns service health status
func (h *DepartmentHandler) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "ok",
		"service":   "department-service",
		"timestamp": time.Now(),
	})
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func getStringValue(ptr *string, defaultValue string) string {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}

func getIntValue(ptr *int, defaultValue int) int {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}

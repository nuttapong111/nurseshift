package handlers

import (
	"strconv"
	"time"

	"nurseshift/department-service/internal/domain/entities"
	"nurseshift/department-service/internal/infrastructure/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// DepartmentHandler handles department-related HTTP requests
type DepartmentHandler struct {
	// departmentUseCase usecases.DepartmentUseCase
	departmentRepo database.DepartmentRepository
}

// NewDepartmentHandler creates a new department handler
func NewDepartmentHandler(repo database.DepartmentRepository) *DepartmentHandler {
	return &DepartmentHandler{
		departmentRepo: repo,
	}
}

// GetDepartments returns list of departments for the authenticated user
func (h *DepartmentHandler) GetDepartments(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "รหัสผู้ใช้ไม่ถูกต้อง",
		})
	}

	// Get departments from database
	departments, err := h.departmentRepo.GetWithStats(c.Context(), userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถดึงข้อมูลแผนกได้",
			"error":   err.Error(),
		})
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

	// Get department from database
	department, err := h.departmentRepo.GetByID(c.Context(), departmentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่พบแผนกที่ต้องการ",
		})
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

	var req entities.CreateDepartmentRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "รหัสผู้ใช้ไม่ถูกต้อง",
		})
	}

	// Create department entity
	department := &entities.Department{
		ID:            uuid.New(),
		Name:          req.Name,
		Description:   req.Description,
		HeadUserID:    req.HeadUserID,
		MaxNurses:     req.MaxNurses,
		MaxAssistants: req.MaxAssistants,
		Settings:      req.Settings,
		IsActive:      true,
		CreatedBy:     &userUUID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save to database
	if err := h.departmentRepo.Create(c.Context(), department); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถสร้างแผนกได้",
			"error":   err.Error(),
		})
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

	var req entities.UpdateDepartmentRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Get existing department
	existingDept, err := h.departmentRepo.GetByID(c.Context(), departmentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่พบแผนกที่ต้องการอัปเดต",
		})
	}

	// Update fields
	if req.Name != nil {
		existingDept.Name = *req.Name
	}
	if req.Description != nil {
		existingDept.Description = req.Description
	}
	if req.HeadUserID != nil {
		existingDept.HeadUserID = req.HeadUserID
	}
	if req.MaxNurses != nil {
		existingDept.MaxNurses = *req.MaxNurses
	}
	if req.MaxAssistants != nil {
		existingDept.MaxAssistants = *req.MaxAssistants
	}
	if req.Settings != nil {
		existingDept.Settings = req.Settings
	}
	if req.IsActive != nil {
		existingDept.IsActive = *req.IsActive
	}
	existingDept.UpdatedAt = time.Now()

	// Save to database
	if err := h.departmentRepo.Update(c.Context(), existingDept); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถอัปเดตแผนกได้",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "อัปเดตแผนกสำเร็จ",
		"data":    existingDept,
	})
}

// DeleteDepartment deletes a department
func (h *DepartmentHandler) DeleteDepartment(c *fiber.Ctx) error {
	departmentIDStr := c.Params("id")
	departmentID, err := uuid.Parse(departmentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "รหัสแผนกไม่ถูกต้อง",
		})
	}

	// Soft delete from database
	if err := h.departmentRepo.Delete(c.Context(), departmentID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถลบแผนกได้",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ลบแผนกสำเร็จ",
	})
}

// GetDepartmentStaff returns staff in a department
func (h *DepartmentHandler) GetDepartmentStaff(c *fiber.Ctx) error {
	departmentIDStr := c.Params("id")
	departmentID, err := uuid.Parse(departmentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "รหัสแผนกไม่ถูกต้อง",
		})
	}

	// Get staff from database
	staff, err := h.departmentRepo.GetStaff(c.Context(), departmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถดึงข้อมูลพนักงานได้",
			"error":   err.Error(),
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	totalPages := (len(staff) + limit - 1) / limit

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลพนักงานสำเร็จ",
		"data": fiber.Map{
			"staff":      staff,
			"total":      len(staff),
			"page":       page,
			"limit":      limit,
			"totalPages": totalPages,
		},
	})
}

// AddDepartmentStaff adds a new staff member to a department
func (h *DepartmentHandler) AddDepartmentStaff(c *fiber.Ctx) error {
	departmentIDStr := c.Params("id")
	departmentID, err := uuid.Parse(departmentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "รหัสแผนกไม่ถูกต้อง",
		})
	}

	var req struct {
		FirstName      string `json:"first_name" validate:"required"`
		LastName       string `json:"last_name" validate:"required"`
		Position       string `json:"position" validate:"required,oneof=nurse assistant"`
		Phone          string `json:"phone,omitempty"`
		Email          string `json:"email,omitempty"`
		DepartmentRole string `json:"department_role" validate:"required,oneof=nurse assistant"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Validate position matches department_role
	if req.Position != req.DepartmentRole {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ตำแหน่งและ role ในแผนกต้องตรงกัน",
		})
	}

	// Create staff member
	staff := &entities.DepartmentStaff{
		ID:           uuid.New(),
		DepartmentID: departmentID,
		Name:         req.FirstName + " " + req.LastName,
		Position:     req.Position,
		Phone:        &req.Phone,
		Email:        &req.Email,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save to database
	if err := h.departmentRepo.CreateStaff(c.Context(), staff); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถเพิ่มพนักงานได้",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "เพิ่มพนักงานสำเร็จ",
		"data":    staff,
	})
}

// DeleteDepartmentStaff deletes a staff member from a department
func (h *DepartmentHandler) DeleteDepartmentStaff(c *fiber.Ctx) error {
	departmentIDStr := c.Params("id")
	staffIDStr := c.Params("staffId")

	departmentID, err := uuid.Parse(departmentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "รหัสแผนกไม่ถูกต้อง",
		})
	}

	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "รหัสพนักงานไม่ถูกต้อง",
		})
	}

	// Delete staff from database
	if err := h.departmentRepo.DeleteStaff(c.Context(), staffID, departmentID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถลบพนักงานได้",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ลบพนักงานสำเร็จ",
	})
}

// GetDepartmentStats returns department statistics
func (h *DepartmentHandler) GetDepartmentStats(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "รหัสผู้ใช้ไม่ถูกต้อง",
		})
	}

	// Get departments with stats from database
	departments, err := h.departmentRepo.GetWithStats(c.Context(), userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถดึงสถิติแผนกได้",
			"error":   err.Error(),
		})
	}

	// Calculate totals
	totalDepartments := len(departments)
	totalEmployees := 0
	var departmentStats []fiber.Map

	for _, dept := range departments {
		totalEmployees += dept.TotalEmployees

		utilizationRate := 0.0
		if dept.Department.MaxNurses > 0 || dept.Department.MaxAssistants > 0 {
			maxTotal := dept.Department.MaxNurses + dept.Department.MaxAssistants
			if maxTotal > 0 {
				utilizationRate = float64(dept.TotalEmployees) / float64(maxTotal) * 100
			}
		}

		departmentStats = append(departmentStats, fiber.Map{
			"departmentName":  dept.Department.Name,
			"totalEmployees":  dept.TotalEmployees,
			"activeEmployees": dept.ActiveEmployees,
			"nurseCount":      dept.NurseCount,
			"assistantCount":  dept.AssistantCount,
			"utilizationRate": utilizationRate,
		})
	}

	stats := fiber.Map{
		"totalDepartments": totalDepartments,
		"totalEmployees":   totalEmployees,
		"departmentStats":  departmentStats,
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

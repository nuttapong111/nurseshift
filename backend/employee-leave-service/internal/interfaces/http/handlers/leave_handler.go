package handlers

import (
	"context"
	"time"

	"nurseshift/employee-leave-service/internal/domain/entities"
	"nurseshift/employee-leave-service/internal/domain/usecases"
	"nurseshift/employee-leave-service/internal/infrastructure/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// LeaveHandler handles leave request HTTP requests
type LeaveHandler struct {
	leaveUseCase usecases.LeaveUseCase
	db           *database.Connection
}

// NewLeaveHandler creates a new leave handler
func NewLeaveHandler(leaveUseCase usecases.LeaveUseCase, db *database.Connection) *LeaveHandler {
	return &LeaveHandler{
		leaveUseCase: leaveUseCase,
		db:           db,
	}
}

// GetLeaves returns all leave requests with filters
func (h *LeaveHandler) GetLeaves(c *fiber.Ctx) error {
	// Get query parameters
	month := c.Query("month")
	employeeId := c.Query("employeeId")
	departmentId := c.Query("departmentId")

	// Build filter
	filter := entities.LeaveRequestFilter{}

	if month != "" {
		filter.Month = &month
	}

	if employeeId != "" {
		if userID, err := uuid.Parse(employeeId); err == nil {
			filter.StaffID = &userID
		}
	}

	if departmentId != "" {
		if deptID, err := uuid.Parse(departmentId); err == nil {
			filter.DepartmentID = &deptID
		}
	}

	// Get leaves from use case
	leaves, err := h.leaveUseCase.GetLeavesByFilter(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch leave requests",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลวันหยุดพนักงานสำเร็จ",
		"data":    leaves,
	})
}

// CreateLeave creates a new leave request
func (h *LeaveHandler) CreateLeave(c *fiber.Ctx) error {
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

	// Parse UUIDs (now staff id)
	userID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid employee ID format",
		})
	}

	deptID, err := uuid.Parse(req.DepartmentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid department ID format",
		})
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid date format. Use YYYY-MM-DD",
		})
	}

	// Create leave request
	leaveReq := entities.LeaveRequestCreate{
		StaffID:      userID,
		DepartmentID: deptID,
		LeaveType:    entities.LeaveTypePersonal, // Default to personal leave
		StartDate:    date,
		EndDate:      date,
		Reason:       &req.Reason,
	}

	// Create leave through use case
	leaveID, err := h.leaveUseCase.CreateLeave(context.Background(), leaveReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create leave request",
			"error":   err.Error(),
		})
	}

	// Get created leave details
	leave, err := h.leaveUseCase.GetLeaveByID(context.Background(), leaveID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Leave created but failed to fetch details",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "สร้างวันหยุดพนักงานสำเร็จ",
		"data":    leave,
	})
}

// GetLeavesByDepartment returns leaves for specific department
func (h *LeaveHandler) GetLeavesByDepartment(c *fiber.Ctx) error {
	departmentId := c.Params("departmentId")

	deptID, err := uuid.Parse(departmentId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid department ID format",
		})
	}

	leaves, err := h.leaveUseCase.GetLeavesByDepartment(context.Background(), deptID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch department leaves",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลวันหยุดพนักงานตามแผนกสำเร็จ",
		"data":    leaves,
	})
}

// GetLeavesByEmployee returns leaves for specific employee
func (h *LeaveHandler) GetLeavesByEmployee(c *fiber.Ctx) error {
	employeeId := c.Params("employeeId")

	userID, err := uuid.Parse(employeeId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid employee ID format",
		})
	}

	leaves, err := h.leaveUseCase.GetLeavesByUser(context.Background(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch employee leaves",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลวันหยุดของพนักงานสำเร็จ",
		"data":    leaves,
	})
}

// UpdateLeave updates leave request information
func (h *LeaveHandler) UpdateLeave(c *fiber.Ctx) error {
	leaveId := c.Params("id")

	leaveID, err := uuid.Parse(leaveId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid leave ID format",
		})
	}

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

	// Build update request
	update := entities.LeaveRequestUpdate{}

	if req.Date != nil {
		if date, err := time.Parse("2006-01-02", *req.Date); err == nil {
			update.StartDate = &date
			update.EndDate = &date
		}
	}

	if req.Reason != nil {
		update.Reason = req.Reason
	}

	// Update leave through use case
	if err := h.leaveUseCase.UpdateLeave(context.Background(), leaveID, update); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update leave request",
			"error":   err.Error(),
		})
	}

	// Get updated leave details
	leave, err := h.leaveUseCase.GetLeaveByID(context.Background(), leaveID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Leave updated but failed to fetch details",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "อัปเดตวันหยุดพนักงานสำเร็จ",
		"data":    leave,
	})
}

// DeleteLeave deletes leave request
func (h *LeaveHandler) DeleteLeave(c *fiber.Ctx) error {
	leaveId := c.Params("id")

	leaveID, err := uuid.Parse(leaveId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid leave ID format",
		})
	}

	if err := h.leaveUseCase.DeleteLeave(context.Background(), leaveID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete leave request",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ลบวันหยุดพนักงานสำเร็จ",
	})
}

// ToggleLeave toggles active status of leave request
func (h *LeaveHandler) ToggleLeave(c *fiber.Ctx) error {
	leaveId := c.Params("id")

	leaveID, err := uuid.Parse(leaveId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid leave ID format",
		})
	}

	if err := h.leaveUseCase.ToggleLeaveStatus(context.Background(), leaveID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to toggle leave status",
			"error":   err.Error(),
		})
	}

	// Get updated leave details
	leave, err := h.leaveUseCase.GetLeaveByID(context.Background(), leaveID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Status toggled but failed to fetch details",
			"error":   err.Error(),
		})
	}

	status := "เปิดใช้งาน"
	if leave.Status == entities.LeaveStatusCancelled {
		status = "ปิดใช้งาน"
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": status + "วันหยุดพนักงานสำเร็จ",
		"data":    leave,
	})
}

// Health returns service health status
func (h *LeaveHandler) Health(c *fiber.Ctx) error {
	// Check database connection
	if err := h.db.GetDB().Ping(); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":    "error",
			"service":   "employee-leave-service",
			"message":   "Database connection failed",
			"error":     err.Error(),
			"timestamp": time.Now(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "ok",
		"service":   "employee-leave-service",
		"message":   "Service and database healthy",
		"timestamp": time.Now(),
	})
}

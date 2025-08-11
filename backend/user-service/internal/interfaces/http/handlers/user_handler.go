package handlers

import (
	"strconv"
	"time"

	"nurseshift/user-service/internal/domain/entities"
	"nurseshift/user-service/internal/domain/usecases"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userUseCase usecases.UserUseCase
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUseCase usecases.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// GetProfile returns current user profile
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	profile, err := h.userUseCase.GetProfile(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถดึงข้อมูลโปรไฟล์ได้",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลโปรไฟล์สำเร็จ",
		"data":    profile,
	})
}

// UpdateProfile updates user profile
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	var req usecases.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	user, err := h.userUseCase.UpdateProfile(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถอัปเดตโปรไฟล์ได้",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "อัปเดตโปรไฟล์สำเร็จ",
		"data":    user,
	})
}

// GetUsers returns paginated list of users (admin/manager only)
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	organizationID := c.Locals("organizationID").(uuid.UUID)

	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	var role *entities.UserRole
	if roleStr := c.Query("role"); roleStr != "" {
		r := entities.UserRole(roleStr)
		role = &r
	}

	var status *entities.UserStatus
	if statusStr := c.Query("status"); statusStr != "" {
		s := entities.UserStatus(statusStr)
		status = &s
	}

	var departmentID *uuid.UUID
	if deptStr := c.Query("departmentId"); deptStr != "" {
		if deptUUID, err := uuid.Parse(deptStr); err == nil {
			departmentID = &deptUUID
		}
	}

	req := &usecases.GetUsersRequest{
		Role:         role,
		Status:       status,
		DepartmentID: departmentID,
		Page:         page,
		Limit:        limit,
	}

	response, err := h.userUseCase.GetUsers(c.Context(), organizationID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถดึงข้อมูลผู้ใช้ได้",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลผู้ใช้สำเร็จ",
		"data":    response,
	})
}

// GetUser returns specific user details
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	requesterID := c.Locals("userID").(uuid.UUID)

	userIDStr := c.Params("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "รหัสผู้ใช้ไม่ถูกต้อง",
		})
	}

	profile, err := h.userUseCase.GetUser(c.Context(), userID, requesterID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่มีสิทธิ์เข้าถึงข้อมูลนี้",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลผู้ใช้สำเร็จ",
		"data":    profile,
	})
}

// SearchUsers searches users by query
func (h *UserHandler) SearchUsers(c *fiber.Ctx) error {
	organizationID := c.Locals("organizationID").(uuid.UUID)

	query := c.Query("q", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	var role *entities.UserRole
	if roleStr := c.Query("role"); roleStr != "" {
		r := entities.UserRole(roleStr)
		role = &r
	}

	var status *entities.UserStatus
	if statusStr := c.Query("status"); statusStr != "" {
		s := entities.UserStatus(statusStr)
		status = &s
	}

	req := &usecases.SearchUsersRequest{
		Query:  query,
		Role:   role,
		Status: status,
		Page:   page,
		Limit:  limit,
	}

	response, err := h.userUseCase.SearchUsers(c.Context(), organizationID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถค้นหาผู้ใช้ได้",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ค้นหาผู้ใช้สำเร็จ",
		"data":    response,
	})
}

// GetUserStats returns user statistics
func (h *UserHandler) GetUserStats(c *fiber.Ctx) error {
	organizationID := c.Locals("organizationID").(uuid.UUID)

	stats, err := h.userUseCase.GetUserStats(c.Context(), organizationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถดึงสstatisticsข้อมูลได้",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงสถิติผู้ใช้สำเร็จ",
		"data":    stats,
	})
}

// UploadAvatar handles avatar upload
func (h *UserHandler) UploadAvatar(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	// For now, we'll just accept a URL
	// In production, this would handle file upload to cloud storage
	var req struct {
		AvatarURL string `json:"avatarUrl" validate:"required,url"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	err := h.userUseCase.UploadAvatar(c.Context(), userID, req.AvatarURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถอัปเดตรูปโปรไฟล์ได้",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "อัปเดตรูปโปรไฟล์สำเร็จ",
	})
}

// Health returns service health status
func (h *UserHandler) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "ok",
		"service":   "user-service",
		"timestamp": time.Now(),
	})
}



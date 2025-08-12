package handlers

import (
	"time"

	"nurseshift/user-service/internal/domain/repositories"
	"nurseshift/user-service/internal/domain/usecases"

	"github.com/gofiber/fiber/v2"
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

// GetProfile returns the authenticated user's profile
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userIDStr := c.Locals("userID").(string)

	// Get user profile
	user, err := h.userUseCase.GetProfile(c.Context(), userIDStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถดึงข้อมูลโปรไฟล์ได้",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลโปรไฟล์สำเร็จ",
		"data":    user,
	})
}

// UpdateProfile updates the authenticated user's profile
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userIDStr := c.Locals("userID").(string)

	// Parse request body
	var req usecases.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถอ่านข้อมูลได้",
			"error":   err.Error(),
		})
	}

	// Update profile
	user, err := h.userUseCase.UpdateProfile(c.Context(), userIDStr, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถอัปเดตโปรไฟล์ได้",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "อัปเดตโปรไฟล์สำเร็จ",
		"data":    user,
	})
}

// GetUsers returns a list of users (admin only)
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userIDStr := c.Locals("userID").(string)

	// Parse query parameters
	var req repositories.GetUsersRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถอ่านพารามิเตอร์ได้",
			"error":   err.Error(),
		})
	}

	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	// Get users
	response, err := h.userUseCase.GetUsers(c.Context(), userIDStr, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถดึงข้อมูลผู้ใช้ได้",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลผู้ใช้สำเร็จ",
		"data":    response,
	})
}

// GetUser returns a specific user by ID
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	_ = c.Locals("userID").(string) // Check if user is authenticated

	// Get user ID from URL parameter
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "กรุณาระบุ ID ผู้ใช้",
		})
	}

	// Get user profile
	user, err := h.userUseCase.GetProfile(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถดึงข้อมูลผู้ใช้ได้",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงข้อมูลผู้ใช้สำเร็จ",
		"data":    user,
	})
}

// SearchUsers searches for users based on query (admin only)
func (h *UserHandler) SearchUsers(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userIDStr := c.Locals("userID").(string)

	// Parse query parameters
	var req repositories.SearchUsersRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถอ่านพารามิเตอร์ได้",
			"error":   err.Error(),
		})
	}

	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	// Search users
	response, err := h.userUseCase.SearchUsers(c.Context(), userIDStr, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถค้นหาผู้ใช้ได้",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "ค้นหาผู้ใช้สำเร็จ",
		"data":    response,
	})
}

// GetUserStats returns user statistics (admin only)
func (h *UserHandler) GetUserStats(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userIDStr := c.Locals("userID").(string)

	// Get user stats
	stats, err := h.userUseCase.GetUserStats(c.Context(), userIDStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถดึงสถิติผู้ใช้ได้",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงสถิติผู้ใช้สำเร็จ",
		"data":    stats,
	})
}

// UploadAvatar uploads a user's avatar
func (h *UserHandler) UploadAvatar(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userIDStr := c.Locals("userID").(string)

	// Parse request body
	var req struct {
		AvatarURL string `json:"avatarUrl"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถอ่านข้อมูลได้",
			"error":   err.Error(),
		})
	}

	if req.AvatarURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "กรุณาระบุ URL รูปภาพ",
		})
	}

	// Upload avatar
	err := h.userUseCase.UploadAvatar(c.Context(), userIDStr, req.AvatarURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถอัปโหลดรูปภาพได้",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "อัปเดตรูปโปรไฟล์สำเร็จ",
	})
}

// SendVerificationEmail sends a verification email to the user
func (h *UserHandler) SendVerificationEmail(c *fiber.Ctx) error {
	// Parse request body
	var req struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถอ่านข้อมูลได้",
			"error":   err.Error(),
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "กรุณาระบุอีเมล",
		})
	}

	// Send verification email
	err := h.userUseCase.SendVerificationEmail(c.Context(), req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถส่งอีเมลยืนยันได้",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "ส่งอีเมลยืนยันแล้ว กรุณาตรวจสอบอีเมลของคุณ",
	})
}

// VerifyEmail verifies the user's email using a token
func (h *UserHandler) VerifyEmail(c *fiber.Ctx) error {
	// Parse request body
	var req struct {
		Token string `json:"token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถอ่านข้อมูลได้",
			"error":   err.Error(),
		})
	}

	if req.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "กรุณาระบุ token",
		})
	}

	// Verify email
	err := h.userUseCase.VerifyEmail(c.Context(), req.Token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถยืนยันอีเมลได้",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "ยืนยันอีเมลสำเร็จ",
	})
}

// CheckEmailVerification checks if a user's email is verified
func (h *UserHandler) CheckEmailVerification(c *fiber.Ctx) error {
	// Get email from URL parameter
	email := c.Params("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "กรุณาระบุอีเมล",
		})
	}

	// Check email verification status
	verified, err := h.userUseCase.IsEmailVerified(c.Context(), email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่สามารถตรวจสอบสถานะการยืนยันอีเมลได้",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "ดึงสถานะการยืนยันอีเมลสำเร็จ",
		"data": fiber.Map{
			"email":    email,
			"verified": verified,
		},
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

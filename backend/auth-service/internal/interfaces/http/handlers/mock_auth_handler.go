package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// MockAuthHandler handles authentication-related HTTP requests with mock data
type MockAuthHandler struct{}

// NewMockAuthHandler creates a new mock auth handler
func NewMockAuthHandler() *MockAuthHandler {
	return &MockAuthHandler{}
}

// MockUser represents a user in mock responses
type MockUser struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Phone     *string   `json:"phone,omitempty"`
	Position  *string   `json:"position,omitempty"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// MockLoginResponse represents a mock login response
type MockLoginResponse struct {
	Status       string    `json:"status"`
	Message      string    `json:"message"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	User         *MockUser `json:"user"`
}

// Login handles user authentication with mock data
func (h *MockAuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Mock login - for demo purposes
	var user *MockUser

	// Check credentials
	if req.Email == "admin@nurseshift.com" && req.Password == "admin123" {
		user = &MockUser{
			ID:        "admin-1",
			Email:     req.Email,
			FirstName: "Admin",
			LastName:  "System",
			Role:      "admin",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	} else if req.Email == "user@nurseshift.com" && req.Password == "user123" {
		user = &MockUser{
			ID:        "user-1",
			Email:     req.Email,
			FirstName: "User",
			LastName:  "Demo",
			Role:      "user",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "อีเมลหรือรหัสผ่านไม่ถูกต้อง",
		})
	}

	// Mock JWT tokens
	accessToken := "mock-jwt-access-token-" + user.ID
	refreshToken := "mock-jwt-refresh-token-" + user.ID
	expiresAt := time.Now().Add(24 * time.Hour)

	return c.Status(fiber.StatusOK).JSON(MockLoginResponse{
		Status:       "success",
		Message:      "เข้าสู่ระบบสำเร็จ",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         user,
	})
}

// Register handles user registration - สมาชิกใหม่จะเป็น role user เสมอ
func (h *MockAuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Mock registration response - บังคับเป็น user role เสมอ
	user := &MockUser{
		ID:        uuid.New().String(),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Position:  req.Position,
		Role:      "user", // บังคับเป็น user เสมอ (ไม่ให้เลือก role)
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "สมัครสมาชิกสำเร็จ ระบบกำหนดให้เป็น role user",
		"user":    user,
	})
}

// RefreshToken handles token refresh
func (h *MockAuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Mock refresh token validation
	if req.RefreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Refresh token ไม่ถูกต้องหรือหมดอายุ",
		})
	}

	// Mock user from refresh token
	user := &MockUser{
		ID:        "user-1",
		Email:     "user@nurseshift.com",
		FirstName: "User",
		LastName:  "Demo",
		Role:      "user",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Generate new tokens
	accessToken := "mock-jwt-access-token-refreshed-" + user.ID
	refreshToken := "mock-jwt-refresh-token-refreshed-" + user.ID
	expiresAt := time.Now().Add(24 * time.Hour)

	return c.Status(fiber.StatusOK).JSON(MockLoginResponse{
		Status:       "success",
		Message:      "ต่ออายุ token สำเร็จ",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         user,
	})
}

// Logout handles user logout
func (h *MockAuthHandler) Logout(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ออกจากระบบสำเร็จ",
	})
}

// LogoutAll handles logout from all devices
func (h *MockAuthHandler) LogoutAll(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ออกจากระบบทุกอุปกรณ์สำเร็จ",
	})
}

// ChangePassword handles password change
func (h *MockAuthHandler) ChangePassword(c *fiber.Ctx) error {
	var req ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "เปลี่ยนรหัสผ่านสำเร็จ",
	})
}

// Me returns current user information
func (h *MockAuthHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	email := c.Locals("email")
	role := c.Locals("role")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"user": fiber.Map{
			"id":    userID,
			"email": email,
			"role":  role,
		},
	})
}

// CreateAdmin creates a new admin user - เฉพาะ admin เท่านั้นที่เรียกได้
func (h *MockAuthHandler) CreateAdmin(c *fiber.Ctx) error {
	// ตรวจสอบว่า user ที่เรียก API นี้เป็น admin หรือไม่
	role := c.Locals("role")
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "ต้องมีสิทธิ์ admin เท่านั้นถึงจะสร้าง admin ใหม่ได้",
		})
	}

	var req CreateAdminRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Mock admin creation response
	admin := &MockUser{
		ID:        uuid.New().String(),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Position:  req.Position,
		Role:      "admin", // สร้างเป็น admin
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "สร้าง admin ใหม่สำเร็จ",
		"user":    admin,
	})
}

// Health returns service health status
func (h *MockAuthHandler) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Auth Service is healthy",
		"timestamp": time.Now(),
	})
}

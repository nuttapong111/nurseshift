package handlers

import (
	"fmt"
	"strings"
	"time"

	"nurseshift/auth-service/internal/domain/usecases"
	"nurseshift/auth-service/internal/infrastructure/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// @title NurseShift Auth Service API
// @version 1.0
// @description Authentication microservice for NurseShift application
// @contact.name NurseShift Team
// @contact.email support@nurseshift.com
// @host localhost:8081
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token in format: Bearer {token}

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authUseCase          usecases.AuthUseCase
	jwtService           usecases.JWTService
	emailService         services.EmailService
	passwordResetService services.PasswordResetService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authUseCase usecases.AuthUseCase, jwtService usecases.JWTService, emailService services.EmailService, passwordResetService services.PasswordResetService) *AuthHandler {
	return &AuthHandler{
		authUseCase:          authUseCase,
		jwtService:           jwtService,
		emailService:         emailService,
		passwordResetService: passwordResetService,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Status       string         `json:"status"`
	Message      string         `json:"message"`
	AccessToken  string         `json:"accessToken"`
	RefreshToken string         `json:"refreshToken"`
	ExpiresAt    time.Time      `json:"expiresAt"`
	User         *usecases.User `json:"user"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=6"`
	FirstName string  `json:"firstName" validate:"required"`
	LastName  string  `json:"lastName" validate:"required"`
	Phone     *string `json:"phone,omitempty"`
	Position  *string `json:"position,omitempty"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// ChangePasswordRequest represents a change password request
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=6"`
}

// ForgotPasswordRequest represents a forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest represents a password reset request
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=6"`
}

// Login handles user authentication
// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User login credentials"
// @Success 200 {object} LoginResponse "Login successful"
// @Failure 400 {object} ErrorResponse "Invalid request data"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Get IP address and user agent
	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	// Use auth usecase for login
	loginResp, err := h.authUseCase.Login(c.Context(), req.Email, req.Password, ipAddress, userAgent)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "อีเมลหรือรหัสผ่านไม่ถูกต้อง",
		})
	}

	return c.Status(fiber.StatusOK).JSON(LoginResponse{
		Status:       "success",
		Message:      "เข้าสู่ระบบสำเร็จ",
		AccessToken:  loginResp.AccessToken,
		RefreshToken: loginResp.RefreshToken,
		ExpiresAt:    loginResp.ExpiresAt,
		User:         loginResp.User,
	})
}

// Register handles user registration - สมาชิกใหม่จะเป็น role user เสมอ
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Convert to usecase request
	registerReq := &usecases.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Position:  req.Position,
	}

	// Use auth usecase for registration
	user, err := h.authUseCase.Register(c.Context(), registerReq)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "การสมัครสมาชิกล้มเหลว",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "สมัครสมาชิกสำเร็จ ระบบกำหนดให้เป็น role user",
		"user":    user,
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// Refresh tokens
	loginResp, err := h.authUseCase.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Refresh token ไม่ถูกต้องหรือหมดอายุ",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(LoginResponse{
		Status:       "success",
		Message:      "ต่ออายุ token สำเร็จ",
		AccessToken:  loginResp.AccessToken,
		RefreshToken: loginResp.RefreshToken,
		ExpiresAt:    loginResp.ExpiresAt,
		User:         loginResp.User,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Get user ID from JWT token
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่พบข้อมูลผู้ใช้",
		})
	}

	// For now, we'll just return success since we're using JWT tokens
	// In a real implementation, you might want to blacklist the token
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ออกจากระบบสำเร็จ",
	})
}

// LogoutAll handles logout from all devices
func (h *AuthHandler) LogoutAll(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	err := h.authUseCase.LogoutAllSessions(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "การออกจากระบบทุกอุปกรณ์ล้มเหลว",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ออกจากระบบทุกอุปกรณ์สำเร็จ",
	})
}

// VerifyTokenRequest represents a token verification request
type VerifyTokenRequest struct {
	Token string `json:"token"`
}

// VerifyToken verifies an access token and returns claims
// @Summary Verify access token
// @Description Verify JWT access token and return claims
// @Tags Authentication
// @Accept json
// @Produce json
// @Param token body VerifyTokenRequest false "Token to verify (if not provided in Authorization header)"
// @Success 200 {object} fiber.Map "Token is valid with claims"
// @Failure 401 {object} ErrorResponse "Invalid token"
// @Router /auth/introspect [post]
func (h *AuthHandler) VerifyToken(c *fiber.Ctx) error {
	// Prefer Authorization header if present
	token := ""
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			token = parts[1]
		}
	}

	// Fallback to body
	if token == "" {
		var req VerifyTokenRequest
		_ = c.BodyParser(&req)
		token = req.Token
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่พบ token สำหรับตรวจสอบ",
		})
	}

	claims, err := h.jwtService.ValidateToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Token ไม่ถูกต้องหรือหมดอายุ",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"active": true,
		"claims": fiber.Map{
			"userId":    claims.UserID,
			"role":      claims.Role,
			"type":      claims.Type,
			"issuedAt":  claims.IssuedAt,
			"expiresAt": claims.ExpiresAt,
		},
	})
}

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	var req ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	userID := c.Locals("userID").(uuid.UUID)

	err := h.authUseCase.ChangePassword(c.Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "การเปลี่ยนรหัสผ่านล้มเหลว",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "เปลี่ยนรหัสผ่านสำเร็จ",
	})
}

// Me returns current user information
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	// For this endpoint, we could add a use case to get current user info
	// For now, we'll return the user ID
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"user": fiber.Map{
			"id": userID,
		},
	})
}

// CreateAdminRequest represents an admin creation request
type CreateAdminRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=6"`
	FirstName string  `json:"firstName" validate:"required"`
	LastName  string  `json:"lastName" validate:"required"`
	Phone     *string `json:"phone,omitempty"`
	Position  *string `json:"position,omitempty"`
}

// CreateAdmin creates a new admin user - เฉพาะ admin เท่านั้นที่เรียกได้
func (h *AuthHandler) CreateAdmin(c *fiber.Ctx) error {
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

	// Convert to usecase request
	registerReq := &usecases.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Position:  req.Position,
	}

	// Use auth usecase for registration
	user, err := h.authUseCase.Register(c.Context(), registerReq)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "การสร้าง admin ล้มเหลว",
			"error":   err.Error(),
		})
	}

	// Note: In a real implementation, you would need to update the user role to admin
	// after creation, or modify the usecase to accept role parameter

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "สร้าง admin ใหม่สำเร็จ",
		"user":    user,
	})
}

// Health returns service health status
func (h *AuthHandler) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Auth Service is healthy",
		"timestamp": time.Now(),
	})
}

// ForgotPassword handles forgot password request
// @Summary Forgot password
// @Description Send password reset token to user's email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Forgot password request"
// @Success 200 {object} fiber.Map "Reset token sent successfully"
// @Failure 400 {object} ErrorResponse "Invalid request data"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// 1. Check if user exists
	user, err := h.authUseCase.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		// Don't reveal if user exists or not for security
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "หากอีเมลนี้มีอยู่ในระบบ เราจะส่งรหัสยืนยันไปให้",
		})
	}

	// 2. Generate reset token
	resetToken, err := h.passwordResetService.GenerateResetToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "เกิดข้อผิดพลาดในการสร้างรหัสยืนยัน",
		})
	}

	// 3. Store reset token with expiration
	err = h.passwordResetService.StoreResetToken(user.ID, resetToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "เกิดข้อผิดพลาดในการบันทึกรหัสยืนยัน",
		})
	}

	// 4. Send email with reset token
	err = h.emailService.SendPasswordResetEmail(user.Email, resetToken)
	if err != nil {
		// Log error but don't reveal to user
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "เกิดข้อผิดพลาดในการส่งอีเมล",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ส่งรหัสยืนยันไปยังอีเมลของคุณแล้ว กรุณาตรวจสอบกล่องจดหมาย",
	})
}

// ResetPassword handles password reset with token
// @Summary Reset password
// @Description Reset password using reset token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Password reset request"
// @Success 200 {object} fiber.Map "Password reset successful"
// @Failure 400 {object} ErrorResponse "Invalid request data"
// @Failure 400 {object} ErrorResponse "Invalid or expired token"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ข้อมูลที่ส่งมาไม่ถูกต้อง",
			"error":   err.Error(),
		})
	}

	// 1. Validate reset token
	if !h.passwordResetService.ValidateResetToken(req.Token) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "รหัสยืนยันไม่ถูกต้องหรือหมดอายุแล้ว",
		})
	}

	// 2. Get user ID from token
	userID, err := h.passwordResetService.GetUserIDByToken(req.Token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "รหัสยืนยันไม่ถูกต้องหรือหมดอายุแล้ว",
		})
	}

	// 3. Update user password
	err = h.authUseCase.UpdatePassword(c.Context(), userID, req.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "เกิดข้อผิดพลาดในการรีเซ็ตรหัสผ่าน",
		})
	}

	// 4. Clear reset token
	err = h.passwordResetService.ClearResetToken(req.Token)
	if err != nil {
		// Log error but don't fail the password reset
		fmt.Printf("Warning: failed to clear reset token: %v\n", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "รีเซ็ตรหัสผ่านสำเร็จ กรุณาเข้าสู่ระบบด้วยรหัสผ่านใหม่",
	})
}

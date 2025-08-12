package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware creates an authentication middleware
type introspectResponse struct {
	Status string `json:"status"`
	Active bool   `json:"active"`
	Claims struct {
		UserID string `json:"userId"`
		Role   string `json:"role"`
	} `json:"claims"`
}

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Authorization header required",
				"message": "กรุณาเข้าสู่ระบบ",
			})
		}

		// Check Bearer token format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Invalid authorization format",
				"message": "รูปแบบ token ไม่ถูกต้อง",
			})
		}

		baseURL := os.Getenv("AUTH_SERVICE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:8081"
		}
		url := strings.TrimRight(baseURL, "/") + "/api/v1/auth/introspect"

		httpClient := &http.Client{Timeout: 3 * time.Second}
		req, _ := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil || resp == nil || resp.StatusCode != http.StatusOK {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Introspection failed",
				"message": "ไม่สามารถตรวจสอบ token ได้",
			})
		}
		defer resp.Body.Close()

		var introspect introspectResponse
		if err := json.NewDecoder(resp.Body).Decode(&introspect); err != nil || !introspect.Active {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Invalid token",
				"message": "Token ไม่ถูกต้องหรือหมดอายุ",
			})
		}

		// Store user information in context
		c.Locals("userID", introspect.Claims.UserID)
		c.Locals("userRole", introspect.Claims.Role)

		return c.Next()
	}
}

// OptionalAuthMiddleware creates an optional authentication middleware
func OptionalAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			return c.Next()
		}

		baseURL := os.Getenv("AUTH_SERVICE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:8081"
		}
		url := strings.TrimRight(baseURL, "/") + "/api/v1/auth/introspect"

		httpClient := &http.Client{Timeout: 3 * time.Second}
		req, _ := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var introspect introspectResponse
			if err := json.NewDecoder(resp.Body).Decode(&introspect); err == nil && introspect.Active {
				c.Locals("userID", introspect.Claims.UserID)
				c.Locals("userRole", introspect.Claims.Role)
			}
		}

		return c.Next()
	}
}

// RoleMiddleware creates a role-based authorization middleware
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("userRole")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Authentication required",
				"message": "กรุณาเข้าสู่ระบบ",
			})
		}

		role := string(userRole.(string))
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "Insufficient permissions",
			"message": "คุณไม่มีสิทธิ์เข้าถึงข้อมูลนี้",
		})
	}
}

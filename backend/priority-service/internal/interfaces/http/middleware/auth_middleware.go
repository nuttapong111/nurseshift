package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// JWTClaims represents the JWT claims structure
type introspectResponse struct {
	Status string `json:"status"`
	Active bool   `json:"active"`
	Claims struct {
		UserID         string `json:"userId"`
		Role           string `json:"role"`
		OrganizationID string `json:"organizationId"`
	} `json:"claims"`
}

// AuthMiddleware validates JWT tokens
func AuthMiddleware(_ string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "ไม่พบ Authorization header",
			})
		}

		// Call auth-service introspection endpoint
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
				"status":  "error",
				"message": "ไม่สามารถตรวจสอบ token ได้",
			})
		}
		defer resp.Body.Close()

		var introspect introspectResponse
		if err := json.NewDecoder(resp.Body).Decode(&introspect); err != nil || !introspect.Active {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Token ไม่ถูกต้องหรือหมดอายุ",
			})
		}

		c.Locals("userID", introspect.Claims.UserID)
		c.Locals("role", introspect.Claims.Role)
		c.Locals("organizationID", introspect.Claims.OrganizationID)

		return c.Next()
	}
}

// AdminOnlyMiddleware ensures only admin users can access
func AdminOnlyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "ต้องมีสิทธิ์ admin เท่านั้น",
			})
		}
		return c.Next()
	}
}

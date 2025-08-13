package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// IntrospectResponse represents the response from auth service introspect endpoint
type IntrospectResponse struct {
	Active bool `json:"active"`
	Claims struct {
		UserID         string `json:"userId"`
		Email          string `json:"email"`
		Role           string `json:"role"`
		OrganizationID string `json:"organizationId"`
	} `json:"claims"`
	Status string `json:"status"`
}

// AuthMiddleware validates JWT tokens by calling the auth service
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Authorization header is required",
			})
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid authorization header format",
			})
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Call auth service to introspect token
		baseURL := os.Getenv("AUTH_SERVICE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:8081"
		}
		url := strings.TrimRight(baseURL, "/") + "/api/v1/auth/introspect"

		// Create HTTP request
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to create introspection request",
			})
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		// Make HTTP request
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to validate token",
			})
		}
		defer resp.Body.Close()

		// Parse response
		var introspect IntrospectResponse
		if err := json.NewDecoder(resp.Body).Decode(&introspect); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to parse introspection response",
			})
		}

		// Check if token is valid
		if !introspect.Active {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid or expired token",
			})
		}

		// Store user information in context
		c.Locals("userID", introspect.Claims.UserID)
		c.Locals("userRole", introspect.Claims.Role)
		c.Locals("organizationID", introspect.Claims.OrganizationID)

		return c.Next()
	}
}

package routes

import (
	"github.com/gofiber/fiber/v2"
)

// AuthUseCase interface for authentication use cases
type AuthUseCase interface {
	// Add methods as needed
}

// JWTService interface for JWT operations
type JWTService interface {
	// Add methods as needed
}

// AuthHandler interface for authentication handlers
type AuthHandler interface {
	// Add methods as needed
}

// AuthMiddleware interface for authentication middleware
type AuthMiddleware interface {
	// Add methods as needed
}

// SetupAuthRoutes sets up authentication-related routes
func SetupAuthRoutes(app *fiber.App, authUseCase AuthUseCase, _ JWTService) {
	// Simple health check endpoint for now
	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Priority service auth routes are available",
		})
	})
}

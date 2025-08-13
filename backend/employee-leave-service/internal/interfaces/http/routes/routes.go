package routes

import (
	"nurseshift/employee-leave-service/internal/interfaces/http/handlers"
	mw "nurseshift/employee-leave-service/internal/interfaces/http/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(app *fiber.App, h *handlers.LeaveHandler) {
	api := app.Group("/api/v1")

	// Leave request routes with authentication
	leaves := api.Group("/leaves")
	leaves.Use(mw.AuthMiddleware())
	{
		leaves.Get("/", h.GetLeaves)
		leaves.Post("/", h.CreateLeave)
		leaves.Get("/departments/:departmentId", h.GetLeavesByDepartment)
		leaves.Get("/employees/:employeeId", h.GetLeavesByEmployee)
		leaves.Put("/:id", h.UpdateLeave)
		leaves.Delete("/:id", h.DeleteLeave)
		leaves.Put("/:id/toggle", h.ToggleLeave)
	}

	// Health check (no auth required)
	app.Get("/health", h.Health)
}

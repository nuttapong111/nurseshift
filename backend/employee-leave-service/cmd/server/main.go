package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nurseshift/employee-leave-service/internal/infrastructure/config"
	"nurseshift/employee-leave-service/internal/interfaces/http/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Starting Employee Leave Service on port %s...\n", cfg.Server.Port)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "NurseShift Employee Leave Service",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:3002",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	if cfg.IsDevelopment() {
		app.Use(logger.New())
	}

	// Initialize handlers
	leaveHandler := handlers.NewEmployeeLeaveHandler()

	// Routes
	api := app.Group("/api/v1")
	leaves := api.Group("/employee-leaves")
	{
		leaves.Get("/", leaveHandler.GetEmployeeLeaves)
		leaves.Post("/", leaveHandler.CreateEmployeeLeave)
		leaves.Get("/departments/:departmentId", leaveHandler.GetLeavesByDepartment)
		leaves.Get("/employees/:employeeId", leaveHandler.GetLeavesByEmployee)
		leaves.Put("/:id", leaveHandler.UpdateEmployeeLeave)
		leaves.Delete("/:id", leaveHandler.DeleteEmployeeLeave)
		leaves.Put("/:id/toggle", leaveHandler.ToggleEmployeeLeave)
	}

	// Health check
	app.Get("/health", leaveHandler.Health)

	// Start server
	go func() {
		fmt.Printf("🚀 Employee Leave Service running on http://localhost:%s\n", cfg.Server.Port)
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			log.Printf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\n🛑 Shutting down Employee Leave Service...")
	app.Shutdown()
	fmt.Println("✅ Employee Leave Service stopped gracefully")
}

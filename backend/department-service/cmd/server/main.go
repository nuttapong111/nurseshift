package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"nurseshift/department-service/internal/infrastructure/config"
	"nurseshift/department-service/internal/infrastructure/database"
	"nurseshift/department-service/internal/interfaces/http/handlers"
	"nurseshift/department-service/internal/interfaces/http/middleware"

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

	fmt.Printf("Starting Department Service on port %s...\n", cfg.Server.Port)

	// Initialize database connection
	dbConn, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Initialize repositories
	deptRepo := database.NewPostgresDepartmentRepository(dbConn.DB, "nurse_shift")

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "NurseShift Department Service",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.CORS.Origins, ","),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: cfg.CORS.Credentials,
	}))

	if cfg.IsDevelopment() {
		app.Use(logger.New())
	}

	// Initialize handlers
	deptHandler := handlers.NewDepartmentHandler(deptRepo)

	// Routes
	api := app.Group("/api/v1")

	// Protected routes (authentication required)
	departments := api.Group("/departments")
	departments.Use(middleware.AuthMiddleware(""))
	{
		departments.Get("/", deptHandler.GetDepartments)
		departments.Post("/", deptHandler.CreateDepartment)
		departments.Get("/stats", deptHandler.GetDepartmentStats)
		departments.Get("/:id", deptHandler.GetDepartment)
		departments.Put("/:id", deptHandler.UpdateDepartment)
		departments.Delete("/:id", deptHandler.DeleteDepartment)
		departments.Get("/:id/staff", deptHandler.GetDepartmentStaff)
		departments.Post("/:id/staff", deptHandler.AddDepartmentStaff)
		departments.Delete("/:id/staff/:staffId", deptHandler.DeleteDepartmentStaff)
	}

	// Health check
	app.Get("/health", deptHandler.Health)

	// Start server in goroutine
	go func() {
		fmt.Printf("🚀 Department Service running on http://localhost:%s\n", cfg.Server.Port)
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			log.Printf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\n🛑 Shutting down Department Service...")
	app.Shutdown()
	fmt.Println("✅ Department Service stopped gracefully")
}

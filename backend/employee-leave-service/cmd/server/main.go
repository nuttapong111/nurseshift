package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nurseshift/employee-leave-service/internal/domain/usecases"
	"nurseshift/employee-leave-service/internal/infrastructure/database"
	"nurseshift/employee-leave-service/internal/infrastructure/repositories"
	"nurseshift/employee-leave-service/internal/interfaces/http/handlers"
	"nurseshift/employee-leave-service/internal/interfaces/http/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load environment variables from config.env
	log.Println("Loading environment variables...")

	// Database connection
	dbConn, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Get schema from environment
	schema := os.Getenv("DB_SCHEMA")
	if schema == "" {
		schema = "nurse_shift"
	}

	// Initialize repository
	leaveRepo := repositories.NewPostgresLeaveRepository(dbConn.GetDB(), schema)

	// Initialize use case
	leaveUseCase := usecases.NewLeaveUseCase(leaveRepo)

	// Initialize handler
	leaveHandler := handlers.NewLeaveHandler(leaveUseCase)

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

	if os.Getenv("ENV") == "development" {
		app.Use(logger.New())
	}

	// Setup routes
	routes.SetupRoutes(app, leaveHandler)

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}

	fmt.Printf("Starting Employee Leave Service on port %s...\n", port)

	// Start server
	go func() {
		fmt.Printf("🚀 Employee Leave Service running on http://localhost:%s\n", port)
		if err := app.Listen(":" + port); err != nil {
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

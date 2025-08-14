package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nurseshift/priority-service/internal/infrastructure/config"
	"nurseshift/priority-service/internal/infrastructure/database"
	"nurseshift/priority-service/internal/interfaces/http/handlers"
	"nurseshift/priority-service/internal/interfaces/http/middleware"

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

	fmt.Printf("Starting Priority Service on port %s...\n", cfg.Server.Port)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "NurseShift Priority Service",
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

	// Initialize DB connection & repository
	conn, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer conn.Close()
	priorityRepo := database.NewPriorityRepository(conn)

	// Initialize handlers
	priorityHandler := handlers.NewPriorityHandler(priorityRepo)

	// Routes
	api := app.Group("/api/v1")

	// Protected routes (authentication required)
	priorities := api.Group("/priorities")
	priorities.Use(middleware.AuthMiddleware(""))
	{
		priorities.Get("/", priorityHandler.GetPriorities)
		priorities.Put("/:id", priorityHandler.UpdatePriority)
		priorities.Put("/:id/setting", priorityHandler.UpdatePrioritySetting)
		priorities.Post("/swap", priorityHandler.SwapPriorityOrder)
	}

	// Health check
	app.Get("/health", priorityHandler.Health)

	// Start server
	go func() {
		fmt.Printf("🚀 Priority Service running on http://localhost:%s\n", cfg.Server.Port)
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			log.Printf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\n🛑 Shutting down Priority Service...")
	app.Shutdown()
	fmt.Println("✅ Priority Service stopped gracefully")
}

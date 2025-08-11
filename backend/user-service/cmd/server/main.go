package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nurseshift/user-service/internal/infrastructure/config"
	"nurseshift/user-service/internal/interfaces/http/handlers"

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

	fmt.Printf("Starting User Service on port %s...\n", cfg.Server.Port)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "NurseShift User Service",
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

	// Initialize handlers (with mock data for now)
	userHandler := handlers.NewUserHandler(nil)

	// Routes
	api := app.Group("/api/v1")
	users := api.Group("/users")
	{
		users.Get("/profile", userHandler.GetProfile)
		users.Put("/profile", userHandler.UpdateProfile)
		users.Post("/avatar", userHandler.UploadAvatar)
		users.Get("/", userHandler.GetUsers)
		users.Get("/search", userHandler.SearchUsers)
		users.Get("/stats", userHandler.GetUserStats)
		users.Get("/:id", userHandler.GetUser)
	}

	// Health check
	app.Get("/health", userHandler.Health)

	// Start server in goroutine
	go func() {
		fmt.Printf("🚀 User Service running on http://localhost:%s\n", cfg.Server.Port)
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			log.Printf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\n🛑 Shutting down User Service...")
	app.Shutdown()
	fmt.Println("✅ User Service stopped gracefully")
}



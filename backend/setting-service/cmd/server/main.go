package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nurseshift/setting-service/internal/infrastructure/config"
	"nurseshift/setting-service/internal/interfaces/http/handlers"

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

	fmt.Printf("Starting Setting Service on port %s...\n", cfg.Server.Port)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "NurseShift Setting Service",
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
	settingHandler := handlers.NewSettingHandler()

	// Routes
	api := app.Group("/api/v1")
	settings := api.Group("/settings")
	{
		settings.Get("/", settingHandler.GetSettings)
		settings.Put("/", settingHandler.UpdateSettings)
	}

	// Health check
	app.Get("/health", settingHandler.Health)

	// Start server
	go func() {
		fmt.Printf("🚀 Setting Service running on http://localhost:%s\n", cfg.Server.Port)
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			log.Printf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\n🛑 Shutting down Setting Service...")
	app.Shutdown()
	fmt.Println("✅ Setting Service stopped gracefully")
}



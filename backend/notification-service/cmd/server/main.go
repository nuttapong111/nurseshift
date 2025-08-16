package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"nurseshift/notification-service/internal/infrastructure/config"
	"nurseshift/notification-service/internal/interfaces/http/handlers"

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

	fmt.Printf("Starting Notification Service on port %s...\n", cfg.Server.Port)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "NurseShift Notification Service",
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
		AllowCredentials: true,
	}))

	if cfg.IsDevelopment() {
		app.Use(logger.New())
	}

	// Initialize handlers
	notificationHandler := handlers.NewNotificationHandler()

	// Routes
	api := app.Group("/api/v1")
	notifications := api.Group("/notifications")
	{
		notifications.Get("/", notificationHandler.GetNotifications)
		notifications.Get("/stats", notificationHandler.GetNotificationStats)
		notifications.Post("/", notificationHandler.CreateNotification)
		notifications.Put("/:id/read", notificationHandler.MarkAsRead)
		notifications.Put("/read-all", notificationHandler.MarkAllAsRead)
		notifications.Delete("/:id", notificationHandler.DeleteNotification)
	}

	// Health check
	app.Get("/health", notificationHandler.Health)

	// Start server
	go func() {
		fmt.Printf("🚀 Notification Service running on http://localhost:%s\n", cfg.Server.Port)
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			log.Printf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\n🛑 Shutting down Notification Service...")
	app.Shutdown()
	fmt.Println("✅ Notification Service stopped gracefully")
}



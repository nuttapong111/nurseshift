package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nurseshift/payment-service/internal/infrastructure/config"
	"nurseshift/payment-service/internal/interfaces/http/handlers"

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

	fmt.Printf("Starting Payment Service on port %s...\n", cfg.Server.Port)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "NurseShift Payment Service",
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
	paymentHandler := handlers.NewPaymentHandler()

	// Routes
	api := app.Group("/api/v1")
	payments := api.Group("/payments")
	{
		payments.Get("/", paymentHandler.GetPayments)
		payments.Post("/", paymentHandler.CreatePayment)
		payments.Get("/stats", paymentHandler.GetPaymentStats)
		payments.Get("/:id", paymentHandler.GetPayment)
		payments.Put("/:id", paymentHandler.UpdatePayment)
		payments.Put("/:id/approve", paymentHandler.ApprovePayment)
		payments.Put("/:id/reject", paymentHandler.RejectPayment)
	}

	// Health check
	app.Get("/health", paymentHandler.Health)

	// Start server
	go func() {
		fmt.Printf("🚀 Payment Service running on http://localhost:%s\n", cfg.Server.Port)
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			log.Printf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\n🛑 Shutting down Payment Service...")
	app.Shutdown()
	fmt.Println("✅ Payment Service stopped gracefully")
}



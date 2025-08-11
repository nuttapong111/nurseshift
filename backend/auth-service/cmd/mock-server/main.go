package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nurseshift/auth-service/internal/interfaces/http/handlers"
	"nurseshift/auth-service/internal/interfaces/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "NurseShift Auth Service (Mock)",
		ServerHeader: "Fiber",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3001,https://localhost:3001",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))
	app.Use(logger.New())

	// JWT Secret for middleware
	jwtSecret := "nurseshift-super-secret-jwt-key-development-only-2024"

	// Initialize mock handlers
	mockAuthHandler := handlers.NewMockAuthHandler()

	// Public routes (no authentication required)
	auth := app.Group("/api/v1/auth")
	{
		auth.Post("/login", mockAuthHandler.Login)
		auth.Post("/register", mockAuthHandler.Register)
		auth.Post("/refresh", mockAuthHandler.RefreshToken)
		auth.Get("/health", mockAuthHandler.Health)
	}

	// Protected routes (authentication required)
	protected := app.Group("/api/v1/auth")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		protected.Post("/logout", mockAuthHandler.Logout)
		protected.Post("/logout-all", mockAuthHandler.LogoutAll)
		protected.Post("/change-password", mockAuthHandler.ChangePassword)
		protected.Get("/me", mockAuthHandler.Me)

		// Admin only routes
		protected.Post("/create-admin", mockAuthHandler.CreateAdmin) // เฉพาะ admin เท่านั้น
	}

	// Health check
	app.Get("/health", mockAuthHandler.Health)

	// Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Start server in goroutine
	go func() {
		fmt.Printf("🚀 Mock Auth Service started on http://localhost:%s\n", port)
		fmt.Println("📚 API Endpoints:")
		fmt.Printf("   - Login: POST http://localhost:%s/api/v1/auth/login\n", port)
		fmt.Printf("   - Register: POST http://localhost:%s/api/v1/auth/register\n", port)
		fmt.Printf("   - Create Admin: POST http://localhost:%s/api/v1/auth/create-admin (admin only)\n", port)
		fmt.Printf("   - Health: GET http://localhost:%s/health\n", port)
		fmt.Println("\n📋 Test Credentials:")
		fmt.Println("   Admin: admin@nurseshift.com / admin123")
		fmt.Println("   User:  user@nurseshift.com / user123")

		if err := app.Listen(":" + port); err != nil {
			log.Printf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until we receive our signal
	<-c

	fmt.Println("\n🛑 Shutting down Mock Auth Service...")

	// Shutdown server
	if err := app.Shutdown(); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	fmt.Println("✅ Mock Auth Service stopped gracefully")
}

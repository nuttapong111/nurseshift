package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"nurseshift/user-service/internal/domain/usecases"
	"nurseshift/user-service/internal/infrastructure/config"
	"nurseshift/user-service/internal/infrastructure/database"
	"nurseshift/user-service/internal/interfaces/http/handlers"
	"nurseshift/user-service/internal/interfaces/http/middleware"

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

	// Initialize database connection
	dbConn, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	fmt.Println("✅ Database connection established successfully")

	// Initialize repository
	userRepo := database.NewPostgresUserRepository(dbConn.DB, cfg.Database.Schema)

	// Initialize use case
	userUseCase := usecases.NewUserUseCase(userRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userUseCase)

	// Import auth middleware
	authMiddleware := middleware.AuthMiddleware()

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
		AllowOrigins:     strings.Join(cfg.CORS.Origins, ","),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: cfg.CORS.Credentials,
	}))

	if cfg.IsDevelopment() {
		app.Use(logger.New(logger.Config{
			Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
		}))
	}

	// Routes
	api := app.Group("/api/v1")
	
	// Public routes (no authentication required)
	{
		api.Post("/users/send-verification-email", userHandler.SendVerificationEmail)
		api.Post("/users/verify-email", userHandler.VerifyEmail)
		api.Get("/users/check-email-verification/:email", userHandler.CheckEmailVerification)
	}
	
	users := api.Group("/users", authMiddleware) // Protected routes
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
		fmt.Printf("📚 API Documentation:\n")
		fmt.Printf("   - Health Check: GET http://localhost:%s/health\n", cfg.Server.Port)
		fmt.Printf("   - Send Verification Email: POST http://localhost:%s/api/v1/users/send-verification-email\n", cfg.Server.Port)
		fmt.Printf("   - Verify Email: POST http://localhost:%s/api/v1/users/verify-email\n", cfg.Server.Port)
		fmt.Printf("   - Check Email Verification: GET http://localhost:%s/api/v1/users/check-email-verification/:email\n", cfg.Server.Port)
		fmt.Printf("   - User Profile: GET http://localhost:%s/api/v1/users/profile\n", cfg.Server.Port)
		fmt.Printf("   - Get Users: GET http://localhost:%s/api/v1/users\n", cfg.Server.Port)
		fmt.Printf("   - Search Users: GET http://localhost:%s/api/v1/users/search\n", cfg.Server.Port)
		fmt.Printf("   - User Stats: GET http://localhost:%s/api/v1/users/stats\n", cfg.Server.Port)
		fmt.Printf("🗄️  Connected to PostgreSQL database: %s\n", cfg.Database.Name)
		
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

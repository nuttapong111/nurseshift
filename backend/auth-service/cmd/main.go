package main

import (
	"log"
	"strconv"
	"time"

	"nurseshift/auth-service/internal/domain/usecases"
	"nurseshift/auth-service/internal/infrastructure/config"
	"nurseshift/auth-service/internal/infrastructure/repositories"
	"nurseshift/auth-service/internal/infrastructure/services"
	"nurseshift/auth-service/internal/interfaces/http/handlers"
	"nurseshift/auth-service/internal/interfaces/http/middleware"
	"nurseshift/auth-service/internal/interfaces/http/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewMockUserRepository() // Using mock for now

	// Initialize services
	jwtService := services.NewJWTService(
		cfg.JWT.Secret,
		time.Duration(cfg.JWT.AccessTokenExpireHours)*time.Hour,
		time.Duration(cfg.JWT.RefreshTokenExpireDays)*24*time.Hour,
	)

	passwordService := services.NewPasswordService(cfg.Security.BCryptCost)

	// Initialize use cases
	authUseCase := usecases.NewAuthUseCase(
		userRepo,
		jwtService,
		passwordService,
		time.Duration(cfg.Security.SessionTimeoutMinutes)*time.Minute,
	)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authUseCase, jwtService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Initialize routes
	routeHandler := routes.NewRouteHandler(authHandler, authMiddleware)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		},
	})

	// Add middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.Origins,
		AllowCredentials: cfg.CORS.Credentials,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
	}))

	// Setup routes
	routeHandler.SetupRoutes(app)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "success",
			"message":   "Auth Service is running",
			"timestamp": time.Now(),
			"version":   "1.0.0",
		})
	})

	// Start server
	port := strconv.Itoa(cfg.Server.Port)
	log.Printf("Starting Auth Service on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

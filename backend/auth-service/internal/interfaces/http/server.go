package http

import (
	"time"

	"nurseshift/auth-service/docs"
	"nurseshift/auth-service/internal/domain/usecases"
	"nurseshift/auth-service/internal/infrastructure/config"
	"nurseshift/auth-service/internal/interfaces/http/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// Server represents the HTTP server
type Server struct {
	app    *fiber.App
	config *config.Config
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config, authUseCase usecases.AuthUseCase, jwtService usecases.JWTService, jwtSecret string) *Server {
	app := fiber.New(fiber.Config{
		AppName:           "NurseShift Auth Service",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		EnablePrintRoutes: cfg.IsDevelopment(),
		ErrorHandler:      errorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(helmet.New())

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     joinStrings(cfg.CORS.Origins),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With,X-Session-ID",
		AllowCredentials: cfg.CORS.Credentials,
		MaxAge:           86400, // 24 hours
	}))

	// Logger middleware (only in development)
	if cfg.IsDevelopment() {
		app.Use(logger.New(logger.Config{
			Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
		}))
	}

	// Routes
	routes.SetupAuthRoutes(app, authUseCase, jwtService, jwtSecret, cfg)

	// Swagger Documentation
	docs.SwaggerInfo.Host = "localhost:8081"
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Global health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":    "ok",
			"service":   "auth-service",
			"timestamp": time.Now(),
		})
	})

	// 404 handler
	app.Use("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "ไม่พบ endpoint ที่ต้องการ",
			"path":    c.Path(),
		})
	})

	return &Server{
		app:    app,
		config: cfg,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.app.Listen(":" + s.config.Server.Port)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

// GetApp returns the fiber app instance
func (s *Server) GetApp() *fiber.App {
	return s.app
}

// errorHandler handles global errors
func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "เกิดข้อผิดพลาดภายในเซิร์ฟเวอร์"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"status":  "error",
		"message": message,
		"error":   err.Error(),
	})
}

// joinStrings joins a slice of strings with comma
func joinStrings(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += "," + strs[i]
	}
	return result
}

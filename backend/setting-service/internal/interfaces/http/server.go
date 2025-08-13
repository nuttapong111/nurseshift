package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// Server represents the HTTP server
type Server struct{ app *fiber.App }

// NewServer creates a new HTTP server
func NewServer() *Server {
	app := fiber.New(fiber.Config{
		AppName:           "NurseShift Auth Service",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		EnablePrintRoutes: true,
		ErrorHandler:      errorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(helmet.New())

	// CORS middleware
	app.Use(cors.New(cors.Config{AllowOrigins: "*", AllowMethods: "GET,POST,PUT,DELETE,OPTIONS", AllowHeaders: "Origin,Content-Type,Accept,Authorization"}))

	// Logger middleware (only in development)
	{
		app.Use(logger.New(logger.Config{
			Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
		}))
	}

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

	return &Server{app: app}
}

// Start starts the HTTP server
func (s *Server) Start(port string) error { return s.app.Listen(":" + port) }

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

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"nurseshift/schedule-service/internal/infrastructure/config"
	dbpkg "nurseshift/schedule-service/internal/infrastructure/database"
	"nurseshift/schedule-service/internal/interfaces/http/handlers"
	"nurseshift/schedule-service/internal/interfaces/http/middleware"

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

	fmt.Printf("Starting Schedule Service on port %s...\n", cfg.Server.Port)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "NurseShift Schedule Service",
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

	// Initialize DB + repository + handlers
	conn, err := dbpkg.NewConnection(cfg)
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}
	defer conn.Close()
	repo := dbpkg.NewScheduleRepository(conn)
	scheduleHandler := handlers.NewScheduleHandler(repo)

	// Routes
	api := app.Group("/api/v1")

	// Protected routes (authentication required)
	schedules := api.Group("/schedules")
	schedules.Use(middleware.AuthMiddleware(""))
	{
		schedules.Get("/", scheduleHandler.GetSchedules)
		schedules.Post("/", scheduleHandler.CreateSchedule)
		schedules.Get("/stats", scheduleHandler.GetScheduleStats)
		schedules.Get("/shifts", scheduleHandler.ListShifts)
		schedules.Get("/available-staff", scheduleHandler.GetAvailableStaff)
		schedules.Post("/optimize-generate", scheduleHandler.OptimizeGenerate)
		schedules.Get("/:id", scheduleHandler.GetSchedule)
		schedules.Put("/:id", scheduleHandler.UpdateSchedule)
		schedules.Delete("/:id", scheduleHandler.DeleteSchedule)

		// TODO: add auto-generate endpoints next
	}

	// Auto-generate routes
	apiAuth := api.Group("/schedules").Use(middleware.AuthMiddleware(""))
	// Route compatibility: point auto-generate to the new optimizer logic
	apiAuth.Post("/auto-generate", scheduleHandler.AutoGenerate) // Enhanced Dynamic Priority Algorithm
	apiAuth.Post("/ai-generate", scheduleHandler.AIGenerate)

	// Calendar meta (working/holiday)
	api.Get("/calendar-meta", scheduleHandler.CalendarMeta)

	// Health check
	app.Get("/health", scheduleHandler.Health)

	// Start server
	go func() {
		fmt.Printf("🚀 Schedule Service running on http://localhost:%s\n", cfg.Server.Port)
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			log.Printf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\n🛑 Shutting down Schedule Service...")
	app.Shutdown()
	fmt.Println("✅ Schedule Service stopped gracefully")
}

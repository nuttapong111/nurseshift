package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"database/sql"
	usecase "nurseshift/setting-service/internal/domain/usecases"
	"nurseshift/setting-service/internal/infrastructure/config"
	repoimpl "nurseshift/setting-service/internal/infrastructure/repositories"
	"nurseshift/setting-service/internal/interfaces/http/handlers"
	settingroutes "nurseshift/setting-service/internal/interfaces/http/routes"

	_ "github.com/lib/pq"

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
		AllowOrigins:     strings.Join(cfg.CORS.Origins, ","),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: cfg.CORS.Credentials,
	}))

	if cfg.IsDevelopment() {
		app.Use(logger.New())
	}

	// Initialize repository/usecase (DB connection via DSN in config of auth-service reused here)
	dsn := cfg.GetDatabaseDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("DB open error: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("DB ping error: %v", err)
	}
	repo := repoimpl.NewPostgresSettingRepository(db, cfg.Database.Schema)
	uc := usecase.NewSettingUseCase(repo)
	settingHandler := handlers.NewSettingHandler(uc)

	// Routes
	settingroutes.SetupRoutes(app, settingHandler)

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

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nurseshift/auth-service/internal/domain/usecases"
	"nurseshift/auth-service/internal/infrastructure/config"
	"nurseshift/auth-service/internal/infrastructure/database"
	"nurseshift/auth-service/internal/infrastructure/services"
	httpServer "nurseshift/auth-service/internal/interfaces/http"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Starting Auth Service on port %s...\n", cfg.Server.Port)
	fmt.Printf("Environment: %s\n", cfg.Server.Env)

	// Initialize database connection
	fmt.Println("🔌 Connecting to PostgreSQL database...")
	dbConn, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	fmt.Println("✅ Database connection established successfully")

	// Initialize real services
	passwordService := services.NewPasswordService(cfg.Security.BcryptCost)
	jwtService := services.NewJWTService(
		cfg.JWT.Secret,
		time.Duration(cfg.JWT.ExpireHours)*time.Hour,
		time.Duration(cfg.JWT.RefreshExpireDays)*24*time.Hour,
	)

	// Initialize PostgreSQL repositories
	userRepo := database.NewPostgresUserRepository(dbConn.DB, cfg.Database.Schema)

	// Initialize use cases
	authUseCase := usecases.NewAuthUseCase(
		userRepo,
		jwtService,
		passwordService,
		time.Duration(cfg.Security.SessionTimeoutMins)*time.Minute,
	)

	// Test database connection by getting user by email
	testEmail := "admin@nurseshift.com"
	user, err := userRepo.GetByEmail(context.Background(), testEmail)
	if err != nil {
		log.Fatalf("Failed to get user by email: %v", err)
	}

	fmt.Printf("✅ Database connection test successful! Found user: %s %s (%s)\n",
		user.FirstName, user.LastName, user.Role)

	// Initialize HTTP server
	server := httpServer.NewServer(cfg, authUseCase, jwtService, cfg.JWT.Secret)

	// Start server in a goroutine
	go func() {
		fmt.Printf("🚀 Auth Service started successfully on http://localhost:%s\n", cfg.Server.Port)
		fmt.Println("📚 API Documentation:")
		fmt.Printf("   - Health Check: GET http://localhost:%s/health\n", cfg.Server.Port)
		fmt.Printf("   - Login: POST http://localhost:%s/api/v1/auth/login\n", cfg.Server.Port)
		fmt.Printf("   - Register: POST http://localhost:%s/api/v1/auth/register\n", cfg.Server.Port)
		fmt.Printf("   - Refresh Token: POST http://localhost:%s/api/v1/auth/refresh\n", cfg.Server.Port)
		fmt.Println("🧪 Test Credentials:")
		fmt.Println("   - Admin: admin@nurseshift.com / admin123")
		fmt.Println("   - User: user@nurseshift.com / user123")
		fmt.Println("   - Test: test@nurseshift.com / test123")

		if err := server.Start(); err != nil {
			log.Printf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until we receive our signal
	<-c

	fmt.Println("\n🛑 Shutting down Auth Service...")

	// Shutdown server
	if err := server.Shutdown(); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	fmt.Println("✅ Auth Service stopped gracefully")
}

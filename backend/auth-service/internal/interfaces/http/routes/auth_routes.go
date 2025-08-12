package routes

import (
	"nurseshift/auth-service/internal/domain/usecases"
	"nurseshift/auth-service/internal/infrastructure/config"
	"nurseshift/auth-service/internal/infrastructure/services"
	"nurseshift/auth-service/internal/interfaces/http/handlers"
	"nurseshift/auth-service/internal/interfaces/http/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes sets up authentication-related routes
func SetupAuthRoutes(app *fiber.App, authUseCase usecases.AuthUseCase, jwtService usecases.JWTService, jwtSecret string, cfg *config.Config) {
	// Initialize services based on configuration
	var emailService services.EmailService
	
	if cfg.Email.Provider == "gmail" && cfg.Email.FromEmail != "" && cfg.Email.FromPassword != "" {
		emailService = services.NewGmailEmailService(cfg.Email.FromEmail, cfg.Email.FromPassword)
	} else {
		// Use mock service if email is not configured
		emailService = services.NewMockEmailService()
	}
	
	passwordResetService := services.NewInMemoryPasswordResetService()
	
	authHandler := handlers.NewAuthHandler(authUseCase, jwtService, emailService, passwordResetService)

	// Public routes (no authentication required)
	auth := app.Group("/api/v1/auth")
	{
		auth.Post("/login", authHandler.Login)
		auth.Post("/register", authHandler.Register)
		auth.Post("/refresh", authHandler.RefreshToken)
		auth.Get("/health", authHandler.Health)
		auth.Post("/introspect", authHandler.VerifyToken)
		auth.Post("/forgot-password", authHandler.ForgotPassword)
		auth.Post("/reset-password", authHandler.ResetPassword)
	}

	// Protected routes (authentication required)
	protected := app.Group("/api/v1/auth", middleware.AuthMiddleware(jwtSecret))
	{
		protected.Post("/logout", authHandler.Logout)
		protected.Post("/logout-all", authHandler.LogoutAll)
		protected.Post("/change-password", authHandler.ChangePassword)
		protected.Get("/me", authHandler.Me)

		// Admin only routes
		protected.Post("/create-admin", authHandler.CreateAdmin) // เฉพาะ admin เท่านั้น
	}
}

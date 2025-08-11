package routes

import (
	"nurseshift/auth-service/internal/domain/usecases"
	"nurseshift/auth-service/internal/interfaces/http/handlers"
	"nurseshift/auth-service/internal/interfaces/http/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes sets up authentication-related routes
func SetupAuthRoutes(app *fiber.App, authUseCase usecases.AuthUseCase, _ usecases.JWTService) {
	authHandler := handlers.NewAuthHandler(authUseCase)

	// Public routes (no authentication required)
	auth := app.Group("/api/v1/auth")
	{
		auth.Post("/login", authHandler.Login)
		auth.Post("/register", authHandler.Register)
		auth.Post("/refresh", authHandler.RefreshToken)
		auth.Get("/health", authHandler.Health)
	}

	// Protected routes (authentication required)
	protected := app.Group("/api/v1/auth", middleware.AuthMiddleware())
	{
		protected.Post("/logout", authHandler.Logout)
		protected.Post("/logout-all", authHandler.LogoutAll)
		protected.Post("/change-password", authHandler.ChangePassword)
		protected.Get("/me", authHandler.Me)
	}
}

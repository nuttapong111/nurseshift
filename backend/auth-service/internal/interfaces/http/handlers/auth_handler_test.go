package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"nurseshift/auth-service/internal/domain/entities"
	"nurseshift/auth-service/internal/domain/usecases"
	"nurseshift/auth-service/internal/infrastructure/repositories"
	"nurseshift/auth-service/internal/infrastructure/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTestAuthHandler() (*AuthHandler, *repositories.MockUserRepository, *services.MockJWTService, *services.MockPasswordService) {
	mockUserRepo := repositories.NewMockUserRepository()
	mockJWTService := services.NewMockJWTService()
	mockPasswordService := services.NewMockPasswordService()

	// Add test users
	adminUser := &entities.User{
		ID:             uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
		Email:          "admin@nurseshift.com",
		PasswordHash:   "mock-hash-admin123",
		FirstName:      "Admin",
		LastName:       "System",
		Role:           entities.RoleAdmin,
		Status:         entities.StatusActive,
		DaysRemaining:  90,
		PackageType:    "enterprise",
		MaxDepartments: 20,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	regularUser := &entities.User{
		ID:             uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
		Email:          "user@nurseshift.com",
		PasswordHash:   "mock-hash-user123",
		FirstName:      "User",
		LastName:       "Demo",
		Role:           entities.RoleUser,
		Status:         entities.StatusActive,
		DaysRemaining:  30,
		PackageType:    "trial",
		MaxDepartments: 2,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mockUserRepo.AddMockUser(adminUser)
	mockUserRepo.AddMockUser(regularUser)

	// Set mock passwords
	mockPasswordService.SetMockPassword("admin@nurseshift.com", "admin123")
	mockPasswordService.SetMockPassword("user@nurseshift.com", "user123")

	authUseCase := usecases.NewAuthUseCase(mockUserRepo, mockJWTService, mockPasswordService, 30*time.Minute)
	handler := NewAuthHandler(authUseCase, mockJWTService)

	return handler, mockUserRepo, mockJWTService, mockPasswordService
}

func TestLogin(t *testing.T) {
	app := fiber.New()
	handler, _, _, _ := setupTestAuthHandler()

	app.Post("/login", handler.Login)

	tests := []struct {
		name           string
		requestBody    LoginRequest
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Valid admin login",
			requestBody: LoginRequest{
				Email:    "admin@nurseshift.com",
				Password: "admin123",
			},
			expectedStatus: fiber.StatusOK,
			expectedError:  false,
		},
		{
			name: "Valid user login",
			requestBody: LoginRequest{
				Email:    "user@nurseshift.com",
				Password: "user123",
			},
			expectedStatus: fiber.StatusOK,
			expectedError:  false,
		},
		{
			name: "Invalid credentials",
			requestBody: LoginRequest{
				Email:    "admin@nurseshift.com",
				Password: "wrongpassword",
			},
			expectedStatus: fiber.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name: "Non-existent user",
			requestBody: LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			expectedStatus: fiber.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError {
				var response LoginResponse
				json.NewDecoder(resp.Body).Decode(&response)
				assert.Equal(t, "success", response.Status)
				assert.NotEmpty(t, response.AccessToken)
				assert.NotEmpty(t, response.RefreshToken)
				assert.NotNil(t, response.User)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	app := fiber.New()
	handler, _, _, _ := setupTestAuthHandler()

	app.Post("/register", handler.Register)

	tests := []struct {
		name           string
		requestBody    RegisterRequest
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Valid registration",
			requestBody: RegisterRequest{
				Email:     "newuser@example.com",
				Password:  "password123",
				FirstName: "New",
				LastName:  "User",
			},
			expectedStatus: fiber.StatusCreated,
			expectedError:  false,
		},
		{
			name: "Missing required fields",
			requestBody: RegisterRequest{
				Email:    "incomplete@example.com",
				Password: "password123",
				// Missing FirstName and LastName
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError {
				var response map[string]interface{}
				json.NewDecoder(resp.Body).Decode(&response)
				assert.Equal(t, "success", response["status"])
				assert.NotNil(t, response["user"])
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	app := fiber.New()
	handler, _, mockJWTService, _ := setupTestAuthHandler()

	app.Post("/verify", handler.VerifyToken)

	// Create a valid token
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	token, _ := mockJWTService.GenerateAccessToken(userID, entities.RoleAdmin)

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "Valid token",
			token:          token,
			expectedStatus: fiber.StatusOK,
			expectedError:  false,
		},
		{
			name:           "Invalid token",
			token:          "invalid-token",
			expectedStatus: fiber.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/verify", nil)
			req.Header.Set("Authorization", "Bearer "+tt.token)

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectedError {
				var response map[string]interface{}
				json.NewDecoder(resp.Body).Decode(&response)
				assert.Equal(t, "success", response["status"])
				assert.Equal(t, true, response["active"])
			}
		})
	}
}

func TestHealth(t *testing.T) {
	app := fiber.New()
	handler, _, _, _ := setupTestAuthHandler()

	app.Get("/health", handler.Health)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Auth Service is healthy", response["message"])
}

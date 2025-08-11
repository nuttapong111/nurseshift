package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the JWT claims structure
type JWTClaims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT tokens
func AuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "ไม่พบ Authorization header",
			})
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "รูปแบบ Authorization header ไม่ถูกต้อง",
			})
		}

		tokenString := tokenParts[1]

		// For mock tokens, parse userID from token string
		if strings.HasPrefix(tokenString, "mock-jwt-access-token-") {
			userID := strings.TrimPrefix(tokenString, "mock-jwt-access-token-")

			// Set mock user data based on userID
			if userID == "admin-1" {
				c.Locals("userID", "admin-1")
				c.Locals("email", "admin@nurseshift.com")
				c.Locals("role", "admin")
			} else if userID == "user-1" {
				c.Locals("userID", "user-1")
				c.Locals("email", "user@nurseshift.com")
				c.Locals("role", "user")
			} else {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"status":  "error",
					"message": "Token ไม่ถูกต้อง",
				})
			}

			return c.Next()
		}

		// Parse and validate real JWT token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Token ไม่ถูกต้องหรือหมดอายุ",
				"error":   err.Error(),
			})
		}

		// Extract claims
		if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
			// Set user information in context
			c.Locals("userID", claims.UserID)
			c.Locals("email", claims.Email)
			c.Locals("role", claims.Role)

			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Token claims ไม่ถูกต้อง",
		})
	}
}

// AdminOnlyMiddleware ensures only admin users can access
func AdminOnlyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "ต้องมีสิทธิ์ admin เท่านั้น",
			})
		}
		return c.Next()
	}
}

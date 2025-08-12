package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Security SecurityConfig
	CORS     CORSConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
	Env  string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	Schema   string
}

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret            string
	ExpireHours       int
	RefreshExpireDays int
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	BcryptCost         int
	SessionTimeoutMins int
	RateLimitMax       int
	RateLimitWindowMin int
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	Origins     []string
	Credentials bool
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load("config.env"); err != nil {
		// It's okay if the file doesn't exist in production
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8081"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "123456"),
			Name:     getEnv("DB_NAME", "nurseshift"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			Schema:   getEnv("DB_SCHEMA", "nurse_shift"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET", "nurseshift-department-service-secret-key-2024"),
			ExpireHours:       getEnvAsInt("JWT_EXPIRE_HOURS", 24),
			RefreshExpireDays: getEnvAsInt("JWT_REFRESH_EXPIRE_DAYS", 7),
		},
		Security: SecurityConfig{
			BcryptCost:         getEnvAsInt("BCRYPT_COST", 12),
			SessionTimeoutMins: getEnvAsInt("SESSION_TIMEOUT_MINUTES", 30),
			RateLimitMax:       getEnvAsInt("RATE_LIMIT_MAX", 100),
			RateLimitWindowMin: getEnvAsInt("RATE_LIMIT_WINDOW_MINUTES", 15),
		},
		CORS: CORSConfig{
			Origins:     strings.Split(getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:3002"), ","),
			Credentials: getEnvAsBool("CORS_CREDENTIALS", true),
		},
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// validate validates the configuration
func (c *Config) validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}
	return nil
}

// GetDatabaseDSN returns the database connection string
func (c *Config) GetDatabaseDSN() string {
	if c.Database.Password == "" {
		return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s",
			c.Database.Host, c.Database.Port, c.Database.User, c.Database.Name, c.Database.SSLMode)
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password,
		c.Database.Name, c.Database.SSLMode)
}

// GetRedisAddr returns the Redis connection address
func (c *Config) GetRedisAddr() string {
	return c.Redis.Host + ":" + c.Redis.Port
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Server.Env == "development"
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

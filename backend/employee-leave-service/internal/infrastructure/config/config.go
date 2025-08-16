package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration for the service
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	CORS     CORSConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Env  string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Schema   string
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	ServiceURL  string
	JWTSecret   string
	ExpireHours int
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	Origins     []string
	Credentials bool
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8088"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "nuttapong2"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "nurseshift"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			Schema:   getEnv("DB_SCHEMA", "nurse_shift"),
		},
		Auth: AuthConfig{
			ServiceURL:  getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
			JWTSecret:   getEnv("JWT_SECRET", "nurseshift-super-secret-jwt-key-development-only-2024"),
			ExpireHours: getEnvAsInt("JWT_EXPIRE_HOURS", 24),
		},
		CORS: CORSConfig{
			Origins:     strings.Split(getEnv("CORS_ORIGINS", "*"), ","),
			Credentials: getEnvAsBool("CORS_CREDENTIALS", true),
		},
	}

	return cfg, nil
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Server.Env == "development"
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

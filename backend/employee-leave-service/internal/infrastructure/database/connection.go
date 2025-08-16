package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// Connection holds the database connection
type Connection struct {
	DB *sql.DB
}

// NewConnection creates a new database connection
func NewConnection() (*Connection, error) {
	// Load environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")
	dbSchema := os.Getenv("DB_SCHEMA")
	databaseURL := os.Getenv("DATABASE_URL")

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbName == "" {
		dbName = "nurseshift"
	}
	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}
	if dbSchema == "" {
		dbSchema = "nurse_shift"
	}

	log.Printf("Connecting to database: %s:%s/%s (schema: %s)", dbHost, dbPort, dbName, dbSchema)

	// Prefer DATABASE_URL if provided (Railway often provides this)
	var dsn string
	if strings.TrimSpace(databaseURL) != "" {
		dsn = databaseURL
		log.Printf("Using DATABASE_URL for connection")
	} else {
		// Build URL DSN to ensure database selection is respected
		userEscaped := url.QueryEscape(dbUser)
		passEscaped := url.QueryEscape(dbPassword)
		hostPort := fmt.Sprintf("%s:%s", dbHost, dbPort)
		dsn = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", userEscaped, passEscaped, hostPort, dbName, dbSSLMode)
	}

	// Open database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Validate connected database (warn only; do not fail in production)
	var curDB string
	if err := db.QueryRow("select current_database()").Scan(&curDB); err == nil {
		if dbName != "" && curDB != dbName {
			log.Printf("Warning: connected to database '%s' while DB_NAME is '%s'", curDB, dbName)
		}
	} else {
		log.Printf("Warning: failed to read current_database(): %v", err)
	}

	// Set search path explicitly
	log.Printf("Setting search_path to: %s, public", dbSchema)
	if _, err := db.Exec(fmt.Sprintf("SET search_path TO %s, public", dbSchema)); err != nil {
		log.Printf("Warning: failed to set search path: %v", err)
		// Continue anyway, we'll use fully qualified names
	} else {
		log.Printf("✅ Search path set successfully to: %s, public", dbSchema)
	}

	// Diagnostics
	var curUser, relExists string
	if err := db.QueryRow("select current_user, coalesce(to_regclass('nurse_shift.leave_requests')::text, 'NULL')").Scan(&curUser, &relExists); err == nil {
		log.Printf("Diagnostics: current_user=%s relation=nurse_shift.leave_requests=%s", curUser, relExists)
	} else {
		log.Printf("Diagnostics query failed: %v", err)
	}

	log.Println("✅ Database connection established successfully")

	return &Connection{DB: db}, nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

// GetDB returns the underlying sql.DB instance
func (c *Connection) GetDB() *sql.DB {
	return c.DB
}

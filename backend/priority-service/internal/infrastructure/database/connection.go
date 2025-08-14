package database

import (
	"database/sql"
	"fmt"
	"time"

	"nurseshift/priority-service/internal/infrastructure/config"

	_ "github.com/lib/pq"
)

// Connection holds database connection and configuration
type Connection struct {
	DB     *sql.DB
	Config *config.Config
}

// NewConnection creates a new database connection
func NewConnection(cfg *config.Config) (*Connection, error) {
	db, err := sql.Open("postgres", cfg.GetDatabaseDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Connection{
		DB:     db,
		Config: cfg,
	}, nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

// Health checks the database connection health
func (c *Connection) Health() error {
	return c.DB.Ping()
}

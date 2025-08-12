package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"nurseshift/auth-service/internal/domain/entities"
	"nurseshift/auth-service/internal/domain/repositories"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// PostgresUserRepository implements UserRepository using PostgreSQL
type PostgresUserRepository struct {
	db     *sql.DB
	schema string
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *sql.DB, schema string) repositories.UserRepository {
	return &PostgresUserRepository{
		db:     db,
		schema: schema,
	}
}

// GetByID retrieves a user by ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	query := fmt.Sprintf(`
		SELECT id, email, password_hash, first_name, last_name, phone, role, status, position,
			   days_remaining, subscription_expires_at, package_type, max_departments, avatar_url,
			   settings, last_login_at, created_at, updated_at
		FROM %s.users 
		WHERE id = $1 AND status = 'active'`, r.schema)

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Phone, &user.Role, &user.Status, &user.Position, &user.DaysRemaining,
		&user.SubscriptionExpiresAt, &user.PackageType, &user.MaxDepartments,
		&user.AvatarURL, &user.Settings, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := fmt.Sprintf(`
		SELECT id, email, password_hash, first_name, last_name, phone, role, status, position,
			   days_remaining, subscription_expires_at, package_type, max_departments, avatar_url,
			   settings, last_login_at, created_at, updated_at
		FROM %s.users 
		WHERE email = $1`, r.schema)

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Phone, &user.Role, &user.Status, &user.Position, &user.DaysRemaining,
		&user.SubscriptionExpiresAt, &user.PackageType, &user.MaxDepartments,
		&user.AvatarURL, &user.Settings, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// Create creates a new user
func (r *PostgresUserRepository) Create(ctx context.Context, user *entities.User) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.users (
			id, email, password_hash, first_name, last_name, phone, role, status, position,
			days_remaining, subscription_expires_at, package_type, max_departments, avatar_url,
			settings, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)`, r.schema)

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.Phone, user.Role, user.Status, user.Position, user.DaysRemaining,
		user.SubscriptionExpiresAt, user.PackageType, user.MaxDepartments,
		user.AvatarURL, user.Settings, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Update updates an existing user
func (r *PostgresUserRepository) Update(ctx context.Context, user *entities.User) error {
	query := fmt.Sprintf(`
		UPDATE %s.users SET 
			email = $2, first_name = $3, last_name = $4, phone = $5, role = $6, status = $7,
			position = $8, days_remaining = $9, subscription_expires_at = $10, package_type = $11,
			max_departments = $12, avatar_url = $13, settings = $14, last_login_at = $15,
			updated_at = $16
		WHERE id = $1`, r.schema)

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.FirstName, user.LastName, user.Phone, user.Role, user.Status,
		user.Position, user.DaysRemaining, user.SubscriptionExpiresAt, user.PackageType,
		user.MaxDepartments, user.AvatarURL, user.Settings, user.LastLoginAt, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// UpdateLastLogin updates the user's last login time
func (r *PostgresUserRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.users SET last_login_at = $2, updated_at = $3
		WHERE id = $1`, r.schema)

	_, err := r.db.ExecContext(ctx, query, userID, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// UpdatePassword updates the user's password
func (r *PostgresUserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error {
	query := fmt.Sprintf(`
		UPDATE %s.users SET password_hash = $2, updated_at = $3
		WHERE id = $1`, r.schema)

	_, err := r.db.ExecContext(ctx, query, userID, hashedPassword, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// EmailExists checks if an email already exists
func (r *PostgresUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS(SELECT 1 FROM %s.users WHERE email = $1)`, r.schema)

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

// CreateSession creates a new user session
func (r *PostgresUserRepository) CreateSession(ctx context.Context, session *entities.UserSession) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.user_sessions (
			id, user_id, token_hash, ip_address, user_agent, created_at, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`, r.schema)

	_, err := r.db.ExecContext(ctx, query,
		session.ID, session.UserID, session.TokenHash, session.IPAddress,
		session.UserAgent, session.CreatedAt, session.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// GetSessionByToken retrieves a session by token hash
func (r *PostgresUserRepository) GetSessionByToken(ctx context.Context, tokenHash string) (*entities.UserSession, error) {
	query := fmt.Sprintf(`
		SELECT id, user_id, token_hash, ip_address, user_agent, created_at, expires_at, revoked_at
		FROM %s.user_sessions 
		WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > $2`, r.schema)

	session := &entities.UserSession{}
	err := r.db.QueryRowContext(ctx, query, tokenHash, time.Now()).Scan(
		&session.ID, &session.UserID, &session.TokenHash, &session.IPAddress,
		&session.UserAgent, &session.CreatedAt, &session.ExpiresAt, &session.RevokedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found or expired")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

// RevokeSession revokes a specific session
func (r *PostgresUserRepository) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.user_sessions SET revoked_at = $2
		WHERE id = $1`, r.schema)

	_, err := r.db.ExecContext(ctx, query, sessionID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	return nil
}

// RevokeAllUserSessions revokes all sessions for a specific user
func (r *PostgresUserRepository) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.user_sessions SET revoked_at = $2
		WHERE user_id = $1 AND revoked_at IS NULL`, r.schema)

	_, err := r.db.ExecContext(ctx, query, userID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to revoke all user sessions: %w", err)
	}

	return nil
}

// CleanExpiredSessions removes expired sessions
func (r *PostgresUserRepository) CleanExpiredSessions(ctx context.Context) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.user_sessions 
		WHERE expires_at < $1`, r.schema)

	_, err := r.db.ExecContext(ctx, query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to clean expired sessions: %w", err)
	}

	return nil
}

// Delete deletes a user (soft delete by setting status to inactive)
func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.users SET status = 'inactive', updated_at = $2
		WHERE id = $1`, r.schema)

	_, err := r.db.ExecContext(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// GetByEmployeeID retrieves a user by employee ID (not implemented in current schema)
func (r *PostgresUserRepository) GetByEmployeeID(ctx context.Context, employeeID string) (*entities.User, error) {
	return nil, fmt.Errorf("employee ID not supported in current schema")
}

// EmployeeIDExists checks if an employee ID exists (not implemented in current schema)
func (r *PostgresUserRepository) EmployeeIDExists(ctx context.Context, employeeID string, organizationID uuid.UUID) (bool, error) {
	return false, fmt.Errorf("employee ID not supported in current schema")
}

// GetOrganizationByID retrieves an organization by ID (not implemented in current schema)
func (r *PostgresUserRepository) GetOrganizationByID(ctx context.Context, id uuid.UUID) (*entities.Organization, error) {
	return nil, fmt.Errorf("organizations not supported in current schema")
}

// GetOrganizationByUser retrieves an organization by user (not implemented in current schema)
func (r *PostgresUserRepository) GetOrganizationByUser(ctx context.Context, userID uuid.UUID) (*entities.Organization, error) {
	return nil, fmt.Errorf("organizations not supported in current schema")
}

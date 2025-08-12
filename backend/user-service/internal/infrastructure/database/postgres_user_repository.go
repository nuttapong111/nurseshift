package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"nurseshift/user-service/internal/domain/entities"
	"nurseshift/user-service/internal/domain/repositories"
	"nurseshift/user-service/internal/infrastructure/services"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// UserRepository interface for user operations
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
	UpdateLastLogin(ctx context.Context, userID string) error
	GetUsers(ctx context.Context, req *repositories.GetUsersRequest) (*repositories.GetUsersResponse, error)
	SearchUsers(ctx context.Context, req *repositories.SearchUsersRequest) (*repositories.SearchUsersResponse, error)
	GetUserStats(ctx context.Context) (*repositories.UserStatsResponse, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	SendVerificationEmail(ctx context.Context, email string) error
	VerifyEmail(ctx context.Context, token string) error
	IsEmailVerified(ctx context.Context, email string) (bool, error)
}

// PostgresUserRepository implements UserRepository using PostgreSQL
type PostgresUserRepository struct {
	db     *sql.DB
	schema string
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *sql.DB, schema string) UserRepository {
	return &PostgresUserRepository{
		db:     db,
		schema: schema,
	}
}

// GetByID retrieves a user by ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT id, email, password_hash, first_name, last_name, phone, role, status, position,
			   days_remaining, subscription_expires_at, package_type, max_departments, avatar_url,
			   settings, last_login_at, created_at, updated_at, email_verified, email_verification_token, email_verification_expires_at
		FROM %s.users 
		WHERE id = $1 AND status = 'active'`, r.schema)

	user := &entities.User{}
	err = r.db.QueryRowContext(ctx, query, userUUID).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Phone, &user.Role, &user.Status, &user.Position, &user.DaysRemaining,
		&user.SubscriptionExpiresAt, &user.PackageType, &user.MaxDepartments,
		&user.AvatarURL, &user.Settings, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
		&user.EmailVerified, &user.EmailVerificationToken, &user.EmailVerificationExpiresAt,
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
			   settings, last_login_at, created_at, updated_at, email_verified, email_verification_token, email_verification_expires_at
		FROM %s.users 
		WHERE email = $1`, r.schema)

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Phone, &user.Role, &user.Status, &user.Position, &user.DaysRemaining,
		&user.SubscriptionExpiresAt, &user.PackageType, &user.MaxDepartments,
		&user.AvatarURL, &user.Settings, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
		&user.EmailVerified, &user.EmailVerificationToken, &user.EmailVerificationExpiresAt,
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
		INSERT INTO %s.users (id, email, password_hash, first_name, last_name, phone, 
							  role, status, position, days_remaining, package_type, max_departments)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`, r.schema)

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.Phone, user.Role, user.Status, user.Position, user.DaysRemaining,
		user.PackageType, user.MaxDepartments,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Update updates an existing user
func (r *PostgresUserRepository) Update(ctx context.Context, user *entities.User) error {
	query := fmt.Sprintf(`
		UPDATE %s.users 
		SET first_name = $2, last_name = $3, phone = $4, position = $5, 
			avatar_url = $6, updated_at = $7
		WHERE id = $1`, r.schema)

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.FirstName, user.LastName, user.Phone, user.Position,
		user.AvatarURL, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// UpdateLastLogin updates the user's last login timestamp
func (r *PostgresUserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	query := fmt.Sprintf(`
		UPDATE %s.users 
		SET last_login_at = $2 
		WHERE id = $1`, r.schema)

	now := time.Now()
	_, err = r.db.ExecContext(ctx, query, userUUID, now)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// GetUsers returns paginated list of users
func (r *PostgresUserRepository) GetUsers(ctx context.Context, req *repositories.GetUsersRequest) (*repositories.GetUsersResponse, error) {
	// Build WHERE clause
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if req.Role != nil {
		whereClause += fmt.Sprintf(" AND role = $%d", argCount)
		args = append(args, *req.Role)
		argCount++
	}

	if req.Status != nil {
		whereClause += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *req.Status)
		argCount++
	}

	// Count total users
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s.users %s`, r.schema, whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Get users with pagination
	offset := (req.Page - 1) * req.Limit
	query := fmt.Sprintf(`
		SELECT id, email, first_name, last_name, phone, role, status, position,
			   days_remaining, subscription_expires_at, package_type, max_departments, avatar_url,
			   settings, last_login_at, created_at, updated_at
		FROM %s.users %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, r.schema, whereClause, argCount, argCount+1)

	args = append(args, req.Limit, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone,
			&user.Role, &user.Status, &user.Position, &user.DaysRemaining,
			&user.SubscriptionExpiresAt, &user.PackageType, &user.MaxDepartments,
			&user.AvatarURL, &user.Settings, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	totalPages := (total + req.Limit - 1) / req.Limit

	return &repositories.GetUsersResponse{
		Users:      users,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// SearchUsers searches users by query
func (r *PostgresUserRepository) SearchUsers(ctx context.Context, req *repositories.SearchUsersRequest) (*repositories.SearchUsersResponse, error) {
	// Build WHERE clause
	whereClause := "WHERE (first_name ILIKE $1 OR last_name ILIKE $1 OR email ILIKE $1)"
	args := []interface{}{"%" + req.Query + "%"}
	argCount := 2

	if req.Role != nil {
		whereClause += fmt.Sprintf(" AND role = $%d", argCount)
		args = append(args, *req.Role)
		argCount++
	}

	if req.Status != nil {
		whereClause += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *req.Status)
		argCount++
	}

	// Count total users
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s.users %s`, r.schema, whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Get users with pagination
	offset := (req.Page - 1) * req.Limit
	query := fmt.Sprintf(`
		SELECT id, email, first_name, last_name, phone, role, status, position,
			   days_remaining, subscription_expires_at, package_type, max_departments, avatar_url,
			   settings, last_login_at, created_at, updated_at
		FROM %s.users %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, r.schema, whereClause, argCount, argCount+1)

	args = append(args, req.Limit, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone,
			&user.Role, &user.Status, &user.Position, &user.DaysRemaining,
			&user.SubscriptionExpiresAt, &user.PackageType, &user.MaxDepartments,
			&user.AvatarURL, &user.Settings, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	totalPages := (total + req.Limit - 1) / req.Limit

	return &repositories.SearchUsersResponse{
		Users:      users,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetUserStats returns user statistics
func (r *PostgresUserRepository) GetUserStats(ctx context.Context) (*repositories.UserStatsResponse, error) {
	query := fmt.Sprintf(`
		SELECT 
			COUNT(*) as total_users,
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active_users,
			COUNT(CASE WHEN status != 'active' THEN 1 END) as inactive_users,
			COUNT(CASE WHEN role = 'admin' THEN 1 END) as admin_count,
			COUNT(CASE WHEN role = 'user' THEN 1 END) as user_count
		FROM %s.users`, r.schema)

	var stats repositories.UserStatsResponse
	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalUsers, &stats.ActiveUsers, &stats.InactiveUsers,
		&stats.AdminCount, &stats.UserCount,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	return &stats, nil
}

// EmailExists checks if an email already exists
func (r *PostgresUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s.users WHERE email = $1)`, r.schema)

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

// Email verification methods
func (r *PostgresUserRepository) SendVerificationEmail(ctx context.Context, email string) error {
	// Check if user exists
	var userID string
	query := fmt.Sprintf("SELECT id FROM %s.users WHERE email = $1", r.schema)
	err := r.db.QueryRowContext(ctx, query, email).Scan(&userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Generate verification token (in real app, use crypto/rand)
	token := fmt.Sprintf("verify_%s_%d", userID, time.Now().Unix())

	// Store verification token in database
	updateQuery := fmt.Sprintf(`
		UPDATE %s.users 
		SET email_verification_token = $1, 
		    email_verification_expires_at = $2,
		    updated_at = $3
		WHERE id = $4
	`, r.schema)

	_, err = r.db.ExecContext(ctx, updateQuery, token, time.Now().Add(24*time.Hour), time.Now(), userID)

	if err != nil {
		return fmt.Errorf("failed to store verification token: %w", err)
	}

	// Send email via email service
	emailService := services.NewSMTPEmailService()
	err = emailService.SendVerificationEmail(email, token)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) VerifyEmail(ctx context.Context, token string) error {
	// In a real implementation, this would:
	// 1. Find the token in the database
	// 2. Check if it's expired
	// 3. Mark the user's email as verified
	// 4. Remove the token

	var userID string
	var expiresAt time.Time

	selectQuery := fmt.Sprintf(`
		SELECT id, email_verification_expires_at 
		FROM %s.users 
		WHERE email_verification_token = $1
	`, r.schema)

	err := r.db.QueryRowContext(ctx, selectQuery, token).Scan(&userID, &expiresAt)

	if err != nil {
		return fmt.Errorf("invalid verification token: %w", err)
	}

	if time.Now().After(expiresAt) {
		return fmt.Errorf("verification token expired")
	}

	// Mark email as verified
	updateQuery := fmt.Sprintf(`
		UPDATE %s.users 
		SET email_verified = true,
		    email_verification_token = NULL,
		    email_verification_expires_at = NULL,
		    updated_at = $1
		WHERE id = $2
	`, r.schema)

	_, err = r.db.ExecContext(ctx, updateQuery, time.Now(), userID)

	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) IsEmailVerified(ctx context.Context, email string) (bool, error) {
	var verified bool
	query := fmt.Sprintf(`
		SELECT COALESCE(email_verified, false) 
		FROM %s.users 
		WHERE email = $1
	`, r.schema)

	err := r.db.QueryRowContext(ctx, query, email).Scan(&verified)

	if err != nil {
		return false, fmt.Errorf("failed to check email verification status: %w", err)
	}

	return verified, nil
}

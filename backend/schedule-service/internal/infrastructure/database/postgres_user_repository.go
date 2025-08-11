package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// UserRole represents user roles
type UserRole string

const (
	RoleAdmin        UserRole = "admin"
	RoleManager      UserRole = "manager"
	RoleNurse        UserRole = "nurse"
	RoleSupervisor   UserRole = "supervisor"
	RoleCoordinator  UserRole = "coordinator"
	RoleAdministrator UserRole = "administrator"
)

// UserStatus represents user status
type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
	StatusSuspended UserStatus = "suspended"
)

// User represents a user entity
type User struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organizationId"`
	EmployeeID     string     `json:"employeeId"`
	Email          string     `json:"email"`
	PasswordHash   string     `json:"-"`
	FirstName      string     `json:"firstName"`
	LastName       string     `json:"lastName"`
	Phone          string     `json:"phone"`
	Role           UserRole   `json:"role"`
	Status         UserStatus `json:"status"`
	Position       string     `json:"position"`
	DateJoined     time.Time  `json:"dateJoined"`
	DateOfBirth    *time.Time `json:"dateOfBirth"`
	AvatarURL      *string    `json:"avatarUrl"`
	LastLoginAt    *time.Time `json:"lastLoginAt"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

// Organization represents an organization entity
type Organization struct {
	ID                    uuid.UUID  `json:"id"`
	Name                  string     `json:"name"`
	Description           *string    `json:"description"`
	Email                 string     `json:"email"`
	Phone                 *string    `json:"phone"`
	Address               *string    `json:"address"`
	Website               *string    `json:"website"`
	LicenseNumber         string     `json:"licenseNumber"`
	SubscriptionExpiresAt *time.Time `json:"subscriptionExpiresAt"`
	PackageType           string     `json:"packageType"`
	MaxUsers              int        `json:"maxUsers"`
	MaxDepartments        int        `json:"maxDepartments"`
	CreatedAt             time.Time  `json:"createdAt"`
	UpdatedAt             time.Time  `json:"updatedAt"`
}

// UserSession represents a user session entity
type UserSession struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"userId"`
	TokenHash  string     `json:"tokenHash"`
	IPAddress  string     `json:"ipAddress"`
	UserAgent  string     `json:"userAgent"`
	CreatedAt  time.Time  `json:"createdAt"`
	ExpiresAt  time.Time  `json:"expiresAt"`
	RevokedAt  *time.Time `json:"revokedAt"`
}

// UserRepository interface for user operations
type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByEmployeeID(ctx context.Context, employeeID string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	GetOrganizationByID(ctx context.Context, id uuid.UUID) (*Organization, error)
	GetOrganizationByUser(ctx context.Context, userID uuid.UUID) (*Organization, error)
	CreateSession(ctx context.Context, session *UserSession) error
	GetSessionByToken(ctx context.Context, tokenHash string) (*UserSession, error)
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
func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := fmt.Sprintf(`
		SELECT id, organization_id, employee_id, email, password_hash, first_name, last_name,
			   phone, role, status, position, date_joined, date_of_birth, avatar_url,
			   last_login_at, created_at, updated_at
		FROM %s.users 
		WHERE id = $1 AND status = 'active'`, r.schema)

	user := &User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.OrganizationID, &user.EmployeeID, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.Phone, &user.Role, &user.Status,
		&user.Position, &user.DateJoined, &user.DateOfBirth, &user.AvatarURL,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
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
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := fmt.Sprintf(`
		SELECT id, organization_id, employee_id, email, password_hash, first_name, last_name,
			   phone, role, status, position, date_joined, date_of_birth, avatar_url,
			   last_login_at, created_at, updated_at
		FROM %s.users 
		WHERE email = $1`, r.schema)

	user := &User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.OrganizationID, &user.EmployeeID, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.Phone, &user.Role, &user.Status,
		&user.Position, &user.DateJoined, &user.DateOfBirth, &user.AvatarURL,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// GetByEmployeeID retrieves a user by employee ID
func (r *PostgresUserRepository) GetByEmployeeID(ctx context.Context, employeeID string) (*User, error) {
	query := fmt.Sprintf(`
		SELECT id, organization_id, employee_id, email, password_hash, first_name, last_name,
			   phone, role, status, position, date_joined, date_of_birth, avatar_url,
			   last_login_at, created_at, updated_at
		FROM %s.users 
		WHERE employee_id = $1`, r.schema)

	user := &User{}
	err := r.db.QueryRowContext(ctx, query, employeeID).Scan(
		&user.ID, &user.OrganizationID, &user.EmployeeID, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.Phone, &user.Role, &user.Status,
		&user.Position, &user.DateJoined, &user.DateOfBirth, &user.AvatarURL,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by employee ID: %w", err)
	}

	return user, nil
}

// Create creates a new user
func (r *PostgresUserRepository) Create(ctx context.Context, user *User) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.users (id, organization_id, employee_id, email, password_hash, 
							  first_name, last_name, phone, role, status, position, date_joined)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`, r.schema)

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.OrganizationID, user.EmployeeID, user.Email, user.PasswordHash,
		user.FirstName, user.LastName, user.Phone, user.Role, user.Status,
		user.Position, user.DateJoined,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Update updates an existing user
func (r *PostgresUserRepository) Update(ctx context.Context, user *User) error {
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
func (r *PostgresUserRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.users 
		SET last_login_at = $2 
		WHERE id = $1`, r.schema)

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, userID, now)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// GetOrganizationByID retrieves an organization by ID
func (r *PostgresUserRepository) GetOrganizationByID(ctx context.Context, id uuid.UUID) (*Organization, error) {
	query := fmt.Sprintf(`
		SELECT id, name, description, email, phone, address, website, license_number,
			   subscription_expires_at, package_type, max_users, max_departments,
			   created_at, updated_at
		FROM %s.organizations 
		WHERE id = $1`, r.schema)

	org := &Organization{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&org.ID, &org.Name, &org.Description, &org.Email, &org.Phone, &org.Address,
		&org.Website, &org.LicenseNumber, &org.SubscriptionExpiresAt, &org.PackageType,
		&org.MaxUsers, &org.MaxDepartments, &org.CreatedAt, &org.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	return org, nil
}

// GetOrganizationByUser retrieves the organization for a specific user
func (r *PostgresUserRepository) GetOrganizationByUser(ctx context.Context, userID uuid.UUID) (*Organization, error) {
	query := fmt.Sprintf(`
		SELECT o.id, o.name, o.description, o.email, o.phone, o.address, o.website, 
			   o.license_number, o.subscription_expires_at, o.package_type, o.max_users, 
			   o.max_departments, o.created_at, o.updated_at
		FROM %s.organizations o
		JOIN %s.users u ON u.organization_id = o.id
		WHERE u.id = $1`, r.schema, r.schema)

	org := &Organization{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&org.ID, &org.Name, &org.Description, &org.Email, &org.Phone, &org.Address,
		&org.Website, &org.LicenseNumber, &org.SubscriptionExpiresAt, &org.PackageType,
		&org.MaxUsers, &org.MaxDepartments, &org.CreatedAt, &org.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("organization not found for user")
		}
		return nil, fmt.Errorf("failed to get organization by user: %w", err)
	}

	return org, nil
}

// CreateSession creates a new user session
func (r *PostgresUserRepository) CreateSession(ctx context.Context, session *UserSession) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.user_sessions (id, user_id, token_hash, ip_address, user_agent, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)`, r.schema)

	_, err := r.db.ExecContext(ctx, query,
		session.ID, session.UserID, session.TokenHash, session.IPAddress, session.UserAgent, session.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// GetSessionByToken retrieves a session by token hash
func (r *PostgresUserRepository) GetSessionByToken(ctx context.Context, tokenHash string) (*UserSession, error) {
	query := fmt.Sprintf(`
		SELECT id, user_id, token_hash, ip_address, user_agent, created_at, expires_at, revoked_at
		FROM %s.user_sessions 
		WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW()`, r.schema)

	session := &UserSession{}
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
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
		UPDATE %s.user_sessions 
		SET revoked_at = NOW() 
		WHERE id = $1`, r.schema)

	_, err := r.db.ExecContext(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	return nil
}

// RevokeAllUserSessions revokes all sessions for a user
func (r *PostgresUserRepository) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.user_sessions 
		SET revoked_at = NOW() 
		WHERE user_id = $1 AND revoked_at IS NULL`, r.schema)

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke all user sessions: %w", err)
	}

	return nil
}

// CleanExpiredSessions removes expired sessions
func (r *PostgresUserRepository) CleanExpiredSessions(ctx context.Context) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.user_sessions 
		WHERE expires_at < NOW() OR revoked_at IS NOT NULL`, r.schema)

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to clean expired sessions: %w", err)
	}

	return nil
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

// EmployeeIDExists checks if an employee ID exists within an organization
func (r *PostgresUserRepository) EmployeeIDExists(ctx context.Context, employeeID string, organizationID uuid.UUID) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s.users WHERE employee_id = $1 AND organization_id = $2)`, r.schema)

	var exists bool
	err := r.db.QueryRowContext(ctx, query, employeeID, organizationID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check employee ID existence: %w", err)
	}

	return exists, nil
}

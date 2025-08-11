package entities

import (
	"time"

	"github.com/google/uuid"
)

// UserSession represents a user session
type UserSession struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	TokenHash string     `json:"token_hash" db:"token_hash"`
	IPAddress *string    `json:"ip_address" db:"ip_address"`
	UserAgent *string    `json:"user_agent" db:"user_agent"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at" db:"revoked_at"`
}

// IsExpired checks if the session has expired
func (s *UserSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsRevoked checks if the session has been revoked
func (s *UserSession) IsRevoked() bool {
	return s.RevokedAt != nil
}

// IsValid checks if the session is valid (not expired and not revoked)
func (s *UserSession) IsValid() bool {
	return !s.IsExpired() && !s.IsRevoked()
}

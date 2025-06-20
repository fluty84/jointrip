package session

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Domain errors
var (
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionExpired     = errors.New("session expired")
	ErrInvalidSessionData = errors.New("invalid session data")
)

// Repository defines the interface for session data persistence
type Repository interface {
	// Create creates a new session
	Create(ctx context.Context, session *UserSession) error

	// GetByID retrieves a session by ID
	GetByID(ctx context.Context, id uuid.UUID) (*UserSession, error)

	// GetByAccessToken retrieves a session by access token
	GetByAccessToken(ctx context.Context, accessToken string) (*UserSession, error)

	// GetByRefreshToken retrieves a session by refresh token
	GetByRefreshToken(ctx context.Context, refreshToken string) (*UserSession, error)

	// GetActiveSessionsByUserID retrieves all active sessions for a user
	GetActiveSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]*UserSession, error)

	// Update updates an existing session
	Update(ctx context.Context, session *UserSession) error

	// Delete deletes a session
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteByUserID deletes all sessions for a user
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error

	// DeactivateExpiredSessions deactivates all expired sessions
	DeactivateExpiredSessions(ctx context.Context) error

	// CountActiveSessionsByUserID counts active sessions for a user
	CountActiveSessionsByUserID(ctx context.Context, userID uuid.UUID) (int, error)
}

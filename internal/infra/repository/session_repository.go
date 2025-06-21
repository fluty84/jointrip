package repository

import (
	"context"
	"database/sql"
	"fmt"

	"jointrip/internal/domain/session"

	"github.com/google/uuid"
)

// SessionRepository implements the session.Repository interface
type SessionRepository struct {
	db *sql.DB
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create creates a new session
func (r *SessionRepository) Create(ctx context.Context, s *session.UserSession) error {
	query := `
		INSERT INTO user_sessions (
			id, user_id, access_token, refresh_token, google_access_token,
			google_refresh_token, expires_at, ip_address, user_agent,
			is_active, created_at, last_used_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)`

	_, err := r.db.ExecContext(ctx, query,
		s.ID, s.UserID, s.AccessToken, s.RefreshToken, s.GoogleAccessToken,
		s.GoogleRefreshToken, s.ExpiresAt, s.IPAddress, s.UserAgent,
		s.IsActive, s.CreatedAt, s.LastUsedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// GetByID retrieves a session by ID
func (r *SessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*session.UserSession, error) {
	query := `
		SELECT id, user_id, access_token, refresh_token, google_access_token,
			   google_refresh_token, expires_at, ip_address, user_agent,
			   is_active, created_at, last_used_at
		FROM user_sessions 
		WHERE id = $1`

	return r.scanSession(r.db.QueryRowContext(ctx, query, id))
}

// GetByAccessToken retrieves a session by access token
func (r *SessionRepository) GetByAccessToken(ctx context.Context, accessToken string) (*session.UserSession, error) {
	query := `
		SELECT id, user_id, access_token, refresh_token, google_access_token,
			   google_refresh_token, expires_at, ip_address, user_agent,
			   is_active, created_at, last_used_at
		FROM user_sessions 
		WHERE access_token = $1`

	return r.scanSession(r.db.QueryRowContext(ctx, query, accessToken))
}

// GetByRefreshToken retrieves a session by refresh token
func (r *SessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*session.UserSession, error) {
	query := `
		SELECT id, user_id, access_token, refresh_token, google_access_token,
			   google_refresh_token, expires_at, ip_address, user_agent,
			   is_active, created_at, last_used_at
		FROM user_sessions 
		WHERE refresh_token = $1`

	return r.scanSession(r.db.QueryRowContext(ctx, query, refreshToken))
}

// GetActiveSessionsByUserID retrieves all active sessions for a user
func (r *SessionRepository) GetActiveSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]*session.UserSession, error) {
	query := `
		SELECT id, user_id, access_token, refresh_token, google_access_token,
			   google_refresh_token, expires_at, ip_address, user_agent,
			   is_active, created_at, last_used_at
		FROM user_sessions 
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at ASC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*session.UserSession
	for rows.Next() {
		s, err := r.scanSessionFromRows(rows)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sessions: %w", err)
	}

	return sessions, nil
}

// Update updates an existing session
func (r *SessionRepository) Update(ctx context.Context, s *session.UserSession) error {
	query := `
		UPDATE user_sessions SET
			access_token = $2, refresh_token = $3, google_access_token = $4,
			google_refresh_token = $5, expires_at = $6, is_active = $7,
			last_used_at = $8
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query,
		s.ID, s.AccessToken, s.RefreshToken, s.GoogleAccessToken,
		s.GoogleRefreshToken, s.ExpiresAt, s.IsActive, s.LastUsedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return session.ErrSessionNotFound
	}

	return nil
}

// Delete deletes a session
func (r *SessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM user_sessions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return session.ErrSessionNotFound
	}

	return nil
}

// DeleteByUserID deletes all sessions for a user
func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user_sessions WHERE user_id = $1`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete sessions by user ID: %w", err)
	}

	return nil
}

// DeactivateExpiredSessions deactivates all expired sessions
func (r *SessionRepository) DeactivateExpiredSessions(ctx context.Context) error {
	query := `
		UPDATE user_sessions 
		SET is_active = false, last_used_at = CURRENT_TIMESTAMP
		WHERE expires_at < CURRENT_TIMESTAMP AND is_active = true`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to deactivate expired sessions: %w", err)
	}

	return nil
}

// CountActiveSessionsByUserID counts active sessions for a user
func (r *SessionRepository) CountActiveSessionsByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM user_sessions 
		WHERE user_id = $1 AND is_active = true AND expires_at > CURRENT_TIMESTAMP`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count active sessions: %w", err)
	}

	return count, nil
}

// scanSession scans a session from a single row
func (r *SessionRepository) scanSession(row *sql.Row) (*session.UserSession, error) {
	s := &session.UserSession{}
	err := row.Scan(
		&s.ID, &s.UserID, &s.AccessToken, &s.RefreshToken, &s.GoogleAccessToken,
		&s.GoogleRefreshToken, &s.ExpiresAt, &s.IPAddress, &s.UserAgent,
		&s.IsActive, &s.CreatedAt, &s.LastUsedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, session.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to scan session: %w", err)
	}

	return s, nil
}

// scanSessionFromRows scans a session from multiple rows
func (r *SessionRepository) scanSessionFromRows(rows *sql.Rows) (*session.UserSession, error) {
	s := &session.UserSession{}
	err := rows.Scan(
		&s.ID, &s.UserID, &s.AccessToken, &s.RefreshToken, &s.GoogleAccessToken,
		&s.GoogleRefreshToken, &s.ExpiresAt, &s.IPAddress, &s.UserAgent,
		&s.IsActive, &s.CreatedAt, &s.LastUsedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan session from rows: %w", err)
	}

	return s, nil
}

package session

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// UserSession represents a user authentication session
type UserSession struct {
	ID                   uuid.UUID `json:"id"`
	UserID               uuid.UUID `json:"user_id"`
	AccessToken          string    `json:"access_token"`
	RefreshToken         string    `json:"refresh_token"`
	GoogleAccessToken    string    `json:"google_access_token"`
	GoogleRefreshToken   string    `json:"google_refresh_token"`
	ExpiresAt            time.Time `json:"expires_at"`
	IPAddress            string    `json:"ip_address"`
	UserAgent            string    `json:"user_agent"`
	IsActive             bool      `json:"is_active"`
	CreatedAt            time.Time `json:"created_at"`
	LastUsedAt           time.Time `json:"last_used_at"`
}

// NewUserSession creates a new user session
func NewUserSession(userID uuid.UUID, accessToken, refreshToken, googleAccessToken, googleRefreshToken string, expiresAt time.Time, ipAddress, userAgent string) (*UserSession, error) {
	if userID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}
	if accessToken == "" {
		return nil, errors.New("access token is required")
	}
	if refreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	now := time.Now()
	session := &UserSession{
		ID:                   uuid.New(),
		UserID:               userID,
		AccessToken:          accessToken,
		RefreshToken:         refreshToken,
		GoogleAccessToken:    googleAccessToken,
		GoogleRefreshToken:   googleRefreshToken,
		ExpiresAt:            expiresAt,
		IPAddress:            ipAddress,
		UserAgent:            userAgent,
		IsActive:             true,
		CreatedAt:            now,
		LastUsedAt:           now,
	}

	return session, nil
}

// UpdateTokens updates the session tokens
func (s *UserSession) UpdateTokens(accessToken, refreshToken string, expiresAt time.Time) {
	s.AccessToken = accessToken
	s.RefreshToken = refreshToken
	s.ExpiresAt = expiresAt
	s.LastUsedAt = time.Now()
}

// UpdateGoogleTokens updates the Google OAuth tokens
func (s *UserSession) UpdateGoogleTokens(googleAccessToken, googleRefreshToken string) {
	s.GoogleAccessToken = googleAccessToken
	s.GoogleRefreshToken = googleRefreshToken
	s.LastUsedAt = time.Now()
}

// UpdateLastUsed updates the last used timestamp
func (s *UserSession) UpdateLastUsed() {
	s.LastUsedAt = time.Now()
}

// Deactivate deactivates the session
func (s *UserSession) Deactivate() {
	s.IsActive = false
	s.LastUsedAt = time.Now()
}

// IsExpired checks if the session is expired
func (s *UserSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsValid checks if the session is valid (active and not expired)
func (s *UserSession) IsValid() bool {
	return s.IsActive && !s.IsExpired()
}

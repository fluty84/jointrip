package session

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserSession(t *testing.T) {
	userID := uuid.New()
	accessToken := "access_token_123"
	refreshToken := "refresh_token_123"
	googleAccessToken := "google_access_token_123"
	googleRefreshToken := "google_refresh_token_123"
	expiresAt := time.Now().Add(time.Hour)
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"

	tests := []struct {
		name               string
		userID             uuid.UUID
		accessToken        string
		refreshToken       string
		googleAccessToken  string
		googleRefreshToken string
		expiresAt          time.Time
		ipAddress          string
		userAgent          string
		expectError        bool
	}{
		{
			name:               "valid session creation",
			userID:             userID,
			accessToken:        accessToken,
			refreshToken:       refreshToken,
			googleAccessToken:  googleAccessToken,
			googleRefreshToken: googleRefreshToken,
			expiresAt:          expiresAt,
			ipAddress:          ipAddress,
			userAgent:          userAgent,
			expectError:        false,
		},
		{
			name:               "missing user ID",
			userID:             uuid.Nil,
			accessToken:        accessToken,
			refreshToken:       refreshToken,
			googleAccessToken:  googleAccessToken,
			googleRefreshToken: googleRefreshToken,
			expiresAt:          expiresAt,
			ipAddress:          ipAddress,
			userAgent:          userAgent,
			expectError:        true,
		},
		{
			name:               "missing access token",
			userID:             userID,
			accessToken:        "",
			refreshToken:       refreshToken,
			googleAccessToken:  googleAccessToken,
			googleRefreshToken: googleRefreshToken,
			expiresAt:          expiresAt,
			ipAddress:          ipAddress,
			userAgent:          userAgent,
			expectError:        true,
		},
		{
			name:               "missing refresh token",
			userID:             userID,
			accessToken:        accessToken,
			refreshToken:       "",
			googleAccessToken:  googleAccessToken,
			googleRefreshToken: googleRefreshToken,
			expiresAt:          expiresAt,
			ipAddress:          ipAddress,
			userAgent:          userAgent,
			expectError:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := NewUserSession(
				tt.userID,
				tt.accessToken,
				tt.refreshToken,
				tt.googleAccessToken,
				tt.googleRefreshToken,
				tt.expiresAt,
				tt.ipAddress,
				tt.userAgent,
			)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, session)
			} else {
				require.NoError(t, err)
				require.NotNil(t, session)

				assert.NotEqual(t, uuid.Nil, session.ID)
				assert.Equal(t, tt.userID, session.UserID)
				assert.Equal(t, tt.accessToken, session.AccessToken)
				assert.Equal(t, tt.refreshToken, session.RefreshToken)
				assert.Equal(t, tt.googleAccessToken, session.GoogleAccessToken)
				assert.Equal(t, tt.googleRefreshToken, session.GoogleRefreshToken)
				assert.Equal(t, tt.expiresAt, session.ExpiresAt)
				assert.Equal(t, tt.ipAddress, session.IPAddress)
				assert.Equal(t, tt.userAgent, session.UserAgent)
				assert.True(t, session.IsActive)
				assert.True(t, time.Since(session.CreatedAt) < time.Second)
				assert.True(t, time.Since(session.LastUsedAt) < time.Second)
			}
		})
	}
}

func TestUserSession_UpdateTokens(t *testing.T) {
	userID := uuid.New()
	session, err := NewUserSession(
		userID,
		"old_access_token",
		"old_refresh_token",
		"google_access_token",
		"google_refresh_token",
		time.Now().Add(time.Hour),
		"192.168.1.1",
		"Mozilla/5.0",
	)
	require.NoError(t, err)

	originalLastUsed := session.LastUsedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(time.Millisecond)

	newAccessToken := "new_access_token"
	newRefreshToken := "new_refresh_token"
	newExpiresAt := time.Now().Add(2 * time.Hour)

	session.UpdateTokens(newAccessToken, newRefreshToken, newExpiresAt)

	assert.Equal(t, newAccessToken, session.AccessToken)
	assert.Equal(t, newRefreshToken, session.RefreshToken)
	assert.Equal(t, newExpiresAt, session.ExpiresAt)
	assert.True(t, session.LastUsedAt.After(originalLastUsed))
}

func TestUserSession_UpdateGoogleTokens(t *testing.T) {
	userID := uuid.New()
	session, err := NewUserSession(
		userID,
		"access_token",
		"refresh_token",
		"old_google_access_token",
		"old_google_refresh_token",
		time.Now().Add(time.Hour),
		"192.168.1.1",
		"Mozilla/5.0",
	)
	require.NoError(t, err)

	originalLastUsed := session.LastUsedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(time.Millisecond)

	newGoogleAccessToken := "new_google_access_token"
	newGoogleRefreshToken := "new_google_refresh_token"

	session.UpdateGoogleTokens(newGoogleAccessToken, newGoogleRefreshToken)

	assert.Equal(t, newGoogleAccessToken, session.GoogleAccessToken)
	assert.Equal(t, newGoogleRefreshToken, session.GoogleRefreshToken)
	assert.True(t, session.LastUsedAt.After(originalLastUsed))
}

func TestUserSession_UpdateLastUsed(t *testing.T) {
	userID := uuid.New()
	session, err := NewUserSession(
		userID,
		"access_token",
		"refresh_token",
		"google_access_token",
		"google_refresh_token",
		time.Now().Add(time.Hour),
		"192.168.1.1",
		"Mozilla/5.0",
	)
	require.NoError(t, err)

	originalLastUsed := session.LastUsedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(time.Millisecond)

	session.UpdateLastUsed()

	assert.True(t, session.LastUsedAt.After(originalLastUsed))
}

func TestUserSession_Deactivate(t *testing.T) {
	userID := uuid.New()
	session, err := NewUserSession(
		userID,
		"access_token",
		"refresh_token",
		"google_access_token",
		"google_refresh_token",
		time.Now().Add(time.Hour),
		"192.168.1.1",
		"Mozilla/5.0",
	)
	require.NoError(t, err)

	assert.True(t, session.IsActive)

	session.Deactivate()

	assert.False(t, session.IsActive)
}

func TestUserSession_IsExpired(t *testing.T) {
	userID := uuid.New()

	// Create session that expires in the future
	futureSession, err := NewUserSession(
		userID,
		"access_token",
		"refresh_token",
		"google_access_token",
		"google_refresh_token",
		time.Now().Add(time.Hour),
		"192.168.1.1",
		"Mozilla/5.0",
	)
	require.NoError(t, err)

	assert.False(t, futureSession.IsExpired())

	// Create session that expired in the past
	pastSession, err := NewUserSession(
		userID,
		"access_token",
		"refresh_token",
		"google_access_token",
		"google_refresh_token",
		time.Now().Add(-time.Hour),
		"192.168.1.1",
		"Mozilla/5.0",
	)
	require.NoError(t, err)

	assert.True(t, pastSession.IsExpired())
}

func TestUserSession_IsValid(t *testing.T) {
	userID := uuid.New()

	// Create valid session (active and not expired)
	validSession, err := NewUserSession(
		userID,
		"access_token",
		"refresh_token",
		"google_access_token",
		"google_refresh_token",
		time.Now().Add(time.Hour),
		"192.168.1.1",
		"Mozilla/5.0",
	)
	require.NoError(t, err)

	assert.True(t, validSession.IsValid())

	// Deactivate session
	validSession.Deactivate()
	assert.False(t, validSession.IsValid())

	// Create expired session
	expiredSession, err := NewUserSession(
		userID,
		"access_token",
		"refresh_token",
		"google_access_token",
		"google_refresh_token",
		time.Now().Add(-time.Hour),
		"192.168.1.1",
		"Mozilla/5.0",
	)
	require.NoError(t, err)

	assert.False(t, expiredSession.IsValid())
}

package auth

import (
	"testing"
	"time"

	"jointrip/internal/infra/config"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTManager_GenerateTokens(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:                "test-secret-key",
			ExpirationHours:       1,
			RefreshExpirationHours: 24,
		},
	}

	jwtManager := NewJWTManager(cfg)
	userID := uuid.New()

	accessToken, refreshToken, expiresAt, err := jwtManager.GenerateTokens(userID)

	require.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	assert.True(t, expiresAt.After(time.Now()))
	assert.True(t, expiresAt.Before(time.Now().Add(2*time.Hour)))

	// Tokens should be different
	assert.NotEqual(t, accessToken, refreshToken)
}

func TestJWTManager_ValidateAccessToken(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:                "test-secret-key",
			ExpirationHours:       1,
			RefreshExpirationHours: 24,
		},
	}

	jwtManager := NewJWTManager(cfg)
	userID := uuid.New()

	// Generate tokens
	accessToken, _, _, err := jwtManager.GenerateTokens(userID)
	require.NoError(t, err)

	// Validate access token
	validatedUserID, err := jwtManager.ValidateAccessToken(accessToken)
	require.NoError(t, err)
	assert.Equal(t, userID, validatedUserID)

	// Test invalid token
	_, err = jwtManager.ValidateAccessToken("invalid-token")
	assert.Error(t, err)

	// Test empty token
	_, err = jwtManager.ValidateAccessToken("")
	assert.Error(t, err)
}

func TestJWTManager_ValidateRefreshToken(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:                "test-secret-key",
			ExpirationHours:       1,
			RefreshExpirationHours: 24,
		},
	}

	jwtManager := NewJWTManager(cfg)
	userID := uuid.New()

	// Generate tokens
	_, refreshToken, _, err := jwtManager.GenerateTokens(userID)
	require.NoError(t, err)

	// Validate refresh token
	validatedUserID, err := jwtManager.ValidateRefreshToken(refreshToken)
	require.NoError(t, err)
	assert.Equal(t, userID, validatedUserID)

	// Test invalid token
	_, err = jwtManager.ValidateRefreshToken("invalid-token")
	assert.Error(t, err)
}

func TestJWTManager_ValidateTokenType(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:                "test-secret-key",
			ExpirationHours:       1,
			RefreshExpirationHours: 24,
		},
	}

	jwtManager := NewJWTManager(cfg)
	userID := uuid.New()

	// Generate tokens
	accessToken, refreshToken, _, err := jwtManager.GenerateTokens(userID)
	require.NoError(t, err)

	// Try to validate access token as refresh token (should fail)
	_, err = jwtManager.ValidateRefreshToken(accessToken)
	assert.Error(t, err)

	// Try to validate refresh token as access token (should fail)
	_, err = jwtManager.ValidateAccessToken(refreshToken)
	assert.Error(t, err)
}

func TestJWTManager_ExpiredToken(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:                "test-secret-key",
			ExpirationHours:       0, // Expire immediately
			RefreshExpirationHours: 0, // Expire immediately
		},
	}

	jwtManager := NewJWTManager(cfg)
	userID := uuid.New()

	// Generate tokens that expire immediately
	accessToken, refreshToken, _, err := jwtManager.GenerateTokens(userID)
	require.NoError(t, err)

	// Wait a bit to ensure expiration
	time.Sleep(time.Millisecond)

	// Validate expired access token (should fail)
	_, err = jwtManager.ValidateAccessToken(accessToken)
	assert.Error(t, err)

	// Validate expired refresh token (should fail)
	_, err = jwtManager.ValidateRefreshToken(refreshToken)
	assert.Error(t, err)
}

func TestJWTManager_GetTokenClaims(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:                "test-secret-key",
			ExpirationHours:       1,
			RefreshExpirationHours: 24,
		},
	}

	jwtManager := NewJWTManager(cfg)
	userID := uuid.New()

	// Generate tokens
	accessToken, refreshToken, _, err := jwtManager.GenerateTokens(userID)
	require.NoError(t, err)

	// Get access token claims
	accessClaims, err := jwtManager.GetTokenClaims(accessToken)
	require.NoError(t, err)
	assert.Equal(t, userID, accessClaims.UserID)
	assert.Equal(t, "access", accessClaims.Type)
	assert.Equal(t, "jointrip", accessClaims.Issuer)
	assert.Equal(t, userID.String(), accessClaims.Subject)

	// Get refresh token claims
	refreshClaims, err := jwtManager.GetTokenClaims(refreshToken)
	require.NoError(t, err)
	assert.Equal(t, userID, refreshClaims.UserID)
	assert.Equal(t, "refresh", refreshClaims.Type)
	assert.Equal(t, "jointrip", refreshClaims.Issuer)
	assert.Equal(t, userID.String(), refreshClaims.Subject)

	// Test invalid token
	_, err = jwtManager.GetTokenClaims("invalid-token")
	assert.Error(t, err)
}

func TestJWTManager_DifferentSecrets(t *testing.T) {
	cfg1 := &config.Config{
		JWT: config.JWTConfig{
			Secret:                "secret-1",
			ExpirationHours:       1,
			RefreshExpirationHours: 24,
		},
	}

	cfg2 := &config.Config{
		JWT: config.JWTConfig{
			Secret:                "secret-2",
			ExpirationHours:       1,
			RefreshExpirationHours: 24,
		},
	}

	jwtManager1 := NewJWTManager(cfg1)
	jwtManager2 := NewJWTManager(cfg2)
	userID := uuid.New()

	// Generate token with first manager
	accessToken, _, _, err := jwtManager1.GenerateTokens(userID)
	require.NoError(t, err)

	// Try to validate with second manager (should fail)
	_, err = jwtManager2.ValidateAccessToken(accessToken)
	assert.Error(t, err)
}

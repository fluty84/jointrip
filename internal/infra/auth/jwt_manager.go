package auth

import (
	"errors"
	"time"

	"jointrip/internal/infra/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents JWT claims
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Type   string    `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// JWTManager handles JWT token operations
type JWTManager struct {
	secretKey                []byte
	accessTokenExpiration    time.Duration
	refreshTokenExpiration   time.Duration
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(cfg *config.Config) *JWTManager {
	return &JWTManager{
		secretKey:                []byte(cfg.JWT.Secret),
		accessTokenExpiration:    cfg.GetJWTExpiration(),
		refreshTokenExpiration:   cfg.GetRefreshTokenExpiration(),
	}
}

// GenerateTokens generates access and refresh tokens for a user
func (j *JWTManager) GenerateTokens(userID uuid.UUID) (accessToken, refreshToken string, expiresAt time.Time, err error) {
	now := time.Now()
	accessExpiresAt := now.Add(j.accessTokenExpiration)
	refreshExpiresAt := now.Add(j.refreshTokenExpiration)

	// Generate access token
	accessClaims := &Claims{
		UserID: userID,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "jointrip",
			Subject:   userID.String(),
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString(j.secretKey)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// Generate refresh token
	refreshClaims := &Claims{
		UserID: userID,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "jointrip",
			Subject:   userID.String(),
		},
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString(j.secretKey)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessToken, refreshToken, accessExpiresAt, nil
}

// ValidateAccessToken validates an access token and returns the user ID
func (j *JWTManager) ValidateAccessToken(tokenString string) (uuid.UUID, error) {
	return j.validateToken(tokenString, "access")
}

// ValidateRefreshToken validates a refresh token and returns the user ID
func (j *JWTManager) ValidateRefreshToken(tokenString string) (uuid.UUID, error) {
	return j.validateToken(tokenString, "refresh")
}

// validateToken validates a token and returns the user ID
func (j *JWTManager) validateToken(tokenString, expectedType string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return j.secretKey, nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	// Validate token type
	if claims.Type != expectedType {
		return uuid.Nil, errors.New("invalid token type")
	}

	// Validate expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return uuid.Nil, errors.New("token expired")
	}

	// Validate not before
	if claims.NotBefore != nil && claims.NotBefore.Time.After(time.Now()) {
		return uuid.Nil, errors.New("token not valid yet")
	}

	return claims.UserID, nil
}

// GetTokenClaims extracts claims from a token without validation (for debugging)
func (j *JWTManager) GetTokenClaims(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

package middleware

import (
	"errors"
	"net/http"
	"strings"

	"jointrip/internal/app/auth"
	"jointrip/internal/domain/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthMiddleware provides authentication middleware
type AuthMiddleware struct {
	authService *auth.Service
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authService *auth.Service) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth middleware that requires authentication
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token required",
			})
			c.Abort()
			return
		}

		user, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Next()
	}
}

// OptionalAuth middleware that optionally authenticates
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token != "" {
			user, err := m.authService.ValidateToken(c.Request.Context(), token)
			if err == nil {
				c.Set("user", user)
				c.Set("user_id", user.ID)
			}
		}
		c.Next()
	}
}

// RequireVerified middleware that requires verified users
func (m *AuthMiddleware) RequireVerified() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		user, ok := userInterface.(*user.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user context",
			})
			c.Abort()
			return
		}

		if !user.IsVerified() {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Account verification required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// extractToken extracts the bearer token from the request
func (m *AuthMiddleware) extractToken(c *gin.Context) string {
	// Try Authorization header first
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	// Try query parameter as fallback
	return c.Query("token")
}

// GetCurrentUser helper function to get current user from context
func GetCurrentUser(c *gin.Context) (*user.User, error) {
	userInterface, exists := c.Get("user")
	if !exists {
		return nil, ErrUserNotInContext
	}

	user, ok := userInterface.(*user.User)
	if !ok {
		return nil, ErrInvalidUserContext
	}

	return user, nil
}

// GetCurrentUserID helper function to get current user ID from context
func GetCurrentUserID(c *gin.Context) (uuid.UUID, error) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, ErrUserNotInContext
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrInvalidUserContext
	}

	return userID, nil
}

// Custom errors
var (
	ErrUserNotInContext   = errors.New("user not in context")
	ErrInvalidUserContext = errors.New("invalid user context")
)

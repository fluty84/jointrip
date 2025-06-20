package handlers

import (
	"net/http"

	"jointrip/internal/app/auth"
	"jointrip/internal/infra/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService *auth.Service
	logger      *logrus.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *auth.Service, logger *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Code  string `json:"code" binding:"required"`
	State string `json:"state"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// GetGoogleAuthURL returns the Google OAuth authorization URL
func (h *AuthHandler) GetGoogleAuthURL(c *gin.Context) {
	state := c.Query("state")
	if state == "" {
		state = uuid.New().String()
	}

	authURL := h.authService.GetGoogleAuthURL(state)

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"state":    state,
	})
}

// Login handles Google OAuth login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Get client info
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Perform login
	response, err := h.authService.LoginWithGoogle(c.Request.Context(), req.Code, ipAddress, userAgent)
	if err != nil {
		h.logger.WithError(err).Error("Failed to login with Google")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication failed",
		})
		return
	}

	h.logger.WithField("user_id", response.User.ID).Info("User logged in successfully")

	c.JSON(http.StatusOK, gin.H{
		"user":          response.User,
		"access_token":  response.AccessToken,
		"refresh_token": response.RefreshToken,
		"expires_at":    response.ExpiresAt,
		"token_type":    "Bearer",
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	response, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.logger.WithError(err).Error("Failed to refresh token")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token refresh failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": response.AccessToken,
		"expires_at":   response.ExpiresAt,
		"token_type":   "Bearer",
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Extract token from request
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Authorization header required",
		})
		return
	}

	// Extract bearer token
	var token string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid authorization header format",
		})
		return
	}

	err := h.authService.Logout(c.Request.Context(), token)
	if err != nil {
		h.logger.WithError(err).Error("Failed to logout")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Logout failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// LogoutAll handles logout from all sessions
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}

	err = h.authService.LogoutAll(c.Request.Context(), userID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to logout from all sessions")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Logout failed",
		})
		return
	}

	h.logger.WithField("user_id", userID).Info("User logged out from all sessions")

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out from all sessions successfully",
	})
}

// GetProfile returns the current user's profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	user, err := middleware.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// ValidateToken validates the current token
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	user, err := middleware.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":   true,
		"user_id": user.ID,
		"email":   user.Email,
	})
}

package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"jointrip/internal/app/auth"
	"jointrip/internal/domain/user"
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

// UpdateProfileRequest represents a profile update request
type UpdateProfileRequest struct {
	FirstName          *string  `json:"first_name,omitempty"`
	LastName           *string  `json:"last_name,omitempty"`
	Bio                *string  `json:"bio,omitempty"`
	Location           *string  `json:"location,omitempty"`
	Phone              *string  `json:"phone,omitempty"`
	Website            *string  `json:"website,omitempty"`
	Languages          []string `json:"languages,omitempty"`
	Interests          []string `json:"interests,omitempty"`
	TravelStyle        *string  `json:"travel_style,omitempty"`
	ProfileVisibility  *string  `json:"profile_visibility,omitempty"`
	EmailNotifications *bool    `json:"email_notifications,omitempty"`
	PushNotifications  *bool    `json:"push_notifications,omitempty"`
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
		"user":         response.User,
		"accessToken":  response.AccessToken,
		"refreshToken": response.RefreshToken,
		"expiresAt":    response.ExpiresAt,
		"tokenType":    "Bearer",
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
		"accessToken": response.AccessToken,
		"expiresAt":   response.ExpiresAt,
		"tokenType":   "Bearer",
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

// UpdateProfile updates the current user's profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	currentUser, err := middleware.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Update user fields if provided
	if req.FirstName != nil {
		currentUser.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		currentUser.LastName = *req.LastName
	}
	if req.Bio != nil {
		currentUser.Bio = *req.Bio
	}
	if req.Location != nil {
		currentUser.Location = *req.Location
	}
	if req.Phone != nil {
		currentUser.Phone = req.Phone
	}
	if req.Website != nil {
		currentUser.Website = *req.Website
	}
	if req.Languages != nil {
		currentUser.Languages = req.Languages
	}
	if req.Interests != nil {
		currentUser.Interests = req.Interests
	}
	if req.TravelStyle != nil {
		travelStyle := user.TravelStyle(*req.TravelStyle)
		currentUser.TravelStyle = &travelStyle
	}
	if req.ProfileVisibility != nil {
		currentUser.ProfileVisibility = user.PrivacyLevel(*req.ProfileVisibility)
	}
	if req.EmailNotifications != nil {
		currentUser.EmailNotifications = *req.EmailNotifications
	}
	if req.PushNotifications != nil {
		currentUser.PushNotifications = *req.PushNotifications
	}

	// Update the user in the database
	err = h.authService.UpdateUser(c.Request.Context(), currentUser)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update user profile")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update profile",
		})
		return
	}

	h.logger.WithField("user_id", currentUser.ID).Info("User profile updated successfully")

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    currentUser,
	})
}

// UploadProfilePhoto handles profile photo upload
func (h *AuthHandler) UploadProfilePhoto(c *gin.Context) {
	currentUser, err := middleware.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}

	// Parse multipart form
	err = c.Request.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse form data",
		})
		return
	}

	// Get file from form
	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No photo file provided",
		})
		return
	}
	defer file.Close()

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File must be an image",
		})
		return
	}

	// Validate file size (5MB max)
	if header.Size > 5<<20 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File size must be less than 5MB",
		})
		return
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := "uploads/profile_photos"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		h.logger.WithError(err).Error("Failed to create uploads directory")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save photo",
		})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s%s", currentUser.ID.String(), ext)
	filePath := filepath.Join(uploadsDir, filename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create destination file")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save photo",
		})
		return
	}
	defer dst.Close()

	// Copy uploaded file to destination
	_, err = io.Copy(dst, file)
	if err != nil {
		h.logger.WithError(err).Error("Failed to copy file")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save photo",
		})
		return
	}

	// Update user's profile photo URL
	photoURL := fmt.Sprintf("/uploads/profile_photos/%s", filename)
	currentUser.ProfilePhotoURL = photoURL

	// Save to database
	err = h.authService.UpdateUser(c.Request.Context(), currentUser)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update user profile photo")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update profile photo",
		})
		return
	}

	h.logger.WithField("user_id", currentUser.ID).Info("Profile photo updated successfully")

	c.JSON(http.StatusOK, gin.H{
		"message":   "Profile photo updated successfully",
		"photo_url": photoURL,
		"user":      currentUser,
	})
}

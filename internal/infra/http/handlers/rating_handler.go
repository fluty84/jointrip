package handlers

import (
	"net/http"

	"jointrip/internal/infra/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// RatingHandler handles user rating-related HTTP requests
type RatingHandler struct {
	logger *logrus.Logger
}

// NewRatingHandler creates a new rating handler
func NewRatingHandler(logger *logrus.Logger) *RatingHandler {
	return &RatingHandler{
		logger: logger,
	}
}

// CreateRatingRequest represents a rating creation request
type CreateRatingRequest struct {
	RatedUserID uuid.UUID `json:"rated_user_id" binding:"required"`
	Rating      int       `json:"rating" binding:"required,min=1,max=5"`
	Review      string    `json:"review,omitempty"`
	TripID      *uuid.UUID `json:"trip_id,omitempty"`
}

// CreateRating creates a new user rating
func (h *RatingHandler) CreateRating(c *gin.Context) {
	currentUser, err := middleware.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}

	var req CreateRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Validate that user is not rating themselves
	if req.RatedUserID == currentUser.ID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot rate yourself",
		})
		return
	}

	// TODO: Implement rating creation logic
	// For now, return success response
	h.logger.WithFields(logrus.Fields{
		"rater_id":     currentUser.ID,
		"rated_user_id": req.RatedUserID,
		"rating":       req.Rating,
	}).Info("Rating created successfully")

	c.JSON(http.StatusCreated, gin.H{
		"message": "Rating created successfully",
		"rating": gin.H{
			"rater_id":     currentUser.ID,
			"rated_user_id": req.RatedUserID,
			"rating":       req.Rating,
			"review":       req.Review,
			"trip_id":      req.TripID,
		},
	})
}

// GetUserRatings returns ratings for a specific user
func (h *RatingHandler) GetUserRatings(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// TODO: Implement get ratings logic
	// For now, return mock data
	h.logger.WithField("user_id", userID).Info("Fetching user ratings")

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"ratings": []gin.H{
			{
				"id":         uuid.New(),
				"rater_id":   uuid.New(),
				"rating":     5,
				"review":     "Great travel companion!",
				"created_at": "2024-01-01T00:00:00Z",
			},
			{
				"id":         uuid.New(),
				"rater_id":   uuid.New(),
				"rating":     4,
				"review":     "Very organized and friendly",
				"created_at": "2024-01-02T00:00:00Z",
			},
		},
		"average_rating": 4.5,
		"total_ratings":  2,
	})
}

// GetMyRatings returns ratings given by the current user
func (h *RatingHandler) GetMyRatings(c *gin.Context) {
	currentUser, err := middleware.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}

	// TODO: Implement get my ratings logic
	// For now, return mock data
	h.logger.WithField("user_id", currentUser.ID).Info("Fetching user's given ratings")

	c.JSON(http.StatusOK, gin.H{
		"user_id": currentUser.ID,
		"ratings": []gin.H{
			{
				"id":           uuid.New(),
				"rated_user_id": uuid.New(),
				"rating":       4,
				"review":       "Good experience traveling together",
				"created_at":   "2024-01-01T00:00:00Z",
			},
		},
		"total_given": 1,
	})
}

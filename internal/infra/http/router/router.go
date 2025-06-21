package router

import (
	"jointrip/internal/app/auth"
	"jointrip/internal/infra/config"
	"jointrip/internal/infra/http/handlers"
	"jointrip/internal/infra/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Router wraps the Gin router with our application routes
type Router struct {
	engine         *gin.Engine
	authHandler    *handlers.AuthHandler
	authMiddleware *middleware.AuthMiddleware
}

// NewRouter creates a new router with all routes configured
func NewRouter(
	cfg *config.Config,
	authService *auth.Service,
	logger *logrus.Logger,
) *Router {
	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()

	// Create middleware
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Add middleware
	engine.Use(gin.Recovery())
	engine.Use(loggingMiddleware.RequestLogger())
	engine.Use(loggingMiddleware.ErrorLogger())

	// CORS middleware
	engine.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService, logger)

	router := &Router{
		engine:         engine,
		authHandler:    authHandler,
		authMiddleware: authMiddleware,
	}

	router.setupRoutes()
	return router
}

// setupRoutes configures all application routes
func (r *Router) setupRoutes() {
	// Health check
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "jointrip-api",
		})
	})

	// API v1 routes
	v1 := r.engine.Group("/api/v1")

	// Authentication routes (public)
	auth := v1.Group("/auth")
	{
		auth.GET("/google/url", r.authHandler.GetGoogleAuthURL)
		auth.POST("/google/login", r.authHandler.Login)
		auth.POST("/refresh", r.authHandler.RefreshToken)
		auth.POST("/logout", r.authHandler.Logout)
	}

	// Protected routes (require authentication)
	protected := v1.Group("/")
	protected.Use(r.authMiddleware.RequireAuth())
	{
		// User profile routes
		protected.GET("/profile", r.authHandler.GetProfile)
		protected.POST("/auth/logout-all", r.authHandler.LogoutAll)
		protected.GET("/auth/validate", r.authHandler.ValidateToken)
	}

	// Optional auth routes (authentication optional)
	optional := v1.Group("/")
	optional.Use(r.authMiddleware.OptionalAuth())
	{
		// Add routes that work with or without authentication
	}

	// Verified user routes (require verified account)
	verified := v1.Group("/")
	verified.Use(r.authMiddleware.RequireAuth())
	verified.Use(r.authMiddleware.RequireVerified())
	{
		// Add routes that require verified users (like creating trips)
	}
}

// GetEngine returns the underlying Gin engine
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

// Run starts the HTTP server
func (r *Router) Run(address string) error {
	return r.engine.Run(address)
}

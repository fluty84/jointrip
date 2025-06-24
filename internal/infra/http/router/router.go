package router

import (
	"io"
	"io/fs"
	"jointrip/internal/app/auth"
	"jointrip/internal/infra/config"
	"jointrip/internal/infra/http/handlers"
	"jointrip/internal/infra/http/middleware"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Router wraps the Gin router with our application routes
type Router struct {
	engine         *gin.Engine
	authHandler    *handlers.AuthHandler
	ratingHandler  *handlers.RatingHandler
	authMiddleware *middleware.AuthMiddleware
	webFS          fs.FS
}

// NewRouter creates a new router with all routes configured
func NewRouter(
	cfg *config.Config,
	authService *auth.Service,
	logger *logrus.Logger,
	webFS fs.FS,
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
	ratingHandler := handlers.NewRatingHandler(logger)

	router := &Router{
		engine:         engine,
		authHandler:    authHandler,
		ratingHandler:  ratingHandler,
		authMiddleware: authMiddleware,
		webFS:          webFS,
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
		protected.PUT("/profile", r.authHandler.UpdateProfile)
		protected.POST("/profile/photo", r.authHandler.UploadProfilePhoto)
		protected.POST("/auth/logout-all", r.authHandler.LogoutAll)
		protected.GET("/auth/validate", r.authHandler.ValidateToken)

		// Rating routes
		protected.POST("/ratings", r.ratingHandler.CreateRating)
		protected.GET("/ratings/my", r.ratingHandler.GetMyRatings)
		protected.GET("/users/:user_id/ratings", r.ratingHandler.GetUserRatings)
	}

	// Optional auth routes (authentication optional)
	optional := v1.Group("/")
	optional.Use(r.authMiddleware.OptionalAuth())
	{
		// Add routes that work with or without authentication
	}

	// Serve React static files
	r.setupStaticRoutes()
}

// setupStaticRoutes configures static file serving for React app
func (r *Router) setupStaticRoutes() {
	// Serve uploaded files
	r.engine.Static("/uploads", "./uploads")

	// Serve assets manually to avoid NoRoute conflicts
	r.engine.GET("/assets/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		// Remove leading slash
		if len(filepath) > 0 && filepath[0] == '/' {
			filepath = filepath[1:]
		}

		file, err := r.webFS.Open("assets/" + filepath)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		defer file.Close()

		// Set appropriate content type
		if strings.HasSuffix(filepath, ".js") {
			c.Header("Content-Type", "application/javascript")
		} else if strings.HasSuffix(filepath, ".css") {
			c.Header("Content-Type", "text/css")
		}

		http.ServeContent(c.Writer, c.Request, filepath, time.Time{}, file.(io.ReadSeeker))
	})

	// Serve specific static files
	r.engine.GET("/vite.svg", func(c *gin.Context) {
		file, err := r.webFS.Open("vite.svg")
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		defer file.Close()

		c.Header("Content-Type", "image/svg+xml")
		http.ServeContent(c.Writer, c.Request, "vite.svg", time.Time{}, file.(io.ReadSeeker))
	})

	// Serve index.html for all non-API, non-asset routes (SPA routing)
	r.engine.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Skip API routes
		if len(path) >= 4 && path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
			return
		}

		// Skip if it's an asset request that wasn't found
		if len(path) >= 7 && path[:7] == "/assets" {
			c.Status(http.StatusNotFound)
			return
		}

		// Serve index.html for SPA routes
		indexFile, err := r.webFS.Open("index.html")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load application"})
			return
		}
		defer indexFile.Close()

		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Status(http.StatusOK)

		// Copy file content to response
		http.ServeContent(c.Writer, c.Request, "index.html", time.Time{}, indexFile.(io.ReadSeeker))
	})
}

// GetEngine returns the underlying Gin engine
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

// Run starts the HTTP server
func (r *Router) Run(address string) error {
	return r.engine.Run(address)
}

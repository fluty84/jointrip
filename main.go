package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"jointrip/internal/app/auth"
	infraAuth "jointrip/internal/infra/auth"
	"jointrip/internal/infra/config"
	"jointrip/internal/infra/database"
	"jointrip/internal/infra/http/router"
	"jointrip/internal/infra/logger"
	"jointrip/internal/infra/repository"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.NewLogger(cfg)
	log.Info("Starting JoinTrip API server")

	// Initialize database connection
	db, err := database.NewConnection(cfg)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	// Run database migrations
	if err := db.RunMigrations("migrations"); err != nil {
		log.WithError(err).Fatal("Failed to run database migrations")
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	sessionRepo := repository.NewSessionRepository(db.DB)

	// Initialize infrastructure services
	jwtManager := infraAuth.NewJWTManager(cfg)
	googleClient := infraAuth.NewGoogleOAuthClient(cfg)

	// Initialize application services
	authService := auth.NewService(
		userRepo,
		sessionRepo,
		googleClient,
		jwtManager,
		cfg.Session.MaxSessionsPerUser,
	)

	// Initialize HTTP router
	httpRouter := router.NewRouter(cfg, authService, log)

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      httpRouter.GetEngine(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.WithField("address", cfg.GetServerAddress()).Info("Starting HTTP server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start HTTP server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Error("Server forced to shutdown")
	} else {
		log.Info("Server shutdown complete")
	}
}

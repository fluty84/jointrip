package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Domain errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUserData   = errors.New("invalid user data")
)

// Repository defines the interface for user data persistence
type Repository interface {
	// Create creates a new user
	Create(ctx context.Context, user *User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)

	// GetByGoogleID retrieves a user by Google ID
	GetByGoogleID(ctx context.Context, googleID string) (*User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*User, error)

	// GetByUsername retrieves a user by username
	GetByUsername(ctx context.Context, username string) (*User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *User) error

	// Delete soft deletes a user
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves users with pagination
	List(ctx context.Context, limit, offset int) ([]*User, error)

	// ExistsByEmail checks if a user exists with the given email
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// ExistsByUsername checks if a user exists with the given username
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}

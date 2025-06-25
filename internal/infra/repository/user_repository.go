package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"jointrip/internal/domain/user"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// UserRepository implements the user.Repository interface
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	query := `
		INSERT INTO users (
			id, google_id, email, username, first_name, last_name, phone,
			date_of_birth, gender, bio, profile_photo_url, google_photo_url,
			reputation_score, privacy_level, is_active,
			last_login, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		)`

	_, err := r.db.ExecContext(ctx, query,
		u.ID, u.GoogleID, u.Email, u.Username, u.FirstName, u.LastName, u.Phone,
		u.DateOfBirth, u.Gender, u.Bio, u.ProfilePhotoURL, u.GooglePhotoURL,
		u.ReputationScore, u.PrivacyLevel, u.IsActive,
		u.LastLogin, u.CreatedAt, u.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				return user.ErrUserAlreadyExists
			}
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	query := `
		SELECT id, google_id, email, username, first_name, last_name, phone,
			   date_of_birth, gender, bio, profile_photo_url, google_photo_url,
			   reputation_score, privacy_level, is_active,
			   last_login, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_active = true`

	return r.scanUser(r.db.QueryRowContext(ctx, query, id))
}

// GetByGoogleID retrieves a user by Google ID
func (r *UserRepository) GetByGoogleID(ctx context.Context, googleID string) (*user.User, error) {
	query := `
		SELECT id, google_id, email, username, first_name, last_name, phone,
			   date_of_birth, gender, bio, profile_photo_url, google_photo_url,
			   reputation_score, privacy_level, is_active,
			   last_login, created_at, updated_at
		FROM users
		WHERE google_id = $1 AND is_active = true`

	return r.scanUser(r.db.QueryRowContext(ctx, query, googleID))
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, google_id, email, username, first_name, last_name, phone,
			   date_of_birth, gender, bio, profile_photo_url, google_photo_url,
			   reputation_score, privacy_level, is_active,
			   last_login, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_active = true`

	return r.scanUser(r.db.QueryRowContext(ctx, query, email))
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	query := `
		SELECT id, google_id, email, username, first_name, last_name, phone,
			   date_of_birth, gender, bio, profile_photo_url, google_photo_url,
			   reputation_score, privacy_level, is_active,
			   last_login, created_at, updated_at
		FROM users
		WHERE username = $1 AND is_active = true`

	return r.scanUser(r.db.QueryRowContext(ctx, query, username))
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	query := `
		UPDATE users SET
			email = $2, username = $3, first_name = $4, last_name = $5, phone = $6,
			date_of_birth = $7, gender = $8, bio = $9, profile_photo_url = $10,
			reputation_score = $11, privacy_level = $12,
			is_active = $13, last_login = $14, updated_at = $15
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query,
		u.ID, u.Email, u.Username, u.FirstName, u.LastName, u.Phone,
		u.DateOfBirth, u.Gender, u.Bio, u.ProfilePhotoURL,
		u.ReputationScore, u.PrivacyLevel,
		u.IsActive, u.LastLogin, u.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return user.ErrUserNotFound
	}

	return nil
}

// UpdateProfile updates user profile fields specifically
func (r *UserRepository) UpdateProfile(ctx context.Context, userID uuid.UUID, profileData map[string]interface{}) error {
	// Build dynamic query based on provided fields
	setParts := []string{}
	args := []interface{}{userID}
	argIndex := 2

	for field, value := range profileData {
		switch field {
		case "first_name", "last_name", "bio", "location", "website", "phone":
			setParts = append(setParts, fmt.Sprintf("%s = $%d", field, argIndex))
			args = append(args, value)
			argIndex++
		case "languages", "interests":
			setParts = append(setParts, fmt.Sprintf("%s = $%d", field, argIndex))
			if strSlice, ok := value.([]string); ok {
				args = append(args, pq.Array(strSlice))
			} else {
				args = append(args, pq.Array([]string{}))
			}
			argIndex++
		case "travel_style", "profile_visibility":
			setParts = append(setParts, fmt.Sprintf("%s = $%d", field, argIndex))
			args = append(args, value)
			argIndex++
		case "email_notifications", "push_notifications":
			setParts = append(setParts, fmt.Sprintf("%s = $%d", field, argIndex))
			args = append(args, value)
			argIndex++
		}
	}

	if len(setParts) == 0 {
		return nil // Nothing to update
	}

	// Add updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $1", strings.Join(setParts, ", "))

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return user.ErrUserNotFound
	}

	return nil
}

// Delete soft deletes a user
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return user.ErrUserNotFound
	}

	return nil
}

// List retrieves users with pagination
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	query := `
		SELECT id, google_id, email, username, first_name, last_name, phone,
			   date_of_birth, gender, bio, profile_photo_url, google_photo_url,
			   reputation_score, privacy_level, is_active,
			   last_login, created_at, updated_at
		FROM users
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		u, err := r.scanUserFromRows(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

// ExistsByEmail checks if a user exists with the given email
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND is_active = true)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence by email: %w", err)
	}

	return exists, nil
}

// ExistsByUsername checks if a user exists with the given username
func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND is_active = true)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence by username: %w", err)
	}

	return exists, nil
}

// scanUser scans a user from a single row
func (r *UserRepository) scanUser(row *sql.Row) (*user.User, error) {
	u := &user.User{}
	err := row.Scan(
		&u.ID, &u.GoogleID, &u.Email, &u.Username, &u.FirstName, &u.LastName, &u.Phone,
		&u.DateOfBirth, &u.Gender, &u.Bio, &u.ProfilePhotoURL, &u.GooglePhotoURL,
		&u.ReputationScore, &u.PrivacyLevel, &u.IsActive,
		&u.LastLogin, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	return u, nil
}

// scanUserFromRows scans a user from multiple rows
func (r *UserRepository) scanUserFromRows(rows *sql.Rows) (*user.User, error) {
	u := &user.User{}
	err := rows.Scan(
		&u.ID, &u.GoogleID, &u.Email, &u.Username, &u.FirstName, &u.LastName, &u.Phone,
		&u.DateOfBirth, &u.Gender, &u.Bio, &u.ProfilePhotoURL, &u.GooglePhotoURL,
		&u.ReputationScore, &u.PrivacyLevel, &u.IsActive,
		&u.LastLogin, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan user from rows: %w", err)
	}

	return u, nil
}

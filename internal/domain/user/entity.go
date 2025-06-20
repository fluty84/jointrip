package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// VerificationStatus represents the user's identity verification status
type VerificationStatus string

const (
	VerificationStatusUnverified VerificationStatus = "unverified"
	VerificationStatusPending    VerificationStatus = "pending"
	VerificationStatusVerified   VerificationStatus = "verified"
	VerificationStatusRejected   VerificationStatus = "rejected"
)

// PrivacyLevel represents the user's profile visibility settings
type PrivacyLevel string

const (
	PrivacyLevelPublic  PrivacyLevel = "public"
	PrivacyLevelFriends PrivacyLevel = "friends"
	PrivacyLevelPrivate PrivacyLevel = "private"
)

// Gender represents the user's gender
type Gender string

const (
	GenderMale      Gender = "male"
	GenderFemale    Gender = "female"
	GenderOther     Gender = "other"
	GenderPreferNot Gender = "prefer_not_to_say"
)

// User represents a user in the system
type User struct {
	ID                 uuid.UUID          `json:"id"`
	GoogleID           string             `json:"google_id"`
	Email              string             `json:"email"`
	Username           string             `json:"username"`
	FirstName          string             `json:"first_name"`
	LastName           string             `json:"last_name"`
	Phone              *string            `json:"phone,omitempty"`
	DateOfBirth        *time.Time         `json:"date_of_birth,omitempty"`
	Gender             *Gender            `json:"gender,omitempty"`
	Bio                string             `json:"bio"`
	ProfilePhotoURL    string             `json:"profile_photo_url"`
	GooglePhotoURL     string             `json:"google_photo_url"`
	VerificationStatus VerificationStatus `json:"verification_status"`
	ReputationScore    float64            `json:"reputation_score"`
	PrivacyLevel       PrivacyLevel       `json:"privacy_level"`
	IsActive           bool               `json:"is_active"`
	LastLogin          *time.Time         `json:"last_login,omitempty"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
}

// NewUser creates a new user from Google OAuth data
func NewUser(googleID, email, firstName, lastName, googlePhotoURL string) (*User, error) {
	if googleID == "" {
		return nil, errors.New("google ID is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if firstName == "" {
		return nil, errors.New("first name is required")
	}

	now := time.Now()
	user := &User{
		ID:                 uuid.New(),
		GoogleID:           googleID,
		Email:              email,
		Username:           generateUsername(firstName, lastName),
		FirstName:          firstName,
		LastName:           lastName,
		GooglePhotoURL:     googlePhotoURL,
		ProfilePhotoURL:    googlePhotoURL, // Initially use Google photo
		VerificationStatus: VerificationStatusUnverified,
		ReputationScore:    0.0,
		PrivacyLevel:       PrivacyLevelPublic,
		IsActive:           true,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	return user, nil
}

// UpdateProfile updates the user's profile information
func (u *User) UpdateProfile(firstName, lastName, bio string, phone *string, dateOfBirth *time.Time, gender *Gender) error {
	if firstName == "" {
		return errors.New("first name is required")
	}
	if lastName == "" {
		return errors.New("last name is required")
	}

	u.FirstName = firstName
	u.LastName = lastName
	u.Bio = bio
	u.Phone = phone
	u.DateOfBirth = dateOfBirth
	u.Gender = gender
	u.UpdatedAt = time.Now()

	return nil
}

// UpdateLastLogin updates the user's last login timestamp
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
	u.UpdatedAt = now
}

// SetVerificationStatus updates the user's verification status
func (u *User) SetVerificationStatus(status VerificationStatus) {
	u.VerificationStatus = status
	u.UpdatedAt = time.Now()
}

// SetPrivacyLevel updates the user's privacy level
func (u *User) SetPrivacyLevel(level PrivacyLevel) {
	u.PrivacyLevel = level
	u.UpdatedAt = time.Now()
}

// Deactivate deactivates the user account
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Activate activates the user account
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// IsVerified returns true if the user is verified
func (u *User) IsVerified() bool {
	return u.VerificationStatus == VerificationStatusVerified
}

// CanCreateTrips returns true if the user can create trips
func (u *User) CanCreateTrips() bool {
	return u.IsActive && u.IsVerified()
}

// generateUsername generates a unique username from first and last name
func generateUsername(firstName, lastName string) string {
	// Simple implementation - in production, you'd want to ensure uniqueness
	return firstName + lastName + uuid.New().String()[:8]
}

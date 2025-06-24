package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// PrivacyLevel represents the user's profile visibility settings
type PrivacyLevel string

const (
	PrivacyLevelPublic  PrivacyLevel = "public"
	PrivacyLevelFriends PrivacyLevel = "friends"
	PrivacyLevelPrivate PrivacyLevel = "private"
)

// TravelStyle represents the user's preferred travel style
type TravelStyle string

const (
	TravelStyleBudget     TravelStyle = "budget"
	TravelStyleMidRange   TravelStyle = "mid-range"
	TravelStyleLuxury     TravelStyle = "luxury"
	TravelStyleBackpacker TravelStyle = "backpacker"
	TravelStyleAdventure  TravelStyle = "adventure"
	TravelStyleCultural   TravelStyle = "cultural"
	TravelStyleRelaxation TravelStyle = "relaxation"
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
	ID              uuid.UUID    `json:"id"`
	GoogleID        string       `json:"google_id"`
	Email           string       `json:"email"`
	Username        string       `json:"username"`
	FirstName       string       `json:"first_name"`
	LastName        string       `json:"last_name"`
	Phone           *string      `json:"phone,omitempty"`
	DateOfBirth     *time.Time   `json:"date_of_birth,omitempty"`
	Gender          *Gender      `json:"gender,omitempty"`
	Bio             string       `json:"bio"`
	Location        string       `json:"location"`
	Website         string       `json:"website"`
	Languages       []string     `json:"languages"`
	Interests       []string     `json:"interests"`
	TravelStyle     *TravelStyle `json:"travel_style,omitempty"`
	ProfilePhotoURL string       `json:"profile_photo_url"`
	GooglePhotoURL  string       `json:"google_photo_url"`

	ReputationScore             float64      `json:"reputation_score"`
	RatingAverage               float64      `json:"rating_average"`
	RatingCount                 int          `json:"rating_count"`
	PrivacyLevel                PrivacyLevel `json:"privacy_level"`
	ProfileVisibility           PrivacyLevel `json:"profile_visibility"`
	EmailNotifications          bool         `json:"email_notifications"`
	PushNotifications           bool         `json:"push_notifications"`
	ProfileCompletionPercentage int          `json:"profile_completion_percentage"`
	IsActive                    bool         `json:"is_active"`
	LastLogin                   *time.Time   `json:"last_login,omitempty"`
	CreatedAt                   time.Time    `json:"created_at"`
	UpdatedAt                   time.Time    `json:"updated_at"`
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
		ID:              uuid.New(),
		GoogleID:        googleID,
		Email:           email,
		Username:        generateUsername(firstName, lastName),
		FirstName:       firstName,
		LastName:        lastName,
		Languages:       []string{},
		Interests:       []string{},
		GooglePhotoURL:  googlePhotoURL,
		ProfilePhotoURL: googlePhotoURL, // Initially use Google photo

		ReputationScore:             0.0,
		RatingAverage:               0.0,
		RatingCount:                 0,
		PrivacyLevel:                PrivacyLevelPublic,
		ProfileVisibility:           PrivacyLevelPublic,
		EmailNotifications:          true,
		PushNotifications:           true,
		ProfileCompletionPercentage: 0,
		IsActive:                    true,
		CreatedAt:                   now,
		UpdatedAt:                   now,
	}

	return user, nil
}

// UpdateProfile updates the user's profile information
func (u *User) UpdateProfile(firstName, lastName, bio, location, website string, phone *string, dateOfBirth *time.Time, gender *Gender, travelStyle *TravelStyle) error {
	if firstName == "" {
		return errors.New("first name is required")
	}
	if lastName == "" {
		return errors.New("last name is required")
	}

	u.FirstName = firstName
	u.LastName = lastName
	u.Bio = bio
	u.Location = location
	u.Website = website
	u.Phone = phone
	u.DateOfBirth = dateOfBirth
	u.Gender = gender
	u.TravelStyle = travelStyle
	u.UpdatedAt = time.Now()

	return nil
}

// UpdateLastLogin updates the user's last login timestamp
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
	u.UpdatedAt = now
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

// CanCreateTrips returns true if the user can create trips
func (u *User) CanCreateTrips() bool {
	return u.IsActive
}

// UpdateLanguages updates the user's language preferences
func (u *User) UpdateLanguages(languages []string) {
	u.Languages = languages
	u.UpdatedAt = time.Now()
}

// AddLanguage adds a language to the user's preferences
func (u *User) AddLanguage(language string) {
	// Check if language already exists
	for _, lang := range u.Languages {
		if lang == language {
			return // Already exists
		}
	}
	u.Languages = append(u.Languages, language)
	u.UpdatedAt = time.Now()
}

// RemoveLanguage removes a language from the user's preferences
func (u *User) RemoveLanguage(language string) {
	for i, lang := range u.Languages {
		if lang == language {
			u.Languages = append(u.Languages[:i], u.Languages[i+1:]...)
			u.UpdatedAt = time.Now()
			return
		}
	}
}

// UpdateInterests updates the user's interests
func (u *User) UpdateInterests(interests []string) {
	u.Interests = interests
	u.UpdatedAt = time.Now()
}

// AddInterest adds an interest to the user's preferences
func (u *User) AddInterest(interest string) {
	// Check if interest already exists
	for _, existingInterest := range u.Interests {
		if existingInterest == interest {
			return // Already exists
		}
	}
	u.Interests = append(u.Interests, interest)
	u.UpdatedAt = time.Now()
}

// RemoveInterest removes an interest from the user's preferences
func (u *User) RemoveInterest(interest string) {
	for i, int := range u.Interests {
		if int == interest {
			u.Interests = append(u.Interests[:i], u.Interests[i+1:]...)
			u.UpdatedAt = time.Now()
			return
		}
	}
}

// UpdateNotificationSettings updates the user's notification preferences
func (u *User) UpdateNotificationSettings(emailNotifications, pushNotifications bool) {
	u.EmailNotifications = emailNotifications
	u.PushNotifications = pushNotifications
	u.UpdatedAt = time.Now()
}

// UpdateRating updates the user's rating information
func (u *User) UpdateRating(average float64, count int) {
	u.RatingAverage = average
	u.RatingCount = count
	u.UpdatedAt = time.Now()
}

// CalculateProfileCompletion calculates and updates the profile completion percentage
func (u *User) CalculateProfileCompletion() {
	totalFields := 15
	completedFields := 0

	// Basic required fields
	if u.Email != "" {
		completedFields++
	}
	if u.FirstName != "" {
		completedFields++
	}
	if u.LastName != "" {
		completedFields++
	}
	if u.ProfilePhotoURL != "" {
		completedFields++
	}

	// Extended profile fields
	if u.Bio != "" {
		completedFields++
	}
	if u.Location != "" {
		completedFields++
	}
	if u.DateOfBirth != nil {
		completedFields++
	}
	if u.Phone != nil && *u.Phone != "" {
		completedFields++
	}
	if len(u.Languages) > 0 {
		completedFields++
	}
	if len(u.Interests) > 0 {
		completedFields++
	}
	if u.TravelStyle != nil {
		completedFields++
	}
	if u.RatingCount > 0 {
		completedFields++
	}
	if u.Website != "" {
		completedFields++
	}

	u.ProfileCompletionPercentage = int((float64(completedFields) / float64(totalFields)) * 100)
	u.UpdatedAt = time.Now()
}

// generateUsername generates a unique username from first and last name
func generateUsername(firstName, lastName string) string {
	// Simple implementation - in production, you'd want to ensure uniqueness
	return firstName + lastName + uuid.New().String()[:8]
}

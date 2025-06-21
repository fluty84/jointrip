package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name           string
		googleID       string
		email          string
		firstName      string
		lastName       string
		googlePhotoURL string
		expectError    bool
	}{
		{
			name:           "valid user creation",
			googleID:       "google123",
			email:          "test@example.com",
			firstName:      "John",
			lastName:       "Doe",
			googlePhotoURL: "https://example.com/photo.jpg",
			expectError:    false,
		},
		{
			name:           "missing google ID",
			googleID:       "",
			email:          "test@example.com",
			firstName:      "John",
			lastName:       "Doe",
			googlePhotoURL: "https://example.com/photo.jpg",
			expectError:    true,
		},
		{
			name:           "missing email",
			googleID:       "google123",
			email:          "",
			firstName:      "John",
			lastName:       "Doe",
			googlePhotoURL: "https://example.com/photo.jpg",
			expectError:    true,
		},
		{
			name:           "missing first name",
			googleID:       "google123",
			email:          "test@example.com",
			firstName:      "",
			lastName:       "Doe",
			googlePhotoURL: "https://example.com/photo.jpg",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.googleID, tt.email, tt.firstName, tt.lastName, tt.googlePhotoURL)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)

				assert.NotEqual(t, "", user.ID.String())
				assert.Equal(t, tt.googleID, user.GoogleID)
				assert.Equal(t, tt.email, user.Email)
				assert.Equal(t, tt.firstName, user.FirstName)
				assert.Equal(t, tt.lastName, user.LastName)
				assert.Equal(t, tt.googlePhotoURL, user.GooglePhotoURL)
				assert.Equal(t, tt.googlePhotoURL, user.ProfilePhotoURL)
				assert.Equal(t, VerificationStatusUnverified, user.VerificationStatus)
				assert.Equal(t, float64(0), user.ReputationScore)
				assert.Equal(t, PrivacyLevelPublic, user.PrivacyLevel)
				assert.True(t, user.IsActive)
				assert.NotEmpty(t, user.Username)
			}
		})
	}
}

func TestUser_UpdateProfile(t *testing.T) {
	user, err := NewUser("google123", "test@example.com", "John", "Doe", "photo.jpg")
	require.NoError(t, err)

	originalUpdatedAt := user.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(time.Millisecond)

	phone := "+1234567890"
	dob := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	gender := GenderMale

	err = user.UpdateProfile("Jane", "Smith", "New bio", &phone, &dob, &gender)
	require.NoError(t, err)

	assert.Equal(t, "Jane", user.FirstName)
	assert.Equal(t, "Smith", user.LastName)
	assert.Equal(t, "New bio", user.Bio)
	assert.Equal(t, &phone, user.Phone)
	assert.Equal(t, &dob, user.DateOfBirth)
	assert.Equal(t, &gender, user.Gender)
	assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
}

func TestUser_UpdateProfile_InvalidData(t *testing.T) {
	user, err := NewUser("google123", "test@example.com", "John", "Doe", "photo.jpg")
	require.NoError(t, err)

	// Test empty first name
	err = user.UpdateProfile("", "Smith", "Bio", nil, nil, nil)
	assert.Error(t, err)

	// Test empty last name
	err = user.UpdateProfile("Jane", "", "Bio", nil, nil, nil)
	assert.Error(t, err)
}

func TestUser_UpdateLastLogin(t *testing.T) {
	user, err := NewUser("google123", "test@example.com", "John", "Doe", "photo.jpg")
	require.NoError(t, err)

	assert.Nil(t, user.LastLogin)

	user.UpdateLastLogin()

	assert.NotNil(t, user.LastLogin)
	assert.True(t, time.Since(*user.LastLogin) < time.Second)
}

func TestUser_SetVerificationStatus(t *testing.T) {
	user, err := NewUser("google123", "test@example.com", "John", "Doe", "photo.jpg")
	require.NoError(t, err)

	assert.Equal(t, VerificationStatusUnverified, user.VerificationStatus)

	user.SetVerificationStatus(VerificationStatusVerified)

	assert.Equal(t, VerificationStatusVerified, user.VerificationStatus)
}

func TestUser_SetPrivacyLevel(t *testing.T) {
	user, err := NewUser("google123", "test@example.com", "John", "Doe", "photo.jpg")
	require.NoError(t, err)

	assert.Equal(t, PrivacyLevelPublic, user.PrivacyLevel)

	user.SetPrivacyLevel(PrivacyLevelPrivate)

	assert.Equal(t, PrivacyLevelPrivate, user.PrivacyLevel)
}

func TestUser_Deactivate(t *testing.T) {
	user, err := NewUser("google123", "test@example.com", "John", "Doe", "photo.jpg")
	require.NoError(t, err)

	assert.True(t, user.IsActive)

	user.Deactivate()

	assert.False(t, user.IsActive)
}

func TestUser_Activate(t *testing.T) {
	user, err := NewUser("google123", "test@example.com", "John", "Doe", "photo.jpg")
	require.NoError(t, err)

	user.Deactivate()
	assert.False(t, user.IsActive)

	user.Activate()

	assert.True(t, user.IsActive)
}

func TestUser_IsVerified(t *testing.T) {
	user, err := NewUser("google123", "test@example.com", "John", "Doe", "photo.jpg")
	require.NoError(t, err)

	assert.False(t, user.IsVerified())

	user.SetVerificationStatus(VerificationStatusVerified)

	assert.True(t, user.IsVerified())
}

func TestUser_CanCreateTrips(t *testing.T) {
	user, err := NewUser("google123", "test@example.com", "John", "Doe", "photo.jpg")
	require.NoError(t, err)

	// Initially cannot create trips (not verified)
	assert.False(t, user.CanCreateTrips())

	// Verify user
	user.SetVerificationStatus(VerificationStatusVerified)
	assert.True(t, user.CanCreateTrips())

	// Deactivate user
	user.Deactivate()
	assert.False(t, user.CanCreateTrips())

	// Reactivate but not verified
	user.Activate()
	user.SetVerificationStatus(VerificationStatusUnverified)
	assert.False(t, user.CanCreateTrips())
}

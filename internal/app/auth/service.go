package auth

import (
	"context"
	"errors"
	"time"

	"jointrip/internal/domain/session"
	"jointrip/internal/domain/user"

	"github.com/google/uuid"
)

// GoogleUserInfo represents user information from Google OAuth
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	User         *user.User `json:"user"`
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresAt    time.Time  `json:"expires_at"`
}

// TokenRefreshResponse represents the response after token refresh
type TokenRefreshResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// GoogleOAuthClient defines the interface for Google OAuth operations
type GoogleOAuthClient interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (accessToken, refreshToken string, err error)
	GetUserInfo(ctx context.Context, accessToken string) (*GoogleUserInfo, error)
}

// JWTManager defines the interface for JWT token operations
type JWTManager interface {
	GenerateTokens(userID uuid.UUID) (accessToken, refreshToken string, expiresAt time.Time, err error)
	ValidateAccessToken(tokenString string) (uuid.UUID, error)
	ValidateRefreshToken(tokenString string) (uuid.UUID, error)
}

// Service provides authentication business logic
type Service struct {
	userRepo     user.Repository
	sessionRepo  session.Repository
	googleClient GoogleOAuthClient
	jwtManager   JWTManager
	maxSessions  int
}

// NewService creates a new authentication service
func NewService(
	userRepo user.Repository,
	sessionRepo session.Repository,
	googleClient GoogleOAuthClient,
	jwtManager JWTManager,
	maxSessions int,
) *Service {
	return &Service{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		googleClient: googleClient,
		jwtManager:   jwtManager,
		maxSessions:  maxSessions,
	}
}

// GetGoogleAuthURL returns the Google OAuth authorization URL
func (s *Service) GetGoogleAuthURL(state string) string {
	return s.googleClient.GetAuthURL(state)
}

// LoginWithGoogle handles Google OAuth login
func (s *Service) LoginWithGoogle(ctx context.Context, code, ipAddress, userAgent string) (*LoginResponse, error) {
	// Exchange code for tokens
	googleAccessToken, googleRefreshToken, err := s.googleClient.ExchangeCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// Get user info from Google
	googleUserInfo, err := s.googleClient.GetUserInfo(ctx, googleAccessToken)
	if err != nil {
		return nil, err
	}

	if !googleUserInfo.VerifiedEmail {
		return nil, errors.New("email not verified with Google")
	}

	// Check if user exists
	existingUser, err := s.userRepo.GetByGoogleID(ctx, googleUserInfo.ID)
	if err != nil && !errors.Is(err, user.ErrUserNotFound) {
		return nil, err
	}

	var currentUser *user.User
	if existingUser != nil {
		// Update existing user's last login
		existingUser.UpdateLastLogin()
		if err := s.userRepo.Update(ctx, existingUser); err != nil {
			return nil, err
		}
		currentUser = existingUser
	} else {
		// Create new user
		newUser, err := user.NewUser(
			googleUserInfo.ID,
			googleUserInfo.Email,
			googleUserInfo.GivenName,
			googleUserInfo.FamilyName,
			googleUserInfo.Picture,
		)
		if err != nil {
			return nil, err
		}

		if err := s.userRepo.Create(ctx, newUser); err != nil {
			return nil, err
		}
		currentUser = newUser
	}

	// Generate JWT tokens
	accessToken, refreshToken, expiresAt, err := s.jwtManager.GenerateTokens(currentUser.ID)
	if err != nil {
		return nil, err
	}

	// Create session
	userSession, err := session.NewUserSession(
		currentUser.ID,
		accessToken,
		refreshToken,
		googleAccessToken,
		googleRefreshToken,
		expiresAt,
		ipAddress,
		userAgent,
	)
	if err != nil {
		return nil, err
	}

	// Check session limit
	if err := s.enforceSessionLimit(ctx, currentUser.ID); err != nil {
		return nil, err
	}

	if err := s.sessionRepo.Create(ctx, userSession); err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:         currentUser,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *Service) RefreshToken(ctx context.Context, refreshTokenString string) (*TokenRefreshResponse, error) {
	// Validate refresh token
	userID, err := s.jwtManager.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	// Get session by refresh token
	userSession, err := s.sessionRepo.GetByRefreshToken(ctx, refreshTokenString)
	if err != nil {
		return nil, err
	}

	if !userSession.IsValid() {
		return nil, errors.New("session is invalid or expired")
	}

	if userSession.UserID != userID {
		return nil, errors.New("token user mismatch")
	}

	// Generate new tokens
	accessToken, newRefreshToken, expiresAt, err := s.jwtManager.GenerateTokens(userID)
	if err != nil {
		return nil, err
	}

	// Update session
	userSession.UpdateTokens(accessToken, newRefreshToken, expiresAt)
	if err := s.sessionRepo.Update(ctx, userSession); err != nil {
		return nil, err
	}

	return &TokenRefreshResponse{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}, nil
}

// Logout logs out a user by deactivating their session
func (s *Service) Logout(ctx context.Context, accessToken string) error {
	userSession, err := s.sessionRepo.GetByAccessToken(ctx, accessToken)
	if err != nil {
		return err
	}

	userSession.Deactivate()
	return s.sessionRepo.Update(ctx, userSession)
}

// LogoutAll logs out a user from all sessions
func (s *Service) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	return s.sessionRepo.DeleteByUserID(ctx, userID)
}

// ValidateToken validates an access token and returns the user
func (s *Service) ValidateToken(ctx context.Context, accessToken string) (*user.User, error) {
	// Validate JWT token
	userID, err := s.jwtManager.ValidateAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	// Get session
	userSession, err := s.sessionRepo.GetByAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	if !userSession.IsValid() {
		return nil, errors.New("session is invalid or expired")
	}

	// Update last used
	userSession.UpdateLastUsed()
	if err := s.sessionRepo.Update(ctx, userSession); err != nil {
		// Log error but don't fail the request
	}

	// Get user
	return s.userRepo.GetByID(ctx, userID)
}

// enforceSessionLimit ensures user doesn't exceed maximum sessions
func (s *Service) enforceSessionLimit(ctx context.Context, userID uuid.UUID) error {
	count, err := s.sessionRepo.CountActiveSessionsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if count >= s.maxSessions {
		// Get oldest sessions and deactivate them
		sessions, err := s.sessionRepo.GetActiveSessionsByUserID(ctx, userID)
		if err != nil {
			return err
		}

		// Sort by creation date and deactivate oldest
		for i := 0; i < len(sessions)-(s.maxSessions-1); i++ {
			sessions[i].Deactivate()
			if err := s.sessionRepo.Update(ctx, sessions[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

// UpdateUser updates a user's profile information
func (s *Service) UpdateUser(ctx context.Context, u *user.User) error {
	return s.userRepo.Update(ctx, u)
}

// UpdateUserProfile updates specific user profile fields
func (s *Service) UpdateUserProfile(ctx context.Context, userID uuid.UUID, profileData map[string]interface{}) error {
	// Check if the repository has the UpdateProfile method
	if repo, ok := s.userRepo.(interface {
		UpdateProfile(ctx context.Context, userID uuid.UUID, profileData map[string]interface{}) error
	}); ok {
		return repo.UpdateProfile(ctx, userID, profileData)
	}

	// Fallback to getting the user and updating normally
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Update fields in memory
	for field, value := range profileData {
		switch field {
		case "first_name":
			if str, ok := value.(string); ok {
				u.FirstName = str
			}
		case "last_name":
			if str, ok := value.(string); ok {
				u.LastName = str
			}
		case "bio":
			if str, ok := value.(string); ok {
				u.Bio = str
			}
		case "location":
			if str, ok := value.(string); ok {
				u.Location = str
			}
		case "website":
			if str, ok := value.(string); ok {
				u.Website = str
			}
		case "phone":
			if str, ok := value.(*string); ok {
				u.Phone = str
			}
		case "languages":
			if strSlice, ok := value.([]string); ok {
				u.Languages = strSlice
			}
		case "interests":
			if strSlice, ok := value.([]string); ok {
				u.Interests = strSlice
			}
		}
	}

	return s.userRepo.Update(ctx, u)
}

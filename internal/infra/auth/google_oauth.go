package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"jointrip/internal/app/auth"
	"jointrip/internal/infra/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleOAuthClient implements the auth.GoogleOAuthClient interface
type GoogleOAuthClient struct {
	config *oauth2.Config
}

// NewGoogleOAuthClient creates a new Google OAuth client
func NewGoogleOAuthClient(cfg *config.Config) *GoogleOAuthClient {
	return &GoogleOAuthClient{
		config: &oauth2.Config{
			ClientID:     cfg.Google.ClientID,
			ClientSecret: cfg.Google.ClientSecret,
			RedirectURL:  cfg.Google.RedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

// GetAuthURL returns the Google OAuth authorization URL
func (g *GoogleOAuthClient) GetAuthURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

// ExchangeCode exchanges an authorization code for tokens
func (g *GoogleOAuthClient) ExchangeCode(ctx context.Context, code string) (accessToken, refreshToken string, err error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return "", "", fmt.Errorf("failed to exchange code: %w", err)
	}

	accessToken = token.AccessToken
	refreshToken = token.RefreshToken

	return accessToken, refreshToken, nil
}

// GetUserInfo retrieves user information from Google
func (g *GoogleOAuthClient) GetUserInfo(ctx context.Context, accessToken string) (*auth.GoogleUserInfo, error) {
	client := &http.Client{}
	
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info, status: %d", resp.StatusCode)
	}

	var userInfo auth.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

// RefreshToken refreshes an access token using a refresh token
func (g *GoogleOAuthClient) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	tokenSource := g.config.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return newToken, nil
}

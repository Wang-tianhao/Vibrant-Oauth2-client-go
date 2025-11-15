package vibrant

import "time"

// TokenResponse represents the OAuth2 token response from Vibrant
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope,omitempty"`
}

// CachedToken represents a cached access token with expiration tracking
type CachedToken struct {
	AccessToken string
	ExpiresAt   time.Time
}

// IsExpired checks if the cached token has expired or will expire soon
// It adds a 60-second buffer to avoid using tokens that are about to expire
func (t *CachedToken) IsExpired() bool {
	return time.Now().Add(60 * time.Second).After(t.ExpiresAt)
}

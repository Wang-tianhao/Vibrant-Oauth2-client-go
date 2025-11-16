package vibrant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

const (
	// VibrantTokenEndpoint is the OAuth2 token endpoint for Vibrant
	VibrantTokenEndpoint = "https://api.vibrant-wellness.com/v1/oauth2/token"

	// Environment variable names
	EnvClientID     = "VIBRANT_CLIENT_ID"
	EnvClientSecret = "VIBRANT_CLIENT_SECRET"
)

// Client represents the Vibrant OAuth2 client
type Client struct {
	clientID     string
	clientSecret string
	httpClient   *http.Client
	cache        *CachedToken
	mu           sync.RWMutex
}

// NewClient creates a new Vibrant OAuth2 client
// It reads credentials from environment variables:
// - VIBRANT_CLIENT_ID
// - VIBRANT_CLIENT_SECRET
func NewClient() (*Client, error) {
	clientID := os.Getenv(EnvClientID)
	if clientID == "" {
		return nil, fmt.Errorf("environment variable %s is not set", EnvClientID)
	}

	clientSecret := os.Getenv(EnvClientSecret)
	if clientSecret == "" {
		return nil, fmt.Errorf("environment variable %s is not set", EnvClientSecret)
	}

	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// GetToken returns a valid access token, either from cache or by fetching a new one
// This is the main function that developers should use
func (c *Client) GetToken() (string, error) {
	// Try to use cached token first (read lock)
	c.mu.RLock()
	if c.cache != nil && !c.cache.IsExpired() {
		token := c.cache.AccessToken
		c.mu.RUnlock()
		return token, nil
	}
	c.mu.RUnlock()

	// Need to fetch new token (write lock)
	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock (another goroutine might have refreshed)
	if c.cache != nil && !c.cache.IsExpired() {
		return c.cache.AccessToken, nil
	}

	// Fetch new token
	token, err := c.fetchToken()
	if err != nil {
		return "", err
	}

	return token, nil
}

// fetchToken fetches a new access token from Vibrant OAuth endpoint
// This method should be called while holding the write lock
func (c *Client) fetchToken() (string, error) {
	// Prepare request body
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)

	// Create request
	req, err := http.NewRequest("POST", VibrantTokenEndpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Cache the token
	c.cache = &CachedToken{
		AccessToken: tokenResp.TokenType + " " + tokenResp.AccessToken,
		ExpiresAt:   time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
	}

	return c.cache.AccessToken, nil
}

// ClearCache clears the cached token, forcing a new token fetch on next GetToken call
func (c *Client) ClearCache() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = nil
}

package credentials

import "errors"

// ErrNotConfigured is returned when no credentials are stored for the requested profile.
var ErrNotConfigured = errors.New("credentials not configured")

// Credentials holds the API key, token, and auth mode for a profile.
type Credentials struct {
	APIKey   string `json:"apiKey"`
	Token    string `json:"token"`
	AuthMode string `json:"authMode"`
}

// Store defines the interface for credential persistence.
type Store interface {
	Get(profile string) (Credentials, error)
	Set(profile string, creds Credentials) error
	Delete(profile string) error
}

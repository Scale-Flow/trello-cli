package auth

import "github.com/brettmcdowell/trello-cli/internal/credentials"

// ClearResult is the response shape for auth clear.
type ClearResult struct {
	Configured bool    `json:"configured"`
	AuthMode   *string `json:"authMode"` // Always null after clear
}

// Clear removes stored credentials for the given profile.
func Clear(store credentials.Store, profile string) (ClearResult, error) {
	if err := store.Delete(profile); err != nil {
		return ClearResult{}, err
	}
	return ClearResult{
		Configured: false,
		AuthMode:   nil,
	}, nil
}

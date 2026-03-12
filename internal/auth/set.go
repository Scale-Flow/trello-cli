package auth

import "github.com/brettmcdowell/trello-cli/internal/credentials"

// SetResult is the response shape for auth set.
type SetResult struct {
	Configured bool        `json:"configured"`
	AuthMode   string      `json:"authMode"`
	Member     interface{} `json:"-"` // Never included — auth set does not validate
}

// Set stores credentials without validating them against the Trello API.
// Validation is deferred to auth status or the first command that requires auth.
func Set(store credentials.Store, profile, apiKey, token string) (SetResult, error) {
	creds := credentials.Credentials{
		APIKey:   apiKey,
		Token:    token,
		AuthMode: "manual",
	}
	if err := store.Set(profile, creds); err != nil {
		return SetResult{}, err
	}
	return SetResult{
		Configured: true,
		AuthMode:   "manual",
	}, nil
}

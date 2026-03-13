package auth

import "github.com/brettmcdowell/trello-cli/internal/credentials"

// SetResult is the response shape for auth set.
type SetResult struct {
	Configured bool        `json:"configured"`
	AuthMode   string      `json:"authMode"`
	Member     interface{} `json:"-"` // Never included — auth set does not validate
}

// SetKeyResult is the response shape for auth set-key.
type SetKeyResult struct {
	Configured bool   `json:"configured"`
	AuthMode   string `json:"authMode"`
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

// SetKey stores or updates only the API key for the given profile.
// If a token already exists, it is preserved so the profile remains configured.
func SetKey(store credentials.Store, profile, apiKey string) (SetKeyResult, error) {
	creds, err := store.Get(profile)
	if err != nil {
		creds = credentials.Credentials{}
	}

	creds.APIKey = apiKey
	if creds.Token == "" {
		creds.AuthMode = "key_only"
	} else if creds.AuthMode == "" || creds.AuthMode == "key_only" {
		creds.AuthMode = "manual"
	}

	if err := store.Set(profile, creds); err != nil {
		return SetKeyResult{}, err
	}

	return SetKeyResult{
		Configured: creds.Token != "",
		AuthMode:   creds.AuthMode,
	}, nil
}

package auth

import (
	"errors"

	"github.com/brettmcdowell/trello-cli/internal/contract"
	"github.com/brettmcdowell/trello-cli/internal/credentials"
)

// RequireAuth loads credentials for the given profile. Returns AUTH_REQUIRED if missing.
func RequireAuth(store credentials.Store, profile string) (credentials.Credentials, error) {
	creds, err := store.Get(profile)
	if err != nil {
		if errors.Is(err, credentials.ErrNotConfigured) {
			return credentials.Credentials{}, contract.NewError(contract.AuthRequired, "not authenticated — run 'trello auth login' or 'trello auth set'")
		}
		return credentials.Credentials{}, err
	}
	return creds, nil
}

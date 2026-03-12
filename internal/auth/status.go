package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/brettmcdowell/trello-cli/internal/contract"
	"github.com/brettmcdowell/trello-cli/internal/credentials"
)

// Member represents a Trello member for auth status responses.
type Member struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
}

// StatusResult is the response shape for auth status.
type StatusResult struct {
	Configured bool    `json:"configured"`
	AuthMode   *string `json:"authMode"`
	Member     *Member `json:"member"`
}

// Status checks the current authentication state.
func Status(ctx context.Context, store credentials.Store, profile, baseURL string) (StatusResult, error) {
	creds, err := store.Get(profile)
	if err != nil {
		if errors.Is(err, credentials.ErrNotConfigured) {
			return StatusResult{
				Configured: false,
				AuthMode:   nil,
				Member:     nil,
			}, nil
		}
		return StatusResult{}, err
	}

	member, err := getMember(ctx, baseURL, creds.APIKey, creds.Token)
	if err != nil {
		return StatusResult{}, err
	}

	authMode := creds.AuthMode
	return StatusResult{
		Configured: true,
		AuthMode:   &authMode,
		Member:     member,
	}, nil
}

func getMember(ctx context.Context, baseURL, apiKey, token string) (*Member, error) {
	url := fmt.Sprintf("%s/1/members/me?key=%s&token=%s", baseURL, apiKey, token)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, contract.NewError(contract.HTTPError, fmt.Sprintf("failed to reach Trello API: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, contract.NewError(contract.AuthInvalid, "Trello rejected the credentials — API key or token is invalid")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, contract.NewError(contract.HTTPError, fmt.Sprintf("Trello API returned status %d", resp.StatusCode))
	}

	var member Member
	if err := json.NewDecoder(resp.Body).Decode(&member); err != nil {
		return nil, contract.NewError(contract.HTTPError, fmt.Sprintf("failed to decode member response: %v", err))
	}

	return &member, nil
}

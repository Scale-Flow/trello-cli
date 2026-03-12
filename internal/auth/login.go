package auth

import (
	"context"
	"fmt"
	"net/url"

	"github.com/brettmcdowell/trello-cli/internal/credentials"
)

const trelloAuthorizeBase = "https://trello.com/1/authorize"

// LoginResult is the response shape for auth login.
type LoginResult struct {
	Configured bool    `json:"configured"`
	AuthMode   *string `json:"authMode"`
	Member     *Member `json:"member"`
}

// BuildAuthorizeURL builds the Trello authorization URL for interactive login.
func BuildAuthorizeURL(apiKey, callbackURL string) string {
	params := url.Values{
		"expiration":    {"never"},
		"name":          {"Trello CLI"},
		"scope":         {"read,write"},
		"response_type": {"token"},
		"key":           {apiKey},
		"return_url":    {callbackURL},
	}
	return trelloAuthorizeBase + "?" + params.Encode()
}

// CompleteLogin validates a captured token and stores credentials.
func CompleteLogin(ctx context.Context, store credentials.Store, profile, apiKey, token, baseURL string) (LoginResult, error) {
	member, err := getMember(ctx, baseURL, apiKey, token)
	if err != nil {
		return LoginResult{}, err
	}

	creds := credentials.Credentials{
		APIKey:   apiKey,
		Token:    token,
		AuthMode: "interactive",
	}
	if err := store.Set(profile, creds); err != nil {
		return LoginResult{}, fmt.Errorf("failed to store credentials: %w", err)
	}

	authMode := "interactive"
	return LoginResult{
		Configured: true,
		AuthMode:   &authMode,
		Member:     member,
	}, nil
}

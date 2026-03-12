package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/auth"
	"github.com/brettmcdowell/trello-cli/internal/credentials"
)

func TestLoginWithToken(t *testing.T) {
	// Mock Trello API for credential validation
	trelloServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/1/members/me" {
			if got := r.URL.Query().Get("key"); got != "test-api-key" {
				t.Errorf("key query = %q, want %q", got, "test-api-key")
			}
			if got := r.URL.Query().Get("token"); got != "captured-token" {
				t.Errorf("token query = %q, want %q", got, "captured-token")
			}
			json.NewEncoder(w).Encode(map[string]string{
				"id":       "member456",
				"username": "loginuser",
				"fullName": "Login User",
			})
			return
		}
		http.Error(w, "not found", 404)
	}))
	defer trelloServer.Close()

	store := credentials.NewMemoryStore()

	// Simulate the login flow by directly providing a token
	result, err := auth.CompleteLogin(
		context.Background(),
		store,
		"default",
		"test-api-key",
		"captured-token",
		trelloServer.URL,
	)
	if err != nil {
		t.Fatalf("CompleteLogin() returned error: %v", err)
	}

	if result.Configured != true {
		t.Errorf("Configured = %v, want true", result.Configured)
	}
	if *result.AuthMode != "interactive" {
		t.Errorf("AuthMode = %q, want %q", *result.AuthMode, "interactive")
	}
	if result.Member == nil {
		t.Fatal("Member should not be nil after login")
	}
	if result.Member.Username != "loginuser" {
		t.Errorf("Member.Username = %q, want %q", result.Member.Username, "loginuser")
	}

	creds, _ := store.Get("default")
	if creds.APIKey != "test-api-key" {
		t.Errorf("stored APIKey = %q, want %q", creds.APIKey, "test-api-key")
	}
	if creds.Token != "captured-token" {
		t.Errorf("stored Token = %q, want %q", creds.Token, "captured-token")
	}
	if creds.AuthMode != "interactive" {
		t.Errorf("stored AuthMode = %q, want %q", creds.AuthMode, "interactive")
	}
}

func TestLoginInvalidToken(t *testing.T) {
	trelloServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("key"); got != "bad-key" {
			t.Errorf("key query = %q, want %q", got, "bad-key")
		}
		if got := r.URL.Query().Get("token"); got != "bad-token" {
			t.Errorf("token query = %q, want %q", got, "bad-token")
		}
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer trelloServer.Close()

	store := credentials.NewMemoryStore()

	_, err := auth.CompleteLogin(
		context.Background(),
		store,
		"default",
		"bad-key",
		"bad-token",
		trelloServer.URL,
	)
	if err == nil {
		t.Fatal("CompleteLogin() should return error for invalid credentials")
	}

	_, getErr := store.Get("default")
	if getErr == nil {
		t.Error("credentials should not be stored after failed login")
	}
}

func TestBuildAuthorizeURL(t *testing.T) {
	rawURL := auth.BuildAuthorizeURL("my-api-key", "http://localhost:12345/callback")

	if rawURL == "" {
		t.Fatal("BuildAuthorizeURL() returned empty string")
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("url.Parse() returned error: %v", err)
	}
	if parsed.Scheme != "https" {
		t.Errorf("scheme = %q, want %q", parsed.Scheme, "https")
	}
	if parsed.Host != "trello.com" {
		t.Errorf("host = %q, want %q", parsed.Host, "trello.com")
	}
	if parsed.Path != "/1/authorize" {
		t.Errorf("path = %q, want %q", parsed.Path, "/1/authorize")
	}
	query := parsed.Query()
	if got := query.Get("key"); got != "my-api-key" {
		t.Errorf("key = %q, want %q", got, "my-api-key")
	}
	if got := query.Get("return_url"); got != "http://localhost:12345/callback" {
		t.Errorf("return_url = %q, want %q", got, "http://localhost:12345/callback")
	}
	if got := query.Get("callback_method"); got != "fragment" {
		t.Errorf("callback_method = %q, want %q", got, "fragment")
	}
	if got := query.Get("response_type"); got != "token" {
		t.Errorf("response_type = %q, want %q", got, "token")
	}
	if got := query.Get("scope"); got != "read,write" {
		t.Errorf("scope = %q, want %q", got, "read,write")
	}
	if got := query.Get("expiration"); got != "never" {
		t.Errorf("expiration = %q, want %q", got, "never")
	}
}

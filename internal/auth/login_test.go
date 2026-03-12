package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/auth"
	"github.com/brettmcdowell/trello-cli/internal/credentials"
)

func TestLoginWithToken(t *testing.T) {
	// Mock Trello API for credential validation
	trelloServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/1/members/me" {
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
	url := auth.BuildAuthorizeURL("my-api-key", "http://localhost:12345/callback")

	if url == "" {
		t.Fatal("BuildAuthorizeURL() returned empty string")
	}
	if !strings.Contains(url, "key=my-api-key") {
		t.Errorf("URL should contain API key, got: %s", url)
	}
	if !strings.Contains(url, "return_url=") {
		t.Errorf("URL should contain return_url, got: %s", url)
	}
}

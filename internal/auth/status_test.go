package auth_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/auth"
	"github.com/brettmcdowell/trello-cli/internal/contract"
	"github.com/brettmcdowell/trello-cli/internal/credentials"
)

func TestStatusConfigured(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/1/members/me" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("key"); got != "key1" {
			t.Errorf("key query = %q, want %q", got, "key1")
		}
		if got := r.URL.Query().Get("token"); got != "tok1" {
			t.Errorf("token query = %q, want %q", got, "tok1")
		}
		json.NewEncoder(w).Encode(map[string]string{
			"id":       "member123",
			"username": "testuser",
			"fullName": "Test User",
		})
	}))
	defer server.Close()

	store := credentials.NewMemoryStore()
	store.Set("default", credentials.Credentials{
		APIKey:   "key1",
		Token:    "tok1",
		AuthMode: "manual",
	})

	result, err := auth.Status(context.Background(), store, "default", server.URL)
	if err != nil {
		t.Fatalf("Status() returned error: %v", err)
	}

	if result.Configured != true {
		t.Errorf("Configured = %v, want true", result.Configured)
	}
	if *result.AuthMode != "manual" {
		t.Errorf("AuthMode = %q, want %q", *result.AuthMode, "manual")
	}
	if result.Member == nil {
		t.Fatal("Member should not be nil when configured")
	}
	if result.Member.ID != "member123" {
		t.Errorf("Member.ID = %q, want %q", result.Member.ID, "member123")
	}
	if result.Member.Username != "testuser" {
		t.Errorf("Member.Username = %q, want %q", result.Member.Username, "testuser")
	}
}

func TestStatusNotConfigured(t *testing.T) {
	store := credentials.NewMemoryStore()

	result, err := auth.Status(context.Background(), store, "default", "")
	if err != nil {
		t.Fatalf("Status() returned error: %v", err)
	}

	if result.Configured != false {
		t.Errorf("Configured = %v, want false", result.Configured)
	}
	if result.AuthMode != nil {
		t.Errorf("AuthMode = %v, want nil", result.AuthMode)
	}
	if result.Member != nil {
		t.Error("Member should be nil when not configured")
	}
}

func TestStatusAuthModeNullJSON(t *testing.T) {
	store := credentials.NewMemoryStore()
	result, _ := auth.Status(context.Background(), store, "default", "")

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json.Marshal() returned error: %v", err)
	}

	var m map[string]any
	json.Unmarshal(data, &m)

	if m["authMode"] != nil {
		t.Errorf("authMode = %v, want null", m["authMode"])
	}
}

func TestStatusInvalidCredentials(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
	}))
	defer server.Close()

	store := credentials.NewMemoryStore()
	store.Set("default", credentials.Credentials{
		APIKey:   "bad-key",
		Token:    "bad-token",
		AuthMode: "manual",
	})

	_, err := auth.Status(context.Background(), store, "default", server.URL)
	if err == nil {
		t.Fatal("Status() should return error for invalid credentials")
	}
	var ce *contract.ContractError
	if !errors.As(err, &ce) {
		t.Fatalf("error should be *ContractError, got %T", err)
	}
	if ce.Code != contract.AuthInvalid {
		t.Errorf("Code = %q, want %q", ce.Code, contract.AuthInvalid)
	}
}

func TestStatusKeyOnlyNotConfigured(t *testing.T) {
	store := credentials.NewMemoryStore()
	store.Set("default", credentials.Credentials{
		APIKey:   "key-only",
		Token:    "",
		AuthMode: "key_only",
	})

	result, err := auth.Status(context.Background(), store, "default", "http://example.com")
	if err != nil {
		t.Fatalf("Status() returned error: %v", err)
	}

	if result.Configured != false {
		t.Errorf("Configured = %v, want false", result.Configured)
	}
	if result.Member != nil {
		t.Errorf("Member = %#v, want nil", result.Member)
	}
	if result.AuthMode == nil || *result.AuthMode != "key_only" {
		t.Fatalf("AuthMode = %v, want key_only", result.AuthMode)
	}
}

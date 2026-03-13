package auth_test

import (
	"encoding/json"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/auth"
	"github.com/brettmcdowell/trello-cli/internal/credentials"
)

func TestSetStoresCredentials(t *testing.T) {
	store := credentials.NewMemoryStore()

	result, err := auth.Set(store, "default", "test-key", "test-token")
	if err != nil {
		t.Fatalf("Set() returned error: %v", err)
	}

	if result.Configured != true {
		t.Errorf("Configured = %v, want true", result.Configured)
	}
	if result.AuthMode != "manual" {
		t.Errorf("AuthMode = %q, want %q", result.AuthMode, "manual")
	}

	creds, err := store.Get("default")
	if err != nil {
		t.Fatalf("store.Get() returned error: %v", err)
	}
	if creds.APIKey != "test-key" {
		t.Errorf("stored APIKey = %q, want %q", creds.APIKey, "test-key")
	}
	if creds.Token != "test-token" {
		t.Errorf("stored Token = %q, want %q", creds.Token, "test-token")
	}
}

func TestSetDoesNotCallTrelloAPI(t *testing.T) {
	store := credentials.NewMemoryStore()
	_, err := auth.Set(store, "default", "any-key", "any-token")
	if err != nil {
		t.Fatalf("Set() should not fail (no API validation): %v", err)
	}
}

func TestSetResultHasNoMemberFieldInJSON(t *testing.T) {
	store := credentials.NewMemoryStore()
	result, _ := auth.Set(store, "default", "k", "t")

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json.Marshal() returned error: %v", err)
	}
	var m map[string]any
	json.Unmarshal(data, &m)
	if _, exists := m["member"]; exists {
		t.Error("auth set response should not include 'member' field in JSON")
	}
}

func TestSetKeyStoresKeyOnlyCredentials(t *testing.T) {
	store := credentials.NewMemoryStore()

	result, err := auth.SetKey(store, "default", "test-key")
	if err != nil {
		t.Fatalf("SetKey() returned error: %v", err)
	}

	if result.Configured != false {
		t.Errorf("Configured = %v, want false", result.Configured)
	}
	if result.AuthMode != "key_only" {
		t.Errorf("AuthMode = %q, want %q", result.AuthMode, "key_only")
	}

	creds, err := store.Get("default")
	if err != nil {
		t.Fatalf("store.Get() returned error: %v", err)
	}
	if creds.APIKey != "test-key" {
		t.Errorf("stored APIKey = %q, want %q", creds.APIKey, "test-key")
	}
	if creds.Token != "" {
		t.Errorf("stored Token = %q, want empty string", creds.Token)
	}
	if creds.AuthMode != "key_only" {
		t.Errorf("stored AuthMode = %q, want %q", creds.AuthMode, "key_only")
	}
}

func TestSetKeyPreservesExistingToken(t *testing.T) {
	store := credentials.NewMemoryStore()
	if err := store.Set("default", credentials.Credentials{
		APIKey:   "old-key",
		Token:    "existing-token",
		AuthMode: "manual",
	}); err != nil {
		t.Fatalf("store.Set() returned error: %v", err)
	}

	result, err := auth.SetKey(store, "default", "new-key")
	if err != nil {
		t.Fatalf("SetKey() returned error: %v", err)
	}

	if result.Configured != true {
		t.Errorf("Configured = %v, want true", result.Configured)
	}
	if result.AuthMode != "manual" {
		t.Errorf("AuthMode = %q, want %q", result.AuthMode, "manual")
	}

	creds, err := store.Get("default")
	if err != nil {
		t.Fatalf("store.Get() returned error: %v", err)
	}
	if creds.APIKey != "new-key" {
		t.Errorf("stored APIKey = %q, want %q", creds.APIKey, "new-key")
	}
	if creds.Token != "existing-token" {
		t.Errorf("stored Token = %q, want %q", creds.Token, "existing-token")
	}
	if creds.AuthMode != "manual" {
		t.Errorf("stored AuthMode = %q, want %q", creds.AuthMode, "manual")
	}
}

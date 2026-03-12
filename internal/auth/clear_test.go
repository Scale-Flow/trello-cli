package auth_test

import (
	"encoding/json"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/auth"
	"github.com/brettmcdowell/trello-cli/internal/credentials"
)

func TestClearRemovesCredentials(t *testing.T) {
	store := credentials.NewMemoryStore()
	store.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	result, err := auth.Clear(store, "default")
	if err != nil {
		t.Fatalf("Clear() returned error: %v", err)
	}

	if result.Configured != false {
		t.Errorf("Configured = %v, want false", result.Configured)
	}
	if result.AuthMode != nil {
		t.Errorf("AuthMode = %v, want nil", result.AuthMode)
	}

	// Verify credentials were removed
	_, err = store.Get("default")
	if err == nil {
		t.Error("credentials should be removed after Clear()")
	}
}

func TestClearWhenNotConfigured(t *testing.T) {
	store := credentials.NewMemoryStore()

	result, err := auth.Clear(store, "default")
	if err != nil {
		t.Fatalf("Clear() returned error: %v", err)
	}

	if result.Configured != false {
		t.Errorf("Configured = %v, want false", result.Configured)
	}
}

func TestClearAuthModeNullJSON(t *testing.T) {
	store := credentials.NewMemoryStore()
	result, _ := auth.Clear(store, "default")

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json.Marshal() returned error: %v", err)
	}

	var m map[string]any
	json.Unmarshal(data, &m)

	if m["authMode"] != nil {
		t.Errorf("authMode = %v, want null in JSON", m["authMode"])
	}
}

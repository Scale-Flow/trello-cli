package main

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/credentials"
)

func setupTestAuth(t *testing.T) {
	t.Helper()
	// Use a memory store for tests
	credStore = credentials.NewMemoryStore()
	apiClient = nil
	// Reset root command output state to avoid cross-test contamination
	rootCmd.SetOut(nil)
	rootCmd.SetArgs(nil)
	if err := authSetCmd.Flags().Set("api-key", ""); err != nil {
		t.Fatalf("failed to reset api-key flag: %v", err)
	}
	if err := authSetCmd.Flags().Set("token", ""); err != nil {
		t.Fatalf("failed to reset token flag: %v", err)
	}
}

func TestAuthSetCommand(t *testing.T) {
	setupTestAuth(t)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"auth", "set", "--api-key", "test-key", "--token", "test-token"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("auth set failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}

	if envelope["ok"] != true {
		t.Errorf("ok = %v, want true", envelope["ok"])
	}

	data := envelope["data"].(map[string]any)
	if data["configured"] != true {
		t.Errorf("configured = %v, want true", data["configured"])
	}
	if data["authMode"] != "manual" {
		t.Errorf("authMode = %v, want manual", data["authMode"])
	}
}

func TestAuthSetMissingFlags(t *testing.T) {
	setupTestAuth(t)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"auth", "set", "--api-key", "test-key"})
	// Missing --token

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("auth set should fail without --token")
	}
}

func TestAuthClearCommand(t *testing.T) {
	setupTestAuth(t)
	// Pre-set credentials
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"auth", "clear"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("auth clear failed: %v", err)
	}

	var envelope map[string]any
	json.Unmarshal(buf.Bytes(), &envelope)

	if envelope["ok"] != true {
		t.Errorf("ok = %v, want true", envelope["ok"])
	}

	data := envelope["data"].(map[string]any)
	if data["configured"] != false {
		t.Errorf("configured = %v, want false", data["configured"])
	}
	if data["authMode"] != nil {
		t.Errorf("authMode = %v, want nil", data["authMode"])
	}
}

func TestAuthLoginReturnsUnsupported(t *testing.T) {
	setupTestAuth(t)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"auth", "login"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("auth login should return error (not yet implemented)")
	}
	// The error handler should produce an UNSUPPORTED error envelope
}

func TestAuthStatusNotConfigured(t *testing.T) {
	setupTestAuth(t)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"auth", "status"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("auth status failed: %v", err)
	}

	var envelope map[string]any
	json.Unmarshal(buf.Bytes(), &envelope)

	if envelope["ok"] != true {
		t.Errorf("ok = %v, want true", envelope["ok"])
	}

	data := envelope["data"].(map[string]any)
	if data["configured"] != false {
		t.Errorf("configured = %v, want false", data["configured"])
	}
}

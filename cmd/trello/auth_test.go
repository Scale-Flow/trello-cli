package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Scale-Flow/trello-cli/internal/auth"
	"github.com/Scale-Flow/trello-cli/internal/credentials"
)

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

func TestAuthSetKeyCommand(t *testing.T) {
	setupTestAuth(t)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"auth", "set-key", "--api-key", "test-key"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("auth set-key failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}

	if envelope["ok"] != true {
		t.Errorf("ok = %v, want true", envelope["ok"])
	}

	data := envelope["data"].(map[string]any)
	if data["configured"] != false {
		t.Errorf("configured = %v, want false", data["configured"])
	}
	if data["authMode"] != "key_only" {
		t.Errorf("authMode = %v, want key_only", data["authMode"])
	}
}

func TestAuthSetKeyMissingFlag(t *testing.T) {
	setupTestAuth(t)

	rootCmd.SetArgs([]string{"auth", "set-key"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("auth set-key should fail without --api-key")
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

func TestAuthLoginCommand(t *testing.T) {
	setupTestAuth(t)

	// Mock device flow to fail so it falls through to browser login
	runAuthLoginWithDeviceFlow = func(_ loginCommandContext, _ credentials.Store, _, _, _ string, _ loginOutputWriter) (auth.LoginResult, error) {
		return auth.LoginResult{}, fmt.Errorf("pairing service unavailable")
	}
	runAuthLogin = func(_ loginCommandContext, _ credentials.Store, _ string, _ string, _ string, _ loginBrowserOpener, _ loginOutputWriter) (auth.LoginResult, error) {
		authMode := "interactive"
		return auth.LoginResult{
			Configured: true,
			AuthMode:   &authMode,
			Member: &auth.Member{
				ID:       "member123",
				Username: "brett",
				FullName: "Brett McDowell",
			},
		}, nil
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"auth", "login"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("auth login failed: %v", err)
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
	if data["authMode"] != "interactive" {
		t.Errorf("authMode = %v, want interactive", data["authMode"])
	}
}

func TestAuthLoginDeviceFlowCommand(t *testing.T) {
	setupTestAuth(t)

	// Mock device flow to succeed
	runAuthLoginWithDeviceFlow = func(_ loginCommandContext, _ credentials.Store, _, _, _ string, _ loginOutputWriter) (auth.LoginResult, error) {
		authMode := "device"
		return auth.LoginResult{
			Configured: true,
			AuthMode:   &authMode,
			Member: &auth.Member{
				ID:       "member123",
				Username: "brett",
				FullName: "Brett McDowell",
			},
		}, nil
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"auth", "login"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("auth login (device flow) failed: %v", err)
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
	if data["authMode"] != "device" {
		t.Errorf("authMode = %v, want device", data["authMode"])
	}
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

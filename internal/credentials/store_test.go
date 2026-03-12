package credentials_test

import (
	"errors"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/credentials"
)

func TestErrNotConfiguredSentinel(t *testing.T) {
	err := credentials.ErrNotConfigured
	if err == nil {
		t.Fatal("ErrNotConfigured should not be nil")
	}
	if !errors.Is(err, credentials.ErrNotConfigured) {
		t.Error("errors.Is should match ErrNotConfigured")
	}
}

func TestCredentialsStruct(t *testing.T) {
	creds := credentials.Credentials{
		APIKey:   "test-key",
		Token:    "test-token",
		AuthMode: "manual",
	}
	if creds.APIKey != "test-key" {
		t.Errorf("APIKey = %q, want %q", creds.APIKey, "test-key")
	}
	if creds.Token != "test-token" {
		t.Errorf("Token = %q, want %q", creds.Token, "test-token")
	}
	if creds.AuthMode != "manual" {
		t.Errorf("AuthMode = %q, want %q", creds.AuthMode, "manual")
	}
}

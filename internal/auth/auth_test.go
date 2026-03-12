package auth_test

import (
	"errors"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/auth"
	"github.com/brettmcdowell/trello-cli/internal/contract"
	"github.com/brettmcdowell/trello-cli/internal/credentials"
)

func TestRequireAuthWithCredentials(t *testing.T) {
	store := credentials.NewMemoryStore()
	store.Set("default", credentials.Credentials{
		APIKey:   "key1",
		Token:    "tok1",
		AuthMode: "manual",
	})

	creds, err := auth.RequireAuth(store, "default")
	if err != nil {
		t.Fatalf("RequireAuth() returned error: %v", err)
	}
	if creds.APIKey != "key1" {
		t.Errorf("APIKey = %q, want %q", creds.APIKey, "key1")
	}
	if creds.Token != "tok1" {
		t.Errorf("Token = %q, want %q", creds.Token, "tok1")
	}
}

func TestRequireAuthWithoutCredentials(t *testing.T) {
	store := credentials.NewMemoryStore()

	_, err := auth.RequireAuth(store, "default")
	if err == nil {
		t.Fatal("RequireAuth() should return error when no credentials")
	}

	var ce *contract.ContractError
	if !errors.As(err, &ce) {
		t.Fatalf("error should be *ContractError, got %T", err)
	}
	if ce.Code != contract.AuthRequired {
		t.Errorf("Code = %q, want %q", ce.Code, contract.AuthRequired)
	}
}

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

func TestMemoryStoreSetAndGet(t *testing.T) {
	store := credentials.NewMemoryStore()
	creds := credentials.Credentials{
		APIKey:   "key1",
		Token:    "tok1",
		AuthMode: "manual",
	}

	if err := store.Set("default", creds); err != nil {
		t.Fatalf("Set() returned error: %v", err)
	}

	got, err := store.Get("default")
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if got.APIKey != "key1" {
		t.Errorf("APIKey = %q, want %q", got.APIKey, "key1")
	}
	if got.Token != "tok1" {
		t.Errorf("Token = %q, want %q", got.Token, "tok1")
	}
	if got.AuthMode != "manual" {
		t.Errorf("AuthMode = %q, want %q", got.AuthMode, "manual")
	}
}

func TestMemoryStoreGetMissing(t *testing.T) {
	store := credentials.NewMemoryStore()
	_, err := store.Get("nonexistent")
	if !errors.Is(err, credentials.ErrNotConfigured) {
		t.Errorf("Get() error = %v, want ErrNotConfigured", err)
	}
}

func TestMemoryStoreDelete(t *testing.T) {
	store := credentials.NewMemoryStore()
	creds := credentials.Credentials{APIKey: "key1", Token: "tok1", AuthMode: "manual"}
	store.Set("default", creds)

	if err := store.Delete("default"); err != nil {
		t.Fatalf("Delete() returned error: %v", err)
	}

	_, err := store.Get("default")
	if !errors.Is(err, credentials.ErrNotConfigured) {
		t.Errorf("Get() after Delete() error = %v, want ErrNotConfigured", err)
	}
}

func TestMemoryStoreDeleteMissing(t *testing.T) {
	store := credentials.NewMemoryStore()
	// Deleting a nonexistent profile should not error
	if err := store.Delete("nonexistent"); err != nil {
		t.Errorf("Delete() on missing profile returned error: %v", err)
	}
}

func TestMemoryStoreOverwrite(t *testing.T) {
	store := credentials.NewMemoryStore()
	store.Set("default", credentials.Credentials{APIKey: "old", Token: "old", AuthMode: "manual"})
	store.Set("default", credentials.Credentials{APIKey: "new", Token: "new", AuthMode: "interactive"})

	got, _ := store.Get("default")
	if got.APIKey != "new" {
		t.Errorf("APIKey = %q, want %q after overwrite", got.APIKey, "new")
	}
	if got.AuthMode != "interactive" {
		t.Errorf("AuthMode = %q, want %q after overwrite", got.AuthMode, "interactive")
	}
}

func TestMemoryStoreMultipleProfiles(t *testing.T) {
	store := credentials.NewMemoryStore()
	store.Set("default", credentials.Credentials{APIKey: "k1", Token: "t1", AuthMode: "manual"})
	store.Set("work", credentials.Credentials{APIKey: "k2", Token: "t2", AuthMode: "interactive"})

	got1, _ := store.Get("default")
	got2, _ := store.Get("work")

	if got1.APIKey != "k1" {
		t.Errorf("default APIKey = %q, want %q", got1.APIKey, "k1")
	}
	if got2.APIKey != "k2" {
		t.Errorf("work APIKey = %q, want %q", got2.APIKey, "k2")
	}
}

func TestKeyringStoreServiceName(t *testing.T) {
	svc := credentials.KeyringServiceName("default")
	if svc != "trello-cli/default" {
		t.Errorf("service name = %q, want %q", svc, "trello-cli/default")
	}

	svc2 := credentials.KeyringServiceName("work")
	if svc2 != "trello-cli/work" {
		t.Errorf("service name = %q, want %q", svc2, "trello-cli/work")
	}
}

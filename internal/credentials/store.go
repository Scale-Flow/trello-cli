package credentials

import (
	"encoding/json"
	"errors"

	"github.com/zalando/go-keyring"
)

// ErrNotConfigured is returned when no credentials are stored for the requested profile.
var ErrNotConfigured = errors.New("credentials not configured")

// Credentials holds the API key, token, and auth mode for a profile.
type Credentials struct {
	APIKey   string `json:"apiKey"`
	Token    string `json:"token"`
	AuthMode string `json:"authMode"`
}

// Store defines the interface for credential persistence.
type Store interface {
	Get(profile string) (Credentials, error)
	Set(profile string, creds Credentials) error
	Delete(profile string) error
}

// MemoryStore is an in-memory Store implementation, primarily for testing.
type MemoryStore struct {
	creds map[string]Credentials
}

// NewMemoryStore creates a new in-memory credential store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{creds: make(map[string]Credentials)}
}

func (m *MemoryStore) Get(profile string) (Credentials, error) {
	cred, ok := m.creds[profile]
	if !ok {
		return Credentials{}, ErrNotConfigured
	}
	return cred, nil
}

func (m *MemoryStore) Set(profile string, creds Credentials) error {
	m.creds[profile] = creds
	return nil
}

func (m *MemoryStore) Delete(profile string) error {
	delete(m.creds, profile)
	return nil
}

const keyringServicePrefix = "trello-cli"

// KeyringServiceName returns the keyring service name for a profile.
func KeyringServiceName(profile string) string {
	return keyringServicePrefix + "/" + profile
}

// KeyringStore persists credentials in the OS keyring.
type KeyringStore struct{}

// NewKeyringStore creates a new keyring-backed credential store.
func NewKeyringStore() *KeyringStore {
	return &KeyringStore{}
}

func (k *KeyringStore) Get(profile string) (Credentials, error) {
	svc := KeyringServiceName(profile)
	data, err := keyring.Get(svc, "credentials")
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return Credentials{}, ErrNotConfigured
		}
		return Credentials{}, err
	}
	var creds Credentials
	if err := json.Unmarshal([]byte(data), &creds); err != nil {
		return Credentials{}, err
	}
	return creds, nil
}

func (k *KeyringStore) Set(profile string, creds Credentials) error {
	svc := KeyringServiceName(profile)
	data, err := json.Marshal(creds)
	if err != nil {
		return err
	}
	return keyring.Set(svc, "credentials", string(data))
}

func (k *KeyringStore) Delete(profile string) error {
	svc := KeyringServiceName(profile)
	err := keyring.Delete(svc, "credentials")
	if errors.Is(err, keyring.ErrNotFound) {
		return nil
	}
	return err
}

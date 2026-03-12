package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/credentials"
)

// ParseEnvelope parses JSON output into a map. Fails the test if invalid.
func ParseEnvelope(t *testing.T, output []byte) map[string]any {
	t.Helper()
	var envelope map[string]any
	if err := json.Unmarshal(output, &envelope); err != nil {
		t.Fatalf("invalid JSON envelope: %v\nraw: %s", err, string(output))
	}
	return envelope
}

// AssertSuccess asserts the envelope has ok: true and returns the data field.
func AssertSuccess(t *testing.T, output []byte) any {
	t.Helper()
	envelope := ParseEnvelope(t, output)
	ok, exists := envelope["ok"]
	if !exists {
		t.Fatal("envelope missing 'ok' field")
	}
	if ok != true {
		t.Errorf("ok = %v, want true", ok)
	}
	return envelope["data"]
}

// AssertError asserts the envelope has ok: false with the expected error code.
func AssertError(t *testing.T, output []byte, expectedCode string) {
	t.Helper()
	envelope := ParseEnvelope(t, output)
	ok, exists := envelope["ok"]
	if !exists {
		t.Fatal("envelope missing 'ok' field")
	}
	if ok != false {
		t.Errorf("ok = %v, want false", ok)
	}
	errObj, exists := envelope["error"]
	if !exists {
		t.Fatal("envelope missing 'error' field")
	}
	errMap, ok2 := errObj.(map[string]any)
	if !ok2 {
		t.Fatal("error field is not an object")
	}
	code, _ := errMap["code"].(string)
	if code != expectedCode {
		t.Errorf("error.code = %q, want %q", code, expectedCode)
	}
}

// MockCredentialStore returns a MemoryStore for use in tests.
// This is a convenience alias — tests should use credentials.NewMemoryStore() directly.
func MockCredentialStore() *credentials.MemoryStore {
	return credentials.NewMemoryStore()
}

// NewTrelloServer creates a test HTTP server with configurable route handlers.
// Routes are keyed by "METHOD /path" (e.g., "GET /1/members/me").
func NewTrelloServer(t *testing.T, routes map[string]http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Method + " " + r.URL.Path
		handler, ok := routes[key]
		if !ok {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		handler(w, r)
	}))
}

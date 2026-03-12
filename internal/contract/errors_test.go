package contract_test

import (
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/contract"
)

func TestContractErrorImplementsError(t *testing.T) {
	err := &contract.ContractError{Code: "TEST", Message: "test message"}
	var _ error = err // compile-time check

	got := err.Error()
	want := "TEST: test message"
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestErrorCodeConstants(t *testing.T) {
	codes := []string{
		contract.AuthRequired,
		contract.AuthInvalid,
		contract.NotFound,
		contract.ValidationError,
		contract.Conflict,
		contract.RateLimited,
		contract.HTTPError,
		contract.FileNotFound,
		contract.Unsupported,
		contract.UnknownError,
	}

	expected := []string{
		"AUTH_REQUIRED",
		"AUTH_INVALID",
		"NOT_FOUND",
		"VALIDATION_ERROR",
		"CONFLICT",
		"RATE_LIMITED",
		"HTTP_ERROR",
		"FILE_NOT_FOUND",
		"UNSUPPORTED",
		"UNKNOWN_ERROR",
	}

	if len(codes) != len(expected) {
		t.Fatalf("expected %d codes, got %d", len(expected), len(codes))
	}

	for i, code := range codes {
		if code != expected[i] {
			t.Errorf("code[%d] = %q, want %q", i, code, expected[i])
		}
	}
}

func TestNewContractError(t *testing.T) {
	err := contract.NewError(contract.ValidationError, "missing required flag")
	ce, ok := err.(*contract.ContractError)
	if !ok {
		t.Fatal("NewError should return *ContractError")
	}
	if ce.Code != "VALIDATION_ERROR" {
		t.Errorf("Code = %q, want %q", ce.Code, "VALIDATION_ERROR")
	}
	if ce.Message != "missing required flag" {
		t.Errorf("Message = %q, want %q", ce.Message, "missing required flag")
	}
}

# Trello CLI Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a cross-platform Go CLI implementing the Trello command contract with TDD, phased by layer.

**Architecture:** Foundation-first approach — contract layer, config, credentials, and auth subsystem built and tested before any resource commands. API client provides an interface that command handlers mock in tests. All commands return JSON envelopes on stdout.

**Tech Stack:** Go 1.23+, Cobra, Viper, go-keyring, net/http, encoding/json, httptest, slog

**Spec:** `docs/superpowers/specs/2026-03-12-trello-cli-implementation-design.md`

---

## File Structure

```
cmd/
  trello/
    main.go              — entry point, ldflags variables, root execute
    root.go              — root command, global flags, output/error rendering
    version.go           — version subcommand
    auth.go              — auth command group (status, login, set, clear)
    boards.go            — boards command group (list, get)
    lists.go             — lists command group (list, create, update, archive, move)
    cards.go             — cards command group (list, get, create, update, move, archive, delete)
    comments.go          — comments command group (list, add, update, delete)
    checklists.go        — checklists command group (list, create, delete, items subgroup)
    attachments.go       — attachments command group (list, add-file, add-url, delete)
    labels.go            — labels command group (list, create, add, remove)
    members.go           — members command group (list, add, remove)
    search.go            — search command group (cards, boards)

internal/
  contract/
    response.go          — Success(), Error(), Render() envelope builders
    response_test.go     — envelope shape and pretty-print tests
    errors.go            — ContractError type, error code constants
    errors_test.go       — error type behavior tests
    validation.go        — RequireFlag, RequireExactlyOne, RequireAtLeastOne, ValidateISO8601, ValidateURL, ValidateFilePath
    validation_test.go   — validator behavior tests
  config/
    config.go            — Viper init, defaults, precedence, env binding
    config_test.go       — precedence, defaults, missing/invalid file tests
  credentials/
    store.go             — Store interface, Credentials struct, KeyringStore, EnvStore
    store_test.go        — mock-based round-trip, missing profile, env fallback tests
  auth/
    auth.go              — RequireAuth helper, shared auth types
    auth_test.go         — auth guard tests
    login.go             — interactive browser flow + copy-paste fallback
    login_test.go        — login flow tests with mocked HTTP + credential store
    set.go               — manual credential entry
    set_test.go          — set validation and storage tests
    status.go            — auth state check via GET /members/me
    status_test.go       — configured/unconfigured state tests
    clear.go             — credential removal
    clear_test.go        — clear behavior tests
  trello/
    client.go            — HTTP client struct, Do(), auth injection, base URL
    client_test.go       — request construction, auth injection, timeout tests
    types.go             — Board, List, Card, Comment, Checklist, CheckItem, Attachment, Label, Member, search result types
    errors.go            — Trello HTTP status → ContractError mapping
    errors_test.go       — error mapping tests (401, 404, 429, 500)
    retry.go             — exponential backoff with jitter, rate-limit detection
    retry_test.go        — retry behavior, mutation safety, max retries tests
    boards.go            — ListBoards, GetBoard methods
    boards_test.go       — httptest-based board method tests
    lists.go             — ListLists, CreateList, UpdateList, ArchiveList, MoveList
    lists_test.go        — httptest-based list method tests
    cards.go             — card CRUD methods
    cards_test.go        — httptest-based card method tests
    comments.go          — comment methods (actions-based)
    comments_test.go     — httptest-based comment method tests
    checklists.go        — checklist + check item methods
    checklists_test.go   — httptest-based checklist method tests
    attachments.go       — attachment methods including multipart upload
    attachments_test.go  — httptest-based attachment method tests
    labels.go            — label methods
    labels_test.go       — httptest-based label method tests
    members.go           — member methods
    members_test.go      — httptest-based member method tests
    search.go            — search methods
    search_test.go       — httptest-based search method tests
  testutil/
    helpers.go           — ParseEnvelope, AssertSuccess, AssertError, MockCredentialStore, NewTrelloServer

go.mod
go.sum
```

---

## Chunk 1: Project Skeleton & Contract Layer

### Task 1: Initialize Go Module and Directory Structure

**Files:**
- Create: `go.mod`
- Create: `cmd/trello/main.go`

- [ ] **Step 1: Initialize Go module**

Run:
```bash
cd /Users/brettmcdowell/Dev/Trello_CLI
go mod init github.com/brettmcdowell/trello-cli
```

- [ ] **Step 2: Create directory structure**

Run:
```bash
mkdir -p cmd/trello
mkdir -p internal/contract
mkdir -p internal/config
mkdir -p internal/credentials
mkdir -p internal/auth
mkdir -p internal/trello
mkdir -p internal/testutil
```

- [ ] **Step 3: Create minimal main.go**

Create `cmd/trello/main.go`:
```go
package main

import "fmt"

func main() {
	fmt.Println("trello cli")
}
```

- [ ] **Step 4: Verify it compiles and runs**

Run: `go run ./cmd/trello/`
Expected: `trello cli`

- [ ] **Step 5: Install dependencies**

Run:
```bash
go get github.com/spf13/cobra
go get github.com/spf13/viper
go get github.com/zalando/go-keyring
```

- [ ] **Step 6: Commit**

```bash
git add go.mod go.sum cmd/trello/main.go
git commit -m "feat: initialize Go module with directory structure and dependencies"
```

---

### Task 2: Contract Error Type and Error Codes

**Files:**
- Create: `internal/contract/errors.go`
- Create: `internal/contract/errors_test.go`

- [x] **Step 1: Write the failing tests**

Create `internal/contract/errors_test.go`:
```go
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
```

- [x] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/contract/ -v`
Expected: compilation failure — `contract` package doesn't exist yet

- [ ] **Step 3: Write minimal implementation**

Create `internal/contract/errors.go`:
```go
package contract

const (
	AuthRequired    = "AUTH_REQUIRED"
	AuthInvalid     = "AUTH_INVALID"
	NotFound        = "NOT_FOUND"
	ValidationError = "VALIDATION_ERROR"
	Conflict        = "CONFLICT"
	RateLimited     = "RATE_LIMITED"
	HTTPError       = "HTTP_ERROR"
	FileNotFound    = "FILE_NOT_FOUND"
	Unsupported     = "UNSUPPORTED"
	UnknownError    = "UNKNOWN_ERROR"
)

// ContractError represents a structured CLI error with a machine-readable code.
type ContractError struct {
	Code    string
	Message string
}

func (e *ContractError) Error() string {
	return e.Code + ": " + e.Message
}

// NewError creates a new ContractError.
func NewError(code, message string) error {
	return &ContractError{Code: code, Message: message}
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/contract/ -v`
Expected: all 3 tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/contract/errors.go internal/contract/errors_test.go
git commit -m "feat: add contract error type and standard error code constants"
```

---

### Task 3: Success Envelope Builder

**Files:**
- Create: `internal/contract/response.go`
- Create: `internal/contract/response_test.go`

- [ ] **Step 1: Write the failing tests**

Create `internal/contract/response_test.go`:
```go
package contract_test

import (
	"encoding/json"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/contract"
)

func TestSuccessEnvelope(t *testing.T) {
	data := map[string]string{"name": "My Board"}
	result, err := contract.Success(data)
	if err != nil {
		t.Fatalf("Success() returned error: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(result, &envelope); err != nil {
		t.Fatalf("Success() produced invalid JSON: %v", err)
	}

	ok, exists := envelope["ok"]
	if !exists {
		t.Fatal("envelope missing 'ok' field")
	}
	if ok != true {
		t.Errorf("ok = %v, want true", ok)
	}

	d, exists := envelope["data"]
	if !exists {
		t.Fatal("envelope missing 'data' field")
	}

	dataMap, ok2 := d.(map[string]any)
	if !ok2 {
		t.Fatal("data is not an object")
	}
	if dataMap["name"] != "My Board" {
		t.Errorf("data.name = %v, want 'My Board'", dataMap["name"])
	}
}

func TestSuccessEnvelopeWithSlice(t *testing.T) {
	data := []map[string]string{{"id": "1"}, {"id": "2"}}
	result, err := contract.Success(data)
	if err != nil {
		t.Fatalf("Success() returned error: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(result, &envelope); err != nil {
		t.Fatalf("Success() produced invalid JSON: %v", err)
	}

	if envelope["ok"] != true {
		t.Errorf("ok = %v, want true", envelope["ok"])
	}

	arr, ok := envelope["data"].([]any)
	if !ok {
		t.Fatal("data is not an array")
	}
	if len(arr) != 2 {
		t.Errorf("data length = %d, want 2", len(arr))
	}
}

func TestSuccessEnvelopeHasNoErrorField(t *testing.T) {
	result, err := contract.Success("hello")
	if err != nil {
		t.Fatalf("Success() returned error: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(result, &envelope); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if _, exists := envelope["error"]; exists {
		t.Error("success envelope should not contain 'error' field")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/contract/ -run TestSuccess -v`
Expected: compilation failure — `Success` function doesn't exist

- [ ] **Step 3: Write minimal implementation**

Create `internal/contract/response.go`:
```go
package contract

import "encoding/json"

type successEnvelope struct {
	OK   bool `json:"ok"`
	Data any  `json:"data"`
}

// Success builds a JSON success envelope.
func Success(data any) ([]byte, error) {
	return json.Marshal(successEnvelope{OK: true, Data: data})
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/contract/ -run TestSuccess -v`
Expected: all 3 tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/contract/response.go internal/contract/response_test.go
git commit -m "feat: add Success envelope builder"
```

---

### Task 4: Error Envelope Builder

**Files:**
- Modify: `internal/contract/response.go`
- Modify: `internal/contract/response_test.go`

- [ ] **Step 1: Write the failing tests**

Append to `internal/contract/response_test.go`:
```go
func TestErrorEnvelope(t *testing.T) {
	result, err := contract.ErrorEnvelope(contract.NotFound, "board not found")
	if err != nil {
		t.Fatalf("ErrorEnvelope() returned error: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(result, &envelope); err != nil {
		t.Fatalf("ErrorEnvelope() produced invalid JSON: %v", err)
	}

	if envelope["ok"] != false {
		t.Errorf("ok = %v, want false", envelope["ok"])
	}

	errObj, exists := envelope["error"]
	if !exists {
		t.Fatal("envelope missing 'error' field")
	}

	errMap, ok := errObj.(map[string]any)
	if !ok {
		t.Fatal("error is not an object")
	}

	if errMap["code"] != "NOT_FOUND" {
		t.Errorf("error.code = %v, want NOT_FOUND", errMap["code"])
	}
	if errMap["message"] != "board not found" {
		t.Errorf("error.message = %v, want 'board not found'", errMap["message"])
	}
}

func TestErrorEnvelopeHasNoDataField(t *testing.T) {
	result, err := contract.ErrorEnvelope(contract.AuthRequired, "not logged in")
	if err != nil {
		t.Fatalf("ErrorEnvelope() returned error: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(result, &envelope); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if _, exists := envelope["data"]; exists {
		t.Error("error envelope should not contain 'data' field")
	}
}

func TestErrorEnvelopeFromContractError(t *testing.T) {
	ce := &contract.ContractError{Code: contract.ValidationError, Message: "missing --board"}
	result, err := contract.ErrorFromContractError(ce)
	if err != nil {
		t.Fatalf("ErrorFromContractError() returned error: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(result, &envelope); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	errMap := envelope["error"].(map[string]any)
	if errMap["code"] != "VALIDATION_ERROR" {
		t.Errorf("error.code = %v, want VALIDATION_ERROR", errMap["code"])
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/contract/ -run TestError -v`
Expected: compilation failure — `ErrorEnvelope` and `ErrorFromContractError` don't exist

- [ ] **Step 3: Write minimal implementation**

Add to `internal/contract/response.go`:
```go
type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorEnvelope struct {
	OK    bool        `json:"ok"`
	Error errorDetail `json:"error"`
}

// ErrorEnvelope builds a JSON error envelope from a code and message.
func ErrorEnvelope(code, message string) ([]byte, error) {
	return json.Marshal(errorEnvelope{
		OK:    false,
		Error: errorDetail{Code: code, Message: message},
	})
}

// ErrorFromContractError builds a JSON error envelope from a ContractError.
func ErrorFromContractError(ce *ContractError) ([]byte, error) {
	return ErrorEnvelope(ce.Code, ce.Message)
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/contract/ -run TestError -v`
Expected: all 3 tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/contract/response.go internal/contract/response_test.go
git commit -m "feat: add ErrorEnvelope and ErrorFromContractError builders"
```

---

### Task 5: Render Function (Pretty/Compact)

**Files:**
- Modify: `internal/contract/response.go`
- Modify: `internal/contract/response_test.go`

- [ ] **Step 1: Write the failing tests**

Append to `internal/contract/response_test.go`:
```go
import (
	"bytes"
)

func TestRenderCompact(t *testing.T) {
	data, _ := contract.Success(map[string]string{"id": "abc"})
	var buf bytes.Buffer
	err := contract.Render(&buf, data, false)
	if err != nil {
		t.Fatalf("Render() returned error: %v", err)
	}

	output := buf.String()
	// Compact: should be a single line ending with newline
	if output[len(output)-1] != '\n' {
		t.Error("Render() output should end with newline")
	}

	// Should be valid JSON
	var envelope map[string]any
	if err := json.Unmarshal([]byte(output), &envelope); err != nil {
		t.Fatalf("Render() produced invalid JSON: %v", err)
	}

	// Should not contain any indentation
	if bytes.Contains([]byte(output), []byte("  ")) {
		t.Error("compact output should not contain indentation")
	}
}

func TestRenderPretty(t *testing.T) {
	data, _ := contract.Success(map[string]string{"id": "abc"})
	var buf bytes.Buffer
	err := contract.Render(&buf, data, true)
	if err != nil {
		t.Fatalf("Render() returned error: %v", err)
	}

	output := buf.String()
	// Pretty: should contain indentation
	if !bytes.Contains([]byte(output), []byte("  ")) {
		t.Error("pretty output should contain indentation")
	}

	// Should be valid JSON
	var envelope map[string]any
	if err := json.Unmarshal([]byte(output), &envelope); err != nil {
		t.Fatalf("Render() produced invalid JSON: %v", err)
	}

	// Should end with newline
	if output[len(output)-1] != '\n' {
		t.Error("Render() output should end with newline")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/contract/ -run TestRender -v`
Expected: compilation failure — `Render` function doesn't exist

- [ ] **Step 3: Write minimal implementation**

Add to `internal/contract/response.go`:
```go
import (
	"encoding/json"
	"io"
)

// Render writes a JSON envelope to the writer. If pretty is true, output is indented.
func Render(w io.Writer, envelope []byte, pretty bool) error {
	if pretty {
		var buf bytes.Buffer
		if err := json.Indent(&buf, envelope, "", "  "); err != nil {
			return err
		}
		buf.WriteByte('\n')
		_, err := w.Write(buf.Bytes())
		return err
	}
	if _, err := w.Write(envelope); err != nil {
		return err
	}
	_, err := w.Write([]byte{'\n'})
	return err
}
```

Update the import block at top of `response.go` to include `"bytes"` and `"io"`.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/contract/ -run TestRender -v`
Expected: all 2 tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/contract/response.go internal/contract/response_test.go
git commit -m "feat: add Render function with pretty-print support"
```

---

### Task 6: Validators — RequireFlag

**Files:**
- Create: `internal/contract/validation.go`
- Create: `internal/contract/validation_test.go`

- [ ] **Step 1: Write the failing tests**

Create `internal/contract/validation_test.go`:
```go
package contract_test

import (
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/contract"
)

func TestRequireFlagPresent(t *testing.T) {
	err := contract.RequireFlag("board", "abc123")
	if err != nil {
		t.Errorf("RequireFlag() with value should return nil, got %v", err)
	}
}

func TestRequireFlagMissing(t *testing.T) {
	err := contract.RequireFlag("board", "")
	if err == nil {
		t.Fatal("RequireFlag() with empty value should return error")
	}

	ce, ok := err.(*contract.ContractError)
	if !ok {
		t.Fatal("RequireFlag() should return *ContractError")
	}
	if ce.Code != contract.ValidationError {
		t.Errorf("Code = %q, want %q", ce.Code, contract.ValidationError)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/contract/ -run TestRequireFlag -v`
Expected: compilation failure — `RequireFlag` doesn't exist

- [ ] **Step 3: Write minimal implementation**

Create `internal/contract/validation.go`:
```go
package contract

import "fmt"

// RequireFlag returns a VALIDATION_ERROR if value is empty.
func RequireFlag(name, value string) error {
	if value == "" {
		return NewError(ValidationError, fmt.Sprintf("--%s is required", name))
	}
	return nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/contract/ -run TestRequireFlag -v`
Expected: all 2 tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/contract/validation.go internal/contract/validation_test.go
git commit -m "feat: add RequireFlag validator"
```

---

### Task 7: Validators — RequireExactlyOne and RequireAtLeastOne

**Files:**
- Modify: `internal/contract/validation.go`
- Modify: `internal/contract/validation_test.go`

- [ ] **Step 1: Write the failing tests**

Append to `internal/contract/validation_test.go`:
```go
func TestRequireExactlyOneWithOne(t *testing.T) {
	err := contract.RequireExactlyOne(map[string]string{
		"board": "abc",
		"list":  "",
	})
	if err != nil {
		t.Errorf("RequireExactlyOne() with one value should return nil, got %v", err)
	}
}

func TestRequireExactlyOneWithNone(t *testing.T) {
	err := contract.RequireExactlyOne(map[string]string{
		"board": "",
		"list":  "",
	})
	if err == nil {
		t.Fatal("RequireExactlyOne() with no values should return error")
	}
	ce := err.(*contract.ContractError)
	if ce.Code != contract.ValidationError {
		t.Errorf("Code = %q, want %q", ce.Code, contract.ValidationError)
	}
}

func TestRequireExactlyOneWithTwo(t *testing.T) {
	err := contract.RequireExactlyOne(map[string]string{
		"board": "abc",
		"list":  "def",
	})
	if err == nil {
		t.Fatal("RequireExactlyOne() with two values should return error")
	}
	ce := err.(*contract.ContractError)
	if ce.Code != contract.ValidationError {
		t.Errorf("Code = %q, want %q", ce.Code, contract.ValidationError)
	}
}

func TestRequireAtLeastOneWithOne(t *testing.T) {
	err := contract.RequireAtLeastOne(map[string]string{
		"name": "new name",
		"pos":  "",
	})
	if err != nil {
		t.Errorf("RequireAtLeastOne() with one value should return nil, got %v", err)
	}
}

func TestRequireAtLeastOneWithNone(t *testing.T) {
	err := contract.RequireAtLeastOne(map[string]string{
		"name": "",
		"pos":  "",
	})
	if err == nil {
		t.Fatal("RequireAtLeastOne() with no values should return error")
	}
	ce := err.(*contract.ContractError)
	if ce.Code != contract.ValidationError {
		t.Errorf("Code = %q, want %q", ce.Code, contract.ValidationError)
	}
}

func TestRequireAtLeastOneWithAll(t *testing.T) {
	err := contract.RequireAtLeastOne(map[string]string{
		"name": "new name",
		"pos":  "top",
	})
	if err != nil {
		t.Errorf("RequireAtLeastOne() with all values should return nil, got %v", err)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/contract/ -run "TestRequireExactly|TestRequireAtLeast" -v`
Expected: compilation failure

- [ ] **Step 3: Write minimal implementation**

Add to `internal/contract/validation.go`:
```go
import (
	"fmt"
	"sort"
	"strings"
)

// RequireExactlyOne returns a VALIDATION_ERROR if not exactly one flag has a non-empty value.
func RequireExactlyOne(flags map[string]string) error {
	var set []string
	var all []string
	for name, value := range flags {
		all = append(all, "--"+name)
		if value != "" {
			set = append(set, "--"+name)
		}
	}
	sort.Strings(all)
	if len(set) == 1 {
		return nil
	}
	return NewError(ValidationError, fmt.Sprintf("exactly one of %s is required", strings.Join(all, ", ")))
}

// RequireAtLeastOne returns a VALIDATION_ERROR if no flags have a non-empty value.
func RequireAtLeastOne(flags map[string]string) error {
	var all []string
	for name, value := range flags {
		all = append(all, "--"+name)
		if value != "" {
			return nil
		}
	}
	sort.Strings(all)
	return NewError(ValidationError, fmt.Sprintf("at least one of %s is required", strings.Join(all, ", ")))
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/contract/ -run "TestRequireExactly|TestRequireAtLeast" -v`
Expected: all 6 tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/contract/validation.go internal/contract/validation_test.go
git commit -m "feat: add RequireExactlyOne and RequireAtLeastOne validators"
```

---

### Task 8: Validators — ValidateISO8601, ValidateURL, ValidateFilePath

**Files:**
- Modify: `internal/contract/validation.go`
- Modify: `internal/contract/validation_test.go`

- [ ] **Step 1: Write the failing tests**

Append to `internal/contract/validation_test.go`:
```go
func TestValidateISO8601Valid(t *testing.T) {
	valid := []string{
		"2026-03-20T00:00:00.000Z",
		"2026-03-20",
		"2026-12-31T23:59:59Z",
	}
	for _, v := range valid {
		if err := contract.ValidateISO8601(v); err != nil {
			t.Errorf("ValidateISO8601(%q) should be valid, got %v", v, err)
		}
	}
}

func TestValidateISO8601Invalid(t *testing.T) {
	invalid := []string{
		"not-a-date",
		"03/20/2026",
		"",
	}
	for _, v := range invalid {
		err := contract.ValidateISO8601(v)
		if err == nil {
			t.Errorf("ValidateISO8601(%q) should be invalid", v)
			continue
		}
		ce := err.(*contract.ContractError)
		if ce.Code != contract.ValidationError {
			t.Errorf("Code = %q, want %q", ce.Code, contract.ValidationError)
		}
	}
}

func TestValidateISO8601Empty(t *testing.T) {
	// Empty string means the flag was not provided — skip validation
	err := contract.ValidateISO8601Optional("")
	if err != nil {
		t.Errorf("ValidateISO8601Optional(\"\") should return nil, got %v", err)
	}
}

func TestValidateURLValid(t *testing.T) {
	err := contract.ValidateURL("https://example.com/file.pdf")
	if err != nil {
		t.Errorf("ValidateURL() with valid URL should return nil, got %v", err)
	}
}

func TestValidateURLInvalid(t *testing.T) {
	invalid := []string{
		"not-a-url",
		"",
		"ftp://nope",
	}
	for _, v := range invalid {
		err := contract.ValidateURL(v)
		if err == nil {
			t.Errorf("ValidateURL(%q) should be invalid", v)
			continue
		}
		ce := err.(*contract.ContractError)
		if ce.Code != contract.ValidationError {
			t.Errorf("Code = %q, want %q", ce.Code, contract.ValidationError)
		}
	}
}

func TestValidateFilePathValid(t *testing.T) {
	// Use a file that definitely exists
	err := contract.ValidateFilePath("/dev/null")
	if err != nil {
		t.Errorf("ValidateFilePath() with existing file should return nil, got %v", err)
	}
}

func TestValidateFilePathMissing(t *testing.T) {
	err := contract.ValidateFilePath("/nonexistent/file/path.txt")
	if err == nil {
		t.Fatal("ValidateFilePath() with missing file should return error")
	}
	ce := err.(*contract.ContractError)
	if ce.Code != contract.FileNotFound {
		t.Errorf("Code = %q, want %q", ce.Code, contract.FileNotFound)
	}
}

func TestValidateFilePathEmpty(t *testing.T) {
	err := contract.ValidateFilePath("")
	if err == nil {
		t.Fatal("ValidateFilePath() with empty path should return error")
	}
	ce := err.(*contract.ContractError)
	if ce.Code != contract.FileNotFound {
		t.Errorf("Code = %q, want %q", ce.Code, contract.FileNotFound)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/contract/ -run "TestValidateISO|TestValidateURL|TestValidateFile" -v`
Expected: compilation failure

- [ ] **Step 3: Write minimal implementation**

Add to `internal/contract/validation.go`:
```go
import (
	"net/url"
	"os"
	"time"
)

// ValidateISO8601 returns a VALIDATION_ERROR if the value is not a valid ISO-8601 date.
func ValidateISO8601(value string) error {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
	}
	for _, f := range formats {
		if _, err := time.Parse(f, value); err == nil {
			return nil
		}
	}
	return NewError(ValidationError, fmt.Sprintf("invalid ISO-8601 date: %q", value))
}

// ValidateISO8601Optional validates an ISO-8601 date only if non-empty.
func ValidateISO8601Optional(value string) error {
	if value == "" {
		return nil
	}
	return ValidateISO8601(value)
}

// ValidateURL returns a VALIDATION_ERROR if the value is not a valid HTTP(S) URL.
func ValidateURL(value string) error {
	if value == "" {
		return NewError(ValidationError, "URL is required")
	}
	u, err := url.Parse(value)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return NewError(ValidationError, fmt.Sprintf("invalid URL: %q", value))
	}
	return nil
}

// ValidateFilePath returns a FILE_NOT_FOUND if the path is empty or does not exist.
func ValidateFilePath(path string) error {
	if path == "" {
		return NewError(FileNotFound, "file path is required")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return NewError(FileNotFound, fmt.Sprintf("file not found: %s", path))
	}
	return nil
}
```

Update the import block to include `"net/url"`, `"os"`, and `"time"`.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/contract/ -run "TestValidateISO|TestValidateURL|TestValidateFile" -v`
Expected: all 8 tests PASS

- [ ] **Step 5: Run full contract test suite**

Run: `go test ./internal/contract/ -v`
Expected: all tests PASS

- [ ] **Step 6: Commit**

```bash
git add internal/contract/validation.go internal/contract/validation_test.go
git commit -m "feat: add ISO-8601, URL, and file path validators"
```

---

### Task 9: Test Helpers

**Files:**
- Create: `internal/testutil/helpers.go`

- [ ] **Step 1: Write test helpers**

Create `internal/testutil/helpers.go`:
```go
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
```

Note: This file references `credentials.Credentials` and `credentials.ErrNotConfigured` which will be created in Chunk 2. For now, this file will not compile until the credentials package exists. The test helpers will be committed but won't be imported until their dependencies are available.

- [ ] **Step 2: Verify the file has no syntax errors** (skip compilation since dependencies don't exist yet)

Run: `gofmt -e internal/testutil/helpers.go`
Expected: formatted output with no errors

- [ ] **Step 3: Commit**

```bash
git add internal/testutil/helpers.go
git commit -m "feat: add test helpers (ParseEnvelope, AssertSuccess, AssertError, mocks)"
```

---

### Task 10: Root Command with Global Flags

**Files:**
- Create: `cmd/trello/root.go`
- Modify: `cmd/trello/main.go`

- [ ] **Step 1: Write root command**

Create `cmd/trello/root.go`:
```go
package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/brettmcdowell/trello-cli/internal/contract"
	"github.com/spf13/cobra"
)

var (
	prettyFlag  bool
	verboseFlag bool
)

var rootCmd = &cobra.Command{
	Use:   "trello",
	Short: "Trello CLI — a machine-friendly Trello interface",
	Long:  "A cross-platform CLI for Trello, designed for coding agents and terminal users. All commands return JSON.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		level := slog.LevelWarn
		if verboseFlag {
			level = slog.LevelDebug
		}
		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
		slog.SetDefault(logger)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&prettyFlag, "pretty", false, "Pretty-print JSON output")
	rootCmd.PersistentFlags().BoolVar(&verboseFlag, "verbose", false, "Enable verbose diagnostic output on stderr")
}

// output writes a success envelope for the given data to stdout.
func output(w io.Writer, data any) error {
	envelope, err := contract.Success(data)
	if err != nil {
		return fmt.Errorf("failed to build success envelope: %w", err)
	}
	return contract.Render(w, envelope, prettyFlag)
}

// handleError writes an error envelope to stdout and returns the appropriate exit code.
func handleError(w io.Writer, err error) int {
	ce, ok := err.(*contract.ContractError)
	if !ok {
		ce = &contract.ContractError{Code: contract.UnknownError, Message: err.Error()}
	}
	envelope, marshalErr := contract.ErrorFromContractError(ce)
	if marshalErr != nil {
		// Last resort: write raw error to stderr
		fmt.Fprintf(os.Stderr, "fatal: %v\n", marshalErr)
		return 1
	}
	contract.Render(w, envelope, prettyFlag)
	return 1
}
```

- [ ] **Step 2: Update main.go to use root command**

Replace `cmd/trello/main.go` with:
```go
package main

import "os"

// Build metadata — set via ldflags at build time.
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(handleError(os.Stdout, err))
	}
}
```

- [ ] **Step 3: Verify it compiles and runs**

Run: `go run ./cmd/trello/ --help`
Expected: help text showing "Trello CLI" with `--pretty` and `--verbose` flags listed

- [ ] **Step 4: Commit**

```bash
git add cmd/trello/main.go cmd/trello/root.go
git commit -m "feat: add root command with --pretty and --verbose global flags"
```

---

### Task 11: Version Command

**Files:**
- Create: `cmd/trello/version.go`
- Create: `cmd/trello/version_test.go`

- [ ] **Step 1: Write the failing test**

Create `cmd/trello/version_test.go`:
```go
package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	// Set test build metadata
	version = "1.0.0-test"
	commit = "abc1234"
	date = "2026-03-12T00:00:00Z"

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"version"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON output: %v\nraw: %s", err, buf.String())
	}

	if envelope["ok"] != true {
		t.Errorf("ok = %v, want true", envelope["ok"])
	}

	data, ok := envelope["data"].(map[string]any)
	if !ok {
		t.Fatal("data is not an object")
	}

	if data["version"] != "1.0.0-test" {
		t.Errorf("version = %v, want '1.0.0-test'", data["version"])
	}
	if data["commit"] != "abc1234" {
		t.Errorf("commit = %v, want 'abc1234'", data["commit"])
	}
	if data["date"] != "2026-03-12T00:00:00Z" {
		t.Errorf("date = %v, want '2026-03-12T00:00:00Z'", data["date"])
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./cmd/trello/ -run TestVersionCommand -v`
Expected: FAIL — no "version" command registered

- [ ] **Step 3: Write minimal implementation**

Create `cmd/trello/version.go`:
```go
package main

import (
	"github.com/spf13/cobra"
)

type versionInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print CLI version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		return output(cmd.OutOrStdout(), versionInfo{
			Version: version,
			Commit:  commit,
			Date:    date,
		})
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./cmd/trello/ -run TestVersionCommand -v`
Expected: PASS

- [ ] **Step 5: Run full test suite**

Run: `go test ./... -v`
Expected: all tests PASS

- [ ] **Step 6: Verify end-to-end from CLI**

Run: `go run ./cmd/trello/ version`
Expected: `{"ok":true,"data":{"version":"dev","commit":"unknown","date":"unknown"}}`

- [ ] **Step 7: Commit**

```bash
git add cmd/trello/version.go cmd/trello/version_test.go
git commit -m "feat: add version command with build metadata"
```

---

**End of Chunk 1**

---

## Chunk 2: Config & Credentials

### Task 12: Config — Viper Initialization and Defaults

**Files:**
- Create: `internal/config/config.go`
- Create: `internal/config/config_test.go`

- [ ] **Step 1: Write the failing tests**

Create `internal/config/config_test.go`:
```go
package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/brettmcdowell/trello-cli/internal/config"
)

func TestDefaults(t *testing.T) {
	cfg := config.Load()

	if cfg.Profile != "default" {
		t.Errorf("Profile = %q, want %q", cfg.Profile, "default")
	}
	if cfg.Pretty != false {
		t.Errorf("Pretty = %v, want false", cfg.Pretty)
	}
	if cfg.Verbose != false {
		t.Errorf("Verbose = %v, want false", cfg.Verbose)
	}
	if cfg.Timeout != 15*time.Second {
		t.Errorf("Timeout = %v, want %v", cfg.Timeout, 15*time.Second)
	}
	if cfg.MaxRetries != 3 {
		t.Errorf("MaxRetries = %d, want 3", cfg.MaxRetries)
	}
	if cfg.RetryMutations != false {
		t.Errorf("RetryMutations = %v, want false", cfg.RetryMutations)
	}
}

func TestEnvOverridesDefaults(t *testing.T) {
	t.Setenv("TRELLO_PROFILE", "work")
	t.Setenv("TRELLO_PRETTY", "true")
	t.Setenv("TRELLO_TIMEOUT", "30s")
	t.Setenv("TRELLO_MAX_RETRIES", "5")
	t.Setenv("TRELLO_RETRY_MUTATIONS", "true")
	t.Setenv("TRELLO_VERBOSE", "true")

	cfg := config.Load()

	if cfg.Profile != "work" {
		t.Errorf("Profile = %q, want %q", cfg.Profile, "work")
	}
	if cfg.Pretty != true {
		t.Errorf("Pretty = %v, want true", cfg.Pretty)
	}
	if cfg.Verbose != true {
		t.Errorf("Verbose = %v, want true", cfg.Verbose)
	}
	if cfg.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want %v", cfg.Timeout, 30*time.Second)
	}
	if cfg.MaxRetries != 5 {
		t.Errorf("MaxRetries = %d, want 5", cfg.MaxRetries)
	}
	if cfg.RetryMutations != true {
		t.Errorf("RetryMutations = %v, want true", cfg.RetryMutations)
	}
}

func TestMissingConfigFileDoesNotError(t *testing.T) {
	// Ensure no config file exists at the default path
	t.Setenv("TRELLO_CONFIG_PATH", "/tmp/nonexistent-trello-config-test.yaml")

	cfg := config.Load()
	// Should still return defaults without error
	if cfg.Profile != "default" {
		t.Errorf("Profile = %q, want %q", cfg.Profile, "default")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/config/ -v`
Expected: compilation failure — `config` package doesn't exist

- [ ] **Step 3: Write minimal implementation**

Create `internal/config/config.go`:
```go
package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Config holds all non-secret CLI configuration.
type Config struct {
	Profile        string        `mapstructure:"profile"`
	Pretty         bool          `mapstructure:"pretty"`
	Verbose        bool          `mapstructure:"verbose"`
	Timeout        time.Duration `mapstructure:"timeout"`
	MaxRetries     int           `mapstructure:"max_retries"`
	RetryMutations bool          `mapstructure:"retry_mutations"`
}

// Load reads configuration from file, environment, and defaults.
// Precedence: env vars > config file > defaults.
func Load() Config {
	v := viper.New()

	// Defaults
	v.SetDefault("profile", "default")
	v.SetDefault("pretty", false)
	v.SetDefault("verbose", false)
	v.SetDefault("timeout", "15s")
	v.SetDefault("max_retries", 3)
	v.SetDefault("retry_mutations", false)

	// Environment binding
	v.SetEnvPrefix("TRELLO")
	v.AutomaticEnv()

	// Config file
	configPath := os.Getenv("TRELLO_CONFIG_PATH")
	if configPath == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			configPath = filepath.Join(home, ".config", "trello-cli", "config.yaml")
		}
	}
	if configPath != "" {
		v.SetConfigFile(configPath)
		_ = v.ReadInConfig() // Ignore error — missing file is fine
	}

	var cfg Config
	_ = v.Unmarshal(&cfg)

	// Parse timeout string into duration
	if cfg.Timeout == 0 {
		dur, err := time.ParseDuration(v.GetString("timeout"))
		if err == nil {
			cfg.Timeout = dur
		} else {
			cfg.Timeout = 15 * time.Second
		}
	}

	return cfg
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/config/ -v`
Expected: all 3 tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/config/config.go internal/config/config_test.go
git commit -m "feat: add config package with Viper defaults and env binding"
```

---

### Task 13: Config — File-Based Configuration

**Files:**
- Modify: `internal/config/config_test.go`

- [ ] **Step 1: Write the failing test**

Append to `internal/config/config_test.go`:
```go
func TestConfigFileOverridesDefaults(t *testing.T) {
	// Create a temp config file
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	content := []byte("profile: personal\npretty: true\ntimeout: 20s\nmax_retries: 1\n")
	if err := os.WriteFile(configPath, content, 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	t.Setenv("TRELLO_CONFIG_PATH", configPath)

	cfg := config.Load()

	if cfg.Profile != "personal" {
		t.Errorf("Profile = %q, want %q", cfg.Profile, "personal")
	}
	if cfg.Pretty != true {
		t.Errorf("Pretty = %v, want true", cfg.Pretty)
	}
	if cfg.Timeout != 20*time.Second {
		t.Errorf("Timeout = %v, want %v", cfg.Timeout, 20*time.Second)
	}
	if cfg.MaxRetries != 1 {
		t.Errorf("MaxRetries = %d, want 1", cfg.MaxRetries)
	}
}

func TestEnvOverridesConfigFile(t *testing.T) {
	// Create a config file with profile=personal
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	content := []byte("profile: personal\n")
	if err := os.WriteFile(configPath, content, 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	t.Setenv("TRELLO_CONFIG_PATH", configPath)
	t.Setenv("TRELLO_PROFILE", "work")

	cfg := config.Load()

	// Env should override config file
	if cfg.Profile != "work" {
		t.Errorf("Profile = %q, want %q (env should override config file)", cfg.Profile, "work")
	}
}
```

Add `"path/filepath"` to the imports.

- [ ] **Step 2: Run tests to verify they pass** (should already work with current implementation)

Run: `go test ./internal/config/ -v`
Expected: all 5 tests PASS

- [ ] **Step 3: Commit**

```bash
git add internal/config/config_test.go
git commit -m "test: add config file loading and precedence tests"
```

---

### Task 14: Credentials — Store Interface and Types

**Files:**
- Create: `internal/credentials/store.go`
- Create: `internal/credentials/store_test.go`

- [ ] **Step 1: Write the failing tests**

Create `internal/credentials/store_test.go`:
```go
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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/credentials/ -v`
Expected: compilation failure

- [ ] **Step 3: Write minimal implementation**

Create `internal/credentials/store.go`:
```go
package credentials

import "errors"

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
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/credentials/ -v`
Expected: all 2 tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/credentials/store.go internal/credentials/store_test.go
git commit -m "feat: add credential Store interface and Credentials type"
```

---

### Task 15: Credentials — MemoryStore (for Testing) and Round-Trip Tests

**Files:**
- Modify: `internal/credentials/store.go`
- Modify: `internal/credentials/store_test.go`

- [ ] **Step 1: Write the failing tests**

Append to `internal/credentials/store_test.go`:
```go
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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/credentials/ -run "TestMemory" -v`
Expected: compilation failure — `NewMemoryStore` doesn't exist

- [ ] **Step 3: Write minimal implementation**

Add to `internal/credentials/store.go`:
```go
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
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/credentials/ -v`
Expected: all 8 tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/credentials/store.go internal/credentials/store_test.go
git commit -m "feat: add MemoryStore implementation with round-trip tests"
```

---

### Task 16: Credentials — KeyringStore

**Files:**
- Modify: `internal/credentials/store.go`
- Modify: `internal/credentials/store_test.go`

- [ ] **Step 1: Write the failing tests**

Append to `internal/credentials/store_test.go`:
```go
func TestKeyringStoreServiceName(t *testing.T) {
	// Verify the service name pattern
	svc := credentials.KeyringServiceName("default")
	if svc != "trello-cli/default" {
		t.Errorf("service name = %q, want %q", svc, "trello-cli/default")
	}

	svc2 := credentials.KeyringServiceName("work")
	if svc2 != "trello-cli/work" {
		t.Errorf("service name = %q, want %q", svc2, "trello-cli/work")
	}
}
```

Note: We do NOT test actual keyring read/write in unit tests — that requires a real OS keyring. The `KeyringStore` implementation delegates to `go-keyring` whose behavior we trust. We test only the service name construction and the mapping of `keyring.ErrNotFound` to `ErrNotConfigured`.

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/credentials/ -run TestKeyring -v`
Expected: compilation failure — `KeyringServiceName` doesn't exist

- [ ] **Step 3: Write minimal implementation**

Add to `internal/credentials/store.go`:
```go
import (
	"encoding/json"
	"errors"

	"github.com/zalando/go-keyring"
)

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
		return nil // Already gone, not an error
	}
	return err
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/credentials/ -run TestKeyring -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/credentials/store.go internal/credentials/store_test.go
git commit -m "feat: add KeyringStore backed by OS keyring via go-keyring"
```

---

### Task 17: Credentials — EnvStore

**Files:**
- Modify: `internal/credentials/store.go`
- Modify: `internal/credentials/store_test.go`

- [ ] **Step 1: Write the failing tests**

Append to `internal/credentials/store_test.go`:
```go
import "os"

func TestEnvStoreGetWithEnvVars(t *testing.T) {
	t.Setenv("TRELLO_API_KEY", "env-key")
	t.Setenv("TRELLO_TOKEN", "env-token")

	store := credentials.NewEnvStore()
	got, err := store.Get("default")
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if got.APIKey != "env-key" {
		t.Errorf("APIKey = %q, want %q", got.APIKey, "env-key")
	}
	if got.Token != "env-token" {
		t.Errorf("Token = %q, want %q", got.Token, "env-token")
	}
	if got.AuthMode != "env" {
		t.Errorf("AuthMode = %q, want %q", got.AuthMode, "env")
	}
}

func TestEnvStoreGetMissingKey(t *testing.T) {
	t.Setenv("TRELLO_API_KEY", "")
	t.Setenv("TRELLO_TOKEN", "")

	store := credentials.NewEnvStore()
	_, err := store.Get("default")
	if !errors.Is(err, credentials.ErrNotConfigured) {
		t.Errorf("Get() error = %v, want ErrNotConfigured", err)
	}
}

func TestEnvStoreGetPartialCreds(t *testing.T) {
	t.Setenv("TRELLO_API_KEY", "env-key")
	t.Setenv("TRELLO_TOKEN", "")

	store := credentials.NewEnvStore()
	_, err := store.Get("default")
	if !errors.Is(err, credentials.ErrNotConfigured) {
		t.Errorf("Get() with only API key error = %v, want ErrNotConfigured", err)
	}
}

func TestEnvStoreSetReturnsError(t *testing.T) {
	store := credentials.NewEnvStore()
	err := store.Set("default", credentials.Credentials{})
	if err == nil {
		t.Error("EnvStore.Set() should return an error (read-only)")
	}
}

func TestEnvStoreDeleteReturnsError(t *testing.T) {
	store := credentials.NewEnvStore()
	err := store.Delete("default")
	if err == nil {
		t.Error("EnvStore.Delete() should return an error (read-only)")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/credentials/ -run TestEnvStore -v`
Expected: compilation failure — `NewEnvStore` doesn't exist

- [ ] **Step 3: Write minimal implementation**

Add to `internal/credentials/store.go`:
```go
import "os"

// ErrReadOnly is returned when Set or Delete is called on a read-only store.
var ErrReadOnly = errors.New("credential store is read-only")

// EnvStore reads credentials from environment variables. It is read-only.
type EnvStore struct{}

// NewEnvStore creates a new environment-variable-backed credential store.
func NewEnvStore() *EnvStore {
	return &EnvStore{}
}

// Get reads credentials from TRELLO_API_KEY and TRELLO_TOKEN environment variables.
// The profile parameter is ignored — env-based auth has no profile concept.
func (e *EnvStore) Get(profile string) (Credentials, error) {
	apiKey := os.Getenv("TRELLO_API_KEY")
	token := os.Getenv("TRELLO_TOKEN")
	if apiKey == "" || token == "" {
		return Credentials{}, ErrNotConfigured
	}
	return Credentials{
		APIKey:   apiKey,
		Token:    token,
		AuthMode: "env",
	}, nil
}

func (e *EnvStore) Set(profile string, creds Credentials) error {
	return ErrReadOnly
}

func (e *EnvStore) Delete(profile string) error {
	return ErrReadOnly
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/credentials/ -v`
Expected: all tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/credentials/store.go internal/credentials/store_test.go
git commit -m "feat: add EnvStore for environment-variable credential fallback"
```

---

### Task 18: Credentials — FallbackStore (Keyring → Env)

**Files:**
- Modify: `internal/credentials/store.go`
- Modify: `internal/credentials/store_test.go`

- [ ] **Step 1: Write the failing tests**

Append to `internal/credentials/store_test.go`:
```go
func TestFallbackStoreUsesFirst(t *testing.T) {
	primary := credentials.NewMemoryStore()
	primary.Set("default", credentials.Credentials{APIKey: "pk", Token: "pt", AuthMode: "manual"})
	fallback := credentials.NewMemoryStore()

	store := credentials.NewFallbackStore(primary, fallback)
	got, err := store.Get("default")
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if got.APIKey != "pk" {
		t.Errorf("APIKey = %q, want %q (should use primary)", got.APIKey, "pk")
	}
}

func TestFallbackStoreUsesSecondWhenFirstMissing(t *testing.T) {
	primary := credentials.NewMemoryStore()
	fallback := credentials.NewMemoryStore()
	fallback.Set("default", credentials.Credentials{APIKey: "fk", Token: "ft", AuthMode: "env"})

	store := credentials.NewFallbackStore(primary, fallback)
	got, err := store.Get("default")
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if got.APIKey != "fk" {
		t.Errorf("APIKey = %q, want %q (should use fallback)", got.APIKey, "fk")
	}
}

func TestFallbackStoreReturnsErrWhenBothMissing(t *testing.T) {
	primary := credentials.NewMemoryStore()
	fallback := credentials.NewMemoryStore()

	store := credentials.NewFallbackStore(primary, fallback)
	_, err := store.Get("default")
	if !errors.Is(err, credentials.ErrNotConfigured) {
		t.Errorf("Get() error = %v, want ErrNotConfigured", err)
	}
}

func TestFallbackStoreSetWritesToPrimary(t *testing.T) {
	primary := credentials.NewMemoryStore()
	fallback := credentials.NewMemoryStore()

	store := credentials.NewFallbackStore(primary, fallback)
	creds := credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"}
	store.Set("default", creds)

	// Should be in primary, not fallback
	got, err := primary.Get("default")
	if err != nil {
		t.Fatalf("primary.Get() returned error: %v", err)
	}
	if got.APIKey != "k" {
		t.Errorf("primary APIKey = %q, want %q", got.APIKey, "k")
	}

	_, err = fallback.Get("default")
	if !errors.Is(err, credentials.ErrNotConfigured) {
		t.Error("fallback should not have the credentials")
	}
}

func TestFallbackStoreDeleteFromPrimary(t *testing.T) {
	primary := credentials.NewMemoryStore()
	primary.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	fallback := credentials.NewMemoryStore()

	store := credentials.NewFallbackStore(primary, fallback)
	store.Delete("default")

	_, err := primary.Get("default")
	if !errors.Is(err, credentials.ErrNotConfigured) {
		t.Error("primary should not have credentials after delete")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/credentials/ -run TestFallback -v`
Expected: compilation failure — `NewFallbackStore` doesn't exist

- [ ] **Step 3: Write minimal implementation**

Add to `internal/credentials/store.go`:
```go
// FallbackStore tries the primary store first, falling back to the secondary.
// Writes always go to the primary store.
type FallbackStore struct {
	primary   Store
	secondary Store
}

// NewFallbackStore creates a store that reads from primary first, then secondary.
func NewFallbackStore(primary, secondary Store) *FallbackStore {
	return &FallbackStore{primary: primary, secondary: secondary}
}

func (f *FallbackStore) Get(profile string) (Credentials, error) {
	creds, err := f.primary.Get(profile)
	if err == nil {
		return creds, nil
	}
	if !errors.Is(err, ErrNotConfigured) {
		return Credentials{}, err
	}
	return f.secondary.Get(profile)
}

func (f *FallbackStore) Set(profile string, creds Credentials) error {
	return f.primary.Set(profile, creds)
}

func (f *FallbackStore) Delete(profile string) error {
	return f.primary.Delete(profile)
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/credentials/ -v`
Expected: all tests PASS

- [ ] **Step 5: Update testutil helpers to use credentials package**

Now that the credentials package exists, verify the `internal/testutil/helpers.go` compiles:

Run: `go build ./internal/testutil/`
Expected: successful compilation (no output)

- [ ] **Step 6: Run full test suite**

Run: `go test ./... -v`
Expected: all tests PASS

- [ ] **Step 7: Commit**

```bash
git add internal/credentials/store.go internal/credentials/store_test.go
git commit -m "feat: add FallbackStore for keyring-to-env credential fallback chain"
```

---

**End of Chunk 2**

---

## Chunk 3: Auth Subsystem

### Task 19: Auth — RequireAuth Guard

**Files:**
- Create: `internal/auth/auth.go`
- Create: `internal/auth/auth_test.go`

- [ ] **Step 1: Write the failing tests**

Create `internal/auth/auth_test.go`:
```go
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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/auth/ -v`
Expected: compilation failure

- [ ] **Step 3: Write minimal implementation**

Create `internal/auth/auth.go`:
```go
package auth

import (
	"errors"

	"github.com/brettmcdowell/trello-cli/internal/contract"
	"github.com/brettmcdowell/trello-cli/internal/credentials"
)

// RequireAuth loads credentials for the given profile. Returns AUTH_REQUIRED if missing.
func RequireAuth(store credentials.Store, profile string) (credentials.Credentials, error) {
	creds, err := store.Get(profile)
	if err != nil {
		if errors.Is(err, credentials.ErrNotConfigured) {
			return credentials.Credentials{}, contract.NewError(contract.AuthRequired, "not authenticated — run 'trello auth login' or 'trello auth set'")
		}
		return credentials.Credentials{}, err
	}
	return creds, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/auth/ -v`
Expected: all 2 tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/auth/auth.go internal/auth/auth_test.go
git commit -m "feat: add RequireAuth guard for centralized auth checking"
```

---

### Task 20: Auth — Set Command Logic

**Files:**
- Create: `internal/auth/set.go`
- Create: `internal/auth/set_test.go`

- [ ] **Step 1: Write the failing tests**

Create `internal/auth/set_test.go`:
```go
package auth_test

import (
	"testing"

	"encoding/json"

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

	// Verify credentials were stored
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
	// Set should store credentials without validation.
	// If it tried to call an API, it would fail here since no HTTP client is provided.
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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/auth/ -run TestSet -v`
Expected: compilation failure

- [ ] **Step 3: Write minimal implementation**

Create `internal/auth/set.go`:
```go
package auth

import "github.com/brettmcdowell/trello-cli/internal/credentials"

// SetResult is the response shape for auth set.
type SetResult struct {
	Configured bool        `json:"configured"`
	AuthMode   string      `json:"authMode"`
	Member     interface{} `json:"-"` // Never included — auth set does not validate
}

// Set stores credentials without validating them against the Trello API.
// Validation is deferred to auth status or the first command that requires auth.
func Set(store credentials.Store, profile, apiKey, token string) (SetResult, error) {
	creds := credentials.Credentials{
		APIKey:   apiKey,
		Token:    token,
		AuthMode: "manual",
	}
	if err := store.Set(profile, creds); err != nil {
		return SetResult{}, err
	}
	return SetResult{
		Configured: true,
		AuthMode:   "manual",
	}, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/auth/ -run TestSet -v`
Expected: all 3 tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/auth/set.go internal/auth/set_test.go
git commit -m "feat: add auth set logic — stores credentials without API validation"
```

---

### Task 21: Auth — Status Command Logic

**Files:**
- Create: `internal/auth/status.go`
- Create: `internal/auth/status_test.go`

- [ ] **Step 1: Write the failing tests**

Create `internal/auth/status_test.go`:
```go
package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/auth"
	"github.com/brettmcdowell/trello-cli/internal/credentials"
)

func TestStatusConfigured(t *testing.T) {
	// Mock Trello API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/1/members/me" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]string{
			"id":       "member123",
			"username": "testuser",
			"fullName": "Test User",
		})
	}))
	defer server.Close()

	store := credentials.NewMemoryStore()
	store.Set("default", credentials.Credentials{
		APIKey: "key1", Token: "tok1", AuthMode: "manual",
	})

	result, err := auth.Status(context.Background(), store, "default", server.URL)
	if err != nil {
		t.Fatalf("Status() returned error: %v", err)
	}

	if result.Configured != true {
		t.Errorf("Configured = %v, want true", result.Configured)
	}
	if *result.AuthMode != "manual" {
		t.Errorf("AuthMode = %q, want %q", *result.AuthMode, "manual")
	}
	if result.Member == nil {
		t.Fatal("Member should not be nil when configured")
	}
	if result.Member.ID != "member123" {
		t.Errorf("Member.ID = %q, want %q", result.Member.ID, "member123")
	}
	if result.Member.Username != "testuser" {
		t.Errorf("Member.Username = %q, want %q", result.Member.Username, "testuser")
	}
}

func TestStatusNotConfigured(t *testing.T) {
	store := credentials.NewMemoryStore()

	result, err := auth.Status(context.Background(), store, "default", "")
	if err != nil {
		t.Fatalf("Status() returned error: %v", err)
	}

	if result.Configured != false {
		t.Errorf("Configured = %v, want false", result.Configured)
	}
	if result.AuthMode != nil {
		t.Errorf("AuthMode = %v, want nil", result.AuthMode)
	}
	if result.Member != nil {
		t.Error("Member should be nil when not configured")
	}
}

func TestStatusAuthModeNullJSON(t *testing.T) {
	store := credentials.NewMemoryStore()
	result, _ := auth.Status(context.Background(), store, "default", "")

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json.Marshal() returned error: %v", err)
	}

	var m map[string]any
	json.Unmarshal(data, &m)

	if m["authMode"] != nil {
		t.Errorf("authMode = %v, want null", m["authMode"])
	}
}

func TestStatusInvalidCredentials(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
	}))
	defer server.Close()

	store := credentials.NewMemoryStore()
	store.Set("default", credentials.Credentials{
		APIKey: "bad-key", Token: "bad-token", AuthMode: "manual",
	})

	_, err := auth.Status(context.Background(), store, "default", server.URL)
	if err == nil {
		t.Fatal("Status() should return error for invalid credentials")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/auth/ -run TestStatus -v`
Expected: compilation failure

- [ ] **Step 3: Write minimal implementation**

Create `internal/auth/status.go`:
```go
package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/brettmcdowell/trello-cli/internal/contract"
	"github.com/brettmcdowell/trello-cli/internal/credentials"
)

// Member represents a Trello member for auth responses.
type Member struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
}

// StatusResult is the response shape for auth status.
type StatusResult struct {
	Configured bool    `json:"configured"`
	AuthMode   *string `json:"authMode"` // *string to marshal null when not configured
	Member     *Member `json:"member"`   // Always present: null when not configured
}

// Status checks the current authentication state.
// If credentials exist, it validates them against the Trello API.
func Status(ctx context.Context, store credentials.Store, profile, baseURL string) (StatusResult, error) {
	creds, err := store.Get(profile)
	if err != nil {
		if errors.Is(err, credentials.ErrNotConfigured) {
			return StatusResult{
				Configured: false,
				AuthMode:   nil,
				Member:     nil,
			}, nil
		}
		return StatusResult{}, err
	}

	// Validate credentials by calling GET /1/members/me
	member, err := getMember(ctx, baseURL, creds.APIKey, creds.Token)
	if err != nil {
		return StatusResult{}, err
	}

	authMode := creds.AuthMode
	return StatusResult{
		Configured: true,
		AuthMode:   &authMode,
		Member:     member,
	}, nil
}

// getMember calls GET /1/members/me to validate credentials and fetch member info.
func getMember(ctx context.Context, baseURL, apiKey, token string) (*Member, error) {
	url := fmt.Sprintf("%s/1/members/me?key=%s&token=%s", baseURL, apiKey, token)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, contract.NewError(contract.HTTPError, fmt.Sprintf("failed to reach Trello API: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, contract.NewError(contract.AuthInvalid, "Trello rejected the credentials — API key or token is invalid")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, contract.NewError(contract.HTTPError, fmt.Sprintf("Trello API returned status %d", resp.StatusCode))
	}

	var member Member
	if err := json.NewDecoder(resp.Body).Decode(&member); err != nil {
		return nil, contract.NewError(contract.HTTPError, fmt.Sprintf("failed to decode member response: %v", err))
	}

	return &member, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/auth/ -run TestStatus -v`
Expected: all 4 tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/auth/status.go internal/auth/status_test.go
git commit -m "feat: add auth status with credential validation via GET /members/me"
```

---

### Task 22: Auth — Clear Command Logic

**Files:**
- Create: `internal/auth/clear.go`
- Create: `internal/auth/clear_test.go`

- [ ] **Step 1: Write the failing tests**

Create `internal/auth/clear_test.go`:
```go
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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/auth/ -run TestClear -v`
Expected: compilation failure

- [ ] **Step 3: Write minimal implementation**

Create `internal/auth/clear.go`:
```go
package auth

import "github.com/brettmcdowell/trello-cli/internal/credentials"

// ClearResult is the response shape for auth clear.
type ClearResult struct {
	Configured bool    `json:"configured"`
	AuthMode   *string `json:"authMode"` // Always null after clear
}

// Clear removes stored credentials for the given profile.
func Clear(store credentials.Store, profile string) (ClearResult, error) {
	if err := store.Delete(profile); err != nil {
		return ClearResult{}, err
	}
	return ClearResult{
		Configured: false,
		AuthMode:   nil,
	}, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/auth/ -run TestClear -v`
Expected: all 3 tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/auth/clear.go internal/auth/clear_test.go
git commit -m "feat: add auth clear — removes credentials and returns null authMode"
```

---

### Task 23: Auth — Login Flow (Callback Server + Copy-Paste Fallback)

**Files:**
- Create: `internal/auth/login.go`
- Create: `internal/auth/login_test.go`

- [ ] **Step 1: Write the failing tests**

Create `internal/auth/login_test.go`:
```go
package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/auth"
	"github.com/brettmcdowell/trello-cli/internal/credentials"
)

func TestLoginWithToken(t *testing.T) {
	// Mock Trello API for credential validation
	trelloServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/1/members/me" {
			json.NewEncoder(w).Encode(map[string]string{
				"id":       "member456",
				"username": "loginuser",
				"fullName": "Login User",
			})
			return
		}
		http.Error(w, "not found", 404)
	}))
	defer trelloServer.Close()

	store := credentials.NewMemoryStore()

	// Simulate the login flow by directly providing a token
	// (bypasses browser — tests the post-token-capture logic)
	result, err := auth.CompleteLogin(
		context.Background(),
		store,
		"default",
		"test-api-key",
		"captured-token",
		trelloServer.URL,
	)
	if err != nil {
		t.Fatalf("CompleteLogin() returned error: %v", err)
	}

	if result.Configured != true {
		t.Errorf("Configured = %v, want true", result.Configured)
	}
	if *result.AuthMode != "interactive" {
		t.Errorf("AuthMode = %q, want %q", *result.AuthMode, "interactive")
	}
	if result.Member == nil {
		t.Fatal("Member should not be nil after login")
	}
	if result.Member.Username != "loginuser" {
		t.Errorf("Member.Username = %q, want %q", result.Member.Username, "loginuser")
	}

	// Verify credentials were stored
	creds, _ := store.Get("default")
	if creds.APIKey != "test-api-key" {
		t.Errorf("stored APIKey = %q, want %q", creds.APIKey, "test-api-key")
	}
	if creds.Token != "captured-token" {
		t.Errorf("stored Token = %q, want %q", creds.Token, "captured-token")
	}
	if creds.AuthMode != "interactive" {
		t.Errorf("stored AuthMode = %q, want %q", creds.AuthMode, "interactive")
	}
}

func TestLoginInvalidToken(t *testing.T) {
	trelloServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer trelloServer.Close()

	store := credentials.NewMemoryStore()

	_, err := auth.CompleteLogin(
		context.Background(),
		store,
		"default",
		"bad-key",
		"bad-token",
		trelloServer.URL,
	)
	if err == nil {
		t.Fatal("CompleteLogin() should return error for invalid credentials")
	}

	// Should not store credentials on failure
	_, getErr := store.Get("default")
	if getErr == nil {
		t.Error("credentials should not be stored after failed login")
	}
}

func TestBuildAuthorizeURL(t *testing.T) {
	url := auth.BuildAuthorizeURL("my-api-key", "http://localhost:12345/callback")

	if url == "" {
		t.Fatal("BuildAuthorizeURL() returned empty string")
	}

	// Should contain the API key and callback URL
	if !strings.Contains(url, "key=my-api-key") {
		t.Errorf("URL should contain API key, got: %s", url)
	}
	if !strings.Contains(url, "return_url=") {
		t.Errorf("URL should contain return_url, got: %s", url)
	}
}

```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/auth/ -run "TestLogin|TestBuild" -v`
Expected: compilation failure

- [ ] **Step 3: Write minimal implementation**

Create `internal/auth/login.go`:
```go
package auth

import (
	"context"
	"fmt"
	"net/url"

	"github.com/brettmcdowell/trello-cli/internal/credentials"
)

const trelloAuthorizeBase = "https://trello.com/1/authorize"

// LoginResult is the response shape for auth login.
type LoginResult struct {
	Configured bool    `json:"configured"`
	AuthMode   *string `json:"authMode"`
	Member     *Member `json:"member"` // Always present after successful login
}

// BuildAuthorizeURL builds the Trello authorization URL for interactive login.
func BuildAuthorizeURL(apiKey, callbackURL string) string {
	params := url.Values{
		"expiration":  {"never"},
		"name":        {"Trello CLI"},
		"scope":       {"read,write"},
		"response_type": {"token"},
		"key":         {apiKey},
		"return_url":  {callbackURL},
	}
	return trelloAuthorizeBase + "?" + params.Encode()
}

// CompleteLogin validates a captured token and stores credentials.
// This is called after the token has been obtained (via callback server or copy-paste).
func CompleteLogin(ctx context.Context, store credentials.Store, profile, apiKey, token, baseURL string) (LoginResult, error) {
	// Validate credentials by calling GET /1/members/me
	member, err := getMember(ctx, baseURL, apiKey, token)
	if err != nil {
		return LoginResult{}, err
	}

	// Store validated credentials
	creds := credentials.Credentials{
		APIKey:   apiKey,
		Token:    token,
		AuthMode: "interactive",
	}
	if err := store.Set(profile, creds); err != nil {
		return LoginResult{}, fmt.Errorf("failed to store credentials: %w", err)
	}

	authMode := "interactive"
	return LoginResult{
		Configured: true,
		AuthMode:   &authMode,
		Member:     member,
	}, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/auth/ -run "TestLogin|TestBuild" -v`
Expected: all 3 tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/auth/login.go internal/auth/login_test.go
git commit -m "feat: add login flow with CompleteLogin and BuildAuthorizeURL"
```

---

### Task 24: Auth — Cobra Command Wiring

**Files:**
- Create: `cmd/trello/auth.go`
- Create: `cmd/trello/auth_test.go`

- [ ] **Step 1: Write the failing tests**

Create `cmd/trello/auth_test.go`:
```go
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
	// Reset root command output state to avoid cross-test contamination
	rootCmd.SetOut(nil)
	rootCmd.SetArgs(nil)
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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./cmd/trello/ -run TestAuth -v`
Expected: compilation failure — no auth commands registered

- [ ] **Step 3: Write minimal implementation**

Create `cmd/trello/auth.go`:
```go
package main

import (
	"github.com/brettmcdowell/trello-cli/internal/auth"
	"github.com/brettmcdowell/trello-cli/internal/contract"
	"github.com/brettmcdowell/trello-cli/internal/credentials"
	"github.com/spf13/cobra"
)

// trelloBaseURL is the Trello API base URL. Variable (not const) so tests can override it.
var trelloBaseURL = "https://api.trello.com"

// credStore is the credential store used by auth commands.
// Overridden in tests with a MemoryStore.
var credStore credentials.Store

func getCredStore() credentials.Store {
	if credStore != nil {
		return credStore
	}
	credStore = credentials.NewFallbackStore(
		credentials.NewKeyringStore(),
		credentials.NewEnvStore(),
	)
	return credStore
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage Trello authentication",
}

var authSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Store API key and token for subsequent commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, _ := cmd.Flags().GetString("api-key")
		token, _ := cmd.Flags().GetString("token")

		if err := contract.RequireFlag("api-key", apiKey); err != nil {
			return err
		}
		if err := contract.RequireFlag("token", token); err != nil {
			return err
		}

		result, err := auth.Set(getCredStore(), "default", apiKey, token)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), result)
	},
}

var authClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove stored credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := auth.Clear(getCredStore(), "default")
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), result)
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication state",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := auth.Status(cmd.Context(), getCredStore(), "default", trelloBaseURL)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), result)
	},
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate via interactive browser login",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Login implementation will be completed when the full interactive
		// flow (callback server + browser launch) is wired up.
		// For now, this returns UNSUPPORTED until the interactive flow is built.
		return contract.NewError(contract.Unsupported, "interactive login not yet implemented — use 'trello auth set' instead")
	},
}

func init() {
	authSetCmd.Flags().String("api-key", "", "Trello API key")
	authSetCmd.Flags().String("token", "", "Trello user token")

	authCmd.AddCommand(authSetCmd)
	authCmd.AddCommand(authClearCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authLoginCmd)
	rootCmd.AddCommand(authCmd)
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./cmd/trello/ -run TestAuth -v`
Expected: all 4 tests PASS

- [ ] **Step 5: Run full test suite**

Run: `go test ./... -v`
Expected: all tests PASS

- [ ] **Step 6: Verify from CLI**

Run: `go run ./cmd/trello/ auth set --api-key test --token test`
Expected: `{"ok":true,"data":{"configured":true,"authMode":"manual"}}`

Run: `go run ./cmd/trello/ auth clear`
Expected: `{"ok":true,"data":{"configured":false,"authMode":null}}`

- [ ] **Step 7: Commit**

```bash
git add cmd/trello/auth.go cmd/trello/auth_test.go
git commit -m "feat: wire auth commands (set, clear, status, login) into Cobra"
```

---

**End of Chunk 3**

---

## Chunk 4: API Client Layer

### Task 25: API Client — Types

**Files:**
- Create: `internal/trello/types.go`

- [ ] **Step 1: Create resource types**

Create `internal/trello/types.go`:
```go
package trello

// Board represents a Trello board.
type Board struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Closed bool   `json:"closed"`
	URL    string `json:"url"`
}

// List represents a Trello list.
type List struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Closed  bool    `json:"closed"`
	IDBoard string  `json:"idBoard"`
	Pos     float64 `json:"pos"`
}

// Card represents a Trello card.
type Card struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Desc    string  `json:"desc"`
	Closed  bool    `json:"closed"`
	IDBoard string  `json:"idBoard"`
	IDList  string  `json:"idList"`
	Due     *string `json:"due"`
	URL     string  `json:"url"`
}

// Comment represents a Trello comment action.
type Comment struct {
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Date          string        `json:"date"`
	MemberCreator MemberCreator `json:"memberCreator"`
	Data          CommentData   `json:"data"`
}

// MemberCreator is the member who created an action.
type MemberCreator struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
}

// CommentData holds the text of a comment action.
type CommentData struct {
	Text string `json:"text"`
}

// Checklist represents a Trello checklist.
type Checklist struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	IDCard     string      `json:"idCard"`
	CheckItems []CheckItem `json:"checkItems"`
}

// CheckItem represents an item in a checklist.
type CheckItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}

// Attachment represents a Trello card attachment.
type Attachment struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	URL      string  `json:"url"`
	Bytes    int     `json:"bytes"`
	MimeType string  `json:"mimeType"`
	Date     string  `json:"date"`
	IsUpload bool    `json:"isUpload"`
}

// Label represents a Trello label.
type Label struct {
	ID      string `json:"id"`
	IDBoard string `json:"idBoard"`
	Name    string `json:"name"`
	Color   string `json:"color"`
}

// Member represents a Trello member.
type Member struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
}

// CardSearchResult wraps search results for cards.
type CardSearchResult struct {
	Query string `json:"query"`
	Cards []Card `json:"cards"`
}

// BoardSearchResult wraps search results for boards.
type BoardSearchResult struct {
	Query  string  `json:"query"`
	Boards []Board `json:"boards"`
}

// UpdateListParams holds optional fields for list updates.
type UpdateListParams struct {
	Name *string  `json:"name,omitempty"`
	Pos  *float64 `json:"pos,omitempty"`
}

// CreateCardParams holds fields for card creation.
type CreateCardParams struct {
	IDList string  `json:"idList"`
	Name   string  `json:"name"`
	Desc   *string `json:"desc,omitempty"`
	Due    *string `json:"due,omitempty"`
}

// UpdateCardParams holds optional fields for card updates.
type UpdateCardParams struct {
	Name    *string `json:"name,omitempty"`
	Desc    *string `json:"desc,omitempty"`
	Due     *string `json:"due,omitempty"`
	Labels  *string `json:"idLabels,omitempty"`
	Members *string `json:"idMembers,omitempty"`
}

// DeleteResult is the response shape for delete operations.
type DeleteResult struct {
	Deleted bool   `json:"deleted"`
	ID      string `json:"id"`
}

// ActionResult is the response shape for void add/remove operations
// (labels add, labels remove, members add, members remove).
type ActionResult struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
```

- [ ] **Step 2: Verify it compiles**

Run: `go build ./internal/trello/`
Expected: successful compilation

- [ ] **Step 3: Commit**

```bash
git add internal/trello/types.go
git commit -m "feat: add Trello resource types matching command spec response fields"
```

---

### Task 26: API Client — Base HTTP Client

**Files:**
- Create: `internal/trello/client.go`
- Create: `internal/trello/client_test.go`

- [ ] **Step 1: Write the failing tests**

Create `internal/trello/client_test.go`:
```go
package trello_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/trello"
)

func TestClientInjectsAuthParams(t *testing.T) {
	var capturedURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL.String()
		json.NewEncoder(w).Encode(map[string]string{"id": "test"})
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "my-key", "my-token", trello.DefaultClientOptions())

	var result map[string]string
	err := client.Get(context.Background(), "/1/members/me", nil, &result)
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}

	// Verify auth params were injected
	if capturedURL == "" {
		t.Fatal("no request was captured")
	}
	if !containsParam(capturedURL, "key=my-key") {
		t.Errorf("URL missing key param: %s", capturedURL)
	}
	if !containsParam(capturedURL, "token=my-token") {
		t.Errorf("URL missing token param: %s", capturedURL)
	}
}

func TestClientGetDecodesResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"id":   "board1",
			"name": "Test Board",
		})
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())

	var board trello.Board
	err := client.Get(context.Background(), "/1/boards/board1", nil, &board)
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}

	if board.ID != "board1" {
		t.Errorf("ID = %q, want %q", board.ID, "board1")
	}
	if board.Name != "Test Board" {
		t.Errorf("Name = %q, want %q", board.Name, "Test Board")
	}
}

func TestClientGetWithQueryParams(t *testing.T) {
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		json.NewEncoder(w).Encode([]any{})
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())

	params := map[string]string{"filter": "open", "fields": "id,name"}
	var result []any
	client.Get(context.Background(), "/1/boards/b1/lists", params, &result)

	if !containsParam(capturedQuery, "filter=open") {
		t.Errorf("query missing filter param: %s", capturedQuery)
	}
	if !containsParam(capturedQuery, "fields=id") {
		t.Errorf("query missing fields param: %s", capturedQuery)
	}
}

func TestClientPostSendsQueryParams(t *testing.T) {
	var capturedQuery string
	var capturedMethod string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedQuery = r.URL.RawQuery
		json.NewEncoder(w).Encode(map[string]string{"id": "new1"})
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())

	params := map[string]string{"name": "New List", "idBoard": "b1"}
	var result map[string]string
	err := client.Post(context.Background(), "/1/lists", params, &result)
	if err != nil {
		t.Fatalf("Post() returned error: %v", err)
	}

	if capturedMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", capturedMethod)
	}
	if !strings.Contains(capturedQuery, "name=New+List") && !strings.Contains(capturedQuery, "name=New%20List") {
		t.Errorf("query missing name param: %s", capturedQuery)
	}
	if !strings.Contains(capturedQuery, "idBoard=b1") {
		t.Errorf("query missing idBoard param: %s", capturedQuery)
	}
}

func TestClientContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		<-r.Context().Done()
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	var result map[string]any
	err := client.Get(ctx, "/1/members/me", nil, &result)
	if err == nil {
		t.Fatal("Get() should return error when context is cancelled")
	}
}

// containsParam checks if a query string contains a key=value pair.
func containsParam(query, param string) bool {
	// Check for exact match with boundary delimiters (?, &, or start/end of string)
	for _, part := range strings.Split(query, "&") {
		if part == param || strings.TrimPrefix(part, "?") == param {
			return true
		}
	}
	return false
}
```

Add `"strings"` to the imports.

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/trello/ -run TestClient -v`
Expected: compilation failure — `trello.NewClient` doesn't exist

- [ ] **Step 3: Write minimal implementation**

Create `internal/trello/client.go`:
```go
package trello

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

// ClientOptions holds configuration for the API client.
type ClientOptions struct {
	Timeout        time.Duration
	MaxRetries     int
	RetryMutations bool
	Verbose        bool
}

// DefaultClientOptions returns sensible defaults.
func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		Timeout:        15 * time.Second,
		MaxRetries:     3,
		RetryMutations: false,
		Verbose:        false,
	}
}

// Client is the Trello API client.
type Client struct {
	baseURL    string
	apiKey     string
	token      string
	httpClient *http.Client
	opts       ClientOptions
}

// NewClient creates a new Trello API client.
func NewClient(baseURL, apiKey, token string, opts ClientOptions) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		token:   token,
		httpClient: &http.Client{
			Timeout: opts.Timeout,
		},
		opts: opts,
	}
}

// Get performs an authenticated GET request and decodes the JSON response.
func (c *Client) Get(ctx context.Context, path string, params map[string]string, result any) error {
	return c.do(ctx, http.MethodGet, path, params, nil, result)
}

// Post performs an authenticated POST request with params as query parameters.
// Trello API expects mutation data as query params, not JSON bodies.
func (c *Client) Post(ctx context.Context, path string, params map[string]string, result any) error {
	return c.do(ctx, http.MethodPost, path, params, result)
}

// Put performs an authenticated PUT request with params as query parameters.
func (c *Client) Put(ctx context.Context, path string, params map[string]string, result any) error {
	return c.do(ctx, http.MethodPut, path, params, result)
}

// Delete performs an authenticated DELETE request.
func (c *Client) Delete(ctx context.Context, path string, result any) error {
	return c.do(ctx, http.MethodDelete, path, nil, result)
}

// PostMultipart performs a multipart file upload POST (for attachments).
func (c *Client) PostMultipart(ctx context.Context, path string, params map[string]string, filePath string, result any) error {
	// Multipart upload implementation — handled in attachments.go
	return nil
}

// buildURL constructs the full URL with auth query params (key, token) and any
// additional request params. Shared by both do() and postMultipartFile().
func (c *Client) buildURL(path string, params map[string]string) string {
	u, _ := url.Parse(c.baseURL + path)
	q := u.Query()
	q.Set("key", c.apiKey)
	q.Set("token", c.token)
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func (c *Client) do(ctx context.Context, method, path string, params map[string]string, result any) error {
	// Use buildURL for centralized auth injection.
	// Trello API expects all data (including mutations) as query params.
	fullURL := c.buildURL(path, params)

	req, err := http.NewRequestWithContext(ctx, method, fullURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Log request if verbose
	start := time.Now()
	if c.opts.Verbose {
		logURL := c.baseURL + path
		slog.Debug("trello request", "method", method, "url", logURL)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if c.opts.Verbose {
		slog.Debug("trello response", "status", resp.StatusCode, "duration", time.Since(start))
	}

	// Handle error status codes
	if resp.StatusCode >= 400 {
		return mapHTTPError(resp)
	}

	// Decode response
	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/trello/ -run TestClient -v`
Expected: compilation failure — `mapHTTPError` doesn't exist yet. Create a stub.

- [ ] **Step 5: Create error mapping stub**

Create `internal/trello/errors.go`:
```go
package trello

import (
	"fmt"
	"net/http"

	"github.com/brettmcdowell/trello-cli/internal/contract"
)

// mapHTTPError converts an HTTP error response to a ContractError.
func mapHTTPError(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return contract.NewError(contract.AuthInvalid, "Trello rejected the credentials")
	case http.StatusNotFound:
		return contract.NewError(contract.NotFound, "resource not found")
	case http.StatusTooManyRequests:
		return contract.NewError(contract.RateLimited, "rate limited by Trello API")
	default:
		return contract.NewError(contract.HTTPError, fmt.Sprintf("Trello API returned status %d", resp.StatusCode))
	}
}
```

- [ ] **Step 6: Run tests to verify they pass**

Run: `go test ./internal/trello/ -run TestClient -v`
Expected: all 5 tests PASS

- [ ] **Step 7: Commit**

```bash
git add internal/trello/client.go internal/trello/client_test.go internal/trello/errors.go
git commit -m "feat: add Trello API client with auth injection and HTTP error mapping"
```

---

### Task 27: API Client — Error Mapping Tests

**Files:**
- Create: `internal/trello/errors_test.go`

- [x] **Step 1: Write the tests**

Create `internal/trello/errors_test.go`:
```go
package trello_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/contract"
	"github.com/brettmcdowell/trello-cli/internal/trello"
)

func TestHTTPError401(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	var result map[string]any
	err := client.Get(context.Background(), "/1/members/me", nil, &result)

	var ce *contract.ContractError
	if !errors.As(err, &ce) {
		t.Fatalf("error should be *ContractError, got %T: %v", err, err)
	}
	if ce.Code != contract.AuthInvalid {
		t.Errorf("Code = %q, want %q", ce.Code, contract.AuthInvalid)
	}
}

func TestHTTPError404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	var result map[string]any
	err := client.Get(context.Background(), "/1/boards/nonexistent", nil, &result)

	var ce *contract.ContractError
	if !errors.As(err, &ce) {
		t.Fatalf("error should be *ContractError, got %T", err)
	}
	if ce.Code != contract.NotFound {
		t.Errorf("Code = %q, want %q", ce.Code, contract.NotFound)
	}
}

func TestHTTPError429(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	// Disable retries for this test
	opts := trello.DefaultClientOptions()
	opts.MaxRetries = 0
	client := trello.NewClient(server.URL, "k", "t", opts)

	var result map[string]any
	err := client.Get(context.Background(), "/1/boards/b1", nil, &result)

	var ce *contract.ContractError
	if !errors.As(err, &ce) {
		t.Fatalf("error should be *ContractError, got %T", err)
	}
	if ce.Code != contract.RateLimited {
		t.Errorf("Code = %q, want %q", ce.Code, contract.RateLimited)
	}
}

func TestHTTPError500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	var result map[string]any
	err := client.Get(context.Background(), "/1/boards/b1", nil, &result)

	var ce *contract.ContractError
	if !errors.As(err, &ce) {
		t.Fatalf("error should be *ContractError, got %T", err)
	}
	if ce.Code != contract.HTTPError {
		t.Errorf("Code = %q, want %q", ce.Code, contract.HTTPError)
	}
}
```

- [x] **Step 2: Run tests to verify they pass**

Run: `go test ./internal/trello/ -run "TestHTTPError" -v`
Expected: all 4 tests PASS

- [x] **Step 3: Commit**

```bash
git add internal/trello/errors_test.go
git commit -m "test: add HTTP error mapping tests (401, 404, 429, 500)"
```

---

### Task 28: API Client — Retry Logic

**Files:**
- Create: `internal/trello/retry.go`
- Create: `internal/trello/retry_test.go`
- Modify: `internal/trello/client.go` (integrate retry into `do()`)

- [ ] **Step 1: Write the failing tests**

Create `internal/trello/retry_test.go`:
```go
package trello_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/trello"
)

func TestRetryOnRateLimit(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n <= 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"id": "b1"})
	}))
	defer server.Close()

	opts := trello.DefaultClientOptions()
	opts.MaxRetries = 3
	client := trello.NewClient(server.URL, "k", "t", opts)

	var result map[string]string
	err := client.Get(context.Background(), "/1/boards/b1", nil, &result)
	if err != nil {
		t.Fatalf("Get() should succeed after retries, got: %v", err)
	}
	if result["id"] != "b1" {
		t.Errorf("id = %q, want %q", result["id"], "b1")
	}
	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("attempts = %d, want 3", atomic.LoadInt32(&attempts))
	}
}

func TestRetryExhausted(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	opts := trello.DefaultClientOptions()
	opts.MaxRetries = 2
	client := trello.NewClient(server.URL, "k", "t", opts)

	var result map[string]string
	err := client.Get(context.Background(), "/1/boards/b1", nil, &result)
	if err == nil {
		t.Fatal("Get() should fail after exhausting retries")
	}
	// 1 initial + 2 retries = 3 total
	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("attempts = %d, want 3", atomic.LoadInt32(&attempts))
	}
}

func TestNoRetryOnMutations(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	opts := trello.DefaultClientOptions()
	opts.MaxRetries = 3
	opts.RetryMutations = false
	client := trello.NewClient(server.URL, "k", "t", opts)

	var result map[string]string
	client.Post(context.Background(), "/1/lists", map[string]string{"name": "test"}, &result) //nolint

	if atomic.LoadInt32(&attempts) != 1 {
		t.Errorf("attempts = %d, want 1 (no retry on mutations)", atomic.LoadInt32(&attempts))
	}
}

func TestRetryOnMutationsWhenEnabled(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n <= 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"id": "new1"})
	}))
	defer server.Close()

	opts := trello.DefaultClientOptions()
	opts.MaxRetries = 3
	opts.RetryMutations = true
	client := trello.NewClient(server.URL, "k", "t", opts)

	var result map[string]string
	err := client.Post(context.Background(), "/1/lists", map[string]string{"name": "test"}, &result)
	if err != nil {
		t.Fatalf("Post() should succeed after retry, got: %v", err)
	}
	if atomic.LoadInt32(&attempts) != 2 {
		t.Errorf("attempts = %d, want 2", atomic.LoadInt32(&attempts))
	}
}

func TestNoRetryOnNon429Errors(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	opts := trello.DefaultClientOptions()
	opts.MaxRetries = 3
	client := trello.NewClient(server.URL, "k", "t", opts)

	var result map[string]string
	client.Get(context.Background(), "/1/boards/nope", nil, &result)

	if atomic.LoadInt32(&attempts) != 1 {
		t.Errorf("attempts = %d, want 1 (no retry on 404)", atomic.LoadInt32(&attempts))
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/trello/ -run TestRetry -v`
Expected: failures — no retry logic exists yet

- [x] **Step 3: Write retry implementation**

Create `internal/trello/retry.go`:
```go
package trello

import (
	"context"
	"math"
	"math/rand"
	"net/http"
	"time"
)

// shouldRetry returns true if the response status code is retryable.
func shouldRetry(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests
}

// isMutation returns true if the HTTP method is not safe (not GET/HEAD).
func isMutation(method string) bool {
	return method != http.MethodGet && method != http.MethodHead
}

// backoff returns the duration to wait before the nth retry (0-indexed).
// Uses exponential backoff with jitter.
func backoff(attempt int) time.Duration {
	base := math.Pow(2, float64(attempt)) * 500 // 500ms, 1s, 2s, 4s...
	jitter := rand.Float64() * base * 0.5        // up to 50% jitter
	return time.Duration(base+jitter) * time.Millisecond
}

// waitForRetry sleeps for the backoff duration, respecting context cancellation.
func waitForRetry(ctx context.Context, attempt int) error {
	d := backoff(attempt)
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
```

- [x] **Step 4: Integrate retry into client.do()**

Modify `internal/trello/client.go` — replace the `do` method with a retry-aware version:

```go
// buildURL constructs the full URL with auth query params (key, token) and any
// additional request params. Shared by both do() and postMultipartFile().
func (c *Client) buildURL(path string, params map[string]string) string {
	u, _ := url.Parse(c.baseURL + path)
	q := u.Query()
	q.Set("key", c.apiKey)
	q.Set("token", c.token)
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func (c *Client) do(ctx context.Context, method, path string, params map[string]string, result any) error {
	fullURL := c.buildURL(path, params)

	maxAttempts := 1 + c.opts.MaxRetries
	if isMutation(method) && !c.opts.RetryMutations {
		maxAttempts = 1
	}

	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			if err := waitForRetry(ctx, attempt-1); err != nil {
				return lastErr
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		start := time.Now()
		if c.opts.Verbose {
			logURL := c.baseURL + path
			slog.Debug("trello request", "method", method, "url", logURL, "attempt", attempt+1)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		if c.opts.Verbose {
			slog.Debug("trello response", "status", resp.StatusCode, "duration", time.Since(start))
		}

		if resp.StatusCode >= 400 {
			resp.Body.Close()
			lastErr = mapHTTPError(resp)
			if shouldRetry(resp.StatusCode) && attempt < maxAttempts-1 {
				continue
			}
			return lastErr
		}

		// Success — decode and close body (no defer in loop)
		if result != nil {
			decodeErr := json.NewDecoder(resp.Body).Decode(result)
			resp.Body.Close()
			if decodeErr != nil {
				return fmt.Errorf("failed to decode response: %w", decodeErr)
			}
		} else {
			resp.Body.Close()
		}
		return nil
	}

	return lastErr
}
```

- [x] **Step 5: Run retry tests to verify they pass**

Run: `go test ./internal/trello/ -run TestRetry -v`
Expected: all 5 tests PASS

- [x] **Step 6: Run full trello package tests**

Run: `go test ./internal/trello/ -v`
Expected: all tests PASS

- [x] **Step 7: Commit**

```bash
git add internal/trello/retry.go internal/trello/retry_test.go internal/trello/client.go
git commit -m "feat: add retry logic with exponential backoff and mutation safety"
```

---

### Task 29: API Client — Client Interface

**Files:**
- Modify: `internal/trello/client.go`

- [x] **Step 1: Define the Client interface for mocking**

Add to `internal/trello/client.go`:
```go
// API defines the interface for Trello API operations.
// Command handlers depend on this interface; tests mock it.
type API interface {
	// Boards
	ListBoards(ctx context.Context) ([]Board, error)
	GetBoard(ctx context.Context, boardID string) (Board, error)
	// Lists
	ListLists(ctx context.Context, boardID string) ([]List, error)
	CreateList(ctx context.Context, boardID, name string) (List, error)
	UpdateList(ctx context.Context, listID string, params UpdateListParams) (List, error)
	ArchiveList(ctx context.Context, listID string) (List, error)
	MoveList(ctx context.Context, listID, boardID string, pos *float64) (List, error)
	// Cards
	ListCardsByBoard(ctx context.Context, boardID string) ([]Card, error)
	ListCardsByList(ctx context.Context, listID string) ([]Card, error)
	GetCard(ctx context.Context, cardID string) (Card, error)
	CreateCard(ctx context.Context, params CreateCardParams) (Card, error)
	UpdateCard(ctx context.Context, cardID string, params UpdateCardParams) (Card, error)
	MoveCard(ctx context.Context, cardID, listID string, pos *float64) (Card, error)
	ArchiveCard(ctx context.Context, cardID string) (Card, error)
	DeleteCard(ctx context.Context, cardID string) error
	// Comments
	ListComments(ctx context.Context, cardID string) ([]Comment, error)
	AddComment(ctx context.Context, cardID, text string) (Comment, error)
	UpdateComment(ctx context.Context, actionID, text string) (Comment, error)
	DeleteComment(ctx context.Context, actionID string) error
	// Checklists
	ListChecklists(ctx context.Context, cardID string) ([]Checklist, error)
	CreateChecklist(ctx context.Context, cardID, name string) (Checklist, error)
	DeleteChecklist(ctx context.Context, checklistID string) error
	AddCheckItem(ctx context.Context, checklistID, name string) (CheckItem, error)
	UpdateCheckItem(ctx context.Context, cardID, itemID, state string) (CheckItem, error)
	DeleteCheckItem(ctx context.Context, checklistID, itemID string) error
	// Attachments
	ListAttachments(ctx context.Context, cardID string) ([]Attachment, error)
	AddFileAttachment(ctx context.Context, cardID, filePath string, name *string) (Attachment, error)
	AddURLAttachment(ctx context.Context, cardID, urlStr string, name *string) (Attachment, error)
	DeleteAttachment(ctx context.Context, cardID, attachmentID string) error
	// Labels
	ListLabels(ctx context.Context, boardID string) ([]Label, error)
	CreateLabel(ctx context.Context, boardID, name, color string) (Label, error)
	AddLabelToCard(ctx context.Context, cardID, labelID string) error
	RemoveLabelFromCard(ctx context.Context, cardID, labelID string) error
	// Members
	ListMembers(ctx context.Context, boardID string) ([]Member, error)
	AddMemberToCard(ctx context.Context, cardID, memberID string) error
	RemoveMemberFromCard(ctx context.Context, cardID, memberID string) error
	// Search
	SearchCards(ctx context.Context, query string) (CardSearchResult, error)
	SearchBoards(ctx context.Context, query string) (BoardSearchResult, error)
	// Auth
	GetMe(ctx context.Context) (Member, error)
}
```

- [x] **Step 2: Verify Client struct can satisfy the interface (compile check)**

Add a compile-time assertion at the bottom of `client.go`:
```go
// Compile-time check that Client implements API.
var _ API = (*Client)(nil)
```

This will fail to compile until all methods are implemented (in Chunk 5-7). For now, add stub methods:

```go
// Stub methods — implemented in resource-specific files (boards.go, lists.go, etc.)
func (c *Client) ListBoards(ctx context.Context) ([]Board, error)                        { return nil, nil }
func (c *Client) GetBoard(ctx context.Context, boardID string) (Board, error)             { return Board{}, nil }
func (c *Client) ListLists(ctx context.Context, boardID string) ([]List, error)           { return nil, nil }
func (c *Client) CreateList(ctx context.Context, boardID, name string) (List, error)      { return List{}, nil }
func (c *Client) UpdateList(ctx context.Context, listID string, params UpdateListParams) (List, error) { return List{}, nil }
func (c *Client) ArchiveList(ctx context.Context, listID string) (List, error)            { return List{}, nil }
func (c *Client) MoveList(ctx context.Context, listID, boardID string, pos *float64) (List, error) { return List{}, nil }
func (c *Client) ListCardsByBoard(ctx context.Context, boardID string) ([]Card, error)    { return nil, nil }
func (c *Client) ListCardsByList(ctx context.Context, listID string) ([]Card, error)      { return nil, nil }
func (c *Client) GetCard(ctx context.Context, cardID string) (Card, error)                { return Card{}, nil }
func (c *Client) CreateCard(ctx context.Context, params CreateCardParams) (Card, error)   { return Card{}, nil }
func (c *Client) UpdateCard(ctx context.Context, cardID string, params UpdateCardParams) (Card, error) { return Card{}, nil }
func (c *Client) MoveCard(ctx context.Context, cardID, listID string, pos *float64) (Card, error) { return Card{}, nil }
func (c *Client) ArchiveCard(ctx context.Context, cardID string) (Card, error)            { return Card{}, nil }
func (c *Client) DeleteCard(ctx context.Context, cardID string) error                     { return nil }
func (c *Client) ListComments(ctx context.Context, cardID string) ([]Comment, error)      { return nil, nil }
func (c *Client) AddComment(ctx context.Context, cardID, text string) (Comment, error)    { return Comment{}, nil }
func (c *Client) UpdateComment(ctx context.Context, actionID, text string) (Comment, error) { return Comment{}, nil }
func (c *Client) DeleteComment(ctx context.Context, actionID string) error                { return nil }
func (c *Client) ListChecklists(ctx context.Context, cardID string) ([]Checklist, error)  { return nil, nil }
func (c *Client) CreateChecklist(ctx context.Context, cardID, name string) (Checklist, error) { return Checklist{}, nil }
func (c *Client) DeleteChecklist(ctx context.Context, checklistID string) error           { return nil }
func (c *Client) AddCheckItem(ctx context.Context, checklistID, name string) (CheckItem, error) { return CheckItem{}, nil }
func (c *Client) UpdateCheckItem(ctx context.Context, cardID, itemID, state string) (CheckItem, error) { return CheckItem{}, nil }
func (c *Client) DeleteCheckItem(ctx context.Context, checklistID, itemID string) error   { return nil }
func (c *Client) ListAttachments(ctx context.Context, cardID string) ([]Attachment, error) { return nil, nil }
func (c *Client) AddFileAttachment(ctx context.Context, cardID, filePath string, name *string) (Attachment, error) { return Attachment{}, nil }
func (c *Client) AddURLAttachment(ctx context.Context, cardID, urlStr string, name *string) (Attachment, error) { return Attachment{}, nil }
func (c *Client) DeleteAttachment(ctx context.Context, cardID, attachmentID string) error { return nil }
func (c *Client) ListLabels(ctx context.Context, boardID string) ([]Label, error)        { return nil, nil }
func (c *Client) CreateLabel(ctx context.Context, boardID, name, color string) (Label, error) { return Label{}, nil }
func (c *Client) AddLabelToCard(ctx context.Context, cardID, labelID string) error        { return nil }
func (c *Client) RemoveLabelFromCard(ctx context.Context, cardID, labelID string) error   { return nil }
func (c *Client) ListMembers(ctx context.Context, boardID string) ([]Member, error)       { return nil, nil }
func (c *Client) AddMemberToCard(ctx context.Context, cardID, memberID string) error      { return nil }
func (c *Client) RemoveMemberFromCard(ctx context.Context, cardID, memberID string) error { return nil }
func (c *Client) SearchCards(ctx context.Context, query string) (CardSearchResult, error) { return CardSearchResult{}, nil }
func (c *Client) SearchBoards(ctx context.Context, query string) (BoardSearchResult, error) { return BoardSearchResult{}, nil }
func (c *Client) GetMe(ctx context.Context) (Member, error)                              { return Member{}, nil }
```

- [x] **Step 3: Verify compilation**

Run: `go build ./internal/trello/`
Expected: successful compilation

- [x] **Step 4: Run full test suite**

Run: `go test ./... -v`
Expected: all tests PASS

- [x] **Step 5: Commit**

```bash
git add internal/trello/client.go
git commit -m "feat: add API interface and stub methods for compile-time verification"
```

---

**End of Chunk 4**

---

## Chunk 5: Resource Commands — Batch A (Boards, Lists, Cards)

Each resource command follows the same TDD recipe established in the design spec:
1. Write command handler test (mock API, assert envelope)
2. Write command handler
3. Write API client test (httptest, assert request/response)
4. Write API client method (replace stub)
5. Run all tests, commit

### Task 30: Boards — API Client Methods

**Files:**
- Create: `internal/trello/boards.go`
- Create: `internal/trello/boards_test.go`

- [ ] **Step 1: Write the failing tests**

Create `internal/trello/boards_test.go`:
```go
package trello_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/trello"
)

func TestListBoards(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/1/members/me/boards" {
			t.Errorf("path = %s, want /1/members/me/boards", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": "b1", "name": "Board One", "desc": "First", "closed": false, "url": "https://trello.com/b/b1"},
			{"id": "b2", "name": "Board Two", "desc": "", "closed": true, "url": "https://trello.com/b/b2"},
		})
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	boards, err := client.ListBoards(context.Background())
	if err != nil {
		t.Fatalf("ListBoards() error: %v", err)
	}
	if len(boards) != 2 {
		t.Fatalf("len = %d, want 2", len(boards))
	}
	if boards[0].ID != "b1" || boards[0].Name != "Board One" {
		t.Errorf("boards[0] = %+v", boards[0])
	}
	if boards[1].Closed != true {
		t.Errorf("boards[1].Closed = %v, want true", boards[1].Closed)
	}
}

func TestGetBoard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/1/boards/b1" {
			t.Errorf("path = %s, want /1/boards/b1", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "b1", "name": "My Board", "desc": "A board", "closed": false, "url": "https://trello.com/b/b1",
		})
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	board, err := client.GetBoard(context.Background(), "b1")
	if err != nil {
		t.Fatalf("GetBoard() error: %v", err)
	}
	if board.ID != "b1" {
		t.Errorf("ID = %q, want %q", board.ID, "b1")
	}
	if board.Name != "My Board" {
		t.Errorf("Name = %q, want %q", board.Name, "My Board")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail** (stubs return nil/empty)

Run: `go test ./internal/trello/ -run "TestListBoards|TestGetBoard" -v`
Expected: FAIL (stubs return empty results)

- [ ] **Step 3: Implement client methods**

Create `internal/trello/boards.go`:
```go
package trello

import (
	"context"
	"fmt"
)

func (c *Client) ListBoards(ctx context.Context) ([]Board, error) {
	var boards []Board
	err := c.Get(ctx, "/1/members/me/boards", nil, &boards)
	return boards, err
}

func (c *Client) GetBoard(ctx context.Context, boardID string) (Board, error) {
	var board Board
	err := c.Get(ctx, fmt.Sprintf("/1/boards/%s", boardID), nil, &board)
	return board, err
}
```

Remove the `ListBoards` and `GetBoard` stubs from `client.go`.

- [ ] **Step 4: Run tests**

Run: `go test ./internal/trello/ -run "TestListBoards|TestGetBoard" -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/trello/boards.go internal/trello/boards_test.go internal/trello/client.go
git commit -m "feat: implement ListBoards and GetBoard API client methods"
```

---

### Task 31: Boards — Cobra Commands

**Files:**
- Create: `cmd/trello/boards.go`
- Create: `cmd/trello/boards_test.go`

- [ ] **Step 1: Write the failing tests**

Create `cmd/trello/boards_test.go`:
```go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/credentials"
	"github.com/brettmcdowell/trello-cli/internal/trello"
)

// mockAPI implements trello.API for command testing.
type mockAPI struct {
	trello.API // embed to satisfy interface; override specific methods
	listBoardsFn func(ctx context.Context) ([]trello.Board, error)
	getBoardFn   func(ctx context.Context, id string) (trello.Board, error)
}

func (m *mockAPI) ListBoards(ctx context.Context) ([]trello.Board, error) {
	if m.listBoardsFn != nil {
		return m.listBoardsFn(ctx)
	}
	return nil, nil
}

func (m *mockAPI) GetBoard(ctx context.Context, id string) (trello.Board, error) {
	if m.getBoardFn != nil {
		return m.getBoardFn(ctx, id)
	}
	return trello.Board{}, nil
}

func TestBoardsListCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		listBoardsFn: func(ctx context.Context) ([]trello.Board, error) {
			return []trello.Board{
				{ID: "b1", Name: "Board One"},
				{ID: "b2", Name: "Board Two"},
			}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"boards", "list"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("boards list failed: %v", err)
	}

	var envelope map[string]any
	json.Unmarshal(buf.Bytes(), &envelope)

	if envelope["ok"] != true {
		t.Errorf("ok = %v, want true", envelope["ok"])
	}
	data, ok := envelope["data"].([]any)
	if !ok {
		t.Fatal("data should be an array")
	}
	if len(data) != 2 {
		t.Errorf("len(data) = %d, want 2", len(data))
	}
}

func TestBoardsListRequiresAuth(t *testing.T) {
	setupTestAuth(t)
	// No credentials set

	rootCmd.SetArgs([]string{"boards", "list"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("boards list should fail without auth")
	}
}

func TestBoardsGetCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		getBoardFn: func(ctx context.Context, id string) (trello.Board, error) {
			if id != "b1" {
				t.Errorf("board ID = %q, want %q", id, "b1")
			}
			return trello.Board{ID: "b1", Name: "My Board"}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"boards", "get", "--board", "b1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("boards get failed: %v", err)
	}

	var envelope map[string]any
	json.Unmarshal(buf.Bytes(), &envelope)

	data := envelope["data"].(map[string]any)
	if data["id"] != "b1" {
		t.Errorf("data.id = %v, want b1", data["id"])
	}
}

func TestBoardsGetMissingFlag(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	rootCmd.SetArgs([]string{"boards", "get"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("boards get should fail without --board")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./cmd/trello/ -run TestBoards -v`
Expected: compilation failure

- [ ] **Step 3: Implement commands**

Create `cmd/trello/boards.go`:
```go
package main

import (
	"github.com/brettmcdowell/trello-cli/internal/auth"
	"github.com/brettmcdowell/trello-cli/internal/contract"
	"github.com/brettmcdowell/trello-cli/internal/trello"
	"github.com/spf13/cobra"
)

// apiClient is the Trello API client. Overridden in tests with a mock.
var apiClient trello.API

func getAPIClient(creds credentials.Credentials) trello.API {
	if apiClient != nil {
		return apiClient
	}
	apiClient = trello.NewClient(trelloBaseURL, creds.APIKey, creds.Token, trello.DefaultClientOptions())
	return apiClient
}

var boardsCmd = &cobra.Command{
	Use:   "boards",
	Short: "Manage Trello boards",
}

var boardsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List visible boards",
	RunE: func(cmd *cobra.Command, args []string) error {
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		boards, err := getAPIClient(creds).ListBoards(cmd.Context())
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), boards)
	},
}

var boardsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a board by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		boardID, _ := cmd.Flags().GetString("board")
		if err := contract.RequireFlag("board", boardID); err != nil {
			return err
		}
		creds, err := auth.RequireAuth(getCredStore(), "default")
		if err != nil {
			return err
		}
		board, err := getAPIClient(creds).GetBoard(cmd.Context(), boardID)
		if err != nil {
			return err
		}
		return output(cmd.OutOrStdout(), board)
	},
}

func init() {
	boardsGetCmd.Flags().String("board", "", "Board ID")

	boardsCmd.AddCommand(boardsListCmd)
	boardsCmd.AddCommand(boardsGetCmd)
	rootCmd.AddCommand(boardsCmd)
}
```

- [ ] **Step 4: Run tests**

Run: `go test ./cmd/trello/ -run TestBoards -v`
Expected: all 4 tests PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/trello/boards.go cmd/trello/boards_test.go
git commit -m "feat: add boards list and boards get commands"
```

---

### Task 32: Lists — API Client Methods

**Files:**
- Create: `internal/trello/lists.go`
- Create: `internal/trello/lists_test.go`

- [ ] **Step 1: Write failing tests**

Create `internal/trello/lists_test.go` with tests for all 5 list methods:
- `TestListLists` — GET /1/boards/{id}/lists, verify path and decoded response
- `TestCreateList` — POST /1/lists, verify query params (idBoard, name)
- `TestUpdateList` — PUT /1/lists/{id}, verify optional params
- `TestArchiveList` — PUT /1/lists/{id}, verify `closed=true` in query
- `TestMoveList` — PUT /1/lists/{id}, verify `idBoard` and optional `pos` in query

- [ ] **Step 2: Run tests to verify they fail**

- [ ] **Step 3: Implement methods**

Create `internal/trello/lists.go`:
```go
package trello

import (
	"context"
	"fmt"
	"strconv"
)

func (c *Client) ListLists(ctx context.Context, boardID string) ([]List, error) {
	var lists []List
	err := c.Get(ctx, fmt.Sprintf("/1/boards/%s/lists", boardID), nil, &lists)
	return lists, err
}

func (c *Client) CreateList(ctx context.Context, boardID, name string) (List, error) {
	var list List
	err := c.Post(ctx, "/1/lists", map[string]string{
		"idBoard": boardID,
		"name":    name,
	}, &list)
	return list, err
}

func (c *Client) UpdateList(ctx context.Context, listID string, params UpdateListParams) (List, error) {
	qp := map[string]string{}
	if params.Name != nil {
		qp["name"] = *params.Name
	}
	if params.Pos != nil {
		qp["pos"] = strconv.FormatFloat(*params.Pos, 'f', -1, 64)
	}
	var list List
	err := c.Put(ctx, fmt.Sprintf("/1/lists/%s", listID), qp, &list)
	return list, err
}

func (c *Client) ArchiveList(ctx context.Context, listID string) (List, error) {
	var list List
	err := c.Put(ctx, fmt.Sprintf("/1/lists/%s/closed", listID), map[string]string{
		"value": "true",
	}, &list)
	return list, err
}

func (c *Client) MoveList(ctx context.Context, listID, boardID string, pos *float64) (List, error) {
	qp := map[string]string{"idBoard": boardID}
	if pos != nil {
		qp["pos"] = strconv.FormatFloat(*pos, 'f', -1, 64)
	}
	var list List
	err := c.Put(ctx, fmt.Sprintf("/1/lists/%s", listID), qp, &list)
	return list, err
}
```

Remove list stubs from `client.go`.

- [ ] **Step 4: Run tests, verify PASS**
- [ ] **Step 5: Commit**

```bash
git add internal/trello/lists.go internal/trello/lists_test.go internal/trello/client.go
git commit -m "feat: implement list API client methods (list, create, update, archive, move)"
```

---

### Task 33: Lists — Cobra Commands

**Files:**
- Create: `cmd/trello/lists.go`
- Create: `cmd/trello/lists_test.go`

Follow the same pattern as boards. All command handlers use the single-auth-check pattern: `creds, err := auth.RequireAuth(...)` then `getAPIClient(creds)`.

Commands and required test cases:
- `lists list --board <id>` — requires `--board`
  - Test: success returns JSON array of lists
  - Test: missing `--board` returns VALIDATION_ERROR
  - Test: auth guard returns AUTH_REQUIRED when no creds
- `lists create --board <id> --name <name>` — requires both
  - Test: success returns created list
  - Test: missing `--board` returns VALIDATION_ERROR
  - Test: missing `--name` returns VALIDATION_ERROR
- `lists update --list <id> [--name] [--pos]` — requires `--list` + at least one mutation
  - Test: success with `--name` returns updated list
  - Test: missing `--list` returns VALIDATION_ERROR
  - Test: no mutation flags returns VALIDATION_ERROR (RequireAtLeastOne)
- `lists archive --list <id>` — requires `--list`
  - Test: success returns list with closed=true
  - Test: missing `--list` returns VALIDATION_ERROR
- `lists move --list <id> --board <id> [--pos]` — requires both `--list` and `--board`
  - Test: success returns moved list
  - Test: missing `--list` returns VALIDATION_ERROR
  - Test: missing `--board` returns VALIDATION_ERROR

- [ ] **Step 1-4: Write tests, implement, verify**

- [ ] **Step 5: Commit**

```bash
git add cmd/trello/lists.go cmd/trello/lists_test.go
git commit -m "feat: add lists commands (list, create, update, archive, move)"
```

---

### Task 34: Cards — API Client Methods

**Files:**
- Create: `internal/trello/cards.go`
- Create: `internal/trello/cards_test.go`

- [ ] **Step 1: Write failing tests** for all 7 card methods:
- `TestListCardsByBoard` — GET /1/boards/{id}/cards
- `TestListCardsByList` — GET /1/lists/{id}/cards
- `TestGetCard` — GET /1/cards/{id}
- `TestCreateCard` — POST /1/cards, verify idList, name, optional desc/due in query
- `TestUpdateCard` — PUT /1/cards/{id}, verify optional params
- `TestMoveCard` — PUT /1/cards/{id}, verify idList and optional pos
- `TestArchiveCard` — PUT /1/cards/{id}, verify closed=true
- `TestDeleteCard` — DELETE /1/cards/{id}

- [ ] **Step 2: Implement methods**

Create `internal/trello/cards.go`:
```go
package trello

import (
	"context"
	"fmt"
	"strconv"
)

func (c *Client) ListCardsByBoard(ctx context.Context, boardID string) ([]Card, error) {
	var cards []Card
	err := c.Get(ctx, fmt.Sprintf("/1/boards/%s/cards", boardID), nil, &cards)
	return cards, err
}

func (c *Client) ListCardsByList(ctx context.Context, listID string) ([]Card, error) {
	var cards []Card
	err := c.Get(ctx, fmt.Sprintf("/1/lists/%s/cards", listID), nil, &cards)
	return cards, err
}

func (c *Client) GetCard(ctx context.Context, cardID string) (Card, error) {
	var card Card
	err := c.Get(ctx, fmt.Sprintf("/1/cards/%s", cardID), nil, &card)
	return card, err
}

func (c *Client) CreateCard(ctx context.Context, params CreateCardParams) (Card, error) {
	qp := map[string]string{
		"idList": params.IDList,
		"name":   params.Name,
	}
	if params.Desc != nil {
		qp["desc"] = *params.Desc
	}
	if params.Due != nil {
		qp["due"] = *params.Due
	}
	var card Card
	err := c.Post(ctx, "/1/cards", qp, &card)
	return card, err
}

func (c *Client) UpdateCard(ctx context.Context, cardID string, params UpdateCardParams) (Card, error) {
	qp := map[string]string{}
	if params.Name != nil {
		qp["name"] = *params.Name
	}
	if params.Desc != nil {
		qp["desc"] = *params.Desc
	}
	if params.Due != nil {
		qp["due"] = *params.Due
	}
	if params.Labels != nil {
		qp["idLabels"] = *params.Labels
	}
	if params.Members != nil {
		qp["idMembers"] = *params.Members
	}
	var card Card
	err := c.Put(ctx, fmt.Sprintf("/1/cards/%s", cardID), qp, &card)
	return card, err
}

func (c *Client) MoveCard(ctx context.Context, cardID, listID string, pos *float64) (Card, error) {
	qp := map[string]string{"idList": listID}
	if pos != nil {
		qp["pos"] = strconv.FormatFloat(*pos, 'f', -1, 64)
	}
	var card Card
	err := c.Put(ctx, fmt.Sprintf("/1/cards/%s", cardID), qp, &card)
	return card, err
}

func (c *Client) ArchiveCard(ctx context.Context, cardID string) (Card, error) {
	var card Card
	err := c.Put(ctx, fmt.Sprintf("/1/cards/%s/closed", cardID), map[string]string{
		"value": "true",
	}, &card)
	return card, err
}

func (c *Client) DeleteCard(ctx context.Context, cardID string) error {
	return c.Delete(ctx, fmt.Sprintf("/1/cards/%s", cardID), nil)
}
```

Remove card stubs from `client.go`.

- [ ] **Step 3: Run tests, verify PASS**
- [ ] **Step 4: Commit**

```bash
git add internal/trello/cards.go internal/trello/cards_test.go internal/trello/client.go
git commit -m "feat: implement card API client methods (list, get, create, update, move, archive, delete)"
```

---

### Task 35: Cards — Cobra Commands

**Files:**
- Create: `cmd/trello/cards.go`
- Create: `cmd/trello/cards_test.go`

All command handlers use the single-auth-check pattern: `creds, err := auth.RequireAuth(...)` then `getAPIClient(creds)`.

Commands, validation patterns, and required test cases:
- `cards list --board <id>|--list <id>` — `RequireExactlyOne({"board": boardID, "list": listID})`
  - Test: success with `--board` returns cards
  - Test: success with `--list` returns cards
  - Test: both `--board` and `--list` returns VALIDATION_ERROR
  - Test: neither flag returns VALIDATION_ERROR
- `cards get --card <id>` — `RequireFlag("card")`
  - Test: success returns card
  - Test: missing `--card` returns VALIDATION_ERROR
- `cards create --list <id> --name <name> [--desc] [--due] [--labels] [--members]`
  - Test: success with required flags returns created card
  - Test: missing `--list` returns VALIDATION_ERROR
  - Test: missing `--name` returns VALIDATION_ERROR
  - Test: invalid `--due` (not ISO8601) returns VALIDATION_ERROR
  - Test: valid `--due` (ISO8601) succeeds
- `cards update --card <id> [--name] [--desc] [--due] [--labels] [--members]`
  - Test: success with `--name` returns updated card
  - Test: missing `--card` returns VALIDATION_ERROR
  - Test: no mutation flags returns VALIDATION_ERROR (RequireAtLeastOne)
  - Test: invalid `--due` returns VALIDATION_ERROR
- `cards move --card <id> --list <id> [--pos]`
  - Test: success returns moved card
  - Test: missing `--card` or `--list` returns VALIDATION_ERROR
- `cards archive --card <id>` — requires `--card`
  - Test: success returns card with closed=true
- `cards delete --card <id>` — returns `DeleteResult{Deleted: true, ID: cardID}` envelope
  - Test: success returns DeleteResult JSON
  - Test: missing `--card` returns VALIDATION_ERROR

- [ ] **Step 1-4: Write tests, implement, verify**
- [ ] **Step 5: Commit**

```bash
git add cmd/trello/cards.go cmd/trello/cards_test.go
git commit -m "feat: add cards commands (list, get, create, update, move, archive, delete)"
```

---

### Task 36: Run Full Test Suite and Commit Batch A

- [ ] **Step 1: Run full test suite**

Run: `go test ./... -v`
Expected: all tests PASS

- [ ] **Step 2: Commit any remaining cleanup**

```bash
git add -A
git commit -m "chore: complete Batch A — boards, lists, cards commands"
```

---

**End of Chunk 5**

---

## Chunk 6: Resource Commands — Batch B (Comments, Checklists, Attachments)

### Task 37: Comments — API Client Methods

**Files:**
- Create: `internal/trello/comments.go`
- Create: `internal/trello/comments_test.go`

Implementation:
```go
func (c *Client) ListComments(ctx context.Context, cardID string) ([]Comment, error) {
	var comments []Comment
	err := c.Get(ctx, fmt.Sprintf("/1/cards/%s/actions", cardID),
		map[string]string{"filter": "commentCard"}, &comments)
	return comments, err
}

func (c *Client) AddComment(ctx context.Context, cardID, text string) (Comment, error) {
	var comment Comment
	err := c.Post(ctx, fmt.Sprintf("/1/cards/%s/actions/comments", cardID),
		map[string]string{"text": text}, &comment)
	return comment, err
}

func (c *Client) UpdateComment(ctx context.Context, actionID, text string) (Comment, error) {
	var comment Comment
	err := c.Put(ctx, fmt.Sprintf("/1/actions/%s/text", actionID),
		map[string]string{"value": text}, &comment)
	return comment, err
}

func (c *Client) DeleteComment(ctx context.Context, actionID string) error {
	return c.Delete(ctx, fmt.Sprintf("/1/actions/%s", actionID), nil)
}
```

- [ ] **Step 1-4: Write tests, implement, verify**
- [ ] **Step 5: Commit**

```bash
git add internal/trello/comments.go internal/trello/comments_test.go internal/trello/client.go
git commit -m "feat: implement comment API client methods"
```

---

### Task 38: Comments — Cobra Commands

**Files:**
- Create: `cmd/trello/comments.go`
- Create: `cmd/trello/comments_test.go`

Commands:
- `comments list --card <id>` — requires `--card`
- `comments add --card <id> --text <text>` — requires both
- `comments update --action <id> --text <text>` — requires both
- `comments delete --action <id>` — requires `--action`

Response shapes for void operations:
- `comments delete` returns `DeleteResult{Deleted: true, ID: actionID}` in the success envelope, same pattern as `cards delete`

- [ ] **Step 1-4: Write tests, implement, verify**
- [ ] **Step 5: Commit**

```bash
git add cmd/trello/comments.go cmd/trello/comments_test.go
git commit -m "feat: add comments commands (list, add, update, delete)"
```

---

### Task 39: Checklists — API Client Methods

**Files:**
- Create: `internal/trello/checklists.go`
- Create: `internal/trello/checklists_test.go`

Implementation:
```go
func (c *Client) ListChecklists(ctx context.Context, cardID string) ([]Checklist, error) {
	var checklists []Checklist
	err := c.Get(ctx, fmt.Sprintf("/1/cards/%s/checklists", cardID), nil, &checklists)
	return checklists, err
}

func (c *Client) CreateChecklist(ctx context.Context, cardID, name string) (Checklist, error) {
	var checklist Checklist
	err := c.Post(ctx, "/1/checklists", map[string]string{
		"idCard": cardID,
		"name":   name,
	}, &checklist)
	return checklist, err
}

func (c *Client) DeleteChecklist(ctx context.Context, checklistID string) error {
	return c.Delete(ctx, fmt.Sprintf("/1/checklists/%s", checklistID), nil)
}

func (c *Client) AddCheckItem(ctx context.Context, checklistID, name string) (CheckItem, error) {
	var item CheckItem
	err := c.Post(ctx, fmt.Sprintf("/1/checklists/%s/checkItems", checklistID),
		map[string]string{"name": name}, &item)
	return item, err
}

func (c *Client) UpdateCheckItem(ctx context.Context, cardID, itemID, state string) (CheckItem, error) {
	var item CheckItem
	err := c.Put(ctx, fmt.Sprintf("/1/cards/%s/checkItem/%s", cardID, itemID),
		map[string]string{"state": state}, &item)
	return item, err
}

func (c *Client) DeleteCheckItem(ctx context.Context, checklistID, itemID string) error {
	return c.Delete(ctx, fmt.Sprintf("/1/checklists/%s/checkItems/%s", checklistID, itemID), nil)
}
```

- [ ] **Step 1-4: Write tests, implement, verify**
- [ ] **Step 5: Commit**

```bash
git add internal/trello/checklists.go internal/trello/checklists_test.go internal/trello/client.go
git commit -m "feat: implement checklist and check item API client methods"
```

---

### Task 40: Checklists — Cobra Commands (including nested items subgroup)

**Files:**
- Create: `cmd/trello/checklists.go`
- Create: `cmd/trello/checklists_test.go`

Commands:
- `checklists list --card <id>`
- `checklists create --card <id> --name <name>`
- `checklists delete --checklist <id>`
- `checklists items add --checklist <id> --name <name>`
- `checklists items update --card <id> --item <id> --state <complete|incomplete>`
- `checklists items delete --checklist <id> --item <id>`

Response shapes for void operations:
- `checklists delete` returns `DeleteResult{Deleted: true, ID: checklistID}`
- `checklists items delete` returns `DeleteResult{Deleted: true, ID: itemID}`

Key validation: `--state` must be `complete` or `incomplete` — use a custom validator:
```go
func validateState(state string) error {
	if state != "complete" && state != "incomplete" {
		return contract.NewError(contract.ValidationError,
			"--state must be 'complete' or 'incomplete'")
	}
	return nil
}
```

Nested command structure:
```go
var checklistsItemsCmd = &cobra.Command{
	Use:   "items",
	Short: "Manage checklist items",
}
// checklistsCmd.AddCommand(checklistsItemsCmd)
// checklistsItemsCmd.AddCommand(checklistsItemsAddCmd, ...)
```

- [ ] **Step 1-4: Write tests, implement, verify**
- [ ] **Step 5: Commit**

```bash
git add cmd/trello/checklists.go cmd/trello/checklists_test.go
git commit -m "feat: add checklists commands with nested items subgroup"
```

---

### Task 41: Attachments — API Client Methods

**Files:**
- Create: `internal/trello/attachments.go`
- Create: `internal/trello/attachments_test.go`

Implementation (including multipart for file upload):
```go
func (c *Client) ListAttachments(ctx context.Context, cardID string) ([]Attachment, error) {
	var attachments []Attachment
	err := c.Get(ctx, fmt.Sprintf("/1/cards/%s/attachments", cardID), nil, &attachments)
	return attachments, err
}

func (c *Client) AddURLAttachment(ctx context.Context, cardID, urlStr string, name *string) (Attachment, error) {
	qp := map[string]string{"url": urlStr}
	if name != nil {
		qp["name"] = *name
	}
	var attachment Attachment
	err := c.Post(ctx, fmt.Sprintf("/1/cards/%s/attachments", cardID), qp, &attachment)
	return attachment, err
}

func (c *Client) AddFileAttachment(ctx context.Context, cardID, filePath string, name *string) (Attachment, error) {
	// Multipart file upload — uses PostMultipart
	qp := map[string]string{}
	if name != nil {
		qp["name"] = *name
	}
	var attachment Attachment
	err := c.postMultipartFile(ctx, fmt.Sprintf("/1/cards/%s/attachments", cardID), filePath, qp, &attachment)
	return attachment, err
}

func (c *Client) DeleteAttachment(ctx context.Context, cardID, attachmentID string) error {
	return c.Delete(ctx, fmt.Sprintf("/1/cards/%s/attachments/%s", cardID, attachmentID), nil)
}

// postMultipartFile handles multipart/form-data file uploads.
// Uses the centralized buildURL method for auth injection (key/token query params)
// and mapHTTPError for error normalization — same as all other API methods.
func (c *Client) postMultipartFile(ctx context.Context, path, filePath string, params map[string]string, result any) error {
	file, err := os.Open(filePath)
	if err != nil {
		return contract.NewError(contract.FileNotFound, fmt.Sprintf("cannot open file: %s", filePath))
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return contract.NewError(contract.UnknownError, fmt.Sprintf("failed to create form file: %v", err))
	}
	if _, err := io.Copy(part, file); err != nil {
		return contract.NewError(contract.UnknownError, fmt.Sprintf("failed to read file: %v", err))
	}

	for k, v := range params {
		writer.WriteField(k, v)
	}
	writer.Close()

	// Use buildURL for centralized auth injection — same as Do()
	fullURL := c.buildURL(path, nil)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return contract.NewError(contract.HTTPError, fmt.Sprintf("upload failed: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return mapHTTPError(resp)
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

// NOTE: The Client must expose a buildURL method that constructs the full URL
// with auth query params (key, token). Extract this from the Do() method
// so both Do() and postMultipartFile() share the same auth injection logic.
```

Additional imports needed: `"bytes"`, `"io"`, `"mime/multipart"`, `"os"`, `"path/filepath"`.

- [ ] **Step 1-4: Write tests, implement, verify**

For `AddFileAttachment` test, create a temp file and verify the server receives multipart data.

- [ ] **Step 5: Commit**

```bash
git add internal/trello/attachments.go internal/trello/attachments_test.go internal/trello/client.go
git commit -m "feat: implement attachment API client methods with multipart upload"
```

---

### Task 42: Attachments — Cobra Commands

**Files:**
- Create: `cmd/trello/attachments.go`
- Create: `cmd/trello/attachments_test.go`

Commands:
- `attachments list --card <id>`
- `attachments add-file --card <id> --path <path> [--name <name>]` — validates file exists via `ValidateFilePath`
- `attachments add-url --card <id> --url <url> [--name <name>]` — validates URL via `ValidateURL`
- `attachments delete --card <id> --attachment <id>`

Response shapes for void operations:
- `attachments delete` returns `DeleteResult{Deleted: true, ID: attachmentID}`

- [ ] **Step 1-4: Write tests, implement, verify**
- [ ] **Step 5: Commit**

```bash
git add cmd/trello/attachments.go cmd/trello/attachments_test.go
git commit -m "feat: add attachments commands (list, add-file, add-url, delete)"
```

---

### Task 43: Run Full Test Suite and Commit Batch B

- [ ] **Step 1: Run full test suite**

Run: `go test ./... -v`
Expected: all tests PASS

- [ ] **Step 2: Commit cleanup**

```bash
git add -A
git commit -m "chore: complete Batch B — comments, checklists, attachments commands"
```

---

**End of Chunk 6**

---

## Chunk 7: Resource Commands — Batch C (Labels, Members, Search)

### Task 44: Labels — API Client Methods + Cobra Commands

**Files:**
- Create: `internal/trello/labels.go`, `internal/trello/labels_test.go`
- Create: `cmd/trello/labels.go`, `cmd/trello/labels_test.go`

API methods:
```go
func (c *Client) ListLabels(ctx context.Context, boardID string) ([]Label, error) {
	var labels []Label
	err := c.Get(ctx, fmt.Sprintf("/1/boards/%s/labels", boardID), nil, &labels)
	return labels, err
}

func (c *Client) CreateLabel(ctx context.Context, boardID, name, color string) (Label, error) {
	var label Label
	err := c.Post(ctx, "/1/labels", map[string]string{
		"idBoard": boardID, "name": name, "color": color,
	}, &label)
	return label, err
}

func (c *Client) AddLabelToCard(ctx context.Context, cardID, labelID string) error {
	return c.Post(ctx, fmt.Sprintf("/1/cards/%s/idLabels", cardID),
		map[string]string{"value": labelID}, nil)
}

func (c *Client) RemoveLabelFromCard(ctx context.Context, cardID, labelID string) error {
	return c.Delete(ctx, fmt.Sprintf("/1/cards/%s/idLabels/%s", cardID, labelID), nil)
}
```

Commands:
- `labels list --board <id>`
- `labels create --board <id> --name <name> --color <color>`
- `labels add --card <id> --label <id>`
- `labels remove --card <id> --label <id>`

Response shapes for void operations:
- `labels add` returns `DeleteResult{Deleted: false, ID: labelID}` — reuse the struct but rename conceptually. Alternatively, define an `ActionResult` type:
```go
// ActionResult is the response shape for void add/remove operations.
type ActionResult struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
```
Add `ActionResult` to `internal/trello/types.go` alongside `DeleteResult`. Use `ActionResult` for `labels add`, `labels remove`, `members add`, `members remove`. Use `DeleteResult` for all delete operations.

- [ ] **Step 1-4: Write tests, implement, verify**
- [ ] **Step 5: Commit**

```bash
git add internal/trello/labels.go internal/trello/labels_test.go cmd/trello/labels.go cmd/trello/labels_test.go internal/trello/client.go
git commit -m "feat: add labels commands (list, create, add, remove)"
```

---

### Task 45: Members — API Client Methods + Cobra Commands

**Files:**
- Create: `internal/trello/members.go`, `internal/trello/members_test.go`
- Create: `cmd/trello/members.go`, `cmd/trello/members_test.go`

API methods:
```go
func (c *Client) ListMembers(ctx context.Context, boardID string) ([]Member, error) {
	var members []Member
	err := c.Get(ctx, fmt.Sprintf("/1/boards/%s/members", boardID), nil, &members)
	return members, err
}

func (c *Client) AddMemberToCard(ctx context.Context, cardID, memberID string) error {
	return c.Post(ctx, fmt.Sprintf("/1/cards/%s/idMembers", cardID),
		map[string]string{"value": memberID}, nil)
}

func (c *Client) RemoveMemberFromCard(ctx context.Context, cardID, memberID string) error {
	return c.Delete(ctx, fmt.Sprintf("/1/cards/%s/idMembers/%s", cardID, memberID), nil)
}

func (c *Client) GetMe(ctx context.Context) (Member, error) {
	var member Member
	err := c.Get(ctx, "/1/members/me", nil, &member)
	return member, err
}
```

Commands:
- `members list --board <id>`
- `members add --card <id> --member <id>`
- `members remove --card <id> --member <id>`

Response shapes for void operations:
- `members add` returns `ActionResult{Success: true, ID: memberID}`
- `members remove` returns `ActionResult{Success: true, ID: memberID}`

All command handlers follow the single-auth-check pattern: call `auth.RequireAuth` once, pass `creds` to `getAPIClient(creds)`.

- [ ] **Step 1-4: Write tests, implement, verify**
- [ ] **Step 5: Commit**

```bash
git add internal/trello/members.go internal/trello/members_test.go cmd/trello/members.go cmd/trello/members_test.go internal/trello/client.go
git commit -m "feat: add members commands (list, add, remove) and GetMe"
```

---

### Task 46: Search — API Client Methods + Cobra Commands

**Files:**
- Create: `internal/trello/search.go`, `internal/trello/search_test.go`
- Create: `cmd/trello/search.go`, `cmd/trello/search_test.go`

API methods:
```go
func (c *Client) SearchCards(ctx context.Context, query string) (CardSearchResult, error) {
	var raw struct {
		Cards []Card `json:"cards"`
	}
	err := c.Get(ctx, "/1/search", map[string]string{
		"query":       query,
		"modelTypes":  "cards",
	}, &raw)
	return CardSearchResult{Query: query, Cards: raw.Cards}, err
}

func (c *Client) SearchBoards(ctx context.Context, query string) (BoardSearchResult, error) {
	var raw struct {
		Boards []Board `json:"boards"`
	}
	err := c.Get(ctx, "/1/search", map[string]string{
		"query":       query,
		"modelTypes":  "boards",
	}, &raw)
	return BoardSearchResult{Query: query, Boards: raw.Boards}, err
}
```

Commands:
- `search cards --query <text>` — returns `{query, cards}` wrapper
- `search boards --query <text>` — returns `{query, boards}` wrapper

- [ ] **Step 1-4: Write tests, implement, verify**
- [ ] **Step 5: Commit**

```bash
git add internal/trello/search.go internal/trello/search_test.go cmd/trello/search.go cmd/trello/search_test.go internal/trello/client.go
git commit -m "feat: add search commands (cards, boards) with query wrapper response"
```

---

### Task 47: Remove All Remaining Stubs + Final Verification

- [ ] **Step 1: Remove all remaining stub methods from client.go**

Verify no stub methods remain — all should be replaced by implementations in resource-specific files.

- [ ] **Step 2: Run full test suite**

Run: `go test ./... -v -count=1`
Expected: all tests PASS

- [ ] **Step 3: Verify compile-time interface check**

Run: `go build ./...`
Expected: successful compilation (the `var _ API = (*Client)(nil)` check passes)

- [ ] **Step 4: Test CLI end-to-end**

Run:
```bash
go run ./cmd/trello/ version
go run ./cmd/trello/ auth set --api-key test --token test
go run ./cmd/trello/ auth status
go run ./cmd/trello/ auth clear
```
Expected: all return valid JSON envelopes

- [ ] **Step 5: Final commit**

```bash
git add -A
git commit -m "chore: complete V1 implementation — all 42 commands with TDD"
```

---

**End of Chunk 7**

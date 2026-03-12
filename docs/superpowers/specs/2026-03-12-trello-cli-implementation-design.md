# Trello CLI Implementation Design

**Date:** 2026-03-12
**Status:** Approved
**Approach:** TDD, phased by layer

---

## Overview

Full TDD implementation of a cross-platform Go CLI that implements the Trello command contract defined in `docs/2026-03-12-trello-cli-command-spec-design.md`, using the technical stack defined in `docs/Stack.md`.

**Key decisions:**
- Go 1.23+ (latest stable)
- Phased by layer: foundation first, then resource commands in batches
- Test layers: `httptest` for API client, interface mocks for command handlers
- Auth login: local callback server with copy-paste fallback

---

## Phase Structure

### Phase 1 — Project Skeleton & Contract Layer

- `go mod init`, Cobra root command, `--pretty` and `--verbose` global flags
- `internal/contract`: success/error envelope builders, error code constants, validation helpers
- `trello version` command (proves Cobra → contract → JSON stdout pipeline end-to-end)
- Tests: envelope construction, error code mapping, validation helpers, version command output

### Phase 2 — Config, Credentials & Auth

- `internal/config`: Viper setup with precedence (flags → env → file → defaults)
- `internal/credentials`: keyring-backed store with interface abstraction (for test mocking + headless fallback)
- `internal/auth`: login flow (local callback server + copy-paste fallback), set, status, clear
- `cmd/trello/auth.go`: all 4 auth subcommands wired through Cobra
- Tests: config precedence, credential store/load/clear (mocked keyring), auth state validation, auth commands with mocked credential store

### Phase 3 — API Client Layer

- `internal/trello/client.go`: base HTTP client with auth injection, timeout, context
- `internal/trello/retry.go`: exponential backoff with jitter, rate-limit detection
- `internal/trello/errors.go`: Trello HTTP → contract error code normalization
- Resource method stubs that Phase 4 fills in
- Tests: request construction, auth header injection, retry behavior, error normalization, rate-limit handling — all via `httptest`

### Phase 4 — Resource Commands (3 batches)

- **Batch A**: Boards, Lists, Cards (core workflow)
- **Batch B**: Comments, Checklists, Attachments (card enrichment)
- **Batch C**: Labels, Members, Search (cross-cutting)

---

## Foundation Layer Design

### `internal/contract` — Response & Validation

```
response.go   — envelope builders
errors.go     — error codes and error type
validation.go — reusable flag/arg validators
```

**Envelope builders:**
- `Success(data any) []byte` — marshals `{"ok": true, "data": ...}`
- `Error(code string, message string) []byte` — marshals `{"ok": false, "error": {"code": ..., "message": ...}}`
- `Render(w io.Writer, envelope []byte, pretty bool)` — writes to stdout, optionally pretty-printed

**Error type:**
- `ContractError` struct with `Code` and `Message` fields
- Implements `error` interface
- Standard code constants: `AUTH_REQUIRED`, `AUTH_INVALID`, `NOT_FOUND`, `VALIDATION_ERROR`, `CONFLICT`, `RATE_LIMITED`, `HTTP_ERROR`, `FILE_NOT_FOUND`, `UNSUPPORTED`, `UNKNOWN_ERROR`

**Validators:**
- `RequireFlag(name, value string) error` — returns `VALIDATION_ERROR` if empty
- `RequireExactlyOne(flags map[string]string) error` — for mutually exclusive flags
- `RequireAtLeastOne(flags map[string]string) error` — for update commands needing at least one mutation flag
- `ValidateISO8601(value string) error` — for `--due` flags
- `ValidateURL(value string) error` — for `--url` flags
- `ValidateFilePath(path string) error` — for `--path` flags, returns `FILE_NOT_FOUND` if missing

### `internal/config` — Viper Setup

```
config.go — initialization, precedence, defaults
```

- Config file location: `~/.config/trello-cli/config.yaml`
- Viper-managed keys: `profile`, `pretty`, `timeout`, `max_retries`, `retry_mutations`, `verbose`
- Precedence: explicit flags → env vars (`TRELLO_*`) → config file → defaults
- Defaults: `pretty: false`, `timeout: 15s`, `max_retries: 3`, `retry_mutations: false`, `verbose: false`, `profile: default`
- No secrets touch this layer

### `internal/credentials` — Secret Store

```
store.go — interface + keyring implementation + env fallback
```

**Interface:**
```go
type Store interface {
    Get(profile string) (Credentials, error)
    Set(profile string, creds Credentials) error
    Delete(profile string) error
}
```

**`Credentials` struct:** `APIKey string`, `Token string`, `AuthMode string`

**Implementations:**
- `KeyringStore` — backed by `go-keyring`, service name `trello-cli/<profile>`
- `EnvStore` — reads `TRELLO_API_KEY` and `TRELLO_TOKEN` from environment (headless/CI fallback)
- Selection: try keyring first, fall back to env if keyring unavailable, never fall back to plaintext file

### `internal/auth` — Auth Subsystem

```
login.go  — interactive browser flow
set.go    — manual credential entry
status.go — auth state check
clear.go  — credential removal
```

**Login flow:**
1. Start local HTTP server on random available port
2. Build Trello `/1/authorize` URL with redirect to `localhost:<port>`
3. Open browser via `os/exec` (`open` on macOS, `xdg-open` on Linux, `start` on Windows)
4. Wait for callback with token (with timeout)
5. If server can't bind or times out → fall back to copy-paste: print the authorize URL, prompt user to paste token
6. Validate credentials via `GET /members/me`
7. Store validated credentials
8. Return success envelope with member info

**Auth guard:** A `RequireAuth` helper that commands call — loads credentials from store, returns `AUTH_REQUIRED` if missing. Single centralized auth check.

### `cmd/trello` — Cobra Wiring

```
main.go    — entry point
root.go    — root command, global flags, output rendering
version.go — version command
auth.go    — auth command group
```

**Root command responsibilities:**
- Bind `--pretty` and `--verbose` global persistent flags
- Initialize config via Viper
- Set up `slog` logger to stderr gated on `--verbose`
- Handle output rendering: commands return data via `output()` helper, errors via `ContractError` returns

**Output flow:** Commands don't write to stdout directly. They call `output(cmd, data)` for success or return `ContractError` for failures. Root's error handler renders error envelopes. Non-contract errors get wrapped as `UNKNOWN_ERROR`.

---

## TDD Strategy & Test Architecture

### Test-First Discipline

Every unit of work follows: write failing test → write minimum code to pass → refactor. No implementation code lands without a test that preceded it.

### Test Layers

**Layer 1 — Contract tests** (`internal/contract/*_test.go`)
- Envelope builders produce valid JSON with correct shape
- Error codes are the exact string constants from the spec
- Validators return correct error codes for each failure case
- Validators return `nil` for valid inputs
- Pretty-print output is valid indented JSON

**Layer 2 — Config tests** (`internal/config/*_test.go`)
- Default values are correct
- Env vars override config file values
- Flags override env vars
- Missing config file doesn't error
- Invalid config file produces a useful error

**Layer 3 — Credential store tests** (`internal/credentials/*_test.go`)
- Uses a mock `Store` interface implementation (not real keyring in tests)
- `Set` → `Get` round-trips correctly
- `Get` on missing profile returns typed error mappable to `AUTH_REQUIRED`
- `Delete` removes credentials
- `EnvStore` reads from environment variables
- Auth mode is stored and retrieved correctly

**Layer 4 — Auth subsystem tests** (`internal/auth/*_test.go`)
- Login flow validates credentials via mocked HTTP
- Login stores credentials on success
- Login returns `AUTH_INVALID` when Trello rejects the token
- `set` validates both flags present
- `status` returns correct shape for configured vs unconfigured states
- `clear` removes credentials and returns `configured: false`
- `RequireAuth` returns `AUTH_REQUIRED` when store is empty
- `RequireAuth` returns loaded credentials when present

**Layer 5 — API client tests** (`internal/trello/*_test.go`)
- All tests use `httptest.NewServer` with canned Trello responses
- Request construction: correct URL paths, query params, auth credentials injected
- Successful responses decoded into correct Go types
- Trello 401 → `AUTH_INVALID`
- Trello 404 → `NOT_FOUND`
- Trello 429 → retry with backoff, then `RATE_LIMITED` if exhausted
- Trello 500 → `HTTP_ERROR`
- Network timeout → `HTTP_ERROR`
- Context cancellation aborts request
- Multipart upload constructs correct `Content-Type` and body
- Retry logic: respects `max_retries`, uses jitter, doesn't retry mutations unless `retry_mutations` is true

**Layer 6 — Command handler tests** (`cmd/trello/*_test.go`)
- Each command tested via Cobra's `Execute()` with captured stdout
- API layer replaced with mock interface
- Tests verify: required flag validation, mutually exclusive flags, correct API method called, success envelope shape, error passthrough, auth guard fires before API call

### Test Helpers

`internal/testutil` package:
- `ParseEnvelope(t, output) map[string]any` — parses JSON output, fails test if invalid
- `AssertSuccess(t, output)` — asserts `ok: true`
- `AssertError(t, output, code)` — asserts `ok: false` with expected error code
- `MockCredentialStore` — in-memory `Store` implementation
- `NewTrelloServer(t, handlers) *httptest.Server` — configurable mock Trello API

### Coverage Expectations

- `internal/contract`: near 100%
- `internal/credentials`: near 100% on interface behavior
- `internal/trello`: high coverage on request construction, error mapping, retry logic
- `cmd/trello`: every command has at least: valid success case, missing required flags case, auth guard case

---

## Resource Command Implementation Pattern

### API Client Interface

```go
type Client interface {
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
    AddURLAttachment(ctx context.Context, cardID, url string, name *string) (Attachment, error)
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
    SearchCards(ctx context.Context, query string) ([]Card, error)
    SearchBoards(ctx context.Context, query string) ([]Board, error)

    // Auth validation
    GetMe(ctx context.Context) (Member, error)
}
```

### Command Handler Pattern

Every command handler follows the same structure:

1. Parse and validate flags (using contract validators)
2. Load auth (using `RequireAuth`)
3. Call API client method
4. Return data for envelope rendering

Commands are thin — no HTTP logic, no JSON marshaling, no direct stdout writes.

### Cobra Wiring Pattern

Each resource group gets one file in `cmd/trello/` that defines the parent command, each subcommand, and an `init()` that registers subcommands and flags.

**Output flow:**
```
RunE → validate flags → requireAuth → apiClient.Method() → output(cmd, data)
     ↘ on error → return ContractError → root error handler → error envelope
```

### Data Types

Resource structs in `internal/trello/types.go` mirror the spec's response fields. One set of types, used end-to-end from API client return values through to envelope construction.

### Adding a New Command (Repeatable Recipe)

1. **Test:** write command test — mock API client, assert correct envelope output
2. **Command:** write Cobra command handler
3. **Test:** write API client test — `httptest` with canned Trello response
4. **Client:** write API client method
5. **Run all tests** — verify green

---

## Resource Command Batches

### Batch A — Core Workflow (Boards, Lists, Cards)

**Build order:**

1. **Boards** (2 commands) — simplest resource, read-only, proves full pipeline
   - `boards list` — GET, returns array
   - `boards get` — GET with ID, returns single object

2. **Lists** (5 commands) — introduces mutations, position handling
   - `lists list`, `lists create`, `lists update`, `lists archive`, `lists move`

3. **Cards** (7 commands) — most complex resource, all patterns
   - `cards list` (mutually exclusive flags), `cards get`, `cards create` (date validation), `cards update` (comma-separated IDs), `cards move`, `cards archive`, `cards delete` (custom response shape)

### Batch B — Card Enrichment (Comments, Checklists, Attachments)

**Build order:**

1. **Comments** (4 commands) — introduces Trello actions model
   - `comments list`, `comments add`, `comments update`, `comments delete`

2. **Checklists** (6 commands) — introduces nested subcommands
   - `checklists list`, `checklists create`, `checklists delete`
   - `checklists items add`, `checklists items update` (enum validation), `checklists items delete`

3. **Attachments** (4 commands) — introduces multipart upload
   - `attachments list`, `attachments add-url`, `attachments add-file` (multipart), `attachments delete`

### Batch C — Cross-Cutting (Labels, Members, Search)

**Build order:**

1. **Labels** (4 commands) — board-level + card-level operations
   - `labels list`, `labels create`, `labels add`, `labels remove`

2. **Members** (3 commands) — same card-level add/remove pattern
   - `members list`, `members add`, `members remove`

3. **Search** (2 commands) — different response shape
   - `search cards`, `search boards`

### Total Command Count

| Group | Commands |
|-------|----------|
| Auth | 4 |
| Boards | 2 |
| Lists | 5 |
| Cards | 7 |
| Comments | 4 |
| Checklists | 6 |
| Attachments | 4 |
| Labels | 4 |
| Members | 3 |
| Search | 2 |
| Version | 1 |
| **Total** | **42** |

---

## Package Layout

```
cmd/
  trello/
    main.go
    root.go
    version.go
    auth.go
    boards.go
    lists.go
    cards.go
    comments.go
    checklists.go
    attachments.go
    labels.go
    members.go
    search.go

internal/
  contract/
    response.go
    errors.go
    validation.go
  config/
    config.go
  credentials/
    store.go
  auth/
    login.go
    set.go
    status.go
    clear.go
  trello/
    client.go
    types.go
    errors.go
    retry.go
    boards.go
    lists.go
    cards.go
    comments.go
    checklists.go
    attachments.go
    labels.go
    members.go
    search.go
  testutil/
    helpers.go

go.mod
go.sum
```

---

## Dependencies

**Required:**
- `github.com/spf13/cobra`
- `github.com/spf13/viper`
- `github.com/zalando/go-keyring`

**Standard library:**
- `net/http`, `encoding/json`, `mime/multipart`, `context`, `time`, `os`, `fmt`, `errors`, `path/filepath`, `log/slog`, `net`, `net/url`, `os/exec`, `runtime`

**Test only:**
- `net/http/httptest`, `testing`, `bytes`, `strings`, `os` (for env manipulation)

# Auth Set-Key Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `trello auth set-key` so the CLI can persist an API key in the credential store without requiring a token up front.

**Architecture:** Extend the existing auth subsystem instead of introducing a second storage path. `set-key` will update the current credential record, preserve any existing token, and introduce a `key_only` partial-auth state that `auth login` can consume but protected commands still reject.

**Tech Stack:** Go, Cobra, existing `internal/auth`, `internal/credentials`, `internal/contract` packages

---

## Chunk 1: Auth Partial-Credential Behavior

### Task 1: Add failing auth package tests for key-only state

**Files:**
- Modify: `internal/auth/auth_test.go`
- Modify: `internal/auth/status_test.go`
- Modify: `internal/auth/set_test.go`

- [ ] **Step 1: Write failing tests showing `RequireAuth` rejects `key_only` credentials**
- [ ] **Step 2: Write failing tests showing `Status` returns `configured: false` for `key_only` credentials**
- [ ] **Step 3: Write failing tests showing `SetKey` stores key-only credentials and preserves an existing token**
- [ ] **Step 4: Run `go test ./internal/auth` and verify the new tests fail for the expected reasons**
- [ ] **Step 5: Commit the failing tests**

### Task 2: Implement key-only auth behavior

**Files:**
- Modify: `internal/auth/auth.go`
- Modify: `internal/auth/status.go`
- Modify: `internal/auth/set.go`

- [ ] **Step 1: Add a `SetKey` result type and implementation in `internal/auth/set.go`**
- [ ] **Step 2: Update `RequireAuth` to treat missing token or `key_only` state as `AUTH_REQUIRED`**
- [ ] **Step 3: Update `Status` to short-circuit to `configured: false` for key-only stored credentials**
- [ ] **Step 4: Run `go test ./internal/auth` and verify the package passes**
- [ ] **Step 5: Commit the auth package implementation**

## Chunk 2: CLI Command Wiring

### Task 3: Add failing CLI tests for `trello auth set-key`

**Files:**
- Modify: `cmd/trello/auth_test.go`

- [ ] **Step 1: Write a failing test for `trello auth set-key --api-key ...` returning a JSON success envelope**
- [ ] **Step 2: Write a failing test for missing `--api-key` validation**
- [ ] **Step 3: Run `go test ./cmd/trello` and verify the tests fail because the command does not exist yet**
- [ ] **Step 4: Commit the failing CLI tests**

### Task 4: Implement the Cobra command

**Files:**
- Modify: `cmd/trello/auth.go`

- [ ] **Step 1: Add `auth set-key` command wiring and `--api-key` flag**
- [ ] **Step 2: Route the command through the new `auth.SetKey` helper**
- [ ] **Step 3: Run `go test ./cmd/trello` and verify the command tests pass**
- [ ] **Step 4: Commit the command implementation**

## Chunk 3: Verification

### Task 5: End-to-end verification

**Files:**
- Modify: `cmd/trello/auth.go`
- Modify: `cmd/trello/auth_test.go`
- Modify: `internal/auth/auth.go`
- Modify: `internal/auth/status.go`
- Modify: `internal/auth/set.go`
- Modify: `internal/auth/auth_test.go`
- Modify: `internal/auth/status_test.go`
- Modify: `internal/auth/set_test.go`

- [ ] **Step 1: Run `go test ./internal/auth ./cmd/trello`**
- [ ] **Step 2: Run `go test ./...`**
- [ ] **Step 3: Run `go build -o bin/trello ./cmd/trello`**
- [ ] **Step 4: Manually smoke-test `./bin/trello auth set-key --api-key test-key`**
- [ ] **Step 5: Commit the finished feature**

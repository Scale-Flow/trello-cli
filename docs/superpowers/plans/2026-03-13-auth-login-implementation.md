# Auth Login Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement a working `trello auth login` flow that opens Trello authorization in the browser, captures the returned token locally, validates it, and stores usable credentials for later commands.

**Architecture:** Keep the browser/callback flow in `internal/auth` so the Cobra command remains thin. Reuse the existing `BuildAuthorizeURL` and `CompleteLogin` helpers, add a small local callback server plus browser-launch abstraction, and source the Trello API key from an existing credential or `TRELLO_API_KEY`.

**Tech Stack:** Go, Cobra, net/http, context, os/exec, existing contract/auth/credentials packages

---

## Chunk 1: Auth Flow Wiring

### Task 1: Add failing login command tests

**Files:**
- Modify: `cmd/trello/auth_test.go`
- Test: `cmd/trello/auth_test.go`

- [ ] **Step 1: Write failing tests for successful interactive login and missing API key**
- [ ] **Step 2: Run the focused auth command tests and verify they fail for the expected reason**
- [ ] **Step 3: Commit the failing tests**

### Task 2: Implement callback capture and browser launch

**Files:**
- Modify: `internal/auth/login.go`
- Test: `internal/auth/login_test.go`

- [ ] **Step 1: Write failing tests for local callback token capture and API-key sourcing behavior**
- [ ] **Step 2: Run the focused auth package tests and verify they fail**
- [ ] **Step 3: Implement the minimal local server, callback HTML/token handoff, and browser-launch hooks**
- [ ] **Step 4: Re-run the focused auth package tests and verify they pass**
- [ ] **Step 5: Commit the auth package implementation**

### Task 3: Wire `trello auth login` to the auth package

**Files:**
- Modify: `cmd/trello/auth.go`
- Modify: `cmd/trello/auth_test.go`

- [ ] **Step 1: Replace the `UNSUPPORTED` stub with the real login execution path**
- [ ] **Step 2: Re-run the focused command tests and verify they pass**
- [ ] **Step 3: Commit the command wiring**

## Chunk 2: Verification

### Task 4: Full verification

**Files:**
- Modify: `cmd/trello/auth.go`
- Modify: `cmd/trello/auth_test.go`
- Modify: `internal/auth/login.go`
- Modify: `internal/auth/login_test.go`

- [ ] **Step 1: Run `go test ./internal/auth ./cmd/trello`**
- [ ] **Step 2: Run `go test ./...`**
- [ ] **Step 3: Manually exercise `go build -o bin/trello ./cmd/trello`**
- [ ] **Step 4: Commit the finished feature**

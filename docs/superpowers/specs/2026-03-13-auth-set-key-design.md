# Auth Set-Key Design

**Date:** 2026-03-13
**Status:** Proposed
**Topic:** Add `trello auth set-key` for storing the Trello API key in the credential store without requiring a token up front.

---

## Overview

Add a new auth command:

```bash
trello auth set-key --api-key <key>
```

This command stores the Trello API key in the same credential record already used by the CLI's keyring-backed auth storage.

The command is meant to support two workflows:

1. Key-first setup for interactive login
2. API key rotation without forcing token re-entry

The chosen behavior is:

- preserve the existing stored token when one already exists for the profile
- store only the API key when no token exists yet
- keep protected Trello commands unauthenticated until both key and token are available

---

## User Experience

### Command

```bash
trello auth set-key --api-key <key>
```

### Intended outcomes

- If the profile already has a stored token, the command updates only the API key and leaves the token intact.
- If the profile has no stored credentials, the command stores a partial credential record containing the API key.
- After storing a key-only record, `trello auth login` can use that stored key without requiring `TRELLO_API_KEY` in the environment.

### Validation

- `--api-key` is required
- missing `--api-key` returns `VALIDATION_ERROR`

### Success shape

```json
{
  "ok": true,
  "data": {
    "configured": false,
    "authMode": "key_only"
  }
}
```

If a full credential set already exists and the token is preserved, the command may return:

```json
{
  "ok": true,
  "data": {
    "configured": true,
    "authMode": "manual"
  }
}
```

This reflects whether the resulting stored record is fully usable for authenticated Trello access.

---

## Data Model

The existing credential record remains the storage unit:

- `apiKey`
- `token`
- `authMode`

No second keyring entry is introduced.

### Auth modes

- `manual` for API key + token stored manually
- `interactive` for API key + token produced by browser login
- `env` for environment-backed fallback reads
- `key_only` for a stored API key with no token yet

`key_only` is a new partial state used only for stored credentials.

---

## Behavior Changes

### `trello auth set-key`

- Load any existing stored credential record for the active profile
- Replace `apiKey`
- Preserve `token` if present
- Preserve `authMode` if a token exists and the record was already complete
- Otherwise write `authMode: "key_only"`

### `trello auth login`

- Continue resolving the API key in this order:
1. explicit argument, if ever added later
2. stored credentials for the profile
3. `TRELLO_API_KEY`

- A stored `key_only` record is valid as an API-key source for interactive login

### `trello auth status`

- If stored credentials contain both API key and token, behave exactly as today
- If credentials are `key_only`, report not configured and do not call Trello

### Protected commands

- Commands that require auth must continue to fail with `AUTH_REQUIRED` unless both API key and token are present
- A `key_only` record is not sufficient authentication

---

## Testing

Add coverage for:

- `auth set-key` stores a key-only record when nothing exists
- `auth set-key` preserves an existing token
- `auth set-key` validates `--api-key`
- `auth login` can source the API key from a key-only stored record
- `auth status` reports not configured for key-only records
- `RequireAuth` treats key-only stored credentials as `AUTH_REQUIRED`

---

## Recommendation

Implement `set-key` as a small extension of the current auth subsystem rather than adding a new storage mechanism. This keeps key rotation simple, keeps interactive login usable without env vars, and avoids duplicating credential state across multiple backends.

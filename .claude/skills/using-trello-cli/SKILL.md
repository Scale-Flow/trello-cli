---
name: using-trello-cli
description: Operate Trello through the local CLI tool. Use this skill whenever the user asks to create, update, move, archive, search, or inspect Trello boards, lists, cards, comments, checklists, attachments, labels, members, or custom fields — or when they mention the Trello CLI, `trello` command, or want to manage Trello from the terminal. Also use this when working in the Trello_CLI repository and needing to exercise the binary.
---

This skill teaches you how to translate natural-language Trello requests into safe, multi-step `trello` CLI workflows. The CLI outputs deterministic JSON — every response is machine-parseable and follows a stable contract.

## Resolve the Binary

Try these in order — use whichever works first:

1. `trello` on `PATH`
2. `./bin/trello` (relative to the repo root)
3. Build from the repo: `go build -o bin/trello ./cmd/trello`

Once resolved, use that path for all subsequent commands in the session.

## Preflight: Check Auth

Before any resource command, verify credentials are configured:

```bash
trello auth status
```

The response follows the JSON envelope. If `ok` is `false` or the auth mode is `key_only`, stop and explain the issue to the user. Do not attempt resource commands without valid auth.

Auth can come from two sources:
- **OS keyring** (stored credentials via `trello auth set` or `trello auth login`)
- **Environment variables** `TRELLO_API_KEY` and `TRELLO_TOKEN`

## Core Workflow: Discover, Mutate, Verify

Every Trello task follows this shape:

1. **Discover** — find the IDs you need (boards, lists, cards) using read commands or search
2. **Mutate** — run the minimum command to accomplish the task
3. **Verify** — re-fetch the resource to confirm the change took effect

This matters because the CLI uses IDs, not names. Never guess an ID — always discover it first.

**Example:** "Create a card in the Doing list on the Marketing board"
1. `trello boards list` → find Marketing board ID
2. `trello lists list --board <board-id>` → find Doing list ID
3. `trello cards create --list <list-id> --name "My card"`
4. `trello cards get --card <card-id>` → confirm creation

## JSON Contract

Every command returns one of two envelopes:

```json
{"ok": true, "data": ...}
{"ok": false, "error": {"code": "...", "message": "..."}}
```

Always branch on `ok`. Read payloads from `data`, errors from `error.code` and `error.message`. The `--pretty` flag changes formatting only, not the schema.

Use compact JSON (no `--pretty`) when piping output or extracting values programmatically. Use `--pretty` when showing results to the user.

## Safety Rules

- **Prefer reads first** when names are ambiguous — list boards/lists before mutating
- **Use IDs after discovery** — never pass names where an ID is expected
- **Never guess IDs** — always discover them from a list or search command
- **Re-fetch after mutations** when confirmation matters
- **Validate file paths** before `attachments add-file`
- **Use ISO-8601** for `--due` and `--date` values
- **`cards list` requires exactly one** of `--board` or `--list`, never both
- **Update commands need at least one** mutation field

## Intent Interpretation

When the user describes a Trello task in natural language, translate it into a multi-step workflow — not a single command. Think about what IDs you need and how to get them.

Common patterns:
- "Move card X to Done" → discover the Done list ID, then `cards move`, then verify
- "Add a checklist to card X" → `checklists create`, add items, list to confirm
- "Find boards about marketing" → `search boards --query marketing` or `boards list` and filter
- "Set priority to High on card X" → discover custom field ID, then `custom-fields items set`

Read `references/task-recipes.md` for complete workflow recipes covering all resource types.

## Command Reference

The CLI has 12 top-level command groups: `auth`, `boards`, `lists`, `cards`, `comments`, `checklists`, `attachments`, `custom-fields`, `labels`, `members`, `search`, `version`.

Read `references/command-digest.md` for the full command surface, flags, and validation rules.

## Error Handling

When a command returns `ok: false`, check the error code:

| Code | Meaning | What to do |
|------|---------|------------|
| `AUTH_REQUIRED` | No credentials configured | Guide user through `auth set` or `auth login` |
| `AUTH_INVALID` | Credentials rejected by Trello | Credentials may be expired — re-authenticate |
| `NOT_FOUND` | Resource ID doesn't exist | Re-discover the ID |
| `VALIDATION_ERROR` | Bad input (missing flag, wrong format) | Fix the command flags |
| `RATE_LIMITED` | Trello API rate limit hit | Wait and retry |
| `CONFLICT` | Resource state conflict | Re-fetch and retry |

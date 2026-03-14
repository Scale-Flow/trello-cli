# Trello CLI LLM Digest

**Repository:** [github.com/Scale-Flow/trello-cli](https://github.com/Scale-Flow/trello-cli)

## Purpose

This repository provides a cross-platform Trello CLI with deterministic JSON output. The CLI is intended for direct terminal use and for machine consumers such as scripts, agents, and future skills.

## Primary Characteristics

- All commands write JSON to `stdout`
- Success and error responses use stable envelopes
- Most resource commands require authentication
- The command tree is organized by Trello resource type
- `--pretty` formats JSON for humans
- `--verbose` sends diagnostics to `stderr`

## Output Contract

Success envelope:

```json
{"ok":true,"data":...}
```

Error envelope:

```json
{"ok":false,"error":{"code":"...","message":"..."}}
```

Important implications:

- Consumers should branch on `ok`
- Success payloads are always under `data`
- Failures are always under `error`
- A newline is always appended
- Pretty-printing changes formatting, not schema

## Error Codes

- `AUTH_REQUIRED`
- `AUTH_INVALID`
- `NOT_FOUND`
- `VALIDATION_ERROR`
- `CONFLICT`
- `RATE_LIMITED`
- `HTTP_ERROR`
- `FILE_NOT_FOUND`
- `UNSUPPORTED`
- `UNKNOWN_ERROR`

## Authentication Model

Auth sources:

- Stored credentials in OS keyring
- Environment variables `TRELLO_API_KEY` and `TRELLO_TOKEN`

Effective read order:

1. Keyring-backed stored credentials
2. Environment credentials

Auth commands:

- `trello auth set --api-key <key> --token <token>`
- `trello auth set-key --api-key <key>`
- `trello auth status`
- `trello auth login`
- `trello auth clear`

Auth modes that may appear in responses:

- `device`
- `manual`
- `interactive`
- `env`
- `key_only`

Interpretation:

- `device`, `manual`, `interactive`, and `env` can represent usable auth depending on whether both key and token exist
- `key_only` means login preparation state, not usable authenticated state
- `auth status` validates credentials against Trello and returns member information when configured

Device flow login (preferred):

- `trello auth login` first attempts the device flow via the Trello Connector Power-Up pairing service
- Displays a pairing code (e.g., `WDJB-MJHT`) for the user to enter in the Power-Up
- No API key or developer setup required — credentials are returned by the pairing service
- Falls back to browser login if the pairing service is unavailable

Browser login fallback:

- Requires a Trello API key from `--api-key`, stored key, or `TRELLO_API_KEY`
- Opens browser authorization flow
- Uses localhost callback `http://localhost:3007/callback`
- If browser launch fails, login URL is printed to `stderr`

### API Key Prerequisite

Users must create a Trello Power-Up before they can get an API key:

1. Go to [trello.com/power-ups/admin](https://trello.com/power-ups/admin)
2. Create a new Power-Up (it serves as the container for API credentials)
3. Go to the API Key tab and generate a key
4. Click the Token hyperlink to authorize and get a token

If a user reports they cannot find their API key or the old `trello.com/app-key` URL doesn't work, direct them through the Power-Up creation flow above.

## Global Flags

- `--pretty`
- `--verbose`

## Command Taxonomy

Top-level commands:

- `auth`
- `boards`
- `lists`
- `cards`
- `comments`
- `checklists`
- `attachments`
- `custom-fields`
- `labels`
- `members`
- `search`
- `version`

### Boards

- `boards list`
- `boards get --board <board-id>`
- `boards create --name <name> [--desc <text>] [--default-lists] [--default-labels] [--organization <org-id>] [--source-board <board-id>]`

### Lists

- `lists list --board <board-id>`
- `lists create --board <board-id> --name <name>`
- `lists update --list <list-id> [--name <name>] [--pos <number>]`
- `lists archive --list <list-id>`
- `lists move --list <list-id> --board <board-id> [--pos <number>]`

### Cards

- `cards list --board <board-id>` or `cards list --list <list-id>`
- `cards get --card <card-id>`
- `cards create --list <list-id> --name <name> [--desc <text>] [--due <iso-8601>]`
- `cards update --card <card-id> [--name] [--desc] [--due] [--labels <csv>] [--members <csv>]`
- `cards move --card <card-id> --list <list-id> [--pos <number>]`
- `cards archive --card <card-id>`
- `cards delete --card <card-id>`

### Comments

- `comments list --card <card-id>`
- `comments add --card <card-id> --text <text>`
- `comments update --action <action-id> --text <text>`
- `comments delete --action <action-id>`

### Checklists

- `checklists list --card <card-id>`
- `checklists create --card <card-id> --name <name>`
- `checklists delete --checklist <checklist-id>`
- `checklists items add --checklist <checklist-id> --name <name>`
- `checklists items update --card <card-id> --item <item-id> --state <complete|incomplete>`
- `checklists items delete --checklist <checklist-id> --item <item-id>`

### Attachments

- `attachments list --card <card-id>`
- `attachments add-file --card <card-id> --path <local-path> [--name <display-name>]`
- `attachments add-url --card <card-id> --url <http-or-https-url> [--name <display-name>]`
- `attachments delete --card <card-id> --attachment <attachment-id>`

### Custom Fields

- `custom-fields list --board <board-id>`
- `custom-fields get --field <field-id>`
- `custom-fields create --board <board-id> --name <name> --type <text|number|date|checkbox|list> [--card-front] [--option <value>...]`
- `custom-fields update --field <field-id> [--name <name>] [--card-front]`
- `custom-fields delete --field <field-id>`
- `custom-fields options list --field <field-id>`
- `custom-fields options add --field <field-id> --text <text> [--color <color>]`
- `custom-fields options update --field <field-id> --option <option-id> [--text <text>] [--color <color>]`
- `custom-fields options delete --field <field-id> --option <option-id>`
- `custom-fields items list --card <card-id>`
- `custom-fields items set --card <card-id> --field <field-id> <exactly one of: --text, --number, --date, --checked, --option>`
- `custom-fields items clear --card <card-id> --field <field-id>`

### Labels

- `labels list --board <board-id>`
- `labels create --board <board-id> --name <name> --color <color>`
- `labels add --card <card-id> --label <label-id>`
- `labels remove --card <card-id> --label <label-id>`

### Members

- `members list --board <board-id>`
- `members add --card <card-id> --member <member-id>`
- `members remove --card <card-id> --member <member-id>`

### Search

- `search cards --query <text>`
- `search boards --query <text>`

### Version

- `version`

## Validation Rules

- Most mutations require IDs rather than names
- `cards list` requires exactly one of `--board` or `--list`
- Update commands require at least one mutation field
- `--due` values must be ISO-8601 compatible
- `attachments add-url` only accepts valid `http` or `https` URLs
- `attachments add-file` requires an existing local file path
- Checklist item state must be `complete` or `incomplete`
- `custom-fields create` requires `--board`, `--name`, and `--type`; `--type` must be one of `text`, `number`, `date`, `checkbox`, `list`; `--option` is only allowed with `--type list`
- `custom-fields update` requires `--field` and at least one mutation flag
- `custom-fields items set` requires exactly one value flag (`--text`, `--number`, `--date`, `--checked`, or `--option`)
- `--date` values must be ISO-8601 compatible

## Recommended Usage Patterns

- Start with discovery commands, then move to mutation commands
- Prefer compact JSON in automation and `--pretty` in interactive use
- Use `auth status` as a health check before larger workflows
- Resolve names to IDs once, then use IDs for subsequent operations
- For complex flows, gather board and list IDs first to reduce ambiguity

Preferred workflow shape:

1. Authenticate
2. Create or discover resource IDs
3. Perform mutations
4. Re-fetch resource state if confirmation is needed

Board bootstrap example:

1. `trello boards create --name "Project Alpha" --default-lists`
2. `trello boards get --board <board-id>`
3. `trello lists list --board <board-id>`

## Common Task Recipes

### Create a Card

1. `trello lists list --board <board-id>`
2. `trello cards create --list <list-id> --name "Write docs" --desc "Initial draft"`
3. `trello cards list --list <list-id>`

### Move a Card

1. `trello lists list --board <board-id>`
2. `trello cards move --card <card-id> --list <destination-list-id> --pos 1`
3. `trello cards get --card <card-id>`

### Add A Comment

1. `trello comments add --card <card-id> --text "Ready for review"`
2. `trello comments list --card <card-id>`

### Attach A File

1. Confirm local file path exists
2. `trello attachments add-file --card <card-id> --path ./brief.pdf --name "Brief"`
3. `trello attachments list --card <card-id>`

### Manage A Checklist

1. `trello checklists create --card <card-id> --name "Release"`
2. `trello checklists items add --checklist <checklist-id> --name "Ship docs"`
3. `trello checklists items update --card <card-id> --item <item-id> --state complete`

### Manage Custom Fields

1. `trello custom-fields list --board <board-id>`
2. `trello custom-fields items list --card <card-id>`
3. `trello custom-fields items set --card <card-id> --field <field-id> --text "value"`
4. `trello custom-fields items list --card <card-id>`

### Search Before Acting

1. `trello search cards --query "documentation"`
2. `trello cards get --card <card-id>`
3. Continue with comments, labels, members, or attachments

## Documentation Map

- Main onboarding: `README.md`
- Guided usage: `docs/getting-started.md`
- Concepts: `docs/concepts/`
- Command reference: `docs/commands/`
- Recipes: `docs/examples/`

## Boundaries For Skill Creation

Use this file as a high-signal summary, not as the only source of truth.

When creating skills from this repo:

- Prefer this file for command discovery, workflow shape, and terminology
- Use the human docs for fuller examples and explanations
- Trust the JSON contract and validation rules here
- Avoid inventing undocumented flags or configuration behavior
- Treat undocumented internal packages as implementation detail unless surfaced in user docs

# Trello CLI Command Specification

**Purpose:** Define the external command contract for a Trello CLI intended for coding-agent use.

**Scope:** This document specifies the what, not the how. It defines command groups, command grammar, arguments, validation rules, response shapes, and Trello REST API mappings. It intentionally excludes implementation details, package layout, internal architecture, and build guidance.

**Primary Goal:** Provide a stable reference for developing a cross-platform CLI and companion Claude Code skill set focused on personal and team task management workflows.

---

## Product Principles

- Commands must be stable, explicit, and machine-friendly.
- All commands must return JSON.
- Resource identifiers remain first-class inputs.
- Command names should reflect user workflows, not raw endpoint names.
- The V1 surface should optimize for boards, lists, cards, comments, attachments, checklists, labels, members, and search.
- Low-value administrative Trello features are out of scope for V1.

---

## Output Contract

All commands return JSON on stdout.

Success shape:

```json
{
  "ok": true,
  "data": {}
}
```

Error shape:

```json
{
  "ok": false,
  "error": {
    "code": "STRING_CODE",
    "message": "Human-readable message"
  }
}
```

### Standard Error Codes

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

### Error Handling Rules

- Invalid or missing required flags must return `VALIDATION_ERROR`.
- A missing Trello resource must return `NOT_FOUND`.
- Trello authentication failures must return `AUTH_INVALID`.
- Missing local auth configuration must return `AUTH_REQUIRED`.
- Attachment file-path failures must return `FILE_NOT_FOUND`.
- Commands unsupported by the Trello API or intentionally excluded from the CLI must return `UNSUPPORTED`.

---

## Authentication Model

Trello API access requires both:

- an API key
- a user token

This CLI specification defines two supported V1 authentication modes:

1. Primary: interactive browser authorization
2. Secondary: manual API key and token entry

The API key is not a standalone alternative authentication mode. It is a required part of both supported flows.

### Primary V1 Auth Experience

The primary user experience must be an interactive login flow:

```bash
trello auth login
```

This command initiates Trello's `/1/authorize` flow to obtain a user token for the local user.

Product requirements:

- It must be the default recommended authentication path in user-facing docs and skills.
- It must produce credentials usable by subsequent CLI commands without further manual token entry.
- It must support permission scopes appropriate for task-management workflows.
- It must support configurable token duration if the product later chooses to expose Trello expiration choices.

Validation:

- The CLI must fail with `AUTH_REQUIRED` if commands requiring Trello access are run before a valid auth flow has completed.

Success schema:

```json
{
  "ok": true,
  "data": {
    "configured": true,
    "authMode": "interactive",
    "member": {
      "id": "string",
      "username": "string",
      "fullName": "string"
    }
  }
}
```

Trello references:
- `GET /members/me`
- Trello authorization guide `/1/authorize` flow

### Secondary V1 Auth Experience

The fallback user experience must support manual credential entry:

```bash
trello auth set --api-key <key> --token <token>
```

This mode is intended for headless environments, CI scenarios, advanced users, or situations where interactive login is not practical.

### Deferred Auth Mode

Basic OAuth 1.0 may be supported in a later version, but it is not required for V1 command compatibility.

If introduced in the future, it should be represented as a distinct auth mode rather than replacing the primary interactive login flow.

---

## Authentication Commands

### `trello auth status`

Returns the current authentication state and, when available, the authenticated member.

Example:

```bash
trello auth status
```

Success schema:

```json
{
  "ok": true,
  "data": {
    "configured": true,
    "authMode": "interactive",
    "member": {
      "id": "string",
      "username": "string",
      "fullName": "string"
    }
  }
}
```

Trello mapping:
- `GET /members/me`

### `trello auth login`

Performs interactive browser authorization and stores the resulting API key and user token for future commands.

Validation:
- User interaction is required

Success schema:

```json
{
  "ok": true,
  "data": {
    "configured": true,
    "authMode": "interactive",
    "member": {
      "id": "string",
      "username": "string",
      "fullName": "string"
    }
  }
}
```

Trello references:
- Trello authorization guide `/1/authorize` flow
- `GET /members/me`

### `trello auth set --api-key <key> --token <token>`

Stores credentials for subsequent CLI calls.

Validation:
- `--api-key` is required
- `--token` is required

Success schema:

```json
{
  "ok": true,
  "data": {
    "configured": true,
    "authMode": "manual"
  }
}
```

### `trello auth clear`

Clears stored credentials.

Success schema:

```json
{
  "ok": true,
  "data": {
    "configured": false,
    "authMode": null
  }
}
```

---

## Board Commands

### `trello boards list`

Lists visible boards for the authenticated member.

Example:

```bash
trello boards list
```

Success schema:

```json
{
  "ok": true,
  "data": [
    {
      "id": "string",
      "name": "string",
      "desc": "string",
      "closed": false,
      "url": "string"
    }
  ]
}
```

Trello mapping:
- `GET /members/{id}/boards`

### `trello boards get --board <board-id>`

Gets a single board.

Validation:
- `--board` is required

Success schema:

```json
{
  "ok": true,
  "data": {
    "id": "string",
    "name": "string",
    "desc": "string",
    "closed": false,
    "url": "string"
  }
}
```

Trello mapping:
- `GET /boards/{id}`

---

## List Commands

### `trello lists list --board <board-id>`

Lists lists in a board.

Validation:
- `--board` is required

Success schema:

```json
{
  "ok": true,
  "data": [
    {
      "id": "string",
      "name": "string",
      "closed": false,
      "idBoard": "string",
      "pos": 12345
    }
  ]
}
```

Trello mapping:
- `GET /boards/{id}/lists`

### `trello lists create --board <board-id> --name <name>`

Creates a list in a board.

Validation:
- `--board` is required
- `--name` is required

Trello mapping:
- `POST /lists`

### `trello lists update --list <list-id> [--name <name>] [--pos <pos>]`

Updates list metadata.

Validation:
- `--list` is required
- At least one mutation flag must be present

Trello mapping:
- `PUT /lists/{id}`

### `trello lists archive --list <list-id>`

Archives a list.

Validation:
- `--list` is required

Trello mapping:
- `PUT /lists/{id}/closed`

### `trello lists move --list <list-id> --board <board-id> [--pos <pos>]`

Moves a list to another board or position.

Validation:
- `--list` is required
- `--board` is required

Trello mapping:
- `PUT /lists/{id}`

---

## Card Commands

### `trello cards list --board <board-id>`

Lists cards in a board.

### `trello cards list --list <list-id>`

Lists cards in a list.

Validation:
- Exactly one of `--board` or `--list` must be provided

Success schema:

```json
{
  "ok": true,
  "data": [
    {
      "id": "string",
      "name": "string",
      "desc": "string",
      "closed": false,
      "idBoard": "string",
      "idList": "string",
      "due": "2026-03-20T00:00:00.000Z",
      "url": "string"
    }
  ]
}
```

Trello mappings:
- `GET /boards/{id}/cards`
- `GET /lists/{id}/cards`

### `trello cards get --card <card-id>`

Gets a single card.

Validation:
- `--card` is required

Trello mapping:
- `GET /cards/{id}`

### `trello cards create --list <list-id> --name <name> [--desc <text>] [--due <iso-date>]`

Creates a card in a list.

Validation:
- `--list` is required
- `--name` is required
- `--due`, when present, must be ISO-8601 compatible

Trello mapping:
- `POST /cards`

### `trello cards update --card <card-id> [--name <name>] [--desc <text>] [--due <iso-date>] [--labels <label-id[,label-id...]>] [--members <member-id[,member-id...]>]`

Updates mutable card fields.

Validation:
- `--card` is required
- At least one mutation flag must be present

Trello mapping:
- `PUT /cards/{id}`

### `trello cards move --card <card-id> --list <list-id> [--pos <pos>]`

Moves a card to another list or position.

Validation:
- `--card` is required
- `--list` is required

Trello mapping:
- `PUT /cards/{id}`

### `trello cards archive --card <card-id>`

Archives a card.

Validation:
- `--card` is required

Trello mapping:
- `PUT /cards/{id}/closed`

### `trello cards delete --card <card-id>`

Deletes a card permanently.

Validation:
- `--card` is required

Success schema:

```json
{
  "ok": true,
  "data": {
    "deleted": true,
    "id": "string"
  }
}
```

Trello mapping:
- `DELETE /cards/{id}`

---

## Comment Commands

Comments are modeled as Trello actions.

### `trello comments list --card <card-id>`

Lists comment actions for a card.

Validation:
- `--card` is required

Success schema:

```json
{
  "ok": true,
  "data": [
    {
      "id": "string",
      "type": "commentCard",
      "date": "2026-03-12T12:00:00.000Z",
      "memberCreator": {
        "id": "string",
        "username": "string",
        "fullName": "string"
      },
      "data": {
        "text": "string"
      }
    }
  ]
}
```

Trello mapping:
- `GET /cards/{id}/actions`

### `trello comments add --card <card-id> --text <text>`

Adds a comment to a card.

Validation:
- `--card` is required
- `--text` is required

Trello mapping:
- `POST /cards/{id}/actions/comments`

### `trello comments update --action <action-id> --text <text>`

Updates an existing comment action.

Validation:
- `--action` is required
- `--text` is required

Trello mapping:
- `PUT /actions/{id}/text`

### `trello comments delete --action <action-id>`

Deletes a comment action.

Validation:
- `--action` is required

Trello mapping:
- `DELETE /actions/{id}`

---

## Checklist Commands

### `trello checklists list --card <card-id>`

Lists checklists on a card.

Validation:
- `--card` is required

Trello mapping:
- `GET /cards/{id}/checklists`

### `trello checklists create --card <card-id> --name <name>`

Creates a checklist on a card.

Validation:
- `--card` is required
- `--name` is required

Trello mapping:
- `POST /checklists`

### `trello checklists delete --checklist <checklist-id>`

Deletes a checklist.

Validation:
- `--checklist` is required

Trello mapping:
- `DELETE /checklists/{id}`

### `trello checklists items add --checklist <checklist-id> --name <name>`

Adds a check item.

Validation:
- `--checklist` is required
- `--name` is required

Trello mapping:
- `POST /checklists/{id}/checkItems`

### `trello checklists items update --card <card-id> --item <item-id> --state <complete|incomplete>`

Updates a check item state.

Validation:
- `--card` is required
- `--item` is required
- `--state` must be `complete` or `incomplete`

Trello mapping:
- `PUT /cards/{idCard}/checkItem/{idCheckItem}`

### `trello checklists items delete --checklist <checklist-id> --item <item-id>`

Deletes a check item.

Validation:
- `--checklist` is required
- `--item` is required

Trello mapping:
- `DELETE /checklists/{idChecklist}/checkItems/{idCheckItem}`

---

## Attachment Commands

### `trello attachments list --card <card-id>`

Lists attachments on a card.

Validation:
- `--card` is required

Trello mapping:
- `GET /cards/{id}/attachments`

### `trello attachments add-file --card <card-id> --path <local-path> [--name <name>]`

Uploads a local file as an attachment.

Validation:
- `--card` is required
- `--path` is required
- `--path` must reference an existing local file

Trello mapping:
- `POST /cards/{id}/attachments`

### `trello attachments add-url --card <card-id> --url <url> [--name <name>]`

Attaches a URL to a card.

Validation:
- `--card` is required
- `--url` is required
- `--url` must be a valid URL

Trello mapping:
- `POST /cards/{id}/attachments`

### `trello attachments delete --card <card-id> --attachment <attachment-id>`

Deletes an attachment from a card.

Validation:
- `--card` is required
- `--attachment` is required

Trello mapping:
- `DELETE /cards/{idCard}/attachments/{idAttachment}`

---

## Label Commands

### `trello labels list --board <board-id>`

Lists labels in a board.

Validation:
- `--board` is required

Trello mapping:
- `GET /boards/{id}/labels`

### `trello labels create --board <board-id> --name <name> --color <color>`

Creates a label in a board.

Validation:
- `--board` is required
- `--name` is required
- `--color` is required

Trello mapping:
- `POST /labels`

### `trello labels add --card <card-id> --label <label-id>`

Adds an existing label to a card.

Validation:
- `--card` is required
- `--label` is required

Trello mapping:
- `POST /cards/{id}/idLabels`

### `trello labels remove --card <card-id> --label <label-id>`

Removes a label from a card.

Validation:
- `--card` is required
- `--label` is required

Trello mapping:
- `DELETE /cards/{id}/idLabels/{idLabel}`

---

## Member Commands

### `trello members list --board <board-id>`

Lists members of a board.

Validation:
- `--board` is required

Trello mapping:
- `GET /boards/{id}/members`

### `trello members add --card <card-id> --member <member-id>`

Adds a member to a card.

Validation:
- `--card` is required
- `--member` is required

Trello mapping:
- `POST /cards/{id}/idMembers`

### `trello members remove --card <card-id> --member <member-id>`

Removes a member from a card.

Validation:
- `--card` is required
- `--member` is required

Trello mapping:
- `DELETE /cards/{id}/idMembers/{idMember}`

---

## Search Commands

### `trello search cards --query <text>`

Searches cards visible to the authenticated member.

Validation:
- `--query` is required

Success schema:

```json
{
  "ok": true,
  "data": {
    "query": "string",
    "cards": [
      {
        "id": "string",
        "name": "string",
        "idBoard": "string",
        "idList": "string",
        "url": "string"
      }
    ]
  }
}
```

Trello mapping:
- `GET /search`

### `trello search boards --query <text>`

Searches boards visible to the authenticated member.

Validation:
- `--query` is required

Trello mapping:
- `GET /search`

---

## Response Field Guidance

V1 responses should include the fields most useful to coding agents:

- Boards: `id`, `name`, `desc`, `closed`, `url`
- Lists: `id`, `name`, `closed`, `idBoard`, `pos`
- Cards: `id`, `name`, `desc`, `closed`, `idBoard`, `idList`, `due`, `url`
- Comments: `id`, `type`, `date`, `memberCreator`, `data.text`
- Checklists: `id`, `name`, `idCard`, `checkItems`
- Attachments: `id`, `name`, `url`, `bytes`, `mimeType`, `date`, `isUpload`
- Labels: `id`, `idBoard`, `name`, `color`
- Members: `id`, `username`, `fullName`

This is a minimum field contract, not a maximum field contract.

---

## V1 Exclusions

The following areas are intentionally out of scope for V1:

- stickers
- reactions
- board stars
- workspace and organization administration
- member profile management
- power-up administration
- board preference administration unrelated to task workflows

---

## Existing CLI Delta

Compared with the reviewed .NET CLI, this specification intentionally adds or changes:

- Adds first-class `search` commands
- Adds first-class `labels` commands
- Adds first-class `members` commands
- Adds comment update and delete via `actions`
- Converts flat flag commands into grouped subcommands
- Standardizes JSON error shapes
- Standardizes validation semantics

---

## Official Reference Sources

- Trello Actions REST API: https://developer.atlassian.com/cloud/trello/rest/api-group-actions/
- Trello Boards REST API: https://developer.atlassian.com/cloud/trello/rest/api-group-boards/
- Trello Cards REST API: https://developer.atlassian.com/cloud/trello/rest/api-group-cards/
- Trello Checklists REST API: https://developer.atlassian.com/cloud/trello/rest/api-group-checklists/
- Trello Labels REST API: https://developer.atlassian.com/cloud/trello/rest/api-group-labels/
- Trello Lists REST API: https://developer.atlassian.com/cloud/trello/rest/api-group-lists/
- Trello Members REST API: https://developer.atlassian.com/cloud/trello/rest/api-group-members/
- Trello Search REST API: https://developer.atlassian.com/cloud/trello/rest/api-group-search/

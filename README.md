# Trello CLI

A cross-platform Go CLI for Trello with machine-friendly JSON output.

Every command writes a JSON envelope to `stdout`, which makes the CLI useful both for terminal users and for automation that needs predictable responses.

## Highlights

- Stable JSON success and error envelopes
- Interactive browser login and manual credential setup
- Broad Trello command coverage for boards, lists, cards, comments, labels, members, checklists, attachments, custom fields, and search
- Single-binary Go CLI with minimal runtime requirements

## Build

```bash
go build -o bin/trello ./cmd/trello
```

Or run it directly during development:

```bash
go run ./cmd/trello version
```

## Quick Start

If you already have a Trello API key and token:

```bash
go run ./cmd/trello auth set --api-key "$TRELLO_API_KEY" --token "$TRELLO_TOKEN"
go run ./cmd/trello auth status --pretty
go run ./cmd/trello boards list --pretty
```

If you want to use interactive login:

```bash
export TRELLO_API_KEY="your-api-key"
go run ./cmd/trello auth login
```

## Output Shape

Success:

```json
{"ok":true,"data":{"version":"dev","commit":"unknown","date":"unknown"}}
```

Error:

```json
{"ok":false,"error":{"code":"VALIDATION_ERROR","message":"--board is required"}}
```

Use `--pretty` on any command to indent the JSON.

## Common Commands

```bash
go run ./cmd/trello boards list
go run ./cmd/trello boards create --name "Project Alpha" --default-lists
go run ./cmd/trello lists list --board <board-id>
go run ./cmd/trello cards list --list <list-id>
go run ./cmd/trello cards create --list <list-id> --name "Follow up with customer"
go run ./cmd/trello search cards --query "customer"
go run ./cmd/trello custom-fields list --board <board-id>
go run ./cmd/trello custom-fields items set --card <card-id> --field <field-id> --text "value"
```

## Documentation

- [Getting Started](docs/getting-started.md)
- [Authentication](docs/concepts/authentication.md)
- [JSON Output](docs/concepts/json-output.md)
- [Errors](docs/concepts/errors.md)
- [Configuration](docs/concepts/configuration.md)
- [Command Reference](docs/commands/README.md)
- [Examples](docs/examples/create-a-card.md)
- [LLM Digest](LLM.md)

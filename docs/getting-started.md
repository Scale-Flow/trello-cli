# Getting Started

## What This CLI Does

`trello` is a JSON-first command-line interface for Trello. It is designed so commands can be piped into scripts, LLM workflows, and other tooling without scraping human-oriented text.

## Build And Verify

```bash
go build -o bin/trello ./cmd/trello
./bin/trello version --pretty
```

You can also replace `./bin/trello` with `go run ./cmd/trello` while developing.

## Authenticate

### Option 1: Manual Credentials

```bash
./bin/trello auth set --api-key "$TRELLO_API_KEY" --token "$TRELLO_TOKEN"
./bin/trello auth status --pretty
```

### Option 2: Interactive Browser Login

Interactive login requires a Trello API key. The CLI can read that key from `TRELLO_API_KEY` or from a previously stored key.

```bash
export TRELLO_API_KEY="your-api-key"
./bin/trello auth login
./bin/trello auth status --pretty
```

The login flow opens a browser and waits for the callback on `http://localhost:3007/callback`.

## Learn The Output Contract

Every command returns one of two shapes:

- Success: `{"ok":true,"data":...}`
- Error: `{"ok":false,"error":{"code":"...","message":"..."}}`

Add `--pretty` to any command for indented JSON.

## Run Your First Workflow

List boards:

```bash
./bin/trello boards list --pretty
```

List lists on a board:

```bash
./bin/trello lists list --board <board-id> --pretty
```

Create a card:

```bash
./bin/trello cards create --list <list-id> --name "Draft docs" --desc "Initial pass"
```

Search cards:

```bash
./bin/trello search cards --query "docs"
```

## Recommended Usage Patterns

- Use `--pretty` when reading output directly in a terminal.
- Omit `--pretty` in scripts and automations to keep output compact.
- Use `auth status` before automation runs that depend on stored credentials.
- Prefer IDs over names in workflows after the discovery step.
- Chain commands as discovery then mutation: `boards list` -> `lists list` -> `cards create` or `cards move`.

## Next Reading

- [Authentication](concepts/authentication.md)
- [JSON Output](concepts/json-output.md)
- [Command Reference](commands/cards.md)

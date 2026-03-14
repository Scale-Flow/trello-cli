# Getting Started

## What This CLI Does

`trello` is a JSON-first command-line interface for Trello. It is designed so commands can be piped into scripts, LLM workflows, and other tooling without scraping human-oriented text.

## Prerequisites: Get Your Trello API Key

The CLI needs a Trello API key and token to access your Trello data. Trello requires you to create a Power-Up to get these credentials:

1. **Create a Trello account** at [trello.com](https://trello.com) if you don't have one
2. **Create a Power-Up** at [trello.com/power-ups/admin](https://trello.com/power-ups/admin) — click "Create new" and fill in the basic details. The Power-Up doesn't need to do anything; it's just the container for your API credentials.
3. **Generate an API key** — in your Power-Up's settings, go to the **API Key** tab and click **Generate a new API Key**
4. **Get a token** — on the same page, click the **Token** hyperlink next to your API key. A permissions screen appears — click **Allow**, then copy the token displayed on the next page.

Keep both your API key and token handy for the authentication step below.

## Install

### Homebrew (macOS / Linux)

```bash
brew tap Scale-Flow/tap
brew install trello-cli
```

### Go Install

```bash
go install github.com/Scale-Flow/trello-cli/cmd/trello@latest
```

### Build From Source

```bash
git clone https://github.com/Scale-Flow/trello-cli.git
cd trello-cli
go build -o bin/trello ./cmd/trello
./bin/trello version --pretty
```

You can also use `go run ./cmd/trello` in place of `./bin/trello` during development.

## Authenticate

### Option 1: Device Flow via Connector Power-Up (Recommended)

The easiest way to authenticate — no API key or developer setup required:

```bash
trello auth login
```

The CLI displays a pairing code (e.g., `WDJB-MJHT`). Enter it in the CLI Connector Power-Up on any Trello board. Once paired, credentials are stored automatically.

### Option 2: Manual Credentials

If you have a Trello API key and token (see [Prerequisites](#prerequisites-get-your-trello-api-key)):

```bash
trello auth set --api-key <your-api-key> --token <your-token>
trello auth status --pretty
```

### Option 3: Interactive Browser Login

Interactive login requires a Trello API key. The CLI can read that key from `TRELLO_API_KEY` or from a previously stored key.

```bash
trello auth set-key --api-key <your-api-key>
trello auth login
trello auth status --pretty
```

If the device flow pairing service is unavailable, `auth login` falls back to a browser flow on `http://localhost:3007/callback`.

## Learn The Output Contract

Every command returns one of two shapes:

- Success: `{"ok":true,"data":...}`
- Error: `{"ok":false,"error":{"code":"...","message":"..."}}`

Add `--pretty` to any command for indented JSON.

## Run Your First Workflow

List boards:

```bash
trello boards list --pretty
```

List lists on a board:

```bash
trello lists list --board <board-id> --pretty
```

Create a card:

```bash
trello cards create --list <list-id> --name "Draft docs" --desc "Initial pass"
```

Search cards:

```bash
trello search cards --query "docs"
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

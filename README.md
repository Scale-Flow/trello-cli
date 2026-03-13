# Trello CLI

[![CI](https://github.com/Scale-Flow/trello-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/Scale-Flow/trello-cli/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/Scale-Flow/trello-cli?include_prereleases)](https://github.com/Scale-Flow/trello-cli/releases)

A CLI for Trello built for AI agents and Claude Code — every command returns structured JSON, making it ideal for autonomous workflows and agent-driven project management.

Includes a ready-to-use **Claude Code skill** so Claude can manage your Trello boards, cards, and lists out of the box.

## Highlights

- **Agent-first design** — stable JSON envelopes on every command for reliable tool use
- **Claude Code skill included** — install once, and Claude can manage Trello autonomously
- Broad command coverage: boards, lists, cards, comments, labels, members, checklists, attachments, custom fields, and search
- Interactive browser login and manual credential setup
- Single-binary Go CLI with minimal runtime requirements

## Installation

### Homebrew (macOS / Linux)

```bash
brew tap Scale-Flow/tap
brew install trello-cli
```

### Go Install

```bash
go install github.com/Scale-Flow/trello-cli/cmd/trello@latest
```

### Download Binary

Download the latest release for your platform from [GitHub Releases](https://github.com/Scale-Flow/trello-cli/releases).

## Claude Code Skill

This repo ships a Claude Code skill that lets Claude manage Trello autonomously — creating cards, moving work across lists, searching boards, and more.

### Install the Skill

Add the skill to your project's `.claude/settings.json`:

```json
{
  "skills": [
    "github:Scale-Flow/trello-cli//using-trello-cli"
  ]
}
```

Or install globally in `~/.claude/settings.json` to use it across all projects.

Once installed, Claude Code will automatically use the skill whenever you ask it to interact with Trello. Try prompts like:

- "Create a card on my Project board in the To Do list"
- "Move all cards labeled 'done' to the Done list"
- "Search for cards mentioning 'bug' across all boards"

> **Note:** The `trello` binary must be installed and authenticated (see below) for the skill to work.

## Prerequisites: Trello API Key

Before using the CLI, you need a Trello API key and token:

1. Log in to [Trello](https://trello.com)
2. Go to the [Power-Up admin page](https://trello.com/power-ups/admin)
3. Create a new Power-Up (it doesn't need to do anything — it's just the container for your API credentials)
4. In your Power-Up's settings, go to the **API Key** tab and click **Generate a new API Key**
5. Click the **Token** hyperlink next to your API key, approve the permissions, and copy the token

See [Getting Started](docs/getting-started.md) for the full walkthrough.

## Quick Start

Store your credentials:

```bash
trello auth set --api-key <your-api-key> --token <your-token>
trello auth status --pretty
```

Or use interactive browser login:

```bash
trello auth set-key --api-key <your-api-key>
trello auth login
```

List your boards:

```bash
trello boards list --pretty
```

## Output Shape

Success:

```json
{"ok":true,"data":{"version":"1.0.0","commit":"abc1234","date":"2026-03-13"}}
```

Error:

```json
{"ok":false,"error":{"code":"VALIDATION_ERROR","message":"--board is required"}}
```

Use `--pretty` on any command to indent the JSON.

## Common Commands

```bash
trello boards list
trello boards create --name "Project Alpha" --default-lists
trello lists list --board <board-id>
trello cards list --list <list-id>
trello cards create --list <list-id> --name "Follow up with customer"
trello search cards --query "customer"
trello custom-fields list --board <board-id>
trello custom-fields items set --card <card-id> --field <field-id> --text "value"
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

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, PR process, and guidelines.

## License

[MIT](LICENSE)

# `trello boards`

Manage Trello boards.

## Subcommands

- `trello boards list`
- `trello boards get --board <board-id>`
- `trello boards create --name <name> [--desc <text>] [--default-lists] [--default-labels] [--organization <org-id>] [--source-board <board-id>]`

## Flags

- `--board`: board ID for `get`
- `--name`: board name for `create`
- `--desc`: board description for `create`
- `--default-lists`: create Trello's default starter lists when creating the board
- `--default-labels`: create Trello's default labels when creating the board
- `--organization`: workspace or organization ID that should own the board
- `--source-board`: source board ID to copy when creating the board

## Examples

```bash
trello boards list --pretty
trello boards get --board <board-id> --pretty
trello boards create --name "Project Alpha" --desc "Delivery tracker" --default-lists --default-labels --pretty
```

## Usage Pattern

Use `boards create` when you want to bootstrap a new board, then `boards list` or `boards get` to capture the resulting board ID before drilling into lists, cards, labels, or members.

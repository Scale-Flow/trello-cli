# Trello CLI Command Digest

Quick-reference for all CLI commands, flags, and validation rules.

## Global Flags

- `--pretty` — format JSON for human reading (does not change schema)
- `--verbose` — send diagnostics to `stderr`

## Auth

| Command | Purpose |
|---------|---------|
| `auth set --api-key <key> --token <token>` | Store both credentials |
| `auth set-key --api-key <key>` | Store API key only (prep for login) |
| `auth status` | Validate credentials, show member info |
| `auth login` | Browser-based OAuth flow (needs API key first) |
| `auth clear` | Remove stored credentials |

Auth modes: `device`, `manual`, `interactive`, `env` = usable. `key_only` = not usable for resource commands.

`auth login` first attempts device flow pairing via the Trello Connector Power-Up. If the pairing service is unavailable, it falls back to browser login on `http://localhost:3007/callback`. If browser launch fails, the URL prints to `stderr`.

## Boards

```
boards list
boards get --board <board-id>
boards create --name <name> [--desc <text>] [--default-lists] [--default-labels]
             [--organization <org-id>] [--source-board <board-id>]
```

## Lists

```
lists list --board <board-id>
lists create --board <board-id> --name <name>
lists update --list <list-id> [--name <name>] [--pos <number>]
lists archive --list <list-id>
lists move --list <list-id> --board <board-id> [--pos <number>]
```

## Cards

```
cards list --board <board-id>       # OR
cards list --list <list-id>         # exactly one required
cards get --card <card-id>
cards create --list <list-id> --name <name> [--desc <text>] [--due <iso-8601>]
cards update --card <card-id> [--name] [--desc] [--due] [--labels <csv>] [--members <csv>]
cards move --card <card-id> --list <list-id> [--pos <number>]
cards archive --card <card-id>
cards delete --card <card-id>
```

## Comments

```
comments list --card <card-id>
comments add --card <card-id> --text <text>
comments update --action <action-id> --text <text>
comments delete --action <action-id>
```

## Checklists

```
checklists list --card <card-id>
checklists create --card <card-id> --name <name>
checklists delete --checklist <checklist-id>
checklists items add --checklist <checklist-id> --name <name>
checklists items update --card <card-id> --item <item-id> --state <complete|incomplete>
checklists items delete --checklist <checklist-id> --item <item-id>
```

## Attachments

```
attachments list --card <card-id>
attachments add-file --card <card-id> --path <local-path> [--name <display-name>]
attachments add-url --card <card-id> --url <http-or-https-url> [--name <display-name>]
attachments delete --card <card-id> --attachment <attachment-id>
```

Validation: `add-file` requires existing local path. `add-url` requires `http://` or `https://`.

## Custom Fields

```
custom-fields list --board <board-id>
custom-fields get --field <field-id>
custom-fields create --board <board-id> --name <name> --type <text|number|date|checkbox|list>
                     [--card-front] [--option <value>...]
custom-fields update --field <field-id> [--name <name>] [--card-front]
custom-fields delete --field <field-id>

custom-fields options list --field <field-id>
custom-fields options add --field <field-id> --text <text> [--color <color>]
custom-fields options update --field <field-id> --option <option-id> [--text <text>] [--color <color>]
custom-fields options delete --field <field-id> --option <option-id>

custom-fields items list --card <card-id>
custom-fields items set --card <card-id> --field <field-id>
                        <exactly one: --text | --number | --date | --checked | --option>
custom-fields items clear --card <card-id> --field <field-id>
```

Validation: `create` requires `--board`, `--name`, `--type`. `--option` only with `--type list`. `items set` requires exactly one value flag. `--date` must be ISO-8601.

## Labels

```
labels list --board <board-id>
labels create --board <board-id> --name <name> --color <color>
labels add --card <card-id> --label <label-id>
labels remove --card <card-id> --label <label-id>
```

## Members

```
members list --board <board-id>
members add --card <card-id> --member <member-id>
members remove --card <card-id> --member <member-id>
```

## Search

```
search cards --query <text>
search boards --query <text>
```

## Version

```
version
```

## Validation Summary

- Most mutations require IDs, not names
- `cards list` needs exactly one of `--board` or `--list`
- Update commands need at least one mutation field
- `--due` and `--date` must be ISO-8601
- Checklist item state must be `complete` or `incomplete`

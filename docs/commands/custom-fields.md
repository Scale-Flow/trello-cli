# `trello custom-fields`

Manage custom field definitions, list-type options, and per-card custom field values.

## Subcommands

### Field Definitions

- `trello custom-fields list --board <board-id>`
- `trello custom-fields get --field <field-id>`
- `trello custom-fields create --board <board-id> --name <name> --type <text|number|date|checkbox|list> [--card-front] [--option <value>...]`
- `trello custom-fields update --field <field-id> [--name <name>] [--card-front]`
- `trello custom-fields delete --field <field-id>`

### List-Type Options

- `trello custom-fields options list --field <field-id>`
- `trello custom-fields options add --field <field-id> --text <text> [--color <color>]`
- `trello custom-fields options update --field <field-id> --option <option-id> [--text <text>] [--color <color>]`
- `trello custom-fields options delete --field <field-id> --option <option-id>`

### Per-Card Values

- `trello custom-fields items list --card <card-id>`
- `trello custom-fields items set --card <card-id> --field <field-id> <exactly one of: --text, --number, --date, --checked, --option>`
- `trello custom-fields items clear --card <card-id> --field <field-id>`

## Rules

- `create` requires `--board`, `--name`, and `--type`
- `--type` must be one of `text`, `number`, `date`, `checkbox`, `list`
- `--option` flag on `create` is only allowed with `--type list`
- `update` requires `--field` and at least one mutation flag (`--name` or `--card-front`)
- `items set` requires exactly one value flag: `--text`, `--number`, `--date`, `--checked`, or `--option`
- `--date` values must be ISO-8601 compatible
- `options update` requires at least one of `--text` or `--color`

## Trello API Mappings

- `custom-fields list` — `GET /boards/{id}/customFields`
- `custom-fields get` — `GET /customFields/{id}`
- `custom-fields create` — `POST /customFields`
- `custom-fields update` — `PUT /customFields/{id}`
- `custom-fields delete` — `DELETE /customFields/{id}`
- `custom-fields options list` — `GET /customFields/{id}/options`
- `custom-fields options add` — `POST /customFields/{id}/options`
- `custom-fields options update` — `PUT /customFields/{id}/options/{idOption}` (Trello uses delete + re-add)
- `custom-fields options delete` — `DELETE /customFields/{id}/options/{idOption}`
- `custom-fields items list` — `GET /cards/{id}/customFieldItems`
- `custom-fields items set` — `PUT /cards/{idCard}/customField/{idCustomField}/item`
- `custom-fields items clear` — `PUT /cards/{idCard}/customField/{idCustomField}/item` (with empty value)

## Examples

```bash
trello custom-fields list --board <board-id>
trello custom-fields create --board <board-id> --name "Priority" --type list --option High --option Medium --option Low
trello custom-fields items list --card <card-id>
trello custom-fields items set --card <card-id> --field <field-id> --text "v2.1"
trello custom-fields items set --card <card-id> --field <field-id> --number 42
trello custom-fields items set --card <card-id> --field <field-id> --checked true
trello custom-fields items clear --card <card-id> --field <field-id>
trello custom-fields options add --field <field-id> --text "Critical" --color red
trello custom-fields options delete --field <field-id> --option <option-id>
```

# Trello CLI Task Recipes

Step-by-step workflows for common Trello tasks. Every recipe follows the **discover, mutate, verify** pattern.

## Preflight: Auth Check

Always run before resource commands:

```bash
trello auth status
```

If `ok` is `false`, help the user authenticate:

### Auth Check (Device Flow — Preferred)

1. `trello auth status` — check if already authenticated
2. If not authenticated:
   - Run `trello auth login`
   - The CLI will print a pairing code (e.g., `WDJB-MJHT`)
   - Tell the user: "Enter this code in your Trello board's CLI Connector Power-Up"
   - Wait for the command to complete
3. Verify: `trello auth status` should show `configured: true, authMode: "device"`

### Auth Check (Manual Fallback)

- `trello auth set --api-key <key> --token <token>` — manual credentials
- `trello auth login` — falls back to browser OAuth if pairing service is unavailable

## Board Discovery

```bash
trello boards list --pretty
```

Note the board ID for subsequent commands. If unsure which board:

```bash
trello search boards --query "project name"
```

## Bootstrap a New Board

```bash
trello boards create --name "Project Alpha" --default-lists --default-labels
# Extract board ID from response
trello lists list --board <board-id>
```

## Create a Card

```bash
# 1. Discover the target list
trello lists list --board <board-id>

# 2. Create the card
trello cards create --list <list-id> --name "Write docs" --desc "Initial draft"

# 3. Verify
trello cards get --card <card-id>
```

With a due date:
```bash
trello cards create --list <list-id> --name "Ship v2" --due "2026-04-01T17:00:00Z"
```

## Move a Card

```bash
# 1. Find the destination list
trello lists list --board <board-id>

# 2. Move
trello cards move --card <card-id> --list <destination-list-id>

# 3. Verify
trello cards get --card <card-id>
```

To set position: add `--pos 1` (top) or `--pos bottom`.

## Update a Card

```bash
trello cards update --card <card-id> --name "New name" --desc "Updated description"
trello cards get --card <card-id>
```

Add labels or members:
```bash
trello cards update --card <card-id> --labels "label1-id,label2-id" --members "member-id"
```

## Add a Comment

```bash
trello comments add --card <card-id> --text "Ready for review"
trello comments list --card <card-id>
```

## Manage Checklists

```bash
# 1. Create checklist
trello checklists create --card <card-id> --name "Release Steps"

# 2. Add items
trello checklists items add --checklist <checklist-id> --name "Update changelog"
trello checklists items add --checklist <checklist-id> --name "Tag release"
trello checklists items add --checklist <checklist-id> --name "Deploy"

# 3. Complete an item
trello checklists items update --card <card-id> --item <item-id> --state complete

# 4. Verify
trello checklists list --card <card-id>
```

## Attach a File

```bash
# 1. Confirm the file exists locally
ls -la ./brief.pdf

# 2. Attach
trello attachments add-file --card <card-id> --path ./brief.pdf --name "Project Brief"

# 3. Verify
trello attachments list --card <card-id>
```

For URL attachments:
```bash
trello attachments add-url --card <card-id> --url "https://example.com/doc" --name "Reference"
```

## Manage Labels

```bash
# List available labels
trello labels list --board <board-id>

# Create a new label
trello labels create --board <board-id> --name "Urgent" --color red

# Add label to card
trello labels add --card <card-id> --label <label-id>

# Remove label from card
trello labels remove --card <card-id> --label <label-id>
```

## Manage Members

```bash
# List board members
trello members list --board <board-id>

# Assign member to card
trello members add --card <card-id> --member <member-id>

# Remove member from card
trello members remove --card <card-id> --member <member-id>
```

## Custom Fields

```bash
# List fields on a board
trello custom-fields list --board <board-id>

# See current values on a card
trello custom-fields items list --card <card-id>

# Set a text field
trello custom-fields items set --card <card-id> --field <field-id> --text "value"

# Set a number field
trello custom-fields items set --card <card-id> --field <field-id> --number 42

# Set a checkbox
trello custom-fields items set --card <card-id> --field <field-id> --checked true

# Set a date field
trello custom-fields items set --card <card-id> --field <field-id> --date "2026-04-01T00:00:00Z"

# Set a list option
trello custom-fields items set --card <card-id> --field <field-id> --option <option-id>

# Clear a field value
trello custom-fields items clear --card <card-id> --field <field-id>

# Verify
trello custom-fields items list --card <card-id>
```

Creating custom field definitions:
```bash
# Text field
trello custom-fields create --board <board-id> --name "Priority" --type text

# List field with options
trello custom-fields create --board <board-id> --name "Status" --type list \
  --option "Not Started" --option "In Progress" --option "Done"

# Show on card front
trello custom-fields create --board <board-id> --name "Points" --type number --card-front
```

## Search Before Acting

When you don't know which board or card to target:

```bash
# Find cards
trello search cards --query "documentation"

# Find boards
trello search boards --query "roadmap"

# Then drill into the result
trello cards get --card <card-id>
trello boards get --board <board-id>
```

## Archive and Delete

```bash
# Archive a card (reversible)
trello cards archive --card <card-id>

# Delete a card (permanent)
trello cards delete --card <card-id>

# Archive a list
trello lists archive --list <list-id>
```

Prefer archiving over deleting unless the user explicitly asks for deletion.

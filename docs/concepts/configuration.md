# Configuration

## User-Facing Configuration Today

The current CLI exposes configuration mainly through:

- Global flags
- Authentication environment variables
- Stored credentials in the OS keyring

## Global Flags

Available on most commands:

- `--pretty`: indent JSON output
- `--verbose`: enable diagnostic logging on `stderr`

Examples:

```bash
trello boards list --pretty
trello cards create --list <list-id> --name "Ship docs" --verbose
```

## Environment Variables

Supported credential environment variables:

- `TRELLO_API_KEY`
- `TRELLO_TOKEN`

These are especially useful for CI jobs, shell sessions, and one-off automation.

```bash
export TRELLO_API_KEY="your-api-key"
export TRELLO_TOKEN="your-token"
trello boards list
```

## Credential Persistence

- `trello auth set` stores credentials for later commands.
- `trello auth set-key` stores only the API key.
- `trello auth clear` removes stored credentials.

## Config File

The CLI reads configuration from `~/.config/trello-cli/config.yaml`. Override the path with `TRELLO_CONFIG_PATH`.

Supported keys:

| Key | Default | Description |
|-----|---------|-------------|
| `pairing_service_url` | `https://trello-connector-production.up.railway.app` | URL of the device flow pairing service |
| `profile` | `default` | Credential profile name |
| `pretty` | `false` | Default pretty-print JSON |
| `verbose` | `false` | Default verbose diagnostics |
| `timeout` | `15s` | HTTP request timeout |
| `max_retries` | `3` | Max retry attempts for failed requests |
| `retry_mutations` | `false` | Whether to retry mutation (write) requests |

Precedence: environment variables (`TRELLO_` prefix) > config file > defaults.

Example `config.yaml`:

```yaml
pretty: true
timeout: 30s
```

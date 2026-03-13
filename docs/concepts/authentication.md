# Authentication

> **First time?** You need a Trello API key before authenticating. See [Prerequisites](../getting-started.md#prerequisites-get-your-trello-api-key) for how to create a Power-Up and generate your key.

## Supported Auth Flows

The CLI supports three practical auth states:

- `manual`: API key and token stored with `trello auth set`
- `interactive`: API key plus browser login token stored with `trello auth login`
- `env`: credentials supplied through `TRELLO_API_KEY` and `TRELLO_TOKEN`

There is also a `key_only` state when only an API key has been stored. That state is not enough to run authenticated commands, but it is enough to support later interactive login.

## Commands

Store an API key and token:

```bash
trello auth set --api-key <api-key> --token <token>
```

Store only an API key:

```bash
trello auth set-key --api-key <api-key>
```

Check current auth state:

```bash
trello auth status --pretty
```

Run interactive browser login:

```bash
export TRELLO_API_KEY=<api-key>
trello auth login
```

Clear stored credentials:

```bash
trello auth clear
```

## Storage Behavior

- Stored credentials are written to the OS keyring when available.
- Environment credentials are read from `TRELLO_API_KEY` and `TRELLO_TOKEN`.
- The CLI uses a fallback chain: keyring first, then environment variables.
- Stored credentials are associated with the `default` profile in the current implementation.

## Interactive Login Details

- The CLI opens the Trello authorization URL in a browser.
- It waits for the token callback on `localhost:3007`.
- If the browser cannot be opened automatically, the CLI prints the authorization URL to `stderr`.
- Interactive login needs an API key from `--api-key`, a stored key, or `TRELLO_API_KEY`.

## Auth Status Response

Typical configured response:

```json
{
  "ok": true,
  "data": {
    "configured": true,
    "authMode": "interactive",
    "member": {
      "id": "member-id",
      "username": "username",
      "fullName": "Full Name"
    }
  }
}
```

Typical not configured response:

```json
{
  "ok": true,
  "data": {
    "configured": false,
    "authMode": null,
    "member": null
  }
}
```

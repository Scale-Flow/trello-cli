# CLI Connector — Privacy Policy

**Effective date:** March 13, 2026

CLI Connector is a Trello Power-Up operated by Scale Flow ("we", "us", "our"). This policy explains what data the Power-Up and its relay server handle.

---

## What We Collect

### Data we do NOT collect

We do not collect, store, or have access to any Trello user personal data, including your email address, username, full name, avatar, bio, or locale.

### Data handled during pairing

When you pair a CLI with your Trello board, the relay server temporarily processes:

- **A one-time pairing code** — an 8-character code used to match your CLI session to the Power-Up.
- **A hashed device code** — a SHA-256 hash used to verify the CLI is polling for the correct session.
- **A Trello OAuth token** — encrypted with AES-256-GCM while in transit through the relay. The token is passed to the CLI and the pairing session is permanently deleted. The server never stores your token long-term.

Pairing sessions expire and are automatically purged after 15 minutes, whether or not pairing completes.

### Telemetry

The relay server logs anonymous event counters (e.g., "pairing started", "pairing completed") with timestamps. These events contain no user identifiers, IP addresses, or personal data.

---

## Where Data Is Stored

The relay server uses an ephemeral database. No data is replicated to third-party services, analytics platforms, or cloud storage.

---

## Data Shared with Third Parties

None. We do not share, sell, or transfer any data to third parties.

---

## Trello API Usage

CLI Connector uses the Trello REST API solely to facilitate the OAuth authorization flow. The Power-Up requests `read,write` scope with a 30-day token expiration. The OAuth token is stored only on your local machine by the CLI — not on our server.

---

## Your Rights

Since we do not store personal data, there is nothing to access, correct, or delete. You can revoke CLI access at any time from the Power-Up settings menu on your Trello board, or by revoking the token in your Trello account settings.

---

## Changes to This Policy

If we update this policy, we will revise the effective date above. Continued use of the Power-Up after changes constitutes acceptance.

---

## Contact

If you have questions about this policy, contact us at:

**Email:** scaleflowsolutions1@gmail.com

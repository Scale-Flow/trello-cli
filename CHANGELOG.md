# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [1.1.0] - 2026-03-13

### Added

- Device flow authentication via Trello Connector Power-Up — `trello auth login` now pairs with the Power-Up automatically, no API key or developer setup required
- New `device` auth mode for credentials obtained through the pairing service
- `pairing_service_url` configuration option (defaults to the hosted pairing service)
- Automatic fallback from device flow to browser-based login when the pairing service is unavailable

### Changed

- `trello auth login` now attempts device flow first before falling back to browser login
- Updated documentation to reflect device flow as the recommended authentication method

## [1.0.0] - 2026-03-13

### Added

- 12 command groups: `auth`, `boards`, `lists`, `cards`, `comments`, `checklists`, `attachments`, `custom-fields`, `labels`, `members`, `search`, `version`
- Three authentication modes: manual credentials, interactive browser login, environment variables
- Stable JSON output contract with success (`{"ok":true,"data":...}`) and error (`{"ok":false,"error":...}`) envelopes
- OS keyring credential storage via `go-keyring`
- Cross-platform support (Linux, macOS, Windows)
- `--pretty` flag for human-readable JSON output
- `--verbose` flag for diagnostic output to stderr
- Comprehensive documentation: getting started guide, concept docs, command reference, workflow examples
- LLM digest (`LLM.md`) for AI agent integration
- Claude Code skill for natural-language Trello workflows
- CI pipeline via GitHub Actions
- Automated releases via goreleaser with Homebrew tap support

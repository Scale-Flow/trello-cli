# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

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

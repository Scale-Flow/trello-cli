# Search Feature (Task 46) Design

## Overview
Task 46 extends the Trello CLI with search commands that mirror the API’s `/1/search` work for cards and boards. The API layer already has `CardSearchResult` and `BoardSearchResult`, so the goal is to add dedicated client helpers that inject the required `modelTypes` query parameters and ensure the CLI surfaces the wrapped response objects without stripping the API’s envelope.

## Approach Options
1. **Separate helpers + dedicated commands (recommended).** Implement `SearchCards` and `SearchBoards` on `trello.Client`, each building the `/1/search` query with the correct `modelTypes`. Cobra adds `search cards` and `search boards` subcommands, each validating `--query` and returning the wrapper directly. This keeps behavior explicit, stays true to the plan, and keeps command tests focused on a single model type.
2. **Single `search` command with `--model` flag.** One command would accept a model selector, and the client helper would accept that flag to build the request. This adds flexibility for future models (members, actions, etc.), but complicates validation (ensuring the modelType is allowed) and makes testing more combinatorial. It also diverges from the Task 46 instructions that explicitly request cards and boards subcommands.
3. **Shared helper pulling all results and filtering post-fetch.** A single helper could request multiple `modelTypes` and filter the response client-side, allowing reuse of the same HTTP call for both commands. That adds parsing complexity and extra data handling, which seems unnecessary given the simplicity of the required result wrappers.

## Chosen Design
- **API layer:** Add `SearchCards(ctx, query)` and `SearchBoards(ctx, query)` to `trello.Client`. Each method builds a `GET /1/search` call with `modelTypes=cards` or `modelTypes=boards`, uses a private struct to capture the slice returned by Trello, and returns `CardSearchResult{Query: query, Cards: raw.Cards}` or `BoardSearchResult{Query: query, Boards: raw.Boards}`. Remove any placeholder/search stubs from `client.go` so the concrete client owns the behavior.
- **Command layer:** Introduce `searchCmd` with two subcommands (`cards` and `boards`). Each subcommand defines `--query` as a required string flag, calls `contract.RequireFlag` on it, and uses the single-auth-check pattern before calling the corresponding client helper. Each command writes the wrapper returned by the helper directly via `output`.
- **Tests:**
  - `internal/trello/search_test.go` will spin up an `httptest` server to verify both helpers send `query` and `modelTypes` params and decode the returned cards/boards into the wrapper result, confirming the returned `Query` field matches the input.
  - `cmd/trello/search_test.go` will use the existing `mockAPI` pattern. Extend `mockAPI` with `searchCardsFn` and `searchBoardsFn`, and add command tests for: successful card search, missing `--query` flag error, successful board search, and missing `--query` error. Each success test verifies the command writes the `CardSearchResult`/`BoardSearchResult` envelope (including the `query` value that matches the flag). The failure tests ensure `contract.RequireFlag` fires before API calls.
- **Test helpers:** Update `cmd/trello/boards_test.go`’s `mockAPI` to include the new function hooks, and extend `setupTestAuth` in `cmd/trello/auth_test.go` to reset the `search cards --query` and `search boards --query` flags so tests are isolated.
- **Workflow:** Follow the mandated TDD loop: write failing tests first (internal and command-level), run them to confirm red, implement the minimal code, then re-run those targeted tests to confirm green. After everything passes, run `go test ./internal/trello/ -run "TestSearchCards|TestSearchBoards" -v` and `go test ./cmd/trello/ -run TestSearch -v` as required verification steps.

## Risks & Notes
- Search commands must continue using the wrapper structures returned by the API so downstream consumers (scripts, other tools) can rely on `query` being echoed.
- Avoid any shared mutable state by resetting new Cobra flags in `setupTestAuth`.
- There is no existing spec-document-reviewer subagent in this environment, so I’ll note here that the formal spec review loop could not be automated; let me know if you’d like me to run a manual review pass or adapt the spec in another way.

## Next Steps
1. After you confirm the design, add the tests described above (red first), wire the API/client/command code, and extend `mockAPI`/`setupTestAuth` for search. 2. Verify the required red/green test commands mentioned in the workflow section. 3. Once the implementation is working, update this plan file’s Task 46 checklist and confirm no stub remains in `internal/trello/client.go`. Let me know if anything in the design needs adjustment before I continue.

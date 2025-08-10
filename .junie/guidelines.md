# Junie Guidelines for sneat-tui

This document guides Junie (and other contributors) on how to make changes to this repository effectively, safely, and consistently.

## Project Overview
- Name: sneat-tui
- Purpose: Text-based User Interface for Sneat.app
- Language: Go (1.21+ recommended)
- UI framework: Bubble Tea (charmbracelet)

## Where Things Are
- main.go: Application entry point
- sneatui/initial.go: Initial screen/model
- sneatui/menu_unsigned.go: Menu when user is not signed in
- Tests: main_test.go, sneatui/*_test.go
- CI: .github/workflows/ci.yml

## Local Development
- Install Go: https://go.dev/dl/
- Run:
  - go run .
- Test:
  - go test ./...
- Lint/Format:
  - go fmt ./...
  - go vet ./...
- Dependencies:
  - Use go mod tidy after adding/removing imports

## Contribution Principles (Very Important)
1. Minimal diffs
   - Make the smallest change necessary to satisfy the issue.
   - Prefer adding small, focused functions or tests to large refactors.
2. Preserve public behavior
   - Do not break CLI behavior or user-visible flows unless the issue requires it.
3. Tests first when bug is reproducible
   - If fixing a bug, try to add/adjust a failing test that demonstrates it, then fix.
4. Keep it idiomatic Go
   - Use standard library where possible; keep code simple and readable.
5. Maintain cross-platform compatibility
   - Avoid OS-specific assumptions. Bubble Tea apps run in common terminals.
6. Document subtle behavior
   - Inline comments for non-obvious logic; update README if user-facing changes occur.

## Code Style
- go fmt, go vet must be clean.
- Prefer small, pure functions where helpful.
- Handle errors explicitly; avoid panics in library-like code.
- Avoid global mutable state; pass dependencies explicitly.

## Testing
- Run full suite: go test ./...
- Add tests alongside the code they cover (e.g., sneatui/feature_test.go).
- Keep tests deterministic and fast.

## CI Expectations
- GitHub Actions workflow should pass.
- If CI fails due to flaky tests or environment, investigate and stabilize tests.

## Common Tasks
- Adding a feature:
  - Add minimal surface to existing models/components in sneatui/.
  - Provide tests for state updates and view rendering if feasible.
- Fixing a bug:
  - Reproduce (preferably with a test), fix, run tests, and explain the root cause in the PR description.

## Dependency Management
- Update modules sparingly; prefer pinned versions that CI passes with.
- After changes: go mod tidy to keep go.mod/go.sum clean.

## Commit and PR Hygiene
- Clear commit messages: what and why.
- Keep PRs small and focused on one issue.

## Security & Privacy
- Do not check in secrets.
- Treat user data paths conservatively; avoid logging sensitive info.

## When in Doubt
- Choose the simplest solution that works.
- Err on the side of minimal changes and strong tests.

Last updated: 2025-08-10 10:01 local time.
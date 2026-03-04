# Repository Guidelines

## Project Structure & Module Organization
`skills-x` is a Go CLI with an interactive TUI plus bundled skill assets.

- `cmd/skills-x/`: CLI entrypoint, Cobra commands, TUI models/views, and i18n wiring.
- `pkg/`: reusable packages (`registry`, `discover`, `gitutil`, `skillvalidator`, `userregistry`, etc.).
- `skills/`: embedded first-party skills shipped with the binary.
- `docs/`: user and contributor docs (for example `docs/tui-guide.md`).
- `npm/`: npm wrapper and installer (`npm/install.js`, package metadata).
- `bin/`: build outputs (generated).

Keep new logic in `pkg/` when reusable, and in `cmd/skills-x/command/...` when command-specific.

## Build, Test, and Development Commands
- `make build`: syncs `skills/` into embed path and builds `bin/skills-x`.
- `make run ARGS="list"`: build and run with CLI arguments.
- `make test`: run all Go tests (`go test -v ./...`).
- `make build-local`: Linux amd64 local build including `local_skills/`.
- `make build-all`: cross-platform binaries in `bin/`.
- `make build-npm`: cross-platform release binaries for npm in `npm/bin/`.
- `make deps`: download and tidy Go modules.

## Coding Style & Naming Conventions
- Follow idiomatic Go; format with `gofmt` before committing.
- Use short, lowercase package names (`pkg/versioncheck`, `pkg/gitutil`).
- Keep command files action-oriented (`init.go`, `update.go`, `list.go`).
- Tests live next to code as `*_test.go`, with table-driven style where practical.
- For locale keys (`cmd/skills-x/i18n/locales/*.yaml`), use stable snake_case identifiers.

## Testing Guidelines
- Primary framework: Go `testing` package.
- Run `make test` before opening a PR.
- Prefer focused unit tests for parsing/validation/registry behavior, then command-level tests for CLI flows.
- Name tests descriptively, e.g. `TestRunAPIGenerate`, `TestValidateSkillPath`.

## Commit & Pull Request Guidelines
Recent history follows Conventional Commit style: `feat:`, `fix:`, `docs:`, `chore:` (example: `fix: starred skills not sorted on initial load in TUI LV3`).

- Use `<type>: <imperative summary>`; keep subject concise and specific.
- PRs should include: purpose, key changes, test evidence (`make test` output summary), and linked issues.
- Add screenshots/GIFs for TUI changes and note any i18n or registry data updates explicitly.

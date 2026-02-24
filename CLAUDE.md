# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Skills-X is a Go-based CLI tool for downloading and managing AI agent skills. It supports multiple AI IDEs (Claude Code, Cursor, Windsurf, CodeBuddy) and distributes skills via npm and Go install.

## Common Commands

```bash
# Build the binary
make build

# Run tests
make test

# Run a single test
go test -v ./cmd/skills-x/command/initcmd/...

# Install to ~/.local/bin
make install-local

# Build cross-platform binaries for npm release
make build-npm

# Build all platforms to bin/
make build-all
```

## Architecture

- **CLI Framework**: Uses Cobra (`github.com/spf13/cobra`)
- **Commands**: `list`, `init`, `tui` (see `cmd/skills-x/command/`)
- **Core Packages**:
  - `pkg/registry/` - Skill registry (registry.yaml)
  - `pkg/discover/` - Skill discovery
  - `pkg/gitutil/` - Git utilities for cloning skills
  - `pkg/versioncheck/` - Version checking against GitHub releases
- **Skills**: Embedded in binary at `cmd/skills-x/skills/skills/` via `go:embed`. Run `make sync-skills` before building to update embedded skills.
- **i18n**: Translations in `cmd/skills-x/i18n/locales/` (en.yaml, zh.yaml)

## Development Workflow

This project uses a plan-first workflow (see `.cursor/rules/development-workflow.mdc`):

1. **Discussion Phase**: When users提出需求, discuss the approach first
2. **Plan Confirmation**: Wait for user to confirm the plan
3. **Execution Phase**: Only start coding after explicit instructions like "开始开发" or "start development"
4. **Documentation**: Save implementation plans to `docs/` before coding

## Key Files

- `cmd/skills-x/main.go` - Entry point, registers commands
- `cmd/skills-x/command/initcmd/init.go` - Download/install skills
- `cmd/skills-x/command/list/list.go` - List available skills
- `cmd/skills-x/tui/` - Interactive TUI for skill management
- `pkg/registry/registry.yaml` - Skill registry definition

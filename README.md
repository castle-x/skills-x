# Skills-X

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue)](https://github.com/castle-x/skills-x)
[![npm](https://img.shields.io/npm/v/skills-x?color=CB3837&logo=npm)](https://www.npmjs.com/package/skills-x)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

> **Note**: For the industry-standard Agent Skills ecosystem, check out [skills.sh](https://skills.sh/) and [`npx skills`](https://github.com/vercel-labs/add-skill) by Vercel Labs.
> This project is just my personal collection for learning purposes.

[中文文档](README_ZH.md)

## Quick Install

```bash
# Using npm (recommended)
npm install -g skills-x

# Or use directly with npx
npx skills-x list

# Using go install
go install github.com/castle-x/skills-x/cmd/skills-x@latest
```

### Update

```bash
# Update to latest version
npm update -g skills-x

# Or reinstall
npm install -g skills-x@latest
```

### Usage

```bash
# Launch interactive TUI (default when no arguments)
skills-x

# Or explicitly
skills-x tui

# List all available skills
skills-x list

# Download all skills
skills-x init --all

# Download specific skill
skills-x init pdf
skills-x init ui-ux-pro-max

# Specify custom target directory
skills-x init pdf --target ./my-skills

# Force overwrite existing skills (skip confirmation)
skills-x init pdf -f
skills-x init --all --force
```

### Interactive TUI

Running `skills-x` with no arguments launches the interactive TUI mode:

1. **Select AI Tool** — Choose your IDE (Claude Code, Cursor, Windsurf, CodeBuddy, Codex)
2. **Select Target** — Install to global or project scope
3. **Select Skills** — Browse, search, select/deselect skills with keyboard
4. **Install/Uninstall** — See real-time progress with per-skill status (✓/✗)

### Setup for AI IDEs (CLI mode)

You can also use the CLI to specify the skills directory directly:

```bash
# For Claude Code
skills-x init --all --target .claude/skills

# For CodeBuddy
skills-x init --all --target .codebuddy/skills

# For Cursor
skills-x init --all --target .cursor/skills

# For Windsurf
skills-x init --all --target .windsurf/skills
```

---

## Directory Structure

```
skills-x/
├── cmd/skills-x/tui/   # Interactive TUI (Bubble Tea + Lipgloss)
├── pkg/products/        # AI IDE product definitions
├── pkg/registry/        # Skill sources registry
├── pkg/gitutil/         # Git clone with caching
└── skills/              # 🏰 Original Skills (embedded in binary)
```

## Collected Skills

> For learning purposes only

```
$ skills-x list

📦 github.com/anthropics/skills (Apache-2.0)
   algorithmic-art           Creating algorithmic art using p5.js...
   brand-guidelines          Apply Anthropic brand colors and typography
   canvas-design             Create visual art in PNG and PDF...
   doc-coauthoring           Collaborative document editing
   docx                      Word document creation, editing and analysis
   frontend-design           Frontend design best practices
   internal-comms            Internal communications templates
   mcp-builder               Generate MCP servers
   pdf                       PDF manipulation - extract, fill forms, merge
   pptx                      PowerPoint presentation creation and editing
   skill-creator             Create new agent skills
   slack-gif-creator         Create animated GIFs optimized for Slack
   theme-factory             Toolkit for styling artifacts with themes
   web-artifacts-builder     Build web artifacts with React
   webapp-testing            Test web applications
   xlsx                      Excel spreadsheet creation, formulas and charts

📦 github.com/remotion-dev/skills (Remotion License)
   remotion             Best practices for Remotion - Video creation...

📦 github.com/vercel-labs/skills (MIT)
   find-skills          Find and install agent skills using the Skills CLI

📦 github.com/vercel-labs/agent-skills (MIT)
   react-best-practices  React and Next.js performance optimization from Vercel Engineering
   react-native-skills   React Native best practices for AI agents
   web-design-guidelines 100+ rules for accessibility, performance and UX

📦 github.com/giuseppe-trisciuoglio/developer-kit (MIT)
   react-patterns        React 19 patterns covering Server Components...
   shadcn-ui             shadcn/ui component library with Radix UI and Tailwind CSS
   tailwind-css-patterns Tailwind CSS utility-first styling patterns...
   typescript-docs       TypeScript documentation with JSDoc, TypeDoc and ADR patterns

📦 github.com/masayuki-kono/agent-skills (MIT)
   implementation-plan   Guide for creating implementation plans

📦 github.com/obra/superpowers (MIT)
   brainstorming                   Brainstorm before any creative work...
   dispatching-parallel-agents     Dispatch 2+ independent tasks to work on in parallel
   executing-plans                 Execute written implementation plans step by step
   finishing-a-development-branch  Complete development branch when tests pass
   receiving-code-review           Handle code review feedback before implementing...
   requesting-code-review          Request code review when completing tasks or features
   subagent-driven-development     Execute plans with independent tasks using subagents
   systematic-debugging            Systematic approach to bugs, test failures...
   test-driven-development         TDD workflow - write tests before implementation
   using-superpowers               How to find and use superpowers skills effectively
   verification-before-completion  Verify work is complete, fixed, or passing...
   writing-plans                   Write implementation plans for multi-step tasks...
   writing-skills                  Create, edit, and verify skills

📦 github.com/tencentcloudbase/skills (MIT)
   ai-model-wechat          AI models in Mini Programs...
   auth-wechat              WeChat Mini Program authentication...
   cloudbase-guidelines     Essential CloudBase development guidelines...
   miniprogram-development  WeChat Mini Program development with CloudBase...
   no-sql-wx-mp-sdk          CloudBase document database for Mini Programs...

📦 github.com/nextlevelbuilder/ui-ux-pro-max-skill (MIT)
   ui-ux-pro-max         UI/UX design intelligence - 50 styles, 97 palettes...

📦 github.com/axtonliu/axton-obsidian-visual-skills (MIT)
   excalidraw-diagram  Generate Excalidraw diagrams from text

📦 github.com/bfollington/terma (CC-BY-SA-4.0)
   strudel             Strudel live-coding music for patterns, rhythms, melodies, and shareable URLs

📦 skills-x (Original)
   baidu-speech-to-text       Baidu speech-to-text (optimized for China mainland)
   go-embedded-spa            Go embedded SPA (single-binary deployment)
   go-i18n                    Go CLI i18n rules (author use)
   minimal-ui-design          Minimal UI design - low-noise, icon-forward
   newapi-deploy-config       New API deploy (host network) and channel configuration
   openclaw-session-header-fix Fix missing session header causing transcript overwrite
   skills-x                   Contribute skills to skills-x collection
   tui-design                 TUI design specification for CLI terminal UI

Total: 55 skills from 12 sources
```

---

## References

- [Claude Official Skills Docs](https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/skills) - Anthropic official documentation
- [Agent Skills Specification](https://agentskills.io/) - Open specification for agent skills
- [AGENTS.md](https://agents.md/) - AI-facing project documentation standard
- [Superpowers](https://github.com/obra/superpowers) - Complete software development workflow for coding agents

---

## 🏰 About X Skills

X skills are **original Skills** by the project author, stored in the `skills/` directory and aligned with common Git repo skill conventions.
They are displayed separately in the list with a ⭐ marker to distinguish them from community skills.

To contribute a new X skill, use the `skills-x` skill for guidance.
After creating a new skill, ask whether to add a `README.md` background summary document.
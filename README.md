# Skills-X

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue)](https://github.com/castle-x/skills-x)
[![npm](https://img.shields.io/npm/v/skills-x?color=CB3837&logo=npm)](https://www.npmjs.com/package/skills-x)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

> **Note**: For the industry-standard Agent Skills ecosystem, check out [skills.sh](https://skills.sh/) and [`npx skills`](https://github.com/vercel-labs/add-skill) by Vercel Labs.
> This project is my personal curated collection.

[中文文档](README_ZH.md)

## Quick Install

```bash
# Using npm (recommended)
npm install -g skills-x

# Using go install
go install github.com/castle-x/skills-x/cmd/skills-x@latest
```

```bash
# Update to latest version
npm install -g skills-x@latest
```

---

## Interactive TUI (Default Mode)

Running `skills-x` with no arguments launches the full interactive TUI — the recommended way to browse, install, update, and uninstall skills.

```bash
skills-x
```

### How the TUI works

The TUI has three levels of navigation:

#### Level 1 — Select AI Tool

Choose which IDE you are managing skills for. Each row shows the number of skills installed globally and at the project level.

```
Skills-X  AI Agent Skills Manager  v0.2.15

  ███████╗██╗  ██╗██╗██╗     ██╗     ███████╗     ██╗  ██╗
  ...

──────────────────────────────────────────────────────────
AI Tool              Global  Project
──────────────────────────────────────────────────────────
❯ Claude Code            45        3
  Cursor                  0        0
  Windsurf                0        0
  Codex                   2        0
  ...
──────────────────────────────────────────────────────────
↑/↓ navigate  |  Enter confirm  |  q quit
```

| Key | Action |
|-----|--------|
| `↑` / `↓` | Move cursor |
| `Enter` | Select tool and go to Level 2 |
| `q` | Quit |

#### Level 2 — Select Install Target

Choose whether to install skills globally (your home config dir) or at the project level (current working directory).

| Key | Action |
|-----|--------|
| `↑` / `↓` | Move cursor |
| `Enter` | Confirm target and go to Level 3 |
| `b` | Go back to Level 1 |
| `q` | Quit |

#### Level 3 — Browse & Select Skills

The main skill management screen. Skills are listed with their current state. Use keyboard shortcuts to search, filter, star, and select skills.

```
Skills For Cursor  /home/user/.cursor/skills
──────────────────────────────────────────────────────────
 输入 / 搜索技能
──────────────────────────────────────────────────────────
[ ]Not installed  [●]Installed  [+]Install  [-]Uninstall  [↑]Update  [↻]Checking
──────────────────────────────────────────────────────────
❯ [●] anthropic/algorithmic-art        2026-01-15  ★
  [●] anthropic/brand-guidelines       2026-01-15
  [ ] anthropic/canvas-design
  [●] superpowers/brainstorming        2026-02-10
  ...

1/53 ↓
Install: 0 | Update: 1 | Uninstall: 0
──────────────────────────────────────────────────────────
Space select | f star | u check update | A select all | Enter confirm | b back | q quit
```

**Navigation & Selection**

| Key | Action |
|-----|--------|
| `↑` / `↓` | Move cursor |
| `Space` | Toggle install / uninstall for the focused skill |
| `A` | Cycle select-all: (1) all Install/Update → (2) all Uninstall → (3) reset |
| `Enter` | Confirm and start installation/uninstall |
| `b` | Go back to Level 2 |
| `q` | Quit |

**Search & Filter**

| Key | Action |
|-----|--------|
| `/` or `、` | Enter search mode |
| `#` | Open the tag picker (interactive category filter) |
| `Esc` | Exit search / cancel tag picker |

When the tag picker opens, navigate with `↑`/`↓` and press `Enter` to apply a tag filter instantly. Available tags:

`#starred` `#featured` `#ai-efficiency` `#planning` `#web-frontend` `#mobile` `#backend` `#testing` `#code-review` `#office` `#design` `#writing` `#media` `#skills`

**Skill Actions**

| Key | Action |
|-----|--------|
| `f` | Star / unstar the focused skill (persists across sessions) |
| `u` | Check for updates on the focused skill (installed skills only) |

- Starred skills show a `★` marker and always sort to the top of the list.
- Starred skills are saved in `~/.config/skills-x/starred.json`.
- Skills with an available update show a `⚠ Update` badge.
- The focused skill's description is shown in the message area.

**Status Bar Legend**

```
[ ] Not installed    [●] Installed    [+] Will install
[-] Will uninstall   [↑] Will update  [↻] Checking for update
```

#### Level 4 — Installation Progress

After confirming, a progress screen shows real-time per-skill status:

```
Installing Skills

Progress: [==================>                     ] 45% (5/11)
Completed: 5 | Failed: 0

  ✓ anthropic/skill-creator
  ✓ superpowers/brainstorming
  ✓ superpowers/test-driven-development
  ...

Press any key to exit
```

---

## CLI Reference

```bash
# List all available skills (non-interactive)
skills-x list

# Install specific skills
skills-x init pdf
skills-x init pdf frontend-design

# Install all skills
skills-x init --all

# Install to a specific directory
skills-x init pdf --target .cursor/skills

# Force overwrite existing skills
skills-x init pdf --force
skills-x init --all --force

# Check for updates
skills-x update --check

# Update a specific skill
skills-x update pdf

# Update all installed skills
skills-x update --all
```

### Target directories by IDE

```bash
# Claude Code (global)
skills-x init --all --target ~/.claude/skills

# Cursor (project)
skills-x init --all --target .cursor/skills

# Windsurf
skills-x init --all --target .windsurf/skills

# CodeBuddy
skills-x init --all --target .codebuddy/skills

# Codex
skills-x init --all --target ~/.codex/skills
```

### Language

```bash
# Switch to English
SKILLS_LANG=en skills-x

# Switch to Chinese (default)
SKILLS_LANG=zh skills-x
```

---

## Collected Skills (49 skills from 13 sources)

```
📦 github.com/anthropics/skills (Apache-2.0)
   algorithmic-art           Creating algorithmic art using p5.js
   brand-guidelines          Apply Anthropic brand colors and typography
   canvas-design             Create visual art in PNG and PDF
   doc-coauthoring           Collaborative document editing
   docx                      Word document creation and editing
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
   xlsx                      Excel spreadsheet creation and formulas

📦 github.com/obra/superpowers (MIT)
   brainstorming                   Brainstorm before any creative work
   dispatching-parallel-agents     Dispatch independent tasks in parallel
   executing-plans                 Execute written implementation plans
   finishing-a-development-branch  Complete development branch after tests pass
   receiving-code-review           Handle code review feedback before implementing
   requesting-code-review          Request code review when completing tasks
   subagent-driven-development     Execute plans with independent tasks
   systematic-debugging            Systematic approach to bugs and failures
   test-driven-development         TDD workflow - write tests before implementation
   using-superpowers               How to find and use superpowers skills
   verification-before-completion  Verify work before claiming done
   writing-plans                   Write implementation plans for multi-step tasks
   writing-skills                  Create, edit, and verify skills

📦 github.com/vercel-labs/agent-skills (MIT)
   react-best-practices   React and Next.js performance optimization
   react-native-skills    React Native best practices
   web-design-guidelines  100+ rules for accessibility, performance, and UX

📦 github.com/vercel-labs/skills (MIT)
   find-skills            Find and install agent skills using the Skills CLI

📦 github.com/giuseppe-trisciuoglio/developer-kit (MIT)
   react-patterns         React 19 patterns - Server Components, Actions, hooks
   shadcn-ui              shadcn/ui with Radix UI and Tailwind CSS
   tailwind-css-patterns  Tailwind CSS utility-first styling patterns
   typescript-docs        TypeScript documentation with JSDoc and TypeDoc

📦 github.com/affaan-m/everything-claude-code (MIT)
   golang-testing         Go testing - table-driven, benchmarks, fuzzing, TDD

📦 github.com/bobmatnyc/claude-mpm-skills (MIT)
   golang-cli-cobra-viper  Go CLI with Cobra and Viper - commands, config, shell completion

📦 github.com/masayuki-kono/agent-skills (MIT)
   implementation-plan    Guide for creating implementation plans

📦 github.com/tencentcloudbase/skills (MIT)
   ai-model-wechat          AI models in Mini Programs (Hunyuan, DeepSeek)
   auth-wechat              WeChat Mini Program authentication
   cloudbase-guidelines     CloudBase development guidelines
   miniprogram-development  WeChat Mini Program development with CloudBase
   no-sql-wx-mp-sdk         CloudBase document database for Mini Programs

📦 github.com/nextlevelbuilder/ui-ux-pro-max-skill (MIT)
   ui-ux-pro-max    UI/UX design intelligence - 50 styles, 21 palettes

📦 github.com/axtonliu/axton-obsidian-visual-skills (MIT)
   excalidraw-diagram  Generate Excalidraw diagrams from text

📦 github.com/remotion-dev/skills (Remotion License)
   remotion   Best practices for Remotion - video creation with React

📦 github.com/bfollington/terma (CC-BY-SA-4.0)
   strudel   Strudel live-coding music - patterns, rhythms, melodies

📦 skills-x (Original)
   baidu-speech-to-text        Baidu speech-to-text (China mainland)
   go-embedded-spa             Go embedded SPA (single-binary deployment)
   go-i18n                     Go CLI i18n rules
   minimal-ui-design           Minimal UI design - low-noise, icon-forward
   newapi-deploy-config        New API deploy and channel configuration
   openclaw-session-header-fix Fix missing session header
   skills-x                    Contribute skills to skills-x collection
   tui-design                  TUI design specification for CLI terminal UI
```

---

## References

- [Claude Official Skills Docs](https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/skills) — Anthropic official documentation
- [Agent Skills Specification](https://agentskills.io/) — Open specification for agent skills
- [Superpowers](https://github.com/obra/superpowers) — Complete development workflow for coding agents

---

## 🏰 About X Skills

X skills are **original skills** by the project author, stored in the `skills/` directory and embedded in the binary.
They are listed separately with a ⭐ marker to distinguish them from community skills.

To contribute a new X skill, use the `skills-x` skill for guidance.

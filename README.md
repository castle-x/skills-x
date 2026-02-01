# Skills-X

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
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

### Setup for AI IDEs

We only provide skills download, you need to specify the skills directory for your AI IDE:

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
├── skills/         # 🏰 Original Skills (embedded in binary)
└── pkg/registry/   # Skill sources registry
```

## Collected Skills

> For learning purposes only

```
$ skills-x list

📦 github.com/anthropics/skills (Apache-2.0)
   algorithmic-art      Creating algorithmic art using p5.js...
   artifacts-builder    Build interactive artifacts with React
   brand-guidelines     Apply Anthropic brand colors and typography
   canvas-design        Create visual art in PNG and PDF...
   doc-coauthoring      Collaborative document editing
   docx                 Word document creation, editing and analysis
   frontend-design      Frontend design best practices
   internal-comms       Internal communications templates
   mcp-builder          Generate MCP servers
   pdf                  PDF manipulation - extract, fill forms, merge
   pptx                 PowerPoint presentation creation and editing
   skill-creator        Create new agent skills
   slack-gif-creator    Create animated GIFs optimized for Slack
   template-skill       Template for creating new skills
   theme-factory        Toolkit for styling artifacts with themes
   web-artifacts-builder Build web artifacts with React
   webapp-testing       Test web applications
   xlsx                 Excel spreadsheet creation, formulas and charts

📦 github.com/remotion-dev/skills (Remotion License)
   remotion             Best practices for Remotion - Video creation...

📦 github.com/vercel-labs/agent-skills (MIT)
   react-best-practices      React and Next.js performance optimization
   react-native-guidelines   React Native best practices for AI agents
   web-design-guidelines     100+ rules for accessibility, performance, UX

📦 skills-x (Original)
   skills-x             Contribute skills to skills-x collection

Total: 23 skills from 4 sources
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
After creating a new skill, ask whether to add a `REAEDME.md` background summary document.
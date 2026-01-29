# Skills-X

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue)](https://github.com/castle-x/skills-x)
[![npm](https://img.shields.io/npm/v/skills-x?color=CB3837&logo=npm)](https://www.npmjs.com/package/skills-x)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

My personal AI Agent Skills collection for quick reference and use.

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
├── castle-x/       # 🏰 My original Skills
├── skills/         # Curated open-source Skills
└── references/     # Skills writing references
```

## Sources

| Project | Link |
|---------|------|
| Claude Official | [anthropics/skills](https://github.com/anthropics/skills) |
| Superpowers | [obra/superpowers](https://github.com/obra/superpowers) |
| Awesome Claude Skills | [ComposioHQ/awesome-claude-skills](https://github.com/ComposioHQ/awesome-claude-skills) |

---

## Skills Directory (52 total)

### 🏰 Castle-X (Original Skills)
| Skill | Purpose |
|-------|---------|
| `skills-x` ⭐ | Contribute skills to skills-x collection |

### 🎨 Creative & Design
| Skill | Purpose |
|-------|---------|
| `ui-ux-pro-max` | **UI/UX Design Intelligence** - 67 styles/96 palettes/57 fonts/25 charts/13 stacks |
| `algorithmic-art` | p5.js generative art, flow fields, particle systems |
| `canvas-design` | Posters, visual art (.png/.pdf) |
| `brand-guidelines` | Anthropic brand styling |
| `theme-factory` | Artifact theme switching (10 presets) |
| `frontend-design` | Frontend design |
| `image-enhancer` | Image upscaling, sharpening, cleanup |

### 📄 Document Processing
| Skill | Purpose |
|-------|---------|
| `pdf` | PDF extract/fill/merge |
| `docx` | Word documents |
| `pptx` | PowerPoint presentations |
| `xlsx` | Excel sheets/formulas/charts |
| `document-skills` | Comprehensive document processing |
| `doc-coauthoring` | Document collaboration |

### 🛠️ Development Tools
| Skill | Purpose |
|-------|---------|
| `mcp-builder` | Build MCP servers |
| `artifacts-builder` | React+Tailwind+shadcn artifacts |
| `web-artifacts-builder` | Complex HTML artifacts |
| `webapp-testing` | Playwright testing |
| `langsmith-fetch` | LangSmith debug tracing |
| `changelog-generator` | Generate changelog from git commits |

### 🔄 Workflows
| Skill | Purpose |
|-------|---------|
| `brainstorming` | Brainstorm before creative work |
| `writing-plans` | Write task plans |
| `executing-plans` | Execute plans |
| `systematic-debugging` | Systematic debugging |
| `test-driven-development` | TDD workflow |
| `verification-before-completion` | Verify before completion |
| `subagent-driven-development` | Subagent-driven development |
| `dispatching-parallel-agents` | Parallel agent dispatching |

### 📝 Git & Code Review
| Skill | Purpose |
|-------|---------|
| `requesting-code-review` | Request CR |
| `receiving-code-review` | Handle CR feedback |
| `finishing-a-development-branch` | Complete branch |
| `using-git-worktrees` | Git Worktree isolation |

### ✍️ Writing
| Skill | Purpose |
|-------|---------|
| `content-research-writer` | Content research writing |
| `internal-comms` | Internal communications/reports |
| `tailored-resume-generator` | Custom resume generation |

### 🔗 Integrations
| Skill | Purpose |
|-------|---------|
| `connect` | Connect 1000+ services |
| `connect-apps` | Gmail/Slack/GitHub etc. |
| `connect-apps-plugin` | App connection plugin |
| `slack-gif-creator` | Slack GIF |

### 📊 Business & Analytics
| Skill | Purpose |
|-------|---------|
| `competitive-ads-extractor` | Competitor ad analysis |
| `developer-growth-analysis` | Developer growth |
| `lead-research-assistant` | Lead research |
| `meeting-insights-analyzer` | Meeting analysis |
| `twitter-algorithm-optimizer` | Tweet optimization |

### 🗂️ File Management
| Skill | Purpose |
|-------|---------|
| `file-organizer` | File organization |
| `invoice-organizer` | Invoice organization/tax prep |

### 🎲 Utilities
| Skill | Purpose |
|-------|---------|
| `video-downloader` | YouTube download |
| `domain-name-brainstormer` | Domain name ideas |
| `raffle-winner-picker` | Raffle picker |

### 🧰 Skills Development
| Skill | Purpose |
|-------|---------|
| `skill-creator` | Create Skills |
| `writing-skills` | Write/validate Skills |
| `skill-share` | Share Skills |
| `template-skill` | Skill template |
| `using-superpowers` | How to use Skills |

---

## References

- `references/claude_official_skills.md` - **Claude Official Skills Docs** (Overview/Quickstart/Best Practices/Enterprise)
- `references/agent_skills.md` - Agent Skills Open Specification (agentskills.io)
- `references/agents_md.md` - **AGENTS.md Specification** (AI-facing project documentation standard)
- `references/superpower.md` - Superpowers documentation

---

## 🏰 About Castle-X Skills

Castle-X skills are **original creations** by the project author, stored in the `castle-x/` directory.
They are displayed separately in the list with a ⭐ marker to distinguish them from community skills.

To contribute a new Castle-X skill, use the `skills-x` skill for guidance.

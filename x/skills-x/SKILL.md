---
name: skills-x
description: Guide for contributing new skills to the skills-x collection. This skill should be used when users want to add new open-source skills from external sources (like agentskills.io or anthropics/skills) to the skills-x repository. It covers the complete workflow from discovery to publishing.
license: MIT
metadata:
  author: x
  version: "1.3"
---

# Skills-X Contribution Guide

This skill provides a standardized workflow for contributing new skills to the skills-x collection.

## When to Use This Skill

- Adding new community skills from external sources (agentskills.io, anthropics/skills, etc.)
- Creating x original skills
- Updating existing skills with new versions
- Validating skill format compliance before submission

## Project Structure Overview

```
skills-x/
├── pkg/registry/        # Skill registry definition
│   └── registry.yaml    # Indexes skills from external sources
├── x/                   # X original skills (自研)
│   └── skills-x/        # This skill (embedded in binary)
├── cmd/skills-x/        # Go source code
│   ├── command/         # CLI commands (list, init)
│   └── i18n/
│       └── locales/     # Language files (zh.yaml, en.yaml)
├── npm/                 # npm package
│   └── package.json     # Version number here
└── Makefile             # Build commands
```

**Key Changes:**
- **No local `skills/` directory** - Skills are fetched directly from external repositories
- **Central registry** - All skill sources defined in `pkg/registry/registry.yaml`
- **Dynamic fetching** - Skills are cloned on-demand, not bundled with binary

---

## ⚠️ Internationalization (i18n) Rules - CRITICAL

**skills-x supports bilingual (Chinese/English) output. Follow these rules strictly:**

### Rule 1: NO Mixing Languages in a Single String

❌ **FORBIDDEN - Never mix Chinese and English in the same string:**
```go
// BAD: Mixed languages
desc = "🔄 套娃! Contribution guide (not for regular use)"
tag = "⭐ 作者自研 Original"
```

✅ **CORRECT - Use separate i18n keys:**
```go
// GOOD: Use i18n.T() to get localized string
desc = i18n.T("list_skillsx_desc")
tag = i18n.T("list_castlex_tag")
```

### Rule 2: All User-Facing Strings Must Use i18n

Any text displayed to users MUST go through the i18n system:

1. **Add keys to both language files:**

`cmd/skills-x/i18n/locales/zh.yaml`:
```yaml
my_message: "这是中文消息"
```

`cmd/skills-x/i18n/locales/en.yaml`:
```yaml
my_message: "This is English message"
```

2. **Use in Go code:**
```go
import "github.com/castle-x/skills-x/cmd/skills-x/i18n"

// Simple string
msg := i18n.T("my_message")

// With format arguments
msg := i18n.Tf("my_format_msg", arg1, arg2)
```

### Rule 3: i18n Key Naming Convention

| Type | Key Prefix | Example |
|------|------------|---------|
| Category names | `cat_` | `cat_creative`, `cat_document` |
| Skill descriptions | `skill_` | `skill_pdf`, `skill_docx` |
| Command descriptions | `cmd_` | `cmd_list_short` |
| List output | `list_` | `list_header`, `list_total` |
| Init output | `init_` | `init_success` |
| Error messages | `err_` | `err_skill_not_found` |

### Rule 4: Adding New Skill Descriptions

When adding a new skill, you MUST add descriptions to BOTH language files:

**Step 1:** Add English description to `en.yaml`:
```yaml
skill_new-skill: "Brief description in English"
```

**Step 2:** Add Chinese description to `zh.yaml`:
```yaml
skill_new-skill: "简短的中文描述"
```

**Step 3:** The code automatically picks up translations via:
```go
desc := i18n.T("skill_" + skillName)
```

### Rule 5: Testing Bilingual Output

**Always test BOTH languages after any UI changes:**

```bash
# Test Chinese
SKILLS_LANG=zh ./bin/skills-x list

# Test English  
SKILLS_LANG=en ./bin/skills-x list
```

### Rule 6: Environment Variable Priority

Language is detected in this order:
1. `SKILLS_LANG` (highest priority, skills-x specific)
2. `LANG` (system locale)
3. `LC_ALL` (system locale)
4. Default: `zh` (Chinese)

---

## Skill Directory Structure Requirements

All skills MUST follow the Agent Skills specification:

```
skill-name/
├── SKILL.md          # Required: Instructions + metadata
├── LICENSE.txt       # Required: License file
├── scripts/          # Optional: Executable code
├── references/       # Optional: Documentation
└── assets/           # Optional: Templates, resources
```

### SKILL.md Format Requirements

The `SKILL.md` file MUST contain YAML frontmatter with required fields:

```yaml
---
name: skill-name        # Required: lowercase, hyphens only, max 64 chars
description: ...        # Required: max 1024 chars, describe what and when
license: MIT            # Optional: license identifier
metadata:               # Optional: additional metadata
  author: example
  version: "1.0"
---
```

#### Name Field Rules

- Length: 1-64 characters
- Characters: lowercase letters, numbers, hyphens only
- Must NOT start or end with hyphen
- Must NOT contain consecutive hyphens (`--`)
- Must match parent directory name

**Valid:** `pdf-processing`, `data-analysis`, `code-review`
**Invalid:** `PDF-Processing`, `-pdf`, `pdf--processing`

#### Description Field Rules

- Length: 1-1024 characters
- Should clearly describe what the skill does AND when to use it
- Include keywords that help AI agents identify relevant tasks

---

## Contributing Community Skills (Open Source Skills)

**NEW WORKFLOW:** To add a new open source skill to the registry, you only need to edit `pkg/registry/registry.yaml`. No manual downloading or local copying required!

### Step 1: Find and Validate Source Skill

Search for skills at:
- https://agentskills.io/
- https://github.com/anthropics/skills
- https://github.com/vercel-labs/agent-skills
- https://github.com/remotion-dev/skills
- Other GitHub repositories with proper skill structure

**Before adding a skill to registry, verify:**
1. The repository has a valid `SKILL.md` file in the skill directory
2. The `SKILL.md` has proper YAML frontmatter with `name` and `description` fields
3. The skill name matches directory name (lowercase, hyphens only)
4. License information is available

### Step 2: Add Skill Source to Registry

Edit `pkg/registry/registry.yaml` and add a new source entry:

```yaml
# Example: Adding a new skill source
new-source-name:
  repo: github.com/owner/repo-name
  license: MIT  # or Apache-2.0, etc.
  skills:
    - name: skill-name
      path: path/to/skill/in/repo  # e.g., "skills/pdf" or "packages/skills/pdf"
      description: "Brief English description"
      description_zh: "简短的中文描述"
```

**Required fields for each source:**
- `repo`: GitHub repository URL (without https://)
- `license`: License type (MIT, Apache-2.0, etc.)
- `skills`: List of skills available from this source

**Required fields for each skill:**
- `name`: Skill name (must match directory name in repo)
- `path`: Path to skill directory within repository
- `description`: English description (max 1024 chars)
- `description_zh`: Chinese description (optional, max 1024 chars)

### Step 3: Verify Skill Installation

Test that the new skill can be installed:

```bash
# Build the binary
make build

# Test listing the skill
./bin/skills-x list --no-fetch | grep "skill-name"

# Test installing the skill
./bin/skills-x init skill-name --target /tmp/test-install
ls /tmp/test-install/skill-name/
```

**Validation checks:**
- [ ] Skill appears in `list` output
- [ ] Skill can be installed with `init`
- [ ] English and Chinese descriptions display correctly
- [ ] License information is accurate

---

## Contributing X Original Skills

X 自研 skills should be placed in `x/` directory, NOT in `skills/`.

### Step 1: Create Skill Directory

```bash
mkdir -p x/<skill-name>
```

### Step 2: Create Required Files

1. Create `SKILL.md` with proper frontmatter:
```yaml
---
name: <skill-name>
description: <what this skill does and when to use it>
license: MIT
metadata:
  author: x
  version: "1.0"
---

# <Skill Name>

<Detailed instructions for the AI agent>
```

2. Add `LICENSE.txt` (copy from project root or create)

### Step 3: Add i18n Translations

Same as community skills - add to both `en.yaml` and `zh.yaml`:

```yaml
skill_new-skill: "Description"
```

Note: X skills are automatically assigned category `x` and marked with `IsX: true`.

---

## Build and Test

```bash
# Build the binary
make build

# Test Chinese output
SKILLS_LANG=zh ./bin/skills-x list | grep "<skill-name>"

# Test English output
SKILLS_LANG=en ./bin/skills-x list | grep "<skill-name>"

# Test downloading the skill
./bin/skills-x init <skill-name> --target /tmp/test-skills
ls /tmp/test-skills/<skill-name>/
```

---

## Release Workflow

### Step 1: Update Version

Increment version in `npm/package.json`:
```json
"version": "0.1.X"  // increment patch version
```

### Step 2: Build for npm

```bash
make build-npm
```

### Step 3: Commit Changes

```bash
git add .
git commit -m "feat: add <skill-name> skill

- Add <skill-name> to skills collection
- Add i18n translations (en/zh)
- Update README"
```

### Step 4: Tag and Push

```bash
git tag -a v0.1.X -m "Add <skill-name> skill"
git push origin main
git push --tags
```

### Step 5: Create GitHub Release

⚠️ **CRITICAL: You MUST upload binary assets to the release!**

GitHub Release without binary assets is useless - users cannot download the tool.

```bash
# Build all platform binaries first
make build-npm

# Create release WITH binary assets (REQUIRED!)
gh release create v0.1.X \
  --title "v0.1.X - Add <skill-name>" \
  --notes "## Added
- New skill: <skill-name>
- Description: <brief description>" \
  npm/bin/skills-x-linux-amd64 \
  npm/bin/skills-x-linux-arm64 \
  npm/bin/skills-x-darwin-amd64 \
  npm/bin/skills-x-darwin-arm64 \
  npm/bin/skills-x-windows-amd64.exe
```

❌ **WRONG - Release without assets:**
```bash
# This creates an EMPTY release - USELESS!
gh release create v0.1.X --title "v0.1.X" --notes "..."
```

✅ **CORRECT - Release with all binary assets:**
```bash
gh release create v0.1.X --title "v0.1.X" --notes "..." \
  npm/bin/skills-x-linux-amd64 \
  npm/bin/skills-x-linux-arm64 \
  npm/bin/skills-x-darwin-amd64 \
  npm/bin/skills-x-darwin-arm64 \
  npm/bin/skills-x-windows-amd64.exe
```

If you forgot to upload assets, use `gh release upload`:
```bash
gh release upload v0.1.X \
  npm/bin/skills-x-* \
  --clobber
```

### Step 6: Publish to npm

```bash
cd npm && npm publish --access public
```

---

## Quick Reference Commands

```bash
# Validate a skill
head -20 skills/<name>/SKILL.md

# Build and test locally (both languages)
make build
SKILLS_LANG=zh ./bin/skills-x list
SKILLS_LANG=en ./bin/skills-x list

# Full release workflow
make build-npm
git add . && git commit -m "feat: add <skill>"
git tag -a v0.1.X -m "Add <skill>"
git push origin main --tags
gh release create v0.1.X --title "v0.1.X" --notes "..." \
  npm/bin/skills-x-linux-amd64 \
  npm/bin/skills-x-linux-arm64 \
  npm/bin/skills-x-darwin-amd64 \
  npm/bin/skills-x-darwin-arm64 \
  npm/bin/skills-x-windows-amd64.exe
cd npm && npm publish --access public
```

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Open source skill not in list | Check registry.yaml entry (repo, path, name fields) |
| X original skill not in list | Check i18n translations and skill directory in `x/` |
| Mixed language output | Ensure ALL strings use `i18n.T()`, no hardcoded text |
| Missing translation | Add keys to BOTH `en.yaml` and `zh.yaml` |
| init fails | Verify SKILL.md exists and has valid frontmatter |
| Windows fails | Ensure using `/` not `\` for embed.FS paths |
| Version mismatch | Check `npm/package.json` version matches build |
| **Release has no assets** | **MUST include binary files when running `gh release create`** |
| **Skill not in README** | **MUST update BOTH `README.md` and `README_ZH.md` with new skill** |

---

## Summary: Skill Contribution Workflows

| Skill Type | Storage Location | Description Source | Workflow |
|------------|------------------|--------------------|----------|
| **Open Source Skills** (from external repositories) | **Remote repositories only** - no local copy | Directly in `registry.yaml` (`description` and `description_zh` fields) | Edit `pkg/registry/registry.yaml` to add source and skills |
| **X Original Skills** (x自研) | `x/<name>/` directory (embedded in binary) | `i18n/locales/` files (`skill_<name>` keys) | 1. Create skill in `x/<name>/`<br>2. Add i18n translations<br>3. Build binary |

**Key Differences:**
- **Open Source Skills**: Descriptions stored in registry.yaml, fetched on-demand from remote repos
- **X Original Skills**: Descriptions stored in i18n files, embedded in binary during build

---

## Checklists for New Skills

### For Open Source Skills (in registry.yaml)

Before submitting a PR for adding new open source skills, verify:

- [ ] Source entry added to `pkg/registry/registry.yaml`
- [ ] `repo` field is correct GitHub repository URL
- [ ] `license` field specifies correct license type
- [ ] `name` field matches skill directory name
- [ ] `path` field points to correct skill location in repo
- [ ] `description` field provides clear English description
- [ ] `description_zh` field provides Chinese translation (optional but recommended)
- [ ] Skill appears in `list` output (use `--no-fetch` flag for testing)
- [ ] Skill can be installed with `init` command

### For X Original Skills (self-developed)

Before submitting a PR for new X skills, verify:

- [ ] Skill directory created in `x/<name>/`
- [ ] `SKILL.md` file with proper YAML frontmatter
- [ ] `skill_<name>` key added to `en.yaml`
- [ ] `skill_<name>` key added to `zh.yaml`
- [ ] No mixed Chinese/English in any single string
- [ ] Tested with `SKILLS_LANG=zh` - shows Chinese
- [ ] Tested with `SKILLS_LANG=en` - shows English
- [ ] Skill appears in `list` output under "skills-x (Original)" section

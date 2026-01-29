---
name: skills-x
description: Guide for contributing new skills to the skills-x collection. This skill should be used when users want to add new open-source skills from external sources (like agentskills.io or anthropics/skills) to the skills-x repository. It covers the complete workflow from discovery to publishing.
license: MIT
metadata:
  author: castle-x
  version: "1.3"
---

# Skills-X Contribution Guide

This skill provides a standardized workflow for contributing new skills to the skills-x collection.

## When to Use This Skill

- Adding new community skills from external sources (agentskills.io, anthropics/skills, etc.)
- Creating castle-x original skills
- Updating existing skills with new versions
- Validating skill format compliance before submission

## Project Structure Overview

```
skills-x/
├── skills/              # Community skills (from agentskills.io, anthropics/skills)
│   ├── pdf/
│   ├── docx/
│   └── ...
├── castle-x/            # Castle-X original skills (自研)
│   └── skills-x/        # This skill
├── cmd/skills-x/        # Go source code
│   ├── skills/
│   │   └── skills.go    # Skill metadata registry
│   └── i18n/
│       └── locales/     # Language files (zh.yaml, en.yaml)
├── npm/                 # npm package
│   └── package.json     # Version number here
└── Makefile             # Build commands
```

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

## Contributing Community Skills

### Step 1: Find and Validate Source Skill

Search for skills at:
- https://agentskills.io/
- https://github.com/anthropics/skills

Before downloading, verify the skill has:
1. A valid `SKILL.md` file with proper YAML frontmatter
2. `name` and `description` fields in frontmatter
3. Name matches directory name
4. Proper license information

**Do NOT download skills that lack proper SKILL.md structure.**

### Step 2: Download to Skills Directory

Download the skill to the root `skills/` directory (NOT `castle-x/`):

```bash
# Clone or download the skill
cp -r /path/to/source-skill skills/<skill-name>

# Or use git sparse-checkout for specific skills
git clone --depth 1 --filter=blob:none --sparse https://github.com/anthropics/skills
cd skills
git sparse-checkout set <skill-name>
```

### Step 3: Validate Skill Structure

Verify the downloaded skill:

```bash
# Check required files exist
ls skills/<skill-name>/SKILL.md

# Verify SKILL.md has proper frontmatter
head -20 skills/<skill-name>/SKILL.md
```

Required validation checks:
- [ ] `SKILL.md` exists
- [ ] YAML frontmatter present (starts with `---`)
- [ ] `name` field matches directory name
- [ ] `description` field is non-empty
- [ ] No uppercase letters in name
- [ ] No consecutive hyphens in name

### Step 4: Add i18n Translations (REQUIRED)

**You MUST add translations for both Chinese and English:**

1. **Edit `cmd/skills-x/i18n/locales/en.yaml`:**
```yaml
# Skill Descriptions
skill_new-skill: "Brief English description"
```

2. **Edit `cmd/skills-x/i18n/locales/zh.yaml`:**
```yaml
# Skill 描述
skill_new-skill: "简短的中文描述"
```

### Step 5: Update skills.go Metadata

Edit `cmd/skills-x/skills/skills.go` to add **category mapping only**:

```go
var skillCategories = map[string]string{
    // ... existing entries
    "new-skill-name": "category",  // creative/document/devtools/workflow/etc.
}
```

**Note:** Descriptions are now in i18n files, NOT in `skillDescriptions` map.

Available categories:
- `creative` - Design, art, UI/UX
- `document` - PDF, DOCX, XLSX, PPTX
- `devtools` - Development tools, MCP, testing
- `workflow` - Processes, debugging, TDD
- `git` - Git operations, code review
- `writing` - Content, communications
- `integration` - External services, APIs
- `business` - Analytics, research
- `files` - File management
- `utility` - General utilities
- `skilldev` - Skill development

### Step 6: Update README Files

Update both `README.md` and `README_ZH.md`:

1. Update skill count if changed
2. Add new skill to the appropriate category table
3. Keep both files in sync

---

## Contributing Castle-X Original Skills

Castle-X 自研 skills should be placed in `castle-x/` directory, NOT in `skills/`.

### Step 1: Create Skill Directory

```bash
mkdir -p castle-x/<skill-name>
```

### Step 2: Create Required Files

1. Create `SKILL.md` with proper frontmatter:
```yaml
---
name: <skill-name>
description: <what this skill does and when to use it>
license: MIT
metadata:
  author: castle-x
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

Note: Castle-X skills are automatically assigned category `castle-x` and marked with `IsCastleX: true`.

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
| Skill not in list | Check `skillCategories` in skills.go and i18n files |
| Mixed language output | Ensure ALL strings use `i18n.T()`, no hardcoded text |
| Missing translation | Add keys to BOTH `en.yaml` and `zh.yaml` |
| init fails | Verify SKILL.md exists and has valid frontmatter |
| Windows fails | Ensure using `/` not `\` for embed.FS paths |
| Version mismatch | Check `npm/package.json` version matches build |
| **Release has no assets** | **MUST include binary files when running `gh release create`** |
| **Skill not in README** | **MUST update BOTH `README.md` and `README_ZH.md` with new skill** |

---

## Summary: Where to Put Skills

| Skill Type | Directory | i18n Key | Category |
|------------|-----------|----------|----------|
| Community (from external sources) | `skills/<name>/` | `skill_<name>` in both yaml files | `skillCategories` map |
| Castle-X (original/自研) | `castle-x/<name>/` | `skill_<name>` in both yaml files | Auto: `castle-x` |

---

## i18n Checklist for New Skills

Before submitting a PR, verify:

- [ ] `skill_<name>` key added to `en.yaml`
- [ ] `skill_<name>` key added to `zh.yaml`
- [ ] No mixed Chinese/English in any single string
- [ ] Tested with `SKILLS_LANG=zh` - shows Chinese
- [ ] Tested with `SKILLS_LANG=en` - shows English
- [ ] Category added to `skillCategories` map

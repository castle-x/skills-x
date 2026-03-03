---
name: skills-x
description: Guide for contributing new skills to the skills-x collection. This skill should be used when users want to add new open-source skills from external sources (like agentskills.io or anthropics/skills) to the skills-x repository. It covers the complete workflow from discovery to publishing.
license: MIT
metadata:
  author: x
  version: "1.5"
---

# Skills-X Contribution Guide

This skill provides a standardized workflow for contributing new skills to the skills-x collection.

## When to Use This Skill

- Adding new community skills from external sources (agentskills.io, anthropics/skills, etc.)
- Creating x original skills
- Updating existing skills with new versions
- Validating skill format compliance before submission
- After creating a new skill, ask whether to generate a README (background summary)

## Project Structure Overview

```
skills-x/
├── pkg/registry/        # Skill registry definition
│   └── registry.yaml    # Indexes skills from external sources
├── skills/              # First-party skill sources (自研)
│   └── skills-x/        # This contribution skill source
├── cmd/skills-x/        # Go source code
│   ├── command/         # CLI commands (list, init)
│   └── i18n/
│       └── locales/     # Language files (zh.yaml, en.yaml)
├── npm/                 # npm package
│   └── package.json     # Version number here
└── Makefile             # Build commands
```

**Key Changes:**
- **Registry-first architecture** - All installable skills are indexed in `pkg/registry/registry.yaml`
- **Merged registry view** - Runtime uses built-in registry + user registry for list/init/update flows
- **Remote source install** - Skills are fetched from repository paths on demand (no embedded skill payloads)

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

For installable skills, descriptions should be defined in `pkg/registry/registry.yaml`:

```yaml
- name: new-skill
  path: skills/new-skill
  description: "Brief English description"
  description_zh: "简短的中文描述"
```

Only add `i18n` keys when introducing new CLI/TUI message keys.

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

### Step 3: Verify Skill Installation (Always)

After updating the registry, ALWAYS build and run list/init tests (no need to ask).

```bash
# Build the binary
make build

# Test listing the skill
./bin/skills-x list | grep "skill-name"

# Test installing the skill
./bin/skills-x init skill-name --target /tmp/test-install
ls /tmp/test-install/skill-name/
```

**Validation checks:**
- [ ] Skill appears in `list` output
- [ ] Skill can be installed with `init`
- [ ] English and Chinese descriptions display correctly
- [ ] License information is accurate

### Step 4: Ask for Release Actions

After tests finish, ask the user whether they want to:
- Commit and push
- Create GitHub release
- Publish to npm

---

## Contributing Self-Developed Skills (自研)

Self-developed skills should be placed in the `skills/` directory at the project root.

### Step 1: Create Skill Directory

```bash
mkdir -p skills/<skill-name>
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

3. Ask the user whether to add a summary document named `README.md`.
   - Purpose: describe the skill’s background, the problem it solves, and the author's goals.
   - Do NOT include any secrets or API keys.

### Step 3: Add i18n Translations

Same as community skills - add to both `en.yaml` and `zh.yaml`:

```yaml
skill_new-skill: "Description"
```

Note: Self-developed skills are treated as normal registry skills. Keep `registry.yaml` entries accurate (`repo`, `path`, tags, descriptions).

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

### Step 3: Pre-Release Testing (CRITICAL)

⚠️ **IMPORTANT: Always run this test before releasing to catch broken or missing skills!**

Test all skills can be installed successfully:

```bash
# Use a clean temporary directory
TEST_DIR=$(mktemp -d)
echo "Testing in: $TEST_DIR"

# Test installing all skills
./bin/skills-x init --all --target "$TEST_DIR"

# Check for failures
if [ $? -ne 0 ]; then
  echo "❌ Some skills failed to install!"
  echo "Review the output above for skills that are:"
  echo "  - Not found in repository"
  echo "  - Have incorrect paths"
  echo "  - Repository no longer exists"
  exit 1
fi

# Clean up
rm -rf "$TEST_DIR"
echo "✅ All skills tested successfully"
```

**If any skills fail:**

1. **Skill not found in repo** (`⚠ 在仓库中未找到 skill 路径`):
   - The skill path in `registry.yaml` is incorrect
   - The skill was removed/renamed in the source repository
   - **Action**: Remove from `pkg/registry/registry.yaml` or fix the path

2. **Repository not accessible**:
   - The repository was deleted or made private
   - **Action**: Remove the entire source from `pkg/registry/registry.yaml`

3. **Clone failed**:
   - Network issue (retry)
   - Repository URL changed
   - **Action**: Update repo URL or remove from registry

**After fixing registry.yaml:**

```bash
# Rebuild and test again
make build-npm
./bin/skills-x init --all --target "$(mktemp -d)"
```

### Step 4: Commit Changes

```bash
git add .
git commit -m "feat: add <skill-name> skill

- Add <skill-name> to skills collection
- Add i18n translations (en/zh)
- Update README"
```

### Step 5: Tag and Push

```bash
git tag -a v0.1.X -m "Add <skill-name> skill"
git push origin main
git push --tags
```

### Step 6: Create GitHub Release

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

### Step 7: Publish to npm

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
./bin/skills-x init --all --target "$(mktemp -d)"  # Test all skills
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
| Self-developed skill not in list | Check `pkg/registry/registry.yaml` source entry, skill path, and source repository accessibility |
| Mixed language output | Ensure ALL strings use `i18n.T()`, no hardcoded text |
| Missing translation | Add keys to BOTH `en.yaml` and `zh.yaml` |
| init fails | Verify SKILL.md exists and has valid frontmatter |
| Windows fails | Ensure registry `path` uses `/` separators and target directories are writable |
| Version mismatch | Check `npm/package.json` version matches build |
| **Release has no assets** | **MUST include binary files when running `gh release create`** |
| **Skill not in README** | **MUST update BOTH `README.md` and `README_ZH.md` with new skill** |

---

## Summary: Skill Contribution Workflows

| Skill Type | Storage Location | Description Source | Workflow |
|------------|------------------|--------------------|----------|
| **Open Source Skills** (from external repositories) | **Remote repositories only** - no local copy | Directly in `registry.yaml` (`description` and `description_zh` fields) | Edit `pkg/registry/registry.yaml` to add source and skills |
| **Self-Developed Skills** (自研) | `skills/<name>/` directory + corresponding registry source entry | Primarily `registry.yaml` (`description` and `description_zh` fields) | 1. Create skill in `skills/<name>/`<br>2. Add/update `pkg/registry/registry.yaml` entry<br>3. Verify list/init/update |

**Key Differences:**
- **Open Source Skills**: Managed by repository/path entries in registry.yaml
- **Self-Developed Skills**: Also managed through registry source entries; no embedded-skill install path

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
- [ ] Skill appears in `list` output
- [ ] Skill can be installed with `init` command

### For Self-Developed Skills (自研)

Before submitting a PR for new self-developed skills, verify:

- [ ] Skill directory created in `skills/<name>/`
- [ ] `SKILL.md` file with proper YAML frontmatter
- [ ] `pkg/registry/registry.yaml` contains correct source/path entry for the skill
- [ ] Registry description fields (`description` / `description_zh`) are complete
- [ ] `skills-x list` shows the skill under `github.com/castle-x/skills-x`
- [ ] `skills-x init <name> --target <tmp>` installs successfully
- [ ] `skills-x update <name> --target <tmp> --check` runs successfully

---
name: skills-x
description: Guide for contributing new skills to the skills-x collection. This skill should be used when users want to add new open-source skills from external sources (like agentskills.io or anthropics/skills) to the skills-x repository. It covers the complete workflow from discovery to publishing.
license: MIT
metadata:
  author: x
  version: "1.7"
---

# Skills-X Contribution Guide

This skill provides a standardized workflow for contributing new skills to the skills-x collection.

## When to Use This Skill

- Adding new community skills from external sources (agentskills.io, anthropics/skills, etc.)
- Creating x original self-developed skills (自研)
- Updating existing skills with new versions
- Validating skill format compliance before submission
- After creating a new skill, ask whether to generate a README (background summary)

---

## ⚡ Core Principle: No Binary Rebuild Needed for Skills

> **Adding or updating skills NEVER requires rebuilding or republishing the skills-x binary.**

The skills-x tool uses a **registry-first architecture**:

- The binary only contains the **registry index** (skill metadata)
- Actual skill content is **fetched from GitHub at install time**
- The registry is **cached locally** at `~/.config/skills-x/registry.yaml`

This means:
- ✅ Add a self-developed skill → push to git → done
- ✅ Add a third-party skill → push to git → done
- ✅ Users get new skills by running `skills-x registry update`
- ❌ No version bump, no binary build, no npm publish needed

**The only time a full release is needed** is when Go source code or CLI behavior changes.

---

## Project Structure Overview

```
skills-x/
├── pkg/registry/        # Skill registry definition
│   └── registry.yaml    # Central index of ALL skills (self-developed + third-party)
├── skills/              # Self-developed skills (自研) source code
│   └── skills-x/        # This contribution skill
├── cmd/skills-x/        # Go source code
│   ├── command/         # CLI commands (list, init, registry)
│   └── i18n/
│       └── locales/     # Language files (zh.yaml, en.yaml)
├── npm/                 # npm package
│   └── package.json     # Version number here
└── Makefile             # Build commands
```

**Architecture summary:**
- `pkg/registry/registry.yaml` — single source of truth for all skills
- `skills/<name>/` — self-developed skill content (committed to this repo)
- Third-party skills — content lives in external repos, only metadata in registry.yaml

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

## User Guide: Keeping Skills Up-to-Date

### The `registry update` Command

```bash
skills-x registry update
```

This command downloads the latest `registry.yaml` from GitHub and caches it locally at `~/.config/skills-x/registry.yaml`.

**When to run it:**
- After any skill has been added/updated and pushed to GitHub
- To see newly contributed skills before upgrading the binary
- Anytime `skills-x list` doesn't show a skill you expect

**How the cache works:**
| State | What `list` / `init` sees |
|-------|--------------------------|
| No cache | Registry embedded in the binary (older) |
| Cache present | Cached registry from GitHub (newer) |
| After `registry update` | Latest registry from GitHub |

> ⚠️ **Important for local development:** After pushing new skills to GitHub, you MUST run `skills-x registry update` before `skills-x list` will show the new skill. The locally built binary still has the old registry embedded until the cache is refreshed.

---

## Contributing Self-Developed Skills (自研)

> **No binary release needed.** Just create the skill, update the registry, push to git.

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

3. Ask the user whether to add a `README.md` (background, problem it solves, author goals — no secrets).

### Step 3: Update Registry

Add the skill to `pkg/registry/registry.yaml` under the `castle-x-skills-x` source:

```yaml
castle-x-skills-x:
  repo: github.com/castle-x/skills-x
  license: MIT
  skills:
    # ... existing skills ...
    - name: <skill-name>
      path: skills/<skill-name>
      tags: [<relevant-tags>]
      description: "Brief English description"
      description_zh: "简短的中文描述"
```

### Step 4: Add i18n Translations

Add to both `cmd/skills-x/i18n/locales/en.yaml` and `zh.yaml`:

```yaml
skill_<skill-name>: "Brief description"
```

### Step 5: Commit and Push

```bash
git add skills/<skill-name>/ pkg/registry/registry.yaml cmd/skills-x/i18n/locales/
git commit -m "feat: add <skill-name> skill"
git push origin main
```

**That’s it — no binary build or npm publish needed!**

### Step 6: Test Locally (After Pushing)

After pushing, verify the new skill is visible:

```bash
# Pull the latest registry from GitHub (REQUIRED to see new skills)
skills-x registry update

# Now the new skill should appear
skills-x list | grep "<skill-name>"

# Try installing it
skills-x init <skill-name> --target /tmp/test-skills
ls /tmp/test-skills/<skill-name>/
```

> ⚠️ **Without `skills-x registry update`, your local binary still shows the old embedded registry and the new skill will NOT appear.**

---

## Contributing Community Skills (Open Source Skills)

> **No binary release needed.** Just validate, update registry, push to git.

### Step 1: Find and Validate the Source Skill

Search for skills at:
- https://agentskills.io/
- https://github.com/anthropics/skills
- https://github.com/vercel-labs/agent-skills
- https://github.com/remotion-dev/skills
- Other GitHub repositories with proper skill structure

**Before adding to registry, verify:**
1. Repository has a valid `SKILL.md` in the skill directory
2. `SKILL.md` has proper YAML frontmatter (`name` + `description`)
3. Skill name matches directory name (lowercase, hyphens only)
4. License is identifiable

### Step 2: Add to Registry

Edit `pkg/registry/registry.yaml`:

```yaml
new-source-name:
  repo: github.com/owner/repo-name
  license: MIT
  skills:
    - name: skill-name
      path: path/to/skill/in/repo
      tags: [<relevant-tags>]
      description: "Brief English description"
      description_zh: "简短的中文描述"
```

### Step 3: Commit and Push

```bash
git add pkg/registry/registry.yaml
git commit -m "feat: add <skill-name> skill from <source>"
git push origin main
```

### Step 4: Test Locally (After Pushing)

```bash
# Pull the latest registry from GitHub (REQUIRED)
skills-x registry update

# Verify the skill appears
skills-x list | grep "<skill-name>"

# Test installing it
skills-x init <skill-name> --target /tmp/test-install
ls /tmp/test-install/<skill-name>/
```

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

> **Skills are fetched from GitHub at install time.** The binary only contains the registry index, not skill content. This means:
> - **Skills-only changes** (add/update skills in registry.yaml) → just commit & push, **no release needed**
> - **Tool changes** (Go code, CLI behavior, registry format) → full release required

### Path A: Skills-only Release (FAST - no version bump needed)

When ONLY `pkg/registry/registry.yaml` or skill content changes:

```bash
# Build and verify
make build
./bin/skills-x init --all --target "$(mktemp -d)"

# Commit and push — done!
git add pkg/registry/registry.yaml
git commit -m "feat: add <skill-name> skill

- Add <skill-name> to registry
- Users can update registry: skills-x registry update"
git push origin main
```

**That's it!** Users get the skill immediately by running:
```bash
skills-x registry update
skills-x list  # See new skill
skills-x init <skill-name>  # Install it
```

### Path B: Tool Release (FULL - version bump required)

When Go source code, CLI behavior, or registry format changes:

#### Step 1: Update Version

Increment version in `npm/package.json`:
```json
"version": "0.1.X"  // increment patch version
```

#### Step 2: Build for npm

```bash
make build-npm
```

#### Step 3: Pre-Release Testing (CRITICAL)

⚠️ **IMPORTANT: Always run this test before releasing to catch broken or missing skills!**

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

#### Step 4: Commit Changes

```bash
git add .
git commit -m "feat: add <skill-name> skill

- Add <skill-name> to skills collection
- Add i18n translations (en/zh)
- Update README"
```

#### Step 5: Tag and Push

```bash
git tag -a v0.1.X -m "Add <skill-name> skill"
git push origin main
git push --tags
```

#### Step 6: Create GitHub Release

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

#### Step 7: Publish to npm

```bash
cd npm && npm publish --access public
```

---

## Quick Reference

```bash
# ─── After adding a skill and pushing ───────────────────────────────────────
# Pull latest registry from GitHub (REQUIRED before testing new skills)
skills-x registry update

# List all skills (now shows newly added skills)
skills-x list
skills-x list | grep "<skill-name>"

# Install a skill
skills-x init <skill-name> --target /tmp/test-skills

# ─── Local dev build & verify ────────────────────────────────────────────────
make build
SKILLS_LANG=zh ./bin/skills-x list
SKILLS_LANG=en ./bin/skills-x list

# ─── Self-developed skill (自研) — no binary release needed ──────────────────
# 1. Create skills/<skill-name>/SKILL.md + LICENSE.txt
# 2. Add entry to pkg/registry/registry.yaml (under castle-x-skills-x)
# 3. Add i18n keys to en.yaml + zh.yaml
git add skills/<skill-name>/ pkg/registry/registry.yaml cmd/skills-x/i18n/locales/
git commit -m "feat: add <skill-name> skill"
git push origin main
# Users get it with: skills-x registry update

# ─── Third-party skill — no binary release needed ────────────────────────────
# 1. Validate SKILL.md in source repo
# 2. Add entry to pkg/registry/registry.yaml
git add pkg/registry/registry.yaml
git commit -m "feat: add <skill-name> skill from <source>"
git push origin main
# Users get it with: skills-x registry update

# ─── Full tool release (Go code changes only) ────────────────────────────────
make build-npm
./bin/skills-x init --all --target "$(mktemp -d)"  # Test all skills
git add . && git commit -m "feat: ..."
git tag -a v0.1.X -m "..."
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

| Skill Type | Storage | Description Source | Release needed? | How users get it |
|------------|---------|-------------------|-----------------|-----------------|
| **Self-Developed** (自研) | `skills/<name>/` in this repo | `registry.yaml` fields | ❌ No | `skills-x registry update` |
| **Third-Party** (open source) | External repo only | `registry.yaml` fields | ❌ No | `skills-x registry update` |
| **Tool change** (Go code) | n/a | n/a | ✅ Yes (full release) | Install new binary version |

---

## Checklists for New Skills

### For Self-Developed Skills (自研)

- [ ] Skill directory created at `skills/<name>/`
- [ ] `SKILL.md` with valid YAML frontmatter (`name`, `description`, `license`)
- [ ] `LICENSE.txt` present
- [ ] Entry added to `pkg/registry/registry.yaml` under `castle-x-skills-x` source
- [ ] `description` (English) and `description_zh` (Chinese) filled in registry.yaml
- [ ] `skill_<name>` key added to `en.yaml` and `zh.yaml`
- [ ] Committed and pushed to `origin/main`
- [ ] `skills-x registry update` run locally
- [ ] `skills-x list | grep "<name>"` shows the skill
- [ ] `skills-x init <name> --target /tmp/test` installs successfully

### For Third-Party / Open Source Skills

- [ ] Source repo has valid `SKILL.md` with frontmatter
- [ ] Skill name matches directory name in source repo
- [ ] Entry added to `pkg/registry/registry.yaml` with correct `repo`, `path`, `license`
- [ ] `description` (English) and `description_zh` (Chinese) filled in registry.yaml
- [ ] Committed and pushed to `origin/main`
- [ ] `skills-x registry update` run locally
- [ ] `skills-x list | grep "<name>"` shows the skill
- [ ] `skills-x init <name> --target /tmp/test` installs successfully

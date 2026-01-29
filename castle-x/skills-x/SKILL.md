---
name: skills-x
description: Guide for contributing new skills to the skills-x collection. This skill should be used when users want to add new open-source skills from external sources (like agentskills.io or anthropics/skills) to the skills-x repository. It covers the complete workflow from discovery to publishing.
license: MIT
metadata:
  author: castle-x
  version: "1.0"
---

# Skills-X Contribution Guide

This skill provides a standardized workflow for contributing new skills to the skills-x collection.

## When to Use This Skill

- Adding new skills from external sources (agentskills.io, anthropics/skills, etc.)
- Updating existing skills with new versions
- Validating skill format compliance before submission

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

## Contribution Workflow

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

Download the skill to the root `skills/` directory:

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

### Step 4: Update skills.go Metadata

Edit `cmd/skills-x/skills/skills.go` to add:

1. **Category mapping** in `skillCategories`:
```go
var skillCategories = map[string]string{
    // ... existing entries
    "new-skill-name": "category",  // creative/document/devtools/workflow/etc.
}
```

2. **Description** in `skillDescriptions`:
```go
var skillDescriptions = map[string]string{
    // ... existing entries
    "new-skill-name": "Brief description for list display",
}
```

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

### Step 5: Update README Files

Update both `README.md` and `README_ZH.md`:

1. Update skill count if changed
2. Add new skill to the appropriate category table
3. Keep both files in sync

### Step 6: Build and Test

```bash
# Build the binary
make build

# Verify new skill appears in list
./bin/skills-x list | grep "<skill-name>"

# Test downloading the skill
./bin/skills-x init <skill-name> --target /tmp/test-skills
ls /tmp/test-skills/<skill-name>/
```

### Step 7: Update Version and Publish

1. **Increment version** in `npm/package.json`:
```json
"version": "0.1.X"  // increment patch version
```

2. **Build for npm**:
```bash
make build-npm
```

3. **Commit changes**:
```bash
git add .
git commit -m "feat: add <skill-name> skill

- Add <skill-name> to skills collection
- Update skills.go metadata
- Update README"
```

4. **Tag and push**:
```bash
git tag -a v0.1.X -m "Add <skill-name> skill"
git push origin main
git push --tags
```

5. **Create GitHub Release**:
```bash
make build-all
gh release create v0.1.X \
  --title "v0.1.X - Add <skill-name>" \
  --notes "## Added
- New skill: <skill-name>
- Description: <brief description>" \
  bin/skills-x-*
```

6. **Publish to npm**:
```bash
cd npm && npm publish --access public
```

## Quick Reference Commands

```bash
# Validate a skill
head -20 skills/<name>/SKILL.md

# Build and test locally
make build && ./bin/skills-x list

# Full release workflow
make build-npm
git add . && git commit -m "feat: add <skill>"
git tag -a v0.1.X -m "Add <skill>"
git push origin main --tags
make build-all
gh release create v0.1.X --title "v0.1.X" --notes "..." bin/skills-x-*
cd npm && npm publish --access public
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Skill not in list | Check `skillCategories` and `skillDescriptions` in skills.go |
| init fails | Verify SKILL.md exists and has valid frontmatter |
| Windows fails | Ensure using `/` not `\` for embed.FS paths |
| Version mismatch | Check `npm/package.json` version matches build |

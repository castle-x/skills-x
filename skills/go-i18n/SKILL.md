---
name: go-i18n
description: Use when defining or updating Go CLI i18n rules in this repo, especially around locale files, env-based language selection, and key naming.
---

# Go CLI i18n

## Overview

This skill documents the Go CLI i18n rules used in this repo, including language detection, locale files, and key naming conventions.

## When to Use

- You add or edit any user-facing CLI string in `cmd/skills-x`
- You need to add a new skill description for list output
- You are troubleshooting mixed language output or missing translations

## Architecture Snapshot

- **i18n runtime:** `cmd/skills-x/i18n/i18n.go`
- **Locale files:** `cmd/skills-x/i18n/locales/zh.yaml`, `cmd/skills-x/i18n/locales/en.yaml`
- **Embed:** `//go:embed locales/*.yaml` (locales are embedded at build time)

## Language Selection

Priority order in `detectLanguage()`:

1. `SKILLS_LANG`
2. `LANG`
3. `LC_ALL`
4. Default: `zh`

Normalization uses prefix matching:
- `zh`, `zh_CN`, `zh_TW` → `zh`
- `en`, `en_US` → `en`
- Unsupported → `zh`

## Using i18n in Go

Initialize once (early in `main`):

```go
import "github.com/castle-x/skills-x/cmd/skills-x/i18n"

i18n.MustInit()
```

Translate strings:

```go
title := i18n.T("list_header")
msg := i18n.Tf("init_success", skillName)
```

Missing keys return the key itself (useful for spotting gaps).

## Key Naming Conventions

Use lowercase + underscores with category prefixes:

- `app_` app metadata
- `cmd_` command descriptions/flags
- `list_` list output
- `init_` init output
- `err_` error messages
- `cat_` category names
- `skill_` skill descriptions
- `meta_` meta warnings

**Rule:** Never mix Chinese and English in the same string. Use two locale keys instead.

## Adding New Strings

1. Add the same key to both locale files:
   - `cmd/skills-x/i18n/locales/zh.yaml`
   - `cmd/skills-x/i18n/locales/en.yaml`
2. Use `i18n.T` / `i18n.Tf` in Go code.
3. Rebuild (`make build`) to embed the updated YAML.
4. Test:
   - `SKILLS_LANG=zh ./bin/skills-x list`
   - `SKILLS_LANG=en ./bin/skills-x list`

### New Skill Description

Add a `skill_<skill-name>` key to both locale files. The list view reads descriptions from i18n, not from the skill file.

## Example

```yaml
# locales/en.yaml
init_success: "Installed: %s"

# locales/zh.yaml
init_success: "安装成功: %s"
```

```go
fmt.Println(i18n.Tf("init_success", skillName))
```

## Quick Reference

| Item | Value |
|------|-------|
| Locale files | `cmd/skills-x/i18n/locales/{zh,en}.yaml` |
| Init | `i18n.MustInit()` |
| Translate | `i18n.T`, `i18n.Tf` |
| Language priority | `SKILLS_LANG > LANG > LC_ALL > zh` |

## Common Mistakes

- Only adding a key in one locale file
- Hardcoding mixed-language strings
- Forgetting to rebuild after YAML changes
- Formatting with `fmt.Sprintf` before calling `i18n.Tf`

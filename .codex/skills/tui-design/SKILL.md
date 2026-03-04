---
name: tui-design
description: Design specification for CLI TUI (Terminal User Interface). This skill provides comprehensive guidelines for implementing interactive terminal UI components, including page layout structure, color schemes, keyboard navigation, and multi-level navigation principles.
license: MIT
metadata:
  author: x
  version: "1.2"
---

# CLI TUI Design Specification

This document provides comprehensive design guidelines for implementing interactive terminal user interface (TUI) applications.

## When to Use This Skill

- Developing or modifying CLI TUI components
- Adding new UI pages or interactions
- Implementing new visual styles or themes
- Creating similar TUI applications following this pattern

---

## 1. Page Layout Structure

Every TUI page follows a consistent four-section layout:

```
┌────────────────────────────────────────────────────────────────┐
│  Header: Title Line + Separator + ASCII Logo + Separator       │
├────────────────────────────────────────────────────────────────┤
│  Info Bar: Context-dependent information (Dynamic)             │
├────────────────────────────────────────────────────────────────┤
│  Content Area: Interactive options/list (Scrollable)           │
│                                                                │
├────────────────────────────────────────────────────────────────┤
│  Footer: Status Bar + Separator + Keyboard Hints               │
└────────────────────────────────────────────────────────────────┘
```

### Section Details

#### 1.1 Header Section (Always Visible)

The header consists of a title line, separator, ASCII logo, and separator with **symmetric spacing**:

```
🚀 Skills-X - AI Agent Skills Manager  v0.2.12
────────────────────────────────────────────────────────────
                                                  <- 1 blank line (from logo constant leading \n)
   ███████╗██╗  ██╗██╗██╗     ...
   ╚══════╝╚═╝  ╚═╝╚═╝╚══════╝...
                                                  <- 1 blank line (from logo constant trailing \n + 1 extra \n)
────────────────────────────────────────────────────────────
```

**Spacing rules:**

- The `logo` constant uses backtick string, which naturally includes a leading `\n` (after opening backtick) and trailing `\n` (before closing backtick)
- After `logoStyle.Render(logo)`, add exactly **one** `\n` — this produces symmetric 1-blank-line padding on both sides
- **Never** add `\n\n` after logo — it creates asymmetric bottom padding

```go
func RenderLogo(version string) string {
    var b strings.Builder
    // Title + version (strip -dirty suffix for display)
    b.WriteString(accentStyle.Render("🚀 Skills-X - AI Agent Skills Manager"))
    if version != "" {
        displayVersion := strings.TrimSuffix(version, "-dirty")
        b.WriteString("  " + hintStyle.Render(displayVersion))
    }
    b.WriteString("\n")
    // Upper separator
    b.WriteString(separatorStyle.Render(strings.Repeat("─", SeparatorWidth)))
    b.WriteString("\n")
    // Logo (constant has leading \n and trailing \n built-in)
    b.WriteString(logoStyle.Render(logo))
    b.WriteString("\n")  // exactly 1 \n — symmetric with leading \n
    // Lower separator
    b.WriteString(separatorStyle.Render(strings.Repeat("─", SeparatorWidth)))
    b.WriteString("\n")
    return b.String()
}
```

**Newline budget between header and content:**

- `RenderLogo` ends with `separator + \n`
- View should **not** add extra `\n` after `RenderLogo()` — content starts immediately

#### 1.2 Info Bar Section (Dynamic)

- Page-specific context: table headers, counts, search box
- Context info uses dim/hint color
- **No working directory display on LV1** — only shown on pages where it adds value (e.g. LV3 skills list)

#### 1.3 Content Area (Interactive)

- **Table View**: Column headers + data rows with visual-width alignment
- **List View**: Scrollable list with checkbox multi-selection
- **Progress View**: Progress bar with percentage and counts

Key features:
- Visual-width-aware column alignment (see Section 3)
- Numeric columns right-aligned with `fmt.Sprintf("%4d")`
- Checkbox (☐/☑) for multi-selection
- Cursor indicator (❯) for current position
- Empty line padding to maintain stable page height

#### 1.4 Footer Section

```
(共 30 个，已安装 12 个，将安装 3 个，将卸载 1 个)

────────────────────────────────────────────────────────────
/ 搜索 | 空格 选择/反选 | A 全选 | Enter 确认安装/卸载 | b 返回 | q 退出
```

**Composition:**

1. **Status bar**: `\n` + counts in dim color, parenthesized
2. **Separator**: `\n` + 60-char `─` line + `\n`
3. **Hints**: Keyboard shortcuts in dim color, pipe-separated

```go
func RenderStatusBar(total, installed, newSelected, deselected int) string {
    return "\n" + hintStyle.Render(
        fmt.Sprintf("(共 %d 个，已安装 %d 个，将安装 %d 个，将卸载 %d 个)",
            total, installed, newSelected, deselected))
}

func RenderHint(hint string) string {
    return "\n" + RenderSeparator() + hintStyle.Render(hint)
}

func RenderSeparator() string {
    return separatorStyle.Render(strings.Repeat("─", SeparatorWidth)) + "\n"
}
```

---

## 2. Style Definitions

### Color Palette

| Color Name | Hex Code  | Variable         | Usage                    |
|------------|-----------|------------------|--------------------------|
| Primary    | `#00D4AA` | `primaryColor`   | Selected items, logo     |
| Secondary  | `#FF6B6B` | `secondaryColor` | Cursor, errors, alerts   |
| Accent     | `#4ECDC4` | `accentColor`    | Title line, highlights   |
| Dim        | `#666666` | `dimColor`       | Hints, separators        |
| White      | `#FFFFFF` | `whiteColor`     | Section titles, headers  |
| Yellow     | `#FFCC00` | `yellowColor`    | Warnings (reserved)      |
| Blue       | `#5599FF` | `blueColor`      | Links (reserved)         |
| Cyan       | `#5A9FB8` | `cyanColor`      | Selectable items         |

### Component Styles

```go
logoStyle       = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
titleStyle      = lipgloss.NewStyle().Foreground(whiteColor).Bold(true)
selectedStyle   = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
normalStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))
selectableStyle = lipgloss.NewStyle().Foreground(cyanColor)
hintStyle       = lipgloss.NewStyle().Foreground(dimColor)
cursorStyle     = lipgloss.NewStyle().Foreground(secondaryColor).Bold(true)
successStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
separatorStyle  = lipgloss.NewStyle().Foreground(dimColor)
searchStyle     = lipgloss.NewStyle().Foreground(whiteColor).Background(lipgloss.Color("#333333"))
```

### Status Indicators

| Symbol | Meaning         | Style         |
|--------|-----------------|---------------|
| `☑`    | Selected        | selectedStyle |
| `☐`    | Not selected    | —             |
| `❯`    | Cursor position | cursorStyle   |

---

## 3. Column Alignment & Visual Width

### The Problem

`fmt.Sprintf` width specifiers count bytes, not terminal columns. This breaks alignment when:
- **ANSI escape codes** (lipgloss styles) add invisible bytes
- **CJK characters** (中文、日本語) occupy 2 terminal columns per rune

### The Solution: Pad First, Style After

**Rule: Always pad plain text to the target visual width, then apply lipgloss style.**

```go
// termWidth returns the visual width of a string in terminal columns.
func termWidth(s string) int {
    w := 0
    for _, r := range s {
        if isCJK(r) {
            w += 2
        } else {
            w++
        }
    }
    return w
}

// padRight pads a string with spaces to reach the target visual width.
func padRight(s string, targetWidth int) string {
    w := termWidth(s)
    for w < targetWidth {
        s += " "
        w++
    }
    return s
}
```

### Usage Pattern

```go
// Header: pad plain text, then style
headerName := padRight("AI 工具", 22)   // "AI 工具" = 7 visual cols → 15 spaces added
b.WriteString(fmt.Sprintf("  %s  %s  %s\n",
    titleStyle.Render(headerName),           // styled AFTER padding
    titleStyle.Render(padRight("全局", 4)),   // "全局" = 4 visual cols → 0 spaces
    titleStyle.Render(padRight("项目", 4))))

// Data rows: same column widths
name := padRight(p.Name, 22)             // "Claude Code" = 11 cols → 11 spaces
globalStr := fmt.Sprintf("%4d", count)   // right-aligned number, 4 chars
b.WriteString(fmt.Sprintf("%s%s  %s  %s\n",
    prefix,                               // "  " or "❯ " (always 2 visual cols)
    style.Render(name),                   // styled AFTER padding
    globalCountStyle.Render(globalStr),   // styled AFTER formatting
    projectCountStyle.Render(projectStr)))
```

### Column Width Standards

| Column       | Visual Width | Alignment | Padding Method                |
|--------------|-------------|-----------|-------------------------------|
| Name         | 22          | Left      | `padRight(text, 22)`          |
| Skill name   | 40          | Left      | `padRight(text, 40)`          |
| Number       | 4           | Right     | `fmt.Sprintf("%4d", n)`       |
| Status (✓/✗) | 1           | Center    | Direct render                 |
| Prefix (❯)   | 2           | Left      | `"  "` or `cursorStyle("❯ ")` |

### Gap Between Columns

Use 2 spaces (`"  "`) between columns in `fmt.Sprintf`:

```go
fmt.Sprintf("%s%s  %s  %s\n", prefix, name, col1, col2)
//                ^^    ^^  two-space gaps
```

---

## 4. Separator Rules

All separators use the same constant width:

```go
const SeparatorWidth = 60
```

| Location             | Character | Width | Color    |
|----------------------|-----------|-------|----------|
| Above ASCII logo     | `─`       | 60    | dimColor |
| Below ASCII logo     | `─`       | 60    | dimColor |
| Below table header   | `─`       | 60    | dimColor |
| Above footer hints   | `─`       | 60    | dimColor |

**All separators use `SeparatorWidth` (60).** Never hardcode different widths.

---

## 5. Keyboard Navigation

### Universal Keys (All Levels)

| Key      | Action           |
|----------|------------------|
| `q`      | Quit application |
| `Ctrl+C` | Force quit       |

### List Navigation

| Key      | Action                          |
|----------|---------------------------------|
| `↑`      | Move cursor up (wrap to bottom) |
| `↓`      | Move cursor down (wrap to top)  |
| `PgUp`   | Page up                         |
| `PgDown` | Page down                       |
| `Home`   | Jump to first item              |
| `End`    | Jump to last item               |

### Selection

| Key       | Action                               |
|-----------|--------------------------------------|
| `Space`   | Toggle selection (cursor stays)      |
| `a` / `A` | Select all visible items             |
| `Enter`   | Confirm selection                    |

### Multi-Level Navigation

| Key          | Action                      |
|--------------|-----------------------------|
| `Esc` / `b`  | Go back to previous level   |
| `Enter`      | Proceed to next level       |

### Search Mode

| Key             | Action                        |
|-----------------|-------------------------------|
| `/` / `Ctrl+F`  | Enter search mode             |
| Typing          | Append to search query        |
| `Backspace`     | Remove last character         |
| `Esc`           | Exit search mode, clear query |

---

## 6. Multi-Level Page Architecture

### Page Flow

```
Level 1           Level 2             Level 3            Level 4
  │                  │                   │                  │
  ▼                  ▼                   ▼                  ▼
Select ────────► Select ──────────► Select ──────────► Install
AI Tool          Target (global/     Skills (multi-     Progress
                 project)            select + search)
```

### Navigation Rules

1. **Forward**: `Enter` confirms and proceeds to next level
2. **Backward**: `Esc` or `b` returns to previous level (restarts the loop)
3. **Quit**: `q` or `Ctrl+C` exits at any level
4. **State preservation**: Going back re-runs the previous level from scratch

### Orchestration Pattern

Alt screen is managed **once** at the top level to avoid flicker between pages:

```go
func RunTUI(opts TUIOptions) error {
    // Enter alt screen ONCE for the entire TUI session
    fmt.Print(EnterAltScreen)
    fmt.Print(HideCursor)
    defer func() {
        fmt.Print(ShowCursor)
        fmt.Print(ExitAltScreen)
    }()
    return runTUIFlow(opts)
}

func runTUIFlow(opts TUIOptions) error {
    // Level 1
    fmt.Print(ClearScreen)
    product, err := RunProductSelect(opts.Version, opts.TargetDir)
    if err != nil || product == nil {
        return err
    }

    // Level 2
    fmt.Print(ClearScreen)
    target, err := RunInstallTargetSelect(product, opts.TargetDir)
    ...

    // Level 3
    fmt.Print(ClearScreen)
    selected, deselected, err := RunSkillsSelect(...)
    if err != nil && err.Error() == "go back" {
        return runTUIFlow(opts)  // recursive restart
    }

    // Level 4
    fmt.Print(ClearScreen)
    RunInstaller(selected, deselected, targetDir)

    // Exit alt screen to show final result on normal terminal
    fmt.Print(ShowCursor)
    fmt.Print(ExitAltScreen)
    fmt.Printf("Completed: %d\n", completed)
    return nil
}
```

---

## 7. Terminal Control

### Two-Layer Rendering Strategy

The rendering uses a two-layer approach to eliminate flicker:

**Layer 1: Alt screen lifecycle (managed by `RunTUI`)**
- Enter alt screen once at the start, exit once at the end
- `ClearScreen` between each page transition (within the alt screen buffer)
- **Never** use `tea.WithAltScreen()` on individual programs — it causes exit/enter cycles that flash the normal terminal

**Layer 2: Per-page rendering (managed by Bubble Tea)**
- Each page runs as `tea.NewProgram(model)` (no WithAltScreen)
- Bubble Tea handles incremental rendering within the page
- `View()` returns complete page content; **never** include `ClearScreen` in View()

```go
// Individual programs: NO WithAltScreen
p := tea.NewProgram(model)  // Bubble Tea renders incrementally

// Page transitions: ClearScreen within the shared alt screen
fmt.Print(ClearScreen)      // clears alt screen buffer between pages
```

**Anti-patterns:**
- `tea.NewProgram(model, tea.WithAltScreen())` per level → flicker on every page switch
- `ClearScreen` inside `View()` → full repaint on every keystroke
- `fmt.Printf(...)` between levels while in alt screen → cursor position corruption

### ANSI Sequences

```go
const ClearScreen    = "\033[2J\033[H"    // Clear screen + cursor home
const EnterAltScreen = "\033[?1049h"      // Enter alternate screen buffer
const ExitAltScreen  = "\033[?1049l"      // Exit alternate screen buffer
const HideCursor     = "\033[?25l"        // Hide cursor
const ShowCursor     = "\033[?25h"        // Show cursor
```

---

## 8. Layout Examples

### Level 1: Product Selection

```
🚀 Skills-X - AI Agent Skills Manager  v0.2.12
────────────────────────────────────────────────────────────

   ███████╗██╗  ██╗██╗██╗     ██╗     ███████╗     ██╗  ██╗
   ...
   ╚══════╝╚═╝  ╚═╝╚═╝╚══════╝╚══════╝╚══════╝     ╚═╝  ╚═╝

────────────────────────────────────────────────────────────
  AI 工具                 全局  项目
────────────────────────────────────────────────────────────
❯ Claude Code               45     0
  Cursor                      0    36
  Windsurf                    0     0

────────────────────────────────────────────────────────────
↑↓ 选择 | Enter 确定 | q 退出
```

### Level 3: Skills Selection

```
🚀 Skills-X - AI Agent Skills Manager  v0.2.12
────────────────────────────────────────────────────────────

   ███████╗...
   ╚══════╝...

────────────────────────────────────────────────────────────

选择要安装的 Skills  Cursor
📁 /path/to/project

 🔍 搜索技能...
                                          已安装
────────────────────────────────────────────────────────────
❯ ☑ anthropic/brainstorming                ✓
  ☑ anthropic/doc-coauthoring              ✓
  ☐ anthropic/frontend-design              ✗
  ...

↑ 3/30 ↓

(共 30 个，已安装 12 个，将安装 3 个，将卸载 1 个)

────────────────────────────────────────────────────────────
/ 搜索 | 空格 选择/反选 | A 全选 | Enter 确认安装/卸载 | b 返回 | q 退出
```

### Level 4: Installation Progress

```
Installing Skills

Progress: [========================================] 100% (9/9)
Completed: 9 | Failed: 0

To uninstall:
   [OK] anthropic/doc-coauthoring
   [OK] anthropic/frontend-design

────────────────────────────────────────────────────────────
Done! 9 completed.
```

---

## 9. Constants

```go
const SeparatorWidth  = 60  // All separators use this width
const DefaultPageSize = 10  // Default items per page in scrollable lists
```

---

## 10. Design Principles Summary

1. **Consistent Layout**: Every page has Header → Info Bar → Content → Footer
2. **Uniform Separators**: All `─` lines use `SeparatorWidth` (60), never hardcode other values
3. **Symmetric Spacing**: Logo has equal blank-line padding above and below
4. **Visual-Width Alignment**: Use `padRight`/`termWidth` for CJK-aware column alignment; pad plain text first, style after
5. **Flicker-Free**: Manage alt screen at RunTUI level (enter once, exit once), use ClearScreen between pages, never ClearScreen in View()
6. **Stable UI**: Fixed column widths, empty line fill, scroll offset tracking
7. **Multi-Level Flow**: Independent programs per level, back/forward navigation

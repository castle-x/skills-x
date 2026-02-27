# TUI Guide

Running `skills-x` with no arguments launches the interactive TUI.

---

## Level 1 — Select AI Tool

```
Skills-X  AI Agent Skills Manager  v0.2.15
────────────────────────────────────────────────────────────

  ███████╗██╗  ██╗██╗██╗     ██╗     ███████╗     ██╗  ██╗
  ██╔════╝██║ ██╔╝██║██║     ██║     ██╔════╝    ╚██╗██╔╝
  ███████╗█████╔╝ ██║██║     ██║     ███████╗█████╗╚███╔╝
  ╚════██║██╔═██╗ ██║██║     ██║     ╚════██║╚════╝██╔██╗
  ███████║██║  ██╗██║███████╗███████╗███████║      ██╔╝ ██╗
  ╚══════╝╚═╝  ╚═╝╚═╝╚══════╝╚══════╝╚══════╝      ╚═╝  ╚═╝

────────────────────────────────────────────────────────────
AI Tool              Global  Project
────────────────────────────────────────────────────────────
❯ Claude Code            45        3
  Cursor                  0        0
  Windsurf                0        0
  Codex                   2        0
  Kimi                    0        0
  CodeBuddy               0        0
  Roo Code                0        0
  Opencode                0        0
  Aider                   0        0
────────────────────────────────────────────────────────────
↑/↓ 选择  |  Enter 确定  |  q 退出
```

Each row shows **Global** (home config dir) and **Project** (current working dir) installed skill count.

`↑`/`↓` move cursor · `Enter` select · `q` quit

---

## Level 2 — Select Install Target

```
Skills For Claude Code
────────────────────────────────────────────────────────────
❯ 全局安装    ~/.claude/skills       (45 installed)
  项目安装    /home/user/myproject   (3 installed)
────────────────────────────────────────────────────────────
↑/↓ 选择  |  Enter 确定  |  b 返回  |  q 退出
```

Choose between global (home config dir) or project-level installation.

`↑`/`↓` move cursor · `Enter` confirm · `b` go back · `q` quit

---

## Level 3 — Browse & Select Skills

```
Skills For Cursor  /home/user/.cursor/skills
────────────────────────────────────────────────────────────
 输入 / 搜索技能
────────────────────────────────────────────────────────────
[ ]未安装  [●]已安装  [+]安装  [-]卸载  [↑]更新  [↻]检测中
────────────────────────────────────────────────────────────
❯ [●] anthropic/algorithmic-art        2026-01-15  ★
  [●] anthropic/brand-guidelines       2026-01-15
  [+] anthropic/canvas-design
  [-] superpowers/brainstorming
  [↑] superpowers/test-driven-development  ⚠ 有新版
  [ ] superpowers/writing-plans
  ...

1/53 ↓
Install: 1 | Update: 1 | Uninstall: 1
────────────────────────────────────────────────────────────
Space 选择  f 收藏  u 检测更新  A 全选  Enter 确定  b 返回  q 退出
```

**Selection**

| Key | Behavior |
|-----|----------|
| `Space` | Toggle: uninstalled → `[+]install` · installed → `[-]uninstall` |
| `A` | Cycle: all install/update → all uninstall → reset |
| `Enter` | Confirm and proceed to installation |

**Search & Filter**

| Key | Behavior |
|-----|----------|
| `/` or `、` | Enter search mode — type skill name |
| `#` | Open tag picker — navigate tags with `↑`/`↓`, `Enter` to apply |
| `Esc` | Exit search or cancel tag picker |

Available tags: `#starred` `#featured` `#ai-efficiency` `#planning` `#web-frontend` `#mobile` `#backend` `#testing` `#code-review` `#office` `#design` `#writing` `#media` `#skills`

**Skill Actions**

| Key | Behavior |
|-----|----------|
| `f` | Star / unstar focused skill — starred skills sort to the top, saved in `~/.config/skills-x/starred.json` |
| `u` | Check for updates on the focused installed skill — shows commit comparison in status area |

---

## Level 4 — Installation Progress

```
Installing Skills

Progress: [===================>                    ] 45% (5/11)
Completed: 5 | Failed: 0

  ✓ anthropic/skill-creator
  ✓ superpowers/brainstorming
  ✓ superpowers/test-driven-development
  ✓ anthropic/brand-guidelines
  ✓ vercel/react-best-practices
  ...

Press any key to exit
```

Real-time progress with per-skill `✓` / `✗` result. The screen persists until you press any key.

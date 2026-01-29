# Skills-X

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue)](https://github.com/castle-x/skills-x)
[![npm](https://img.shields.io/npm/v/skills-x?color=CB3837&logo=npm)](https://www.npmjs.com/package/skills-x)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

我的 AI Agent Skills 个人收藏，方便查找和使用。

## 快速安装

```bash
# 使用 npm（推荐）
npm install -g skills-x

# 或直接使用 npx
npx skills-x list

# 使用 go install
go install github.com/castle-x/skills-x/cmd/skills-x@latest
```

### 更新

```bash
# 更新到最新版本
npm update -g skills-x

# 或重新安装
npm install -g skills-x@latest
```

### 使用方法

```bash
# 查看所有可用 skills
skills-x list

# 下载全部 skills
skills-x init --all

# 下载指定 skill
skills-x init pdf
skills-x init ui-ux-pro-max

# 指定自定义目标目录
skills-x init pdf --target .claude/skills

# 强制覆盖已存在的 skills（跳过确认）
skills-x init pdf -f
skills-x init --all --force
```

### AI IDE 配置

我们仅提供 skills 下载，请指定你的 AI IDE 的 skills 目录：

```bash
# Claude Code
skills-x init --all --target .claude/skills

# CodeBuddy
skills-x init --all --target .codebuddy/skills

# Cursor
skills-x init --all --target .cursor/skills

# Windsurf
skills-x init --all --target .windsurf/skills
```

---

## 目录结构

```
skills-x/
├── castle-x/       # 🏰 我的原创 Skills
├── skills/         # 收录的开源 Skills
└── references/     # Skills 编写参考资料
```

## 来源

| 项目 | 链接 |
|------|------|
| Claude Official | [anthropics/skills](https://github.com/anthropics/skills) |
| Superpowers | [obra/superpowers](https://github.com/obra/superpowers) |
| Awesome Claude Skills | [ComposioHQ/awesome-claude-skills](https://github.com/ComposioHQ/awesome-claude-skills) |

---

## Skills 目录 (53个)

### 🏰 Castle-X (作者自研 Skills)
| Skill | 用途 |
|-------|------|
| `skills-x` ⭐ | 向 skills-x 集合贡献新 skill |

### 🎨 创意设计
| Skill | 用途 |
|-------|------|
| `ui-ux-pro-max` | **UI/UX 设计智能** - 67风格/96色板/57字体/25图表/13技术栈，设计系统生成 |
| `algorithmic-art` | p5.js 生成艺术、流场、粒子系统 |
| `canvas-design` | 海报、视觉艺术 (.png/.pdf) |
| `brand-guidelines` | Anthropic 品牌风格 |
| `theme-factory` | 工件主题切换 (10种预设) |
| `frontend-design` | 前端设计 |
| `image-enhancer` | 图像放大、锐化、清理 |
| `remotion` | **AI 视频编程** - 使用 React + Remotion 制作视频 |

### 📄 文档处理
| Skill | 用途 |
|-------|------|
| `pdf` | PDF 提取/填写/合并 |
| `docx` | Word 文档 |
| `pptx` | PPT 演示文稿 |
| `xlsx` | Excel 表格/公式/图表 |
| `document-skills` | 综合文档处理 |
| `doc-coauthoring` | 文档协作 |

### 🛠️ 开发工具
| Skill | 用途 |
|-------|------|
| `mcp-builder` | 构建 MCP 服务器 |
| `artifacts-builder` | React+Tailwind+shadcn 工件 |
| `web-artifacts-builder` | 复杂 HTML 工件 |
| `webapp-testing` | Playwright 测试 |
| `langsmith-fetch` | LangSmith 调试追踪 |
| `changelog-generator` | Git 提交生成更新日志 |

### 🔄 工作流程
| Skill | 用途 |
|-------|------|
| `brainstorming` | 创意工作前头脑风暴 |
| `writing-plans` | 编写任务计划 |
| `executing-plans` | 执行计划 |
| `systematic-debugging` | 系统化调试 |
| `test-driven-development` | TDD 流程 |
| `verification-before-completion` | 完成前验证 |
| `subagent-driven-development` | 子代理开发 |
| `dispatching-parallel-agents` | 并行代理调度 |

### 📝 Git & 代码审查
| Skill | 用途 |
|-------|------|
| `requesting-code-review` | 请求 CR |
| `receiving-code-review` | 处理 CR 反馈 |
| `finishing-a-development-branch` | 完成分支 |
| `using-git-worktrees` | Git Worktree 隔离 |

### ✍️ 写作
| Skill | 用途 |
|-------|------|
| `content-research-writer` | 内容研究写作 |
| `internal-comms` | 内部沟通/报告 |
| `tailored-resume-generator` | 定制简历 |

### 🔗 集成
| Skill | 用途 |
|-------|------|
| `connect` | 连接 1000+ 服务 |
| `connect-apps` | Gmail/Slack/GitHub 等 |
| `connect-apps-plugin` | 应用连接插件 |
| `slack-gif-creator` | Slack GIF |

### 📊 商业分析
| Skill | 用途 |
|-------|------|
| `competitive-ads-extractor` | 竞品广告分析 |
| `developer-growth-analysis` | 开发者增长 |
| `lead-research-assistant` | 客户研究 |
| `meeting-insights-analyzer` | 会议分析 |
| `twitter-algorithm-optimizer` | 推文优化 |

### 🗂️ 文件管理
| Skill | 用途 |
|-------|------|
| `file-organizer` | 文件整理 |
| `invoice-organizer` | 发票整理/报税 |

### 🎲 实用工具
| Skill | 用途 |
|-------|------|
| `video-downloader` | YouTube 下载 |
| `domain-name-brainstormer` | 域名灵感 |
| `raffle-winner-picker` | 抽奖 |

### 🧰 Skills 开发
| Skill | 用途 |
|-------|------|
| `skill-creator` | 创建 Skill |
| `writing-skills` | 编写/验证 Skill |
| `skill-share` | 分享 Skill |
| `template-skill` | Skill 模板 |
| `using-superpowers` | 如何使用 Skills |

---

## 参考资料

- `references/claude_official_skills.md` - **Claude 官方 Skills 文档**（概述/快速开始/最佳实践/企业级）
- `references/agent_skills.md` - Agent Skills 开放规范 (agentskills.io)
- `references/agents_md.md` - **AGENTS.md 规范**（面向 AI 的项目说明标准）
- `references/superpower.md` - Superpowers 文档

---

## 🏰 关于 Castle-X Skills

Castle-X skills 是**项目作者的原创作品**，存放在 `castle-x/` 目录下。
它们在列表中单独显示，并带有 ⭐ 标记以区分社区 skills。

要贡献新的 Castle-X skill，请参考 `skills-x` skill 的指南。

# Skills-X

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue)](https://github.com/castle-x/skills-x)
[![npm](https://img.shields.io/npm/v/skills-x?color=CB3837&logo=npm)](https://www.npmjs.com/package/skills-x)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

> **提示**: 业内已有成熟的 Agent Skills 生态，推荐使用 [skills.sh](https://skills.sh/) 和 Vercel Labs 的 [`npx skills`](https://github.com/vercel-labs/add-skill)。
> 本项目只是我的个人收藏，学习使用。

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
├── skills/         # 🏰 原创 Skills（嵌入到二进制中）
└── pkg/registry/   # Skill 来源注册表
```

## 收藏的 Skills

> 仅供学习使用

```
$ skills-x list

📦 github.com/anthropics/skills (Apache-2.0)
   algorithmic-art      使用 p5.js 创建算法艺术...
   artifacts-builder    使用 React 构建交互式 artifacts
   brand-guidelines     应用 Anthropic 品牌色彩和排版
   canvas-design        创建 PNG 和 PDF 视觉艺术...
   doc-coauthoring      协作文档编辑
   docx                 Word 文档创建、编辑和分析
   frontend-design      前端设计最佳实践
   internal-comms       内部沟通模板
   mcp-builder          生成 MCP 服务器
   pdf                  PDF 操作 - 提取、填写表单、合并
   pptx                 PowerPoint 演示文稿创建和编辑
   skill-creator        创建新的 agent skills
   slack-gif-creator    创建针对 Slack 优化的动画 GIF
   template-skill       创建新 skills 的模板
   theme-factory        使用主题样式化 artifacts 的工具包
   web-artifacts-builder 使用 React 构建 Web artifacts
   webapp-testing       测试 Web 应用程序
   xlsx                 Excel 电子表格创建、公式和图表

📦 github.com/remotion-dev/skills (Remotion License)
   remotion             Remotion 最佳实践 - 使用 React 创建视频

📦 github.com/vercel-labs/agent-skills (MIT)
   react-best-practices      React 和 Next.js 性能优化指南
   react-native-guidelines   AI agents 的 React Native 最佳实践
   web-design-guidelines     100+ 条可访问性、性能和用户体验规则

📦 skills-x (Original)
   skills-x             向 skills-x 贡献 skills

共 23 个 skills，来自 4 个源
```

---

## 参考资料

- [Claude 官方 Skills 文档](https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/skills) - Anthropic 官方文档
- [Agent Skills 规范](https://agentskills.io/) - Agent Skills 开放规范
- [AGENTS.md](https://agents.md/) - 面向 AI 的项目说明标准
- [Superpowers](https://github.com/obra/superpowers) - 编码代理的完整开发工作流

---

## 🏰 关于 X Skills

X skills 是**项目作者的原创作品**，存放在 `x/` 目录下。
它们在列表中单独显示，并带有 ⭐ 标记以区分社区 skills。

要贡献新的 X skill，请参考 `skills-x` skill 的指南。

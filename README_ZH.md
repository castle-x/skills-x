# Skills-X

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue)](https://github.com/castle-x/skills-x)
[![npm](https://img.shields.io/npm/v/skills-x?color=CB3837&logo=npm)](https://www.npmjs.com/package/skills-x)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

> **提示**：业内已有成熟的 Agent Skills 生态，推荐使用 [skills.sh](https://skills.sh/) 和 Vercel Labs 的 [`npx skills`](https://github.com/vercel-labs/add-skill)。
> 本项目是我的个人精选收藏。

[English](README.md)

## 快速安装

```bash
# 使用 npm（推荐）
npm install -g skills-x

# 使用 go install
go install github.com/castle-x/skills-x/cmd/skills-x@latest
```

```bash
# 更新到最新版本
npm install -g skills-x@latest
```

---

## 交互式 TUI（默认模式）

直接运行 `skills-x` 即可进入交互式 TUI，这是浏览、安装、更新、卸载 skills 的推荐方式。

```bash
skills-x
```

**主要特性：**
- 四级页面导航：选择 IDE → 安装位置 → 浏览技能 → 安装进度
- 按名称搜索（`/`）或按标签筛选（`#`），支持交互式标签选择器
- 收藏技能（`f`）— 持久保存至 `~/.config/skills-x/starred.json`，始终排列在列表最前
- 检测更新（`u`）— 显示已安装技能的版本对比信息
- 双语界面 — 根据 `SKILLS_LANG` 自动切换中英文

> 详细的页面截图与操作说明见 [docs/tui-guide-zh.md](docs/tui-guide-zh.md)

---

## CLI 参考

```bash
# 查看所有可用 skills（非交互模式）
skills-x list

# 安装指定 skill
skills-x init pdf
skills-x init pdf frontend-design

# 安装全部 skills
skills-x init --all

# 安装到指定目录
skills-x init pdf --target .cursor/skills

# 强制覆盖已存在的 skills
skills-x init pdf --force
skills-x init --all --force

# 检测可用更新（不执行更新）
skills-x update --check

# 更新指定 skill
skills-x update pdf

# 更新全部已安装 skills
skills-x update --all
```

### 各 IDE 目标目录

```bash
# Claude Code（全局）
skills-x init --all --target ~/.claude/skills

# Cursor（项目级）
skills-x init --all --target .cursor/skills

# Windsurf
skills-x init --all --target .windsurf/skills

# CodeBuddy
skills-x init --all --target .codebuddy/skills

# Codex
skills-x init --all --target ~/.codex/skills
```

### 语言切换

```bash
# 切换为英文
SKILLS_LANG=en skills-x

# 切换为中文（默认）
SKILLS_LANG=zh skills-x
```

---

## 收藏的 Skills（最新总数请运行 `skills-x list` 查看）

```
📦 github.com/anthropics/skills (Apache-2.0)
   algorithmic-art           使用 p5.js 创建算法艺术
   brand-guidelines          应用 Anthropic 品牌色彩和排版
   canvas-design             创建 PNG 和 PDF 视觉艺术
   doc-coauthoring           协作文档编辑
   docx                      Word 文档创建和编辑
   frontend-design           前端设计最佳实践
   internal-comms            内部沟通模板
   mcp-builder               生成 MCP 服务器
   pdf                       PDF 操作 - 提取、填写表单、合并
   pptx                      PowerPoint 演示文稿创建和编辑
   skill-creator             创建新的 agent skills
   slack-gif-creator         创建针对 Slack 优化的动画 GIF
   theme-factory             使用主题样式化 artifacts
   web-artifacts-builder     使用 React 构建 Web artifacts
   webapp-testing            测试 Web 应用程序
   xlsx                      Excel 电子表格创建和公式

📦 github.com/obra/superpowers (MIT)
   brainstorming                   创意工作前的头脑风暴
   dispatching-parallel-agents     分派独立任务并行处理
   executing-plans                 逐步执行书面实施计划
   finishing-a-development-branch  测试通过后完成开发分支
   receiving-code-review           实施前处理代码审查反馈
   requesting-code-review          完成任务时请求代码审查
   subagent-driven-development     使用子代理执行独立任务
   systematic-debugging            系统化调试 bug 和异常
   test-driven-development         TDD 工作流 - 先写测试再实现
   using-superpowers               如何有效使用 superpowers 技能
   verification-before-completion  完成前验证工作
   writing-plans                   编写多步骤任务实施计划
   writing-skills                  创建、编辑和验证 skills

📦 github.com/vercel-labs/agent-skills (MIT)
   react-best-practices   React 和 Next.js 性能优化指南
   react-native-skills    React Native 最佳实践
   web-design-guidelines  100+ 条可访问性、性能和用户体验规则

📦 github.com/vercel-labs/skills (MIT)
   find-skills   使用 Skills CLI 查找并安装技能

📦 github.com/vercel/next.js (MIT)
   cache-components  Next.js 缓存组件最佳实践（PPR 与缓存工作流）
   update-docs       基于源码变更更新 Next.js 文档

📦 github.com/google-gemini/gemini-cli (Apache-2.0)
   code-reviewer  审查本地改动或 PR，并输出结构化反馈
   pr-creator     基于模板与预检流程创建高质量 Pull Request

📦 github.com/langgenius/dify (Other)
   frontend-code-review  面向前端的代码审查（质量、性能、业务逻辑）

📦 github.com/Shubhamsaboo/awesome-llm-apps (MIT)
   fullstack-developer  使用 React/Next.js、Node.js 与数据库进行全栈开发

📦 github.com/facebook/react (MIT)
   fix  在 CI 前执行格式化与 lint 检查，快速修复代码质量问题

📦 github.com/giuseppe-trisciuoglio/developer-kit (MIT)
   react-patterns         React 19 模式 - 服务器组件、Actions、hooks
   shadcn-ui              shadcn/ui 组件库 - Radix UI + Tailwind CSS
   tailwind-css-patterns  Tailwind CSS 实用优先样式模式
   typescript-docs        TypeScript 文档 - JSDoc 和 TypeDoc

📦 github.com/affaan-m/everything-claude-code (MIT)
   golang-testing   Go 测试模式 - 表格驱动测试、基准、模糊测试、TDD

📦 github.com/bobmatnyc/claude-mpm-skills (MIT)
   golang-cli-cobra-viper  Go CLI 开发 - Cobra + Viper 命令结构、配置管理与 Shell 补全

📦 github.com/masayuki-kono/agent-skills (MIT)
   implementation-plan  编写实施/开发方案的指南

📦 github.com/tencentcloudbase/skills (MIT)
   ai-model-wechat          小程序 AI 模型 - 混元和 DeepSeek，支持流式响应
   auth-wechat              微信小程序认证 - 自动注入 OPENID/UNIONID
   cloudbase-guidelines     腾讯云开发指南
   miniprogram-development  微信小程序开发 - 免登录认证、AI模型、部署发布
   no-sql-wx-mp-sdk         小程序文档数据库 - 增删改查、复杂查询

📦 github.com/nextlevelbuilder/ui-ux-pro-max-skill (MIT)
   ui-ux-pro-max   UI/UX 设计智能 - 50 种风格、21 色板

📦 github.com/axtonliu/axton-obsidian-visual-skills (MIT)
   excalidraw-diagram  根据文本生成 Excalidraw 图

📦 github.com/remotion-dev/skills (Remotion License)
   remotion   Remotion 最佳实践 - 使用 React 创建视频

📦 github.com/bfollington/terma (CC-BY-SA-4.0)
   strudel   Strudel 现场编程音乐：节奏、旋律与可分享链接

📦 github.com/castle-x/skills-x (MIT)
   baidu-speech-to-text  百度语音识别 - 语音转文本（国内环境优化）
   go-embedded-spa       Go 内嵌 SPA（单二进制部署）
   go-i18n               Go CLI 多语言规则（作者自用）
   minimal-ui-design     极简 UI 设计 - 低噪声、图标优先
   skills-x              向 skills-x 贡献 skills
   tui-design            TUI 终端交互界面设计规范
```

---

## 参考资料

- [Claude 官方 Skills 文档](https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/skills) — Anthropic 官方文档
- [Agent Skills 规范](https://agentskills.io/) — Agent Skills 开放规范
- [Superpowers](https://github.com/obra/superpowers) — 编码代理的完整开发工作流

---

## 🏰 关于 X Skills

X skills 是**项目作者的原创 Skills**，维护在 `skills/` 目录，并通过内置注册表源（`github.com/castle-x/skills-x`）分发。
它们现在与其他注册表 skills 使用同一套安装/更新流程。

要贡献新的 X skill，请参考 `skills-x` skill 的指南。若需私有分叉或自定义来源，可使用 `skills-x registry add`。

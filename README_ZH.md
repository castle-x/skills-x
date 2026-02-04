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
   algorithmic-art           使用 p5.js 创建算法艺术...
   brand-guidelines          应用 Anthropic 品牌色彩和排版
   canvas-design             创建 PNG 和 PDF 视觉艺术...
   doc-coauthoring           协作文档编辑
   docx                      Word 文档创建、编辑和分析
   frontend-design           前端设计最佳实践
   internal-comms            内部沟通模板
   mcp-builder               生成 MCP 服务器
   pdf                       PDF 操作 - 提取、填写表单、合并
   pptx                      PowerPoint 演示文稿创建和编辑
   skill-creator             创建新的 agent skills
   slack-gif-creator         创建针对 Slack 优化的动画 GIF
   theme-factory             使用主题样式化 artifacts 的工具包
   web-artifacts-builder     使用 React 构建 Web artifacts
   webapp-testing            测试 Web 应用程序
   xlsx                      Excel 电子表格创建、公式和图表

📦 github.com/remotion-dev/skills (Remotion License)
   remotion             Remotion 最佳实践 - 使用 React 创建视频

📦 github.com/vercel-labs/skills (MIT)
   find-skills          使用 Skills CLI 查找并安装技能

📦 github.com/vercel-labs/agent-skills (MIT)
   react-best-practices  React 和 Next.js 性能优化指南
   react-native-skills   AI agents 的 React Native 最佳实践
   web-design-guidelines 100+ 条可访问性、性能和用户体验规则

📦 github.com/giuseppe-trisciuoglio/developer-kit (MIT)
   react-patterns        React 19 模式 - 服务器组件、Actions、hooks、Suspense 和 TypeScript
   shadcn-ui             shadcn/ui 组件库 - 基于 Radix UI 和 Tailwind CSS
   tailwind-css-patterns Tailwind CSS 实用优先样式模式 - 响应式设计
   typescript-docs       TypeScript 文档 - JSDoc、TypeDoc 和 ADR 模式

📦 github.com/masayuki-kono/agent-skills (MIT)
   implementation-plan   编写实施/开发方案的指南

📦 github.com/obra/superpowers (MIT)
   brainstorming                   创意工作前的头脑风暴
   dispatching-parallel-agents     分派 2+ 个独立任务并行处理
   executing-plans                 逐步执行书面实施计划
   finishing-a-development-branch  测试通过后完成开发分支
   receiving-code-review           在实施更改前处理代码审查反馈
   requesting-code-review          完成任务或功能时请求代码审查
   subagent-driven-development     使用子代理执行独立任务的计划
   systematic-debugging            系统化调试 - 处理 bug、测试失败、异常行为
   test-driven-development         TDD 工作流 - 先写测试再实现
   using-superpowers               如何有效查找和使用 superpowers 技能
   verification-before-completion  完成前验证工作已完成、修复或通过
   writing-plans                   根据规格编写多步骤任务的实施计划
   writing-skills                  创建、编辑和验证 skills

📦 github.com/tencentcloudbase/skills (MIT)
   ai-model-wechat          小程序 AI 模型 - 混元和 DeepSeek，支持流式响应
   auth-wechat              微信小程序认证 - 自动注入 OPENID/UNIONID
   cloudbase-guidelines     腾讯云开发指南 - Web、小程序、后端服务开发规范
   miniprogram-development  微信小程序开发 - 免登录认证、AI模型集成、部署发布
   no-sql-wx-mp-sdk          小程序文档数据库 - 增删改查、复杂查询、分页、地理位置

📦 github.com/nextlevelbuilder/ui-ux-pro-max-skill (MIT)
   ui-ux-pro-max         UI/UX 设计智能 - 50 种风格、97 色板、57 字体、9 技术栈

📦 skills-x (Original)
   baidu-speech-to-text       百度语音识别 - 语音转文本（国内环境优化）
   go-embedded-spa            Go 内嵌 SPA（单二进制部署）
   go-i18n                    Go CLI 多语言规则（作者自用）
   minimal-ui-design          极简 UI 设计 - 低噪声、图标优先
   newapi-deploy-config       New API 部署与模型配置（Host 网络）
   openclaw-session-header-fix 修复 openclaw session 覆写问题
   skills-x                   向 skills-x 贡献 skills

共 52 个 skills，来自 10 个源
```

---

## 参考资料

- [Claude 官方 Skills 文档](https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/skills) - Anthropic 官方文档
- [Agent Skills 规范](https://agentskills.io/) - Agent Skills 开放规范
- [AGENTS.md](https://agents.md/) - 面向 AI 的项目说明标准
- [Superpowers](https://github.com/obra/superpowers) - 编码代理的完整开发工作流

---

## 🏰 关于 X Skills

X skills 是**项目作者的原创 Skill**，存放在 `skills/` 目录下，并已适配业内 Git 仓库的 skill 规范（使用 skills 目录）。
它们在列表中单独显示，并带有 ⭐ 标记以区分社区 skills。

要贡献新的 X skill，请参考 `skills-x` skill 的指南。
创建新 skill 后，需要询问是否添加 `REAEDME.md` 背景总结文档，用于记录问题、复盘方案，并沉淀为后续可复用的工具。

# AGENTS.md 规范

> 来源: https://agents.md | https://github.com/agentsmd/agents.md
> 一种专门为 AI 编程代理设计的项目说明文档标准

---

## 一、什么是 AGENTS.md？

AGENTS.md 是一种简单、开放的标准格式，旨在为 AI 编程代理提供指导文档。**已被超过 60,000 个开源项目使用**。

- **核心定位**：把它想象成**专门给 AI 看的 README 文件**
- **目的**：为 AI 编程代理提供一个 dedicated（专门）且 predictable（可预测）的地方，以提供上下文和指令

### 与 README.md 的区别

| 文件 | 目标读者 | 内容特点 |
|:---|:---|:---|
| `README.md` | 人类开发者 | 快速入门、项目描述、贡献指南 |
| `AGENTS.md` | AI 代理 | 构建步骤、测试细节、代码规范、详细指令 |

### 为什么需要 AGENTS.md？

- **分离关注点**：将给 AI 的详细指令与给人类看的文档分开，保持 README 简洁
- **提供详细上下文**：包含对人类来说可能过于冗余的技术细节
- **通用性**：采用开放格式，不绑定任何特定供应商

---

## 二、兼容性

AGENTS.md 被广泛的 AI 编码工具和生态系统所支持：

| 工具/平台 | 支持情况 |
|:---|:---|
| OpenAI Codex | ✅ |
| Google Jules / Gemini CLI | ✅ |
| Cursor | ✅ |
| GitHub Copilot (Coding agent) | ✅ |
| Aider | ✅ |
| Devin (Cognition) | ✅ |
| VS Code | ✅ |
| Warp | ✅ |
| Zed | ✅ |

**只需维护一份 AGENTS.md 文件，即可被多种 AI 工具识别和使用。**

---

## 三、规范说明

### 1. 文件放置规则

- **位置**：将文件放置在仓库的**根目录**
- **文件名**：`AGENTS.md`（大写）

### 2. Monorepo 支持

对于大型代码仓库，可以在子项目目录中放置嵌套的 `AGENTS.md` 文件：

```
project/
├── AGENTS.md              # 根目录规则
├── packages/
│   ├── frontend/
│   │   └── AGENTS.md      # frontend 专属规则
│   └── backend/
│       └── AGENTS.md      # backend 专属规则
```

**优先级规则**：
- 离被编辑文件**最近**的 `AGENTS.md` 拥有最高优先级
- 用户在聊天中输入的显式指令会**覆盖**文件中的所有内容

### 3. 内容规范

使用标准的 Markdown 格式，**没有必填字段**，可根据项目需求自定义。

**推荐的章节：**

| 章节 | 说明 |
|:---|:---|
| Setup commands | 依赖安装、启动开发服务器、运行测试 |
| Code style | 语言模式、引号风格、分号使用、设计模式 |
| Testing instructions | 如何运行测试、测试通过标准 |
| Project overview | 简要介绍项目背景 |
| Dev environment tips | 特定于该项目的工具链使用技巧 |
| Security considerations | 安全相关注意事项 |
| PR/Commit 指南 | 提交信息或 PR 的格式要求 |

---

## 四、示例

### 基础示例

```markdown
# AGENTS.md

## Dev environment tips
- Use `pnpm dlx turbo run where <project_name>` to jump to a package.
- Run `pnpm install --filter <project_name>` to add the package to workspace.
- Use `pnpm create vite@latest <project_name> -- --template react-ts` for new React + Vite package.

## Testing instructions
- Find CI plan in `.github/workflows` folder.
- Run `pnpm turbo run test --filter <project_name>` to run checks.
- Fix any test or type errors until the whole suite is green.
- Add or update tests for the code you change, even if nobody asked.

## PR instructions
- Title format: `[<project_name>] <Title>`
- Always run `pnpm lint` and `pnpm test` before committing.
```

### Next.js 项目示例

```markdown
# AGENTS Guidelines

## 1. Use the Development Server, **not** `npm run build`

* **Always use `npm run dev`** while iterating on the application.
* **Do _not_ run `npm run build` inside the agent session.** It disables hot reload.

## 2. Keep Dependencies in Sync

If you add or update dependencies:
1. Update the appropriate lockfile.
2. Re-start the development server.

## 3. Coding Conventions

* Prefer TypeScript (`.tsx`/`.ts`) for new components.
* Co-locate component-specific styles in the same folder.

## 4. Useful Commands

| Command | Purpose |
|---------|---------|
| `npm run dev` | Start dev server with HMR |
| `npm run lint` | Run ESLint checks |
| `npm run test` | Execute the test suite |
| `npm run build` | **Production build – do not run during agent sessions** |
```

---

## 五、工具配置

部分工具可能需要简单配置：

### Aider

在 `.aider.conf.yml` 中：
```yaml
read: AGENTS.md
```

### Gemini CLI

在 `.gemini/settings.json` 中：
```json
{
  "contextFileName": "AGENTS.md"
}
```

---

## 六、组织与管理

- **开源协作**：源于 AI 软件开发生态系统的协作努力（OpenAI, Google, Cursor 等）
- **监管机构**：由 Linux Foundation 旗下的 **Agentic AI Foundation** 托管和维护

---

## 七、与 Agent Skills 的对比

| 特性 | AGENTS.md | Agent Skills (SKILL.md) |
|:---|:---|:---|
| **定位** | 项目级 AI 指导文档 | 可复用的技能包 |
| **作用范围** | 单个项目/仓库 | 跨项目通用 |
| **内容** | 项目构建/测试/风格指南 | 特定任务的详细指令 |
| **触发方式** | AI 自动读取 | 按需加载/激活 |
| **复杂度** | 简单（纯 Markdown） | 可包含脚本/资源 |
| **使用场景** | 告诉 AI 如何在这个项目工作 | 教 AI 完成某类任务 |

**两者可以互补使用**：
- `AGENTS.md` 描述项目特定的上下文
- `SKILL.md` 提供通用的任务能力

---

## 八、相关资源

- **官网**: https://agents.md
- **GitHub**: https://github.com/agentsmd/agents.md
- **托管组织**: Agentic AI Foundation (Linux Foundation)

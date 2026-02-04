# Claude 官方 Agent Skills 文档

> 整理自 Claude 官方文档和博客
> 来源: platform.claude.com, code.claude.com, claude.com/blog/skills

---

## 一、概述 (Overview)

### 什么是 Agent Skills？

Skills 是 Claude 的一项功能，允许它通过加载特定的文件夹来执行特定任务。这些文件夹包含了**指令、脚本和资源**。Claude 仅在处理与当前任务相关的内容时才会访问这些技能。

简而言之，Skills 就像是让 Claude 获取专业知识的"入职材料"，将其塑造成特定领域的专家。

### 核心特点

1. **可组合性**：技能可以叠加使用，Claude 会自动识别需要哪些技能并协调使用
2. **可移植性**：技能在所有平台上使用相同的格式，构建一次即可在 Claude 应用、Claude Code 和 API 之间通用
3. **高效性**：仅在需要时加载所需的最少信息和文件，保持运行速度
4. **强大功能**：技能可以包含可执行代码，用于处理传统编程比生成 Token 更可靠的任务

### 本质理解

**Skills 本质上不是可执行代码**，而是专业的**提示词模板**。它们通过向对话上下文中注入特定领域的指令来改变 Claude 的行为。

- **Skill Tool（大写 S）**：元工具，管理所有技能
- **skills（小写 s）**：具体的技能（如 pdf, skill-creator）

---

## 二、快速开始 (Quickstart)

### 创建第一个 Skill

#### 1. 创建目录

```bash
# 个人技能（所有项目通用）
mkdir -p ~/.claude/skills/explain-code

# 项目技能（仅此项目）
mkdir -p .claude/skills/explain-code
```

#### 2. 编写 SKILL.md

文件由 YAML 前置元数据和 Markdown 说明组成：

```markdown
---
name: explain-code
description: 使用视觉图表和类比解释代码。用于解释代码如何工作、教授代码库或用户询问"这是如何工作的？"时。
---

解释代码时，始终包括：
1. **以类比开始**：将代码与日常生活进行比较
2. **绘制图表**：使用 ASCII 艺术展示流程、结构或关系
3. **逐步演练**：解释每一步发生了什么
4. **强调陷阱**：常见的错误或误解是什么？
```

### 技能存放位置

| 位置 | 路径 | 适用范围 |
|:---|:---|:---|
| **个人** | `~/.claude/skills/<skill-name>/SKILL.md` | 你的所有项目 |
| **项目** | `.claude/skills/<skill-name>/SKILL.md` | 仅此项目 |
| **插件** | `<plugin>/skills/<skill-name>/SKILL.md` | 启用插件的位置 |
| **嵌套** | `packages/frontend/.claude/skills/` | 单仓库中的特定包 |

*注：项目技能会覆盖同名的个人技能。*

---

## 三、最佳实践 (Best Practices)

### SKILL.md 结构

每个技能定义在 `SKILL.md` 文件中，包含两部分：

#### Frontmatter（YAML 头部）

配置如何运行：

| 字段 | 说明 |
|:---|:---|
| `name` | 技能名称（仅小写字母、数字、连字符） |
| `description` | 技能的作用及何时使用（**关键字段**，帮助 Claude 决定何时调用） |
| `disable-model-invocation` | 设为 `true` 防止自动调用，仅允许手动触发 |
| `user-invocable` | 设为 `false` 从 `/` 菜单中隐藏 |
| `allowed-tools` | 限制技能激活时可使用的工具（如 `Read, Grep`） |
| `context` | 设为 `fork` 在子代理中运行 |
| `agent` | 指定子代理类型（如 `Explore`, `Plan`） |
| `model` | 指定使用的模型（默认继承会话模型） |

#### Markdown 内容

具体的指令，告诉 Claude 做什么。

### 资源捆绑结构

为保持提示词简洁，技能支持捆绑外部资源：

```
my-skill/
├── SKILL.md          # 必需：指令 + 元数据
├── scripts/          # Python/Bash 脚本，由 Claude 通过 Bash 执行
├── references/       # 文档，Claude 通过 Read 工具读取
└── assets/           # 模板和二进制文件，仅引用路径不加载内容
```

### 高级功能

#### 参数传递

使用 `$ARGUMENTS` 占位符接收命令行传入的参数：

```markdown
运行 `/fix-issue 123`，内容中的 `$ARGUMENTS` 会被替换为 `123`
```

#### 动态上下文注入

使用 `` !`command` `` 语法在发送给 Claude 之前执行 Shell 命令：

```markdown
---
name: pr-summary
description: 总结拉取请求的更改
---

## PR 上下文
- PR 差异: !`gh pr diff`
- PR 评论: !`gh pr view --comments`

## 你的任务
总结此拉取请求...
```

### 关键设计模式

| 模式 | 说明 |
|:---|:---|
| **脚本自动化** | 将复杂任务交给脚本处理 |
| **读取-处理-写入** | 最基础的文件转换模式 |
| **搜索-分析-报告** | 用于代码库分析 |
| **向导式多步工作流** | 在每一步等待用户确认 |
| **模板生成** | 基于 Assets 中的模板生成结构化输出 |

### 故障排除

| 问题 | 解决方案 |
|:---|:---|
| 技能未触发 | 检查 `description` 是否包含用户常说的关键词；尝试直接用 `/skill-name` 调用 |
| 触发过于频繁 | 将描述写得更具体，或添加 `disable-model-invocation: true` |
| Claude 看不到所有技能 | 技能描述受上下文字符预算限制（默认 15,000 字符），可通过 `SLASH_COMMAND_TOOL_CHAR_BUDGET` 环境变量增加 |

---

## 四、企业级功能 (Enterprise)

### 产品支持

Skills 已集成到 Claude 的全线产品中：

#### Claude 应用
- **适用对象**：Pro, Max, Team 和 Enterprise 用户
- **功能**：Claude 会根据任务自动调用相关技能，用户无需手动选择
- **注意**：Team 和 Enterprise 用户需要管理员先在组织范围内启用该功能

#### Claude 开发者平台 (API)
- 通过 Messages API 请求添加技能
- 使用 `/v1/skills` 端点进行版本控制和管理
- 技能运行需要"代码执行工具"测试版提供的 secure environment
- 支持读取和生成 Excel、PowerPoint、Word 文档和可填写的 PDF

#### Claude Code
- 通过插件从 `anthropics/skills` 市场安装技能
- Claude 会自动加载相关技能
- 团队可以通过版本控制共享技能
- 也可手动安装到 `~/.claude/skills`

### 组织级管理

- **组织范围内的技能管理**：管理员可以统一管理和分发技能
- **合作伙伴构建的技能目录**：提供预构建的专业技能
- **开放标准**：Agent Skills 发布为跨平台可移植的开放标准

### 安全建议

由于该功能赋予 Claude 执行代码的权限，应注意：
- 只使用受信任的技能来源
- 审核技能中的脚本内容
- 使用 `allowed-tools` 限制权限

---

## 五、内部架构原理

### 技能选择机制

- 完全基于 LLM 的**推理能力**，而非算法匹配
- 代码层面没有嵌入、分类器或正则匹配
- 所有可用技能的描述会被格式化并放入 Skill Tool 的描述中
- Claude 读取后根据语义理解决定调用哪个技能

### 执行流程

1. **发现与加载**：系统从用户设置、项目设置、插件等路径扫描技能
2. **注入上下文**：
   - **对话上下文**：通过插入新的用户消息注入详细指令
   - **执行上下文**：修改工具权限和模型选择

### 双消息注入机制

为平衡用户透明度和指令清晰度，技能执行时会注入两条消息：

1. **可见消息**：包含简短的 XML 元数据，用户可见，提供状态指示
   ```xml
   <command-message>The "pdf" skill is loading</command-message>
   ```

2. **隐藏消息**：包含完整的技能提示词（500-5000 字），用户不可见，发送给 API 以指导 Claude 行为

---

## 六、合作伙伴案例

| 合作伙伴 | 使用场景 |
|:---|:---|
| **Box** | 将存储的文件转换为符合组织标准的 PPT、Excel 和 Word 文档 |
| **Canva** | 定制代理并扩展其功能，使设计工作流更加自动化 |
| **Notion** | 减少在复杂任务上的提示词调整，提供更可预测的结果 |
| **管理会计团队** | 自动化处理多个电子表格、捕捉异常并生成报告，将一天的工作缩短至一小时 |

---

## 七、相关资源

- **官方博客**: https://claude.com/blog/skills
- **GitHub 示例**: https://github.com/anthropics/skills
- **Claude Code 文档**: https://code.claude.com/docs/skills
- **开放规范**: https://agentskills.io/

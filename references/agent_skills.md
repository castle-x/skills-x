# Agent Skills 完整指南

> 来源: https://agentskills.io/

---

## 一、概述 (Overview)

### 什么是 Agent Skills？

Agent Skills 是一种简单、开放的格式，用于赋予 AI 智能体新的能力和专业知识。它由包含指令、脚本和资源的文件夹组成，智能体可以发现并使用这些内容，从而更准确、高效地完成任务。

### 为什么需要 Agent Skills？

尽管智能体的能力日益增强，但它们往往缺乏可靠完成实际工作所需的上下文。Agent Skills 通过以下方式解决这个问题：

- **按需加载**：为智能体提供程序性知识以及公司、团队和用户特定的上下文
- **扩展能力**：智能体可以根据当前任务加载不同的技能集来扩展自身功能
- **多方受益**：
  - **技能开发者**：一次构建能力，即可部署到多种智能体产品上
  - **兼容的智能体**：支持最终用户直接赋予智能体开箱即用的新能力
  - **团队与企业**：将组织知识捕获为可移植、可版本控制的软件包

### Agent Skills 能实现什么？

1. **领域专业知识**：将专业知识打包为可复用的指令，例如从法律审查流程到数据分析管道
2. **新能力**：赋予智能体全新的能力（例如创建演示文稿、构建 MCP 服务器、分析数据集）
3. **可重复的工作流**：将多步骤任务转化为一致且可审计的工作流
4. **互操作性**：在不同的、兼容技能的智能体产品中复用同一个技能

### 开发与采用情况

- **开源开发**：该标准最初由 Anthropic 开发，现已作为开放标准发布
- **广泛支持**：Agent Skills 得到了领先 AI 开发工具的支持，生态系统开放，欢迎社区贡献

---

## 二、什么是 Skills (What are Skills?)

### 核心定义

Agent Skills 是一种轻量级、开放的标准格式，旨在通过专门的知识和工作流程来扩展 AI 智能体的能力。

### 基本结构

在最核心的层面上，一个"技能"就是一个包含 `SKILL.md` 文件的文件夹。

```text
my-skill/
├── SKILL.md          # 必需：指令 + 元数据
├── scripts/          # 可选：可执行代码
├── references/       # 可选：文档
└── assets/           # 可选：模板、资源
```

- **SKILL.md** (必需)：包含元数据（名称、描述）和告诉智能体如何执行特定任务的指令
- **scripts/** (可选)：可执行代码
- **references/** (可选)：参考文档
- **assets/** (可选)：模板或资源

### 运作机制：渐进式披露 (Progressive Disclosure)

技能使用"渐进式披露"技术来高效管理上下文，确保智能体保持快速响应。该机制分为三个阶段：

1. **发现 (Discovery)**：启动时，智能体仅加载每个可用技能的**名称**和**描述**。这足以让智能体知道该技能何时可能与当前任务相关。
2. **激活 (Activation)**：当任务与技能的描述匹配时，智能体会将完整的 `SKILL.md` 指令读入上下文中。
3. **执行 (Execution)**：智能体遵循指令，并根据需要加载引用的文件或执行捆绑的代码。

### 关键优势

- **自文档化**：作者或用户可以直接阅读 `SKILL.md` 来理解其功能，便于审核和改进
- **可扩展**：技能的复杂度可以从纯文本指令扩展到包含可执行代码、资源和模板
- **可移植**：技能只是文件，因此易于编辑、版本控制和共享

---

## 三、规范 (Specification)

### 1. 目录结构

一个技能至少包含一个名为 `SKILL.md` 的文件：

```text
skill-name/
└── SKILL.md          # 必需文件
```

可选目录：`scripts/`、`references/`、`assets/`

### 2. SKILL.md 文件格式

`SKILL.md` 文件必须包含 **YAML 前置元数据**，后跟 **Markdown 内容**。

#### 最小示例：

```yaml
---
name: skill-name
description: A description of what this skill does and when to use it.
---
```

#### 包含可选字段的示例：

```yaml
---
name: pdf-processing
description: Extract text and tables from PDF files, fill forms, merge documents.
license: Apache-2.0
metadata:
  author: example-org
  version: "1.0"
---
```

### 3. 字段详细说明

| 字段 | 是否必需 | 约束条件 |
| :--- | :--- | :--- |
| **name** | 是 | 最多 64 个字符。仅限小写字母、数字和连字符。不能以连字符开头或结尾。必须与父目录名称匹配。 |
| **description** | 是 | 最多 1024 个字符。非空。描述技能的功能以及何时使用它。 |
| **license** | 否 | 许可证名称或对捆绑许可证文件的引用。 |
| **compatibility** | 否 | 最多 500 个字符。指示环境要求（如目标产品、系统包、网络访问等）。 |
| **metadata** | 否 | 任意键值映射，用于存储额外的元数据。 |
| **allowed-tools** | 否 | 预批准可用工具的空格分隔列表。（实验性功能） |

### 4. 字段规则详解

#### name 字段
- 长度：1-64 个字符
- 字符：仅限 Unicode 小写字母数字字符和连字符
- 禁止：以 `-` 开头或结尾，包含连续连字符 (`--`)
- 必须与父目录名一致

**有效示例：** `pdf-processing`, `data-analysis`, `code-review`

**无效示例：** `PDF-Processing` (含大写), `-pdf` (以连字符开头), `pdf--processing` (含连续连字符)

#### description 字段
- 长度：1-1024 个字符
- 内容：应描述技能的功能和使用场景，并包含有助于代理识别相关任务的关键词

**好例子：** "Extracts text and tables from PDF files, fills PDF forms, and merges multiple PDFs. Use when working with PDF documents or when the user mentions PDFs, forms, or document extraction."

**坏例子：** "Helps with PDFs."

### 5. 可选目录说明

#### scripts/
包含代理可以运行的可执行代码：
- 应自包含或清晰记录依赖项
- 包含有用的错误消息
- 优雅地处理边缘情况

#### references/
包含代理在需要时可以阅读的额外文档：
- `REFERENCE.md` - 详细的技术参考
- `FORMS.md` - 表单模板或结构化数据格式
- 特定领域的文件（如 `finance.md`, `legal.md` 等）

#### assets/
包含静态资源：
- 模板（文档模板、配置模板）
- 图片（图表、示例）
- 数据文件（查找表、模式）

### 6. 渐进式披露与上下文管理

为了高效利用上下文，技能应按以下结构组织：

1. **元数据**（约 100 个令牌）：所有技能在启动时加载 `name` 和 `description`
2. **指令**（建议少于 5000 个令牌）：激活技能时加载完整的 `SKILL.md` 正文
3. **资源**（按需）：仅在需要时加载文件

**建议：** 保持主 `SKILL.md` 在 500 行以内。将详细的参考资料移至单独的文件。

### 7. 文件引用

引用技能中的其他文件时，请使用相对于技能根目录的路径：

```markdown
See [the reference guide](references/REFERENCE.md) for details.

Run the extraction script:
scripts/extract.py
```

保持文件引用在 `SKILL.md` 之下的一级深度。避免深层嵌套的引用链。

### 8. 验证

使用 `skills-ref` 参考库来验证技能：

```bash
skills-ref validate ./my-skill
```

---

## 四、集成 Skills (Integrate Skills)

### 1. 集成方法概览

主要有两种集成 Agent Skills 的方式：

- **基于文件系统的代理:** 在计算机环境（bash/unix）中运行，功能最强大。当模型发出 shell 命令（如 `cat /path/to/my-skill/SKILL.md`）时激活技能，并通过 shell 命令访问打包的资源。
- **基于工具的代理:** 在没有专用计算机环境的情况下运行。通过实现工具来允许模型触发技能和访问打包资产。

### 2. 集成步骤概述

一个兼容技能的智能体需要完成以下 5 个步骤：

1. **发现技能:** 在配置的目录中扫描
2. **加载元数据:** 在启动时解析每个技能的名称和描述
3. **匹配任务:** 将用户任务与相关技能进行匹配
4. **激活技能:** 通过加载完整指令来激活技能
5. **执行与访问:** 根据需要执行脚本并访问资源

### 3. 技能发现

- 技能是包含 `SKILL.md` 文件的文件夹
- 智能体应扫描配置的目录以查找有效的技能

### 4. 加载与解析元数据

启动时仅解析 `SKILL.md` 文件的前言，以保持初始上下文使用量较低。

```python
function parseMetadata(skillPath):
    content = readFile(skillPath + "/SKILL.md")
    frontmatter = extractYAMLFrontmatter(content)

    return {
        name: frontmatter.name,
        description: frontmatter.description,
        path: skillPath
    }
```

### 5. 注入上下文

将技能元数据包含在系统提示中，以便模型知道有哪些可用技能。

**推荐格式 (针对 Claude) - 使用 XML 格式：**

```xml
<available_skills>
  <skill>
    <name>pdf-processing</name>
    <description>Extracts text and tables from PDF files...</description>
    <location>/path/to/skills/pdf-processing/SKILL.md</location>
  </skill>
  ...
</available_skills>
```

**注意：**
- 对于基于文件系统的代理，必须包含 `location` 字段（SKILL.md 的绝对路径）
- 对于基于工具的代理，可以省略
- 尽量保持元数据简洁，每个技能增加约 50-100 个 token

### 6. 安全考虑

脚本执行会带来安全风险，建议采取以下措施：

- **沙箱化:** 在隔离环境中运行脚本
- **白名单:** 仅执行来自受信任技能的脚本
- **确认机制:** 在运行潜在危险操作之前询问用户
- **日志记录:** 记录所有脚本执行以供审计

### 7. 参考实现

使用 `skills-ref` 库，它提供了用于处理技能的 Python 实用程序和 CLI：

```bash
# 验证技能目录
skills-ref validate <path>

# 生成代理提示词 XML
skills-ref to-prompt <path>...
```

---

## 五、SKILL.md 模板

```markdown
---
name: my-skill-name
description: 描述这个技能做什么以及何时使用它。包含关键词以帮助智能体识别相关任务。
license: MIT
metadata:
  author: your-name
  version: "1.0"
---

# 技能名称

## 何时使用此技能

描述触发此技能的场景和关键词...

## 工作流程

### 步骤 1: ...

详细说明...

### 步骤 2: ...

详细说明...

## 示例

### 输入示例

...

### 输出示例

...

## 注意事项

- 边缘情况处理
- 常见错误
- 最佳实践
```

---

## 六、相关资源

- **官方网站**: https://agentskills.io/
- **GitHub 示例**: https://github.com/anthropics/skills
- **参考库**: skills-ref (Python)

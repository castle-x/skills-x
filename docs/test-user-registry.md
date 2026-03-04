# 用户本地注册表模块 — 测试用例文档

> 测试仓库：`affaan-m/everything-claude-code`
> 测试文件：
> - `pkg/skillvalidator/validator_test.go`
> - `pkg/userregistry/userregistry_test.go`

---

## 总览

| 包 | 测试函数 | 子用例 | 覆盖率 | 类型 |
|---|---|---|---|---|
| `pkg/skillvalidator` | 7 | 26 | 77.7% | 单元 + 集成 |
| `pkg/userregistry` | 16 | 22 | 92.1% | 纯单元 |

---

## 一、skillvalidator 测试用例

### 1.1 ParseInput — 用户输入解析（13 组，纯单元）

测试 `ParseInput()` 函数对各种用户输入格式的解析能力。

| # | 用例 | 输入 | 期望 Kind | 期望 Repo | 期望 SkillHint |
|---|---|---|---|---|---|
| 1 | owner/repo 短格式 | `affaan-m/everything-claude-code` | RepoScan | `github.com/affaan-m/everything-claude-code` | (空) |
| 2 | owner/repo/skill 三段式 | `affaan-m/everything-claude-code/golang-testing` | SingleSkill | `github.com/affaan-m/everything-claude-code` | `golang-testing` |
| 3 | github.com 前缀兼容 | `github.com/affaan-m/everything-claude-code` | RepoScan | `github.com/affaan-m/everything-claude-code` | (空) |
| 4 | github.com 前缀带路径 | `github.com/affaan-m/everything-claude-code/skills/golang-testing` | SingleSkill | `github.com/affaan-m/everything-claude-code` | `skills/golang-testing` |
| 5 | https URL 自动剥离 | `https://github.com/affaan-m/everything-claude-code` | RepoScan | `github.com/affaan-m/everything-claude-code` | (空) |
| 6 | https URL tree/main 路径 | `https://github.com/.../tree/main/skills/golang-testing` | SingleSkill | `github.com/affaan-m/everything-claude-code` | `skills/golang-testing` |
| 7 | 绝对本地路径 | `/home/user/skills/my-skill` | Local | `/home/user/skills/my-skill` | (空) |
| 8 | 相对本地路径 | `./skills/go-i18n` | Local | `./skills/go-i18n` | (空) |
| 9 | 家目录路径 | `~/my-skills/test` | Local | `~/my-skills/test` | (空) |
| 10 | 前后空白裁剪 | `  affaan-m/everything-claude-code  ` | RepoScan | `github.com/affaan-m/everything-claude-code` | (空) |
| 11 | 尾部斜杠剥离 | `affaan-m/everything-claude-code/golang-testing/` | SingleSkill | `github.com/affaan-m/everything-claude-code` | `golang-testing` |
| 12 | http 前缀剥离 | `http://github.com/affaan-m/everything-claude-code` | RepoScan | `github.com/affaan-m/everything-claude-code` | (空) |
| 13 | blob/main 路径剥离 | `https://github.com/.../blob/main/skills/golang-testing` | SingleSkill | `github.com/affaan-m/everything-claude-code` | `skills/golang-testing` |

### 1.2 Validate — 本地技能校验（8 组，纯单元）

在 `t.TempDir()` 中构造 SKILL.md，测试 `Validate()` 函数的校验逻辑。

| # | 用例 | 场景 | 期望结果 |
|---|---|---|---|
| 1 | 合法 skill | 含完整 frontmatter (name/description/license) | Valid=true，字段解析正确 |
| 2 | LICENSE.txt 缺失警告 | 有 SKILL.md 但无 LICENSE.txt | Valid=true，Warnings 包含提示 |
| 3 | LICENSE.txt 存在 | 补充 LICENSE.txt 后校验 | 无 LICENSE.txt 相关警告 |
| 4 | 缺少 SKILL.md | 空目录 | Valid=false |
| 5 | 无 frontmatter | SKILL.md 只有正文无 `---` 块 | Valid=false |
| 6 | 非法 name 格式 | name 含大写/下划线 (`INVALID_NAME`) | Valid=false |
| 7 | 空 description | description 为空字符串 | Valid=false |
| 8 | 路径不存在 | 传入不存在的路径 | Valid=false |

### 1.3 parseFrontmatter — YAML 前置块解析（5 组，纯单元）

测试 `parseFrontmatter()` 内部函数对各种 SKILL.md 内容的解析容错。

| # | 用例 | 场景 | 期望 |
|---|---|---|---|
| 1 | 标准 frontmatter | `---\nname: x\ndescription: y\nlicense: MIT\n---` | 全字段正确解析 |
| 2 | 前导空行 | 文件头有空行，随后才是 `---` | 正常解析 |
| 3 | 无 frontmatter | 纯 Markdown 内容 | 返回空结构体 |
| 4 | 空文件 | 0 字节文件 | 返回空结构体 |
| 5 | 引号包裹描述 | `description: "Line one. Line two."` | 正确解析含标点的描述 |

### 1.4 Discover — 仓库级技能扫描（集成测试，需网络）

> 运行 `go test -short` 时自动跳过

| # | 用例 | 断言 |
|---|---|---|
| 1 | 扫描 `affaan-m/everything-claude-code` 全量技能 | 返回 ≥ 1 个技能 |
| 2 | 结果中包含 `golang-testing` | Name 匹配、Description 非空、Valid=true |
| 3 | 所有技能至少有 name | 无空 Name 存在 |

### 1.5 FindSkill — 仓库内定向查找（集成测试，需网络）

> 运行 `go test -short` 时自动跳过

| # | 用例 | 输入 | 期望 |
|---|---|---|---|
| 1 | 按名称查找 | `skillHint="golang-testing"` | 找到，Name/Description/Valid 正确 |
| 2 | 按显式路径查找 | `skillHint="skills/golang-testing"` | 找到，Name 正确 |
| 3 | 查找不存在的技能 | `skillHint="this-skill-does-not-exist-xyz"` | 返回 nil，无 error |

---

## 二、userregistry 测试用例

> 所有用例通过 `t.Setenv("XDG_CONFIG_HOME", tmpDir)` 隔离，不污染真实配置。

### 2.1 FilePath — 路径定位（1 组）

| # | 用例 | 断言 |
|---|---|---|
| 1 | XDG_CONFIG_HOME 重定向 | `FilePath()` = `$tmpDir/skills-x/user-registry.yaml` |

### 2.2 Load — 加载与持久化（2 组）

| # | 用例 | 断言 |
|---|---|---|
| 1 | 文件不存在时加载 | 返回空注册表，IsEmpty=true，TotalSkillCount=0 |
| 2 | Add → Save → Load 完整回环 | 重新 Load 后数据一致 (name, description) |

### 2.3 Add — 添加技能（6 组）

| # | 用例 | 输入 | 断言 |
|---|---|---|---|
| 1 | GitHub 源名推导 | repo=`github.com/affaan-m/everything-claude-code` | SourceName=`affaan-m-everything-claude-code` |
| 2 | 本地绝对路径源名 | repo=`/home/user/my-skills` | SourceName=`local` |
| 3 | 本地相对路径源名 | repo=`./my-skills` | SourceName=`local` |
| 4 | 同源多技能聚合 | 同一 repo 添加 2 个 skill | 同一 SourceEntry 下 Skills 长度=2 |
| 5 | 重复名称检测 | 添加两次 `golang-testing` | 第二次返回 "already exists" 错误 |
| 6 | 大小写不敏感去重 | 先加 `golang-testing`，再加 `Golang-Testing` | 第二次返回错误 |

### 2.4 Add — 冲突检测（1 组）

| # | 用例 | 输入 | 断言 |
|---|---|---|---|
| 1 | 与内建注册表冲突 | builtinNames 含 `golang-testing→[anthropic]` | ConflictSources=`["anthropic"]`，但技能仍然添加成功（用户优先） |

### 2.5 Remove — 删除技能（4 组）

| # | 用例 | 操作 | 断言 |
|---|---|---|---|
| 1 | 基本删除 | 添加后 Remove | IsEmpty=true，Source key 也被清理 |
| 2 | 大小写不敏感 | 添加 `golang-testing`，Remove `Golang-Testing` | 删除成功 |
| 3 | 技能不存在 | Remove 不存在的名称 | 返回 "not found" 错误 |
| 4 | 仅删除目标技能 | 添加 2 个 skill，删除其中 1 个 | 剩余 1 个，Source 保留 |

### 2.6 ListAll — 列表顺序（1 组）

| # | 用例 | 操作 | 断言 |
|---|---|---|---|
| 1 | 插入顺序保持 | 依次添加 3 个不同源的 skill | ListAll 输出顺序与插入顺序一致 |

### 2.7 IsEmpty / TotalSkillCount — 边界（1 组）

| # | 用例 | 场景 | 断言 |
|---|---|---|---|
| 1 | 空 Source 条目 | Source 存在但 Skills 列表为空 | IsEmpty=true，TotalSkillCount=0 |

### 2.8 deriveSourceName — 源名推导（6 组）

| # | 输入 | 期望输出 |
|---|---|---|
| 1 | `github.com/affaan-m/everything-claude-code` | `affaan-m-everything-claude-code` |
| 2 | `github.com/anthropics/courses` | `anthropics-courses` |
| 3 | `/home/user/skills` | `local` |
| 4 | `./my-skills` | `local` |
| 5 | `~/my-skills` | `local` |
| 6 | `custom.gitlab.com/foo/bar` | `custom-gitlab-com-foo-bar` |

### 2.9 Persistence — YAML 完整性验证（1 组）

| # | 用例 | 断言 |
|---|---|---|
| 1 | YAML 文件内容完整性 | 原始 YAML 包含 skill 名、source key、中文描述；重新 Load 后 DescriptionZh/License 一致 |

---

## 运行命令

```bash
# 仅运行单元测试（跳过集成测试）
go test -v -short ./pkg/skillvalidator/... ./pkg/userregistry/...

# 运行全部测试（含网络集成测试）
go test -v ./pkg/skillvalidator/... ./pkg/userregistry/...

# 查看覆盖率
go test -cover ./pkg/skillvalidator/... ./pkg/userregistry/...

# 生成覆盖率 HTML 报告
go test -coverprofile=coverage.out ./pkg/skillvalidator/... ./pkg/userregistry/...
go tool cover -html=coverage.out
```

---

## 测试隔离说明

| 策略 | 适用包 | 方法 |
|---|---|---|
| 文件系统隔离 | `userregistry` | `t.Setenv("XDG_CONFIG_HOME", tmpDir)` 重定向配置目录 |
| 临时目录 | `skillvalidator` | `t.TempDir()` 构造测试用 SKILL.md |
| 网络隔离 | `skillvalidator` (集成) | `testing.Short()` 跳过，`go test -short` 可避免网络依赖 |

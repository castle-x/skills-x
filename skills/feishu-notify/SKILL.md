---
name: feishu-notify
description: |
  帮助用户将飞书群机器人 Webhook 通知接入 Claude Code，使得每次任务完成后自动发送飞书消息卡片通知。

  触发场景（当用户提到以下任意内容时激活此 skill）：
  - "飞书通知"、"接入飞书"、"飞书提醒"
  - "task 完成通知"、"任务完成提醒"
  - "claude 完成提醒"、"claude 通知"
  - "Stop hook 通知"、"hook 通知"
  - "发飞书消息"、"飞书机器人"
  - "任务结束通知"、"配置通知"
---

# Feishu Notify — Claude Code 飞书通知 Skill

每当 Claude Code 完成任务（触发 Stop Hook）时，自动向飞书群发送一条消息卡片通知，内容包含会话 ID、完成时间、停止原因、最后使用的工具、工作目录、主机名、系统运行时间，以及 Claude 最后一条回复的摘要（前 300 字符）。

---

## 执行流程

### Step 1：引导用户获取飞书群机器人 Webhook 地址

首先询问用户是否已有飞书群机器人 Webhook 地址。如果没有，引导用户按以下步骤操作：

```
1. 打开飞书，进入任意群组（或新建一个专用的「Claude 通知」群）
2. 点击右上角「群设置」图标 → 「群机器人」→ 「添加机器人」
3. 选择「自定义机器人」
4. 填写机器人名称（如：Claude Code 助手）
5. 点击「添加」，复制生成的 Webhook 地址
   格式示例：https://open.feishu.cn/open-apis/bot/v2/hook/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

获取到 Webhook 地址后，继续执行 Step 2。

---

### Step 2：创建目录并写入通知脚本

创建 `~/.claude/` 目录，并将以下脚本写入 `~/.claude/notify-feishu.sh`。

**将脚本中的 `__WEBHOOK_URL__` 替换为用户提供的真实 Webhook 地址。**

```bash
#!/usr/bin/env bash
# 飞书通知脚本 - Claude Code Stop Hook
# 保存路径：~/.claude/notify-feishu.sh

WEBHOOK_URL="__WEBHOOK_URL__"

# 从 stdin 读取 Claude Code Hook 传入的 JSON
INPUT=$(cat)

# 提取 Hook 字段
SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // "unknown"')
STOP_HOOK_REASON=$(echo "$INPUT" | jq -r '.stop_hook_reason // empty')
TOOL_NAME=$(echo "$INPUT" | jq -r '.tool_name // empty')
TRANSCRIPT_PATH=$(echo "$INPUT" | jq -r '.transcript_path // empty')

# 提取环境信息
HOSTNAME=$(hostname 2>/dev/null || echo "unknown")
WORK_DIR=$(pwd)
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')
UPTIME=$(uptime -p 2>/dev/null || uptime | sed 's/.*up/up/')

SHORT_SESSION="${SESSION_ID:0:12}"

# 提取最后一条 assistant 文本回复（截取前 300 字符）
LAST_REPLY=""
if [[ -n "$TRANSCRIPT_PATH" && -f "$TRANSCRIPT_PATH" ]]; then
  LAST_REPLY=$(grep '"type":"assistant"' "$TRANSCRIPT_PATH" \
    | while IFS= read -r line; do
        echo "$line" | jq -r '
          select(.message.content[]?.type == "text")
          | [.message.content[] | select(.type == "text") | .text]
          | join("")' 2>/dev/null
      done \
    | tail -1 \
    | cut -c1-300)
  ORIGINAL_LEN=$(grep '"type":"assistant"' "$TRANSCRIPT_PATH" \
    | while IFS= read -r line; do
        echo "$line" | jq -r '
          select(.message.content[]?.type == "text")
          | [.message.content[] | select(.type == "text") | .text]
          | join("")' 2>/dev/null
      done | tail -1 | wc -c)
  [[ "$ORIGINAL_LEN" -gt 300 ]] && LAST_REPLY="${LAST_REPLY}..."
fi

# ⚠️ 关键：用数组 + printf '%s\n' 产生真实换行符
# 不能用 "\n" 字符串拼接，否则飞书会渲染为字面 \n 而不是换行
LINES=()
LINES+=("📋 **会话 ID**：\`${SHORT_SESSION}...\`")
LINES+=("⏰ **完成时间**：${TIMESTAMP}")
[[ -n "$STOP_HOOK_REASON" ]] && LINES+=("📝 **停止原因**：${STOP_HOOK_REASON}")
[[ -n "$TOOL_NAME" ]] && LINES+=("🔧 **最后工具**：${TOOL_NAME}")
LINES+=("📂 **工作目录**：\`${WORK_DIR}\`")
LINES+=("🖥️ **主机名称**：${HOSTNAME}")
LINES+=("🔌 **系统运行**：${UPTIME}")
if [[ -n "$LAST_REPLY" ]]; then
  LINES+=("")
  LINES+=("💬 **最后回复**：")
  LINES+=("${LAST_REPLY}")
fi

# printf '%s\n' 将数组每个元素转为真实换行的多行字符串
MD_CONTENT=$(printf '%s\n' "${LINES[@]}")

# ⚠️ 关键：用 jq -n --arg 传入多行字符串，jq 会正确转义为合法 JSON
# 不能用 echo/printf 手动拼接 JSON，否则换行符会破坏 JSON 结构
CARD_JSON=$(jq -n \
  --arg md "$MD_CONTENT" \
  '{
    msg_type: "interactive",
    card: {
      schema: "2.0",
      config: { wide_screen_mode: true },
      header: {
        template: "green",
        title: { tag: "plain_text", content: "✅ Claude Code 任务完成" }
      },
      body: {
        elements: [{ tag: "markdown", content: $md }]
      }
    }
  }')

# 发送通知，静默失败不影响 Claude Code 正常运行
curl -s -X POST "$WEBHOOK_URL" \
  -H 'Content-Type: application/json' \
  -d "$CARD_JSON" \
  > /dev/null 2>&1 || true
```

执行以下命令创建目录并写入脚本（替换 `<YOUR_WEBHOOK_URL>` 为真实地址）：

```bash
mkdir -p ~/.claude
# 写入脚本后赋予执行权限
chmod +x ~/.claude/notify-feishu.sh
```

---

### Step 3：配置 Claude Code settings.json

配置文件路径：`~/.claude/settings.json`（全局生效，对所有项目有效）。

> ⚠️ **合并注意**：如果该文件已存在，必须只追加 `hooks.Stop` 字段，不能覆盖文件中已有的其他配置（如 `mcpServers`、`permissions` 等）。

**如果文件不存在**，直接写入：

```json
{
  "hooks": {
    "Stop": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "bash ~/.claude/notify-feishu.sh",
            "timeout": 15,
            "statusMessage": "发送飞书通知..."
          }
        ]
      }
    ]
  }
}
```

**如果文件已存在**，使用 `jq` 合并（安全方式，不丢失已有配置）：

```bash
# 备份原文件
cp ~/.claude/settings.json ~/.claude/settings.json.bak

# 用 jq 合并，追加 Stop hook（不覆盖其他字段）
jq '.hooks.Stop += [{"hooks": [{"type": "command", "command": "bash ~/.claude/notify-feishu.sh", "timeout": 15, "statusMessage": "发送飞书通知..."}]}]' \
  ~/.claude/settings.json > /tmp/settings_merged.json \
  && mv /tmp/settings_merged.json ~/.claude/settings.json
```

---

### Step 4：测试发送通知

写入配置后，立即发送一条测试通知，验证 Webhook 是否有效：

```bash
echo '{"session_id":"test-session-abcdefg","stop_hook_reason":"task_complete","tool_name":"Bash","transcript_path":""}' \
  | bash ~/.claude/notify-feishu.sh
```

如果飞书群收到一条绿色标题的消息卡片，则说明配置成功。

如果没有收到，检查：
1. Webhook 地址是否正确粘贴（注意末尾不要有多余空格）
2. `jq` 和 `curl` 是否已安装（`which jq && which curl`）
3. 网络是否可以访问 `open.feishu.cn`

---

### Step 5：向用户确认完成

配置完成后，告知用户：

```
✅ 飞书通知配置完成！

从现在起，每次 Claude Code 完成任务时，飞书群将自动收到通知卡片，包含：
- 会话 ID（前12位）
- 完成时间
- 停止原因
- 最后使用的工具
- 工作目录 / 主机名 / 系统运行时间
- Claude 最后一条回复摘要（前 300 字符）

配置文件位置：
- 通知脚本：~/.claude/notify-feishu.sh
- Hook 配置：~/.claude/settings.json
```

---

## 注意事项

### ⚠️ 换行符问题（最常见的坑）

飞书消息卡片的 Markdown 内容必须包含**真实换行符**（`\n` 字节），不能是字面字符串 `\n`。

**正确做法**：
```bash
# 1. 用数组收集每行内容
LINES=()
LINES+=("第一行内容")
LINES+=("第二行内容")

# 2. 用 printf '%s\n' 将数组转为真实多行字符串
MD_CONTENT=$(printf '%s\n' "${LINES[@]}")

# 3. 用 jq --arg 传入，jq 会将真实换行转义为合法 JSON 中的 \n
CARD_JSON=$(jq -n --arg md "$MD_CONTENT" '{ ..., content: $md }')
```

**错误做法**（会导致飞书渲染出字面 `\n`）：
```bash
# ❌ 不要用字符串拼接
MD_CONTENT="第一行\n第二行"

# ❌ 不要手动拼接 JSON 字符串
CARD_JSON='{"content": "第一行\n第二行"}'
```

### ⚠️ settings.json 配置位置说明

| 位置 | 路径 | 作用范围 |
|------|------|---------|
| 全局配置 | `~/.claude/settings.json` | 对所有项目生效（推荐） |
| 项目配置 | `.claude/settings.json`（项目根目录） | 仅对当前项目生效 |

> 注意：Claude Code 官方全局配置路径可能因版本而异，也可能在 `~/.claude/settings.json`。请以实际环境为准，通过 `claude config` 或文档确认。

### ⚠️ settings.json 合并方式

- **不要直接覆盖**已有的 `settings.json`
- 使用 `jq` 合并，确保 `mcpServers`、`permissions`、其他 `hooks` 等已有配置不丢失
- 合并前务必备份（`.bak` 文件）

### ⚠️ 脚本依赖

脚本依赖以下工具，配置前请确认已安装：
- `jq`：JSON 处理（`apt install jq` / `brew install jq`）
- `curl`：HTTP 请求（通常系统自带）

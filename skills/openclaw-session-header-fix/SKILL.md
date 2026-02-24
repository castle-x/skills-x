---
name: openclaw-session-header-fix
description: 修复 OpenClaw 会话 JSONL 缺失 session header 导致的会话覆盖与上下文丢失问题。
license: MIT
metadata:
  author: x
  version: "1.0"
  tags:
    - openclaw
    - session
    - transcript
    - repair
---

# OpenClaw 会话 Header 修复

## 问题现象

- session 文件行数始终很少，且每次对话后像是被覆盖
- `sessions_history` 只看到 assistant 消息
- `openclaw doctor` 提示 sessions missing transcripts

## 根因

Session JSONL 第一行缺少 `type=session` 的 header。
OpenClaw 在加载时判定文件无效并重写，导致历史丢失。

## 处理步骤

1. 定位 sessionId 与文件
   - `cat ~/.openclaw/agents/<agentId>/sessions/sessions.json`
2. 检查 header
   - `head -1 <sessionId>.jsonl` 必须包含 `"type":"session"`
3. 备份并补写 header（建议先确保 gateway 空闲）

```bash
SESSION_ID="..."
FILE="/root/.openclaw/agents/<agentId>/sessions/${SESSION_ID}.jsonl"
FIRST_LINE=$(head -1 "$FILE")
if ! echo "$FIRST_LINE" | grep -q '"type":"session"'; then
  BACKUP="$FILE.bak.$(date +%s)"
  cp "$FILE" "$BACKUP"
  HEADER=$(printf '{"type":"session","version":3,"id":"%s","timestamp":"%s","cwd":"/root/openclaw"}' \
    "$SESSION_ID" "$(date -Iseconds)")
  { echo "$HEADER"; cat "$BACKUP"; } > "$FILE"
fi
```

## 验证

- 发送新消息后，`wc -l <sessionId>.jsonl` 行数应持续增长
- `grep '"role":"user"' <sessionId>.jsonl` 能看到用户消息

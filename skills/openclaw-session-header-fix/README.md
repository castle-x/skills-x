# OpenClaw 会话 Header 修复

## 简介

该技能用于修复 OpenClaw 会话 JSONL 缺失 session header 导致的会话覆盖与上下文丢失问题。

## 适用场景

- session 文件行数异常偏少、对话后被覆盖
- `sessions_history` 仅有 assistant 消息
- `openclaw doctor` 提示 sessions missing transcripts

## 处理概览

1. 定位 session 文件
2. 校验首行是否包含 `"type":"session"`
3. 备份并补写 header
4. 发送新消息后验证文件增长

## 目录结构

```
openclaw-session-header-fix/
├── SKILL.md
└── LICENSE.txt
```

## 使用提示

建议在 gateway 空闲时操作，并先做文件备份。

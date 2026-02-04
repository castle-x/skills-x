# 百度语音识别（音频转文本）

## 简介

该技能用于将用户语音消息（ogg/opus）转换为文本，面向国内服务器 + 代理环境优化，支持普通话、英语、粤语和四川话。

## 适用场景

- OpenClaw 在国内服务器运行，需要把语音消息转文字
- 代理环境下访问海外服务，但百度 API 需要直连

## 内容要点

- 通过 wrapper 脚本绕过代理访问百度 API
- 支持多语言识别与极速版
- 语音文件来源路径与处理流程
- 常见错误原因与处理方式

## 配置要求

运行前需配置环境变量（请勿写入仓库）：

```
export BAIDU_APP_ID="your_app_id"
export BAIDU_API_KEY="your_api_key"
export BAIDU_SECRET_KEY="your_secret_key"
```

## 目录结构

```
baidu-speech-to-text/
├── SKILL.md
└── LICENSE.txt
```

## 使用提示

不要直接调用 Python 脚本，优先使用 shell wrapper；涉及的 API 配置需妥善保管，避免泄露。

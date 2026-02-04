---
name: baidu-speech-to-text
description: 百度语音识别 - 将语音消息转换为文本。支持中文普通话、英语、粤语、四川话。专为国内服务器环境优化，自动绕过代理访问百度 API。
description_en: Baidu Speech Recognition - Convert voice messages to text. Supports Mandarin, English, Cantonese, and Sichuan dialect. Optimized for China mainland servers with automatic proxy bypass for Baidu API.
license: MIT
metadata:
  author: x
  version: "1.0"
  tags:
    - speech
    - voice
    - transcription
    - baidu
    - china
---

# 百度语音识别 - 音频转文本

## 概述

本 skill 用于将用户发送的语音消息（ogg/opus 格式）转换为文本。使用百度语音识别 API，专为**国内服务器 + 代理环境**优化。

## 环境背景

### 问题场景
- **国内服务器**通过 `proxychains4` 代理访问海外服务（Discord、WhatsApp 等）
- 但**百度 API 是国内服务**，通过代理访问反而会失败（SSL 错误）
- 需要在代理环境下**选择性绕过**百度 API

### 解决方案
- 使用 wrapper 脚本，通过 `unset LD_PRELOAD` 绕过 proxychains 注入
- 脚本内部清除代理环境变量，确保直连百度 API

## 脚本文件

| 脚本 | 路径 | 用途 |
|------|------|------|
| `ogg_to_text.sh` | `/root/.openclaw/workspace/scripts/ogg_to_text.sh` | **推荐** - 简单易用 |
| `speech_to_text.sh` | `/root/.openclaw/workspace/scripts/speech_to_text.sh` | 完整功能版 |
| `baidu_speech_to_text.py` | `/root/.openclaw/workspace/scripts/baidu_speech_to_text.py` | Python 主脚本 |

## 使用方法

### 基本用法（中文普通话）

```bash
/root/.openclaw/workspace/scripts/ogg_to_text.sh <音频文件路径>
```

### 示例

```bash
# 转换用户语音消息
/root/.openclaw/workspace/scripts/ogg_to_text.sh /root/.openclaw/media/inbound/xxxxx.ogg
```

### 支持的语言

| 参数 | 语言 |
|------|------|
| `--lang zh` | 中文普通话（默认）|
| `--lang en` | 英语 |
| `--lang cantonese` | 粤语 |
| `--lang sichuan` | 四川话 |

### 使用极速版（更快，仅中文）

```bash
/root/.openclaw/workspace/scripts/speech_to_text.sh <音频文件> --pro
```

## 音频文件位置

用户通过 WhatsApp/Discord 发送的语音消息保存在：
```
/root/.openclaw/media/inbound/
```

文件格式通常为 `.ogg`（Opus 编码），脚本会自动转换为 PCM 格式。

## 工作流程

当用户发送语音消息并请求转文本时：

1. 确定语音文件路径（`/root/.openclaw/media/inbound/*.ogg`）
2. 执行转换：
   ```bash
   /root/.openclaw/workspace/scripts/ogg_to_text.sh /root/.openclaw/media/inbound/<文件名>.ogg
   ```
3. 返回识别结果给用户

## 性能参考

| 阶段 | 耗时 | 说明 |
|------|------|------|
| 获取 token | ~100ms | 可缓存优化 |
| ogg 转 pcm | ~150-220ms | ffmpeg 转换 |
| API 调用 | ~400-600ms | **主要耗时** |
| **总计** | **~700-1000ms** | |

## API 配置

使用环境变量提供百度 API 账号信息（请勿写入仓库）：

```
export BAIDU_APP_ID="your_app_id"
export BAIDU_API_KEY="your_api_key"
export BAIDU_SECRET_KEY="your_secret_key"
```

端点：
- http://vop.baidu.com/server_api (标准版)
- https://vop.baidu.com/pro_api (极速版)

## 错误处理

| 错误 | 原因 | 解决方案 |
|------|------|---------|
| SSL 错误 | 代理影响 | 确保使用 wrapper 脚本（.sh），不要直接调用 Python |
| 识别结果为空 | 静音或无语音 | 告知用户音频可能没有语音内容 |
| 3301 错误 | 音频质量差 | 请用户重新录制 |
| 3303 错误 | 语音过长 | 需要分段处理 |

## 注意事项

1. **必须使用 shell wrapper 脚本**，不要直接调用 Python 脚本
2. 脚本已内置代理绕过逻辑，在 proxychains 环境下可正常工作
3. 短语音识别限制约 60 秒，超长音频需分段

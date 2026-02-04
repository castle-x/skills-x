#!/bin/bash
# OGG 音频自动转文本处理器（百度语音识别）
# 用法: ./ogg_to_text.sh <ogg文件路径>

# ========== 重要：绕过 proxychains 代理 ==========
unset LD_PRELOAD
unset http_proxy https_proxy HTTP_PROXY HTTPS_PROXY all_proxy ALL_PROXY
export no_proxy="aip.baidubce.com,vop.baidu.com,baidubce.com,baidu.com"
export NO_PROXY="$no_proxy"
# =================================================

OGG_FILE="$1"

if [ -z "$OGG_FILE" ]; then
    echo "用法: $0 <ogg文件路径>"
    exit 1
fi

if [ ! -f "$OGG_FILE" ]; then
    echo "文件不存在: $OGG_FILE"
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "========================================"
echo "🎵 OGG 音频转文本处理器"
echo "========================================"
echo "文件: $OGG_FILE"
echo "========================================"

# 调用百度语音识别（使用绕过代理的方式）
exec /usr/bin/python3 "$SCRIPT_DIR/baidu_speech_to_text.py" "$OGG_FILE"

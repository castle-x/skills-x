#!/bin/bash
#
# 百度语音识别 wrapper 脚本
# 绕过 proxychains 代理，直接访问百度 API
#

# 清除 LD_PRELOAD（禁用 proxychains 注入）
unset LD_PRELOAD

# 清除代理环境变量
unset http_proxy https_proxy HTTP_PROXY HTTPS_PROXY all_proxy ALL_PROXY socks_proxy SOCKS_PROXY

# 设置 no_proxy
export no_proxy="aip.baidubce.com,vop.baidu.com,baidubce.com,baidu.com,localhost,127.0.0.1"
export NO_PROXY="$no_proxy"

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 调用 Python 脚本
exec /usr/bin/python3 "$SCRIPT_DIR/baidu_speech_to_text.py" "$@"

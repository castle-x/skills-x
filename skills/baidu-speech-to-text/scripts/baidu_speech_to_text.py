#!/usr/bin/env python3
"""
百度语音识别 - 语音转文本
文档: https://cloud.baidu.com/doc/SPEECH/s/qlcirqhz0

支持的音频格式: pcm, wav, amr, m4a (需要16k采样率)
本脚本会自动将其他格式转换为 pcm 16k
"""

import os
import sys

# ============== 禁用代理（百度 API 是国内服务，不需要代理）==============
# 清除代理环境变量，避免 proxychains 影响
for proxy_var in ['http_proxy', 'https_proxy', 'HTTP_PROXY', 'HTTPS_PROXY', 
                  'all_proxy', 'ALL_PROXY', 'socks_proxy', 'SOCKS_PROXY']:
    os.environ.pop(proxy_var, None)
# 设置 no_proxy 包含百度域名
os.environ['no_proxy'] = 'aip.baidubce.com,vop.baidu.com,baidubce.com,baidu.com'
os.environ['NO_PROXY'] = os.environ['no_proxy']

import json
import base64
import subprocess
import urllib.request
import urllib.parse
import time
from pathlib import Path

# ============== 耗时统计 ==============
_timings = {}

def _start_timer(name):
    _timings[name] = {'start': time.time()}

def _end_timer(name):
    if name in _timings:
        _timings[name]['end'] = time.time()
        _timings[name]['duration'] = _timings[name]['end'] - _timings[name]['start']

def _print_timings():
    print("\n⏱️  耗时统计:")
    print("-" * 40)
    total = 0
    for name, t in _timings.items():
        if 'duration' in t:
            duration_ms = t['duration'] * 1000
            total += t['duration']
            print(f"  {name}: {duration_ms:.0f}ms")
    print("-" * 40)
    print(f"  总计: {total*1000:.0f}ms ({total:.2f}s)")
    print("-" * 40)

# ============== 百度 API 配置（从环境变量读取） ==============
APP_ID = os.environ.get("BAIDU_APP_ID")
API_KEY = os.environ.get("BAIDU_API_KEY")
SECRET_KEY = os.environ.get("BAIDU_SECRET_KEY")

if not APP_ID or not API_KEY or not SECRET_KEY:
    print("✗ 缺少环境变量：BAIDU_APP_ID / BAIDU_API_KEY / BAIDU_SECRET_KEY")
    print("  请先导出环境变量再运行脚本。")
    sys.exit(1)

# API 端点
TOKEN_URL = "https://aip.baidubce.com/oauth/2.0/token"
# 短语音识别标准版 (支持中文普通话、英语等)
ASR_URL = "http://vop.baidu.com/server_api"
# 短语音识别极速版 (仅支持中文普通话，识别更快)
ASR_PRO_URL = "https://vop.baidu.com/pro_api"


def get_access_token():
    """获取百度 API access_token"""
    _start_timer("1.获取token")
    
    params = {
        "grant_type": "client_credentials",
        "client_id": API_KEY,
        "client_secret": SECRET_KEY
    }
    
    url = TOKEN_URL + "?" + urllib.parse.urlencode(params)
    
    try:
        with urllib.request.urlopen(url, timeout=10) as response:
            result = json.loads(response.read().decode())
            _end_timer("1.获取token")
            if "access_token" in result:
                print(f"✓ 获取 access_token 成功")
                return result["access_token"]
            else:
                print(f"✗ 获取 token 失败: {result}")
                return None
    except Exception as e:
        _end_timer("1.获取token")
        print(f"✗ 请求 token 出错: {e}")
        return None


def convert_audio_to_pcm(input_file, output_file=None):
    """
    使用 ffmpeg 将音频转换为百度支持的 PCM 格式
    - 采样率: 16000 Hz
    - 声道: 单声道
    - 位深: 16bit
    """
    _start_timer("2.ogg转pcm")
    
    if output_file is None:
        output_file = str(Path(input_file).with_suffix('.pcm'))
    
    cmd = [
        "ffmpeg", "-y",  # 覆盖输出文件
        "-i", input_file,
        "-f", "s16le",   # PCM signed 16-bit little-endian
        "-acodec", "pcm_s16le",
        "-ar", "16000",  # 采样率 16kHz
        "-ac", "1",      # 单声道
        output_file
    ]
    
    print(f"转换音频: {input_file} -> {output_file}")
    
    try:
        result = subprocess.run(cmd, capture_output=True, text=True)
        _end_timer("2.ogg转pcm")
        if result.returncode == 0:
            print(f"✓ 音频转换成功")
            return output_file
        else:
            print(f"✗ ffmpeg 转换失败: {result.stderr}")
            return None
    except FileNotFoundError:
        _end_timer("2.ogg转pcm")
        print("✗ ffmpeg 未安装，请先安装: dnf install -y ffmpeg 或 apt install -y ffmpeg")
        return None


def speech_to_text(audio_file, use_pro=False, dev_pid=1537):
    """
    调用百度语音识别 API
    
    参数:
        audio_file: 音频文件路径 (pcm/wav/amr/m4a)
        use_pro: 是否使用极速版 (仅支持中文普通话)
        dev_pid: 语言模型
            - 1537: 中文普通话 (标准版默认)
            - 1737: 英语
            - 1637: 粤语
            - 1837: 四川话
            - 80001: 中文普通话 (极速版)
    
    返回:
        识别结果文本，失败返回 None
    """
    # 获取 token
    access_token = get_access_token()
    if not access_token:
        return None
    
    # 检查文件格式，如果不是 pcm/wav/amr/m4a 则转换
    file_ext = Path(audio_file).suffix.lower()
    if file_ext not in ['.pcm', '.wav', '.amr', '.m4a']:
        print(f"音频格式 {file_ext} 需要转换...")
        pcm_file = convert_audio_to_pcm(audio_file)
        if not pcm_file:
            return None
        audio_file = pcm_file
        file_ext = '.pcm'
    
    # 读取音频文件
    _start_timer("3.读取文件")
    with open(audio_file, 'rb') as f:
        audio_data = f.read()
    _end_timer("3.读取文件")
    
    audio_length = len(audio_data)
    print(f"音频大小: {audio_length} bytes")
    
    # 检查文件大小限制 (短语音识别限制约 60 秒，文件不超过 10MB)
    if audio_length > 10 * 1024 * 1024:
        print("✗ 音频文件过大 (>10MB)，请使用音频文件转写 API")
        return None
    
    # Base64 编码
    _start_timer("4.Base64编码")
    audio_base64 = base64.b64encode(audio_data).decode('utf-8')
    _end_timer("4.Base64编码")
    
    # 构建请求参数
    format_map = {'.pcm': 'pcm', '.wav': 'wav', '.amr': 'amr', '.m4a': 'm4a'}
    
    params = {
        "format": format_map.get(file_ext, 'pcm'),
        "rate": 16000,
        "channel": 1,
        "cuid": f"openclaw_{APP_ID}",
        "token": access_token,
        "speech": audio_base64,
        "len": audio_length,
        "dev_pid": 80001 if use_pro else dev_pid  # 极速版使用 80001
    }
    
    # 选择 API 端点
    api_url = ASR_PRO_URL if use_pro else ASR_URL
    print(f"调用 {'极速版' if use_pro else '标准版'} API: {api_url}")
    
    # 发送请求
    headers = {"Content-Type": "application/json"}
    request_data = json.dumps(params).encode('utf-8')
    
    _start_timer("5.API调用")
    try:
        req = urllib.request.Request(api_url, data=request_data, headers=headers)
        with urllib.request.urlopen(req, timeout=30) as response:
            result = json.loads(response.read().decode())
            _end_timer("5.API调用")
            
            if result.get("err_no") == 0:
                text = "".join(result.get("result", []))
                print(f"✓ 识别成功!")
                return text
            else:
                err_no = result.get("err_no")
                err_msg = result.get("err_msg", "未知错误")
                print(f"✗ 识别失败 [{err_no}]: {err_msg}")
                
                # 常见错误提示
                if err_no == 3301:
                    print("  提示: 音频质量过差或格式不正确")
                elif err_no == 3302:
                    print("  提示: 鉴权失败，请检查 API Key 和 Secret Key")
                elif err_no == 3303:
                    print("  提示: 语音过长，请使用音频文件转写 API")
                elif err_no == 3304:
                    print("  提示: 请求参数错误")
                elif err_no == 3305:
                    print("  提示: 其他客户端错误")
                
                return None
                
    except Exception as e:
        _end_timer("5.API调用")
        print(f"✗ API 请求出错: {e}")
        return None


def main():
    """主函数"""
    import argparse
    
    parser = argparse.ArgumentParser(
        description="百度语音识别 - 语音转文本",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  %(prog)s audio.ogg                    # 识别音频文件
  %(prog)s audio.ogg --pro              # 使用极速版
  %(prog)s audio.ogg --lang en          # 识别英语
  %(prog)s audio.ogg --lang cantonese   # 识别粤语
        """
    )
    parser.add_argument("audio_file", help="音频文件路径")
    parser.add_argument("--pro", action="store_true", help="使用极速版 (仅支持中文普通话)")
    parser.add_argument("--lang", default="zh", choices=["zh", "en", "cantonese", "sichuan"],
                        help="语言: zh=中文普通话, en=英语, cantonese=粤语, sichuan=四川话")
    parser.add_argument("--output", "-o", help="输出文件路径 (不指定则输出到终端)")
    
    args = parser.parse_args()
    
    # 检查文件是否存在
    if not os.path.exists(args.audio_file):
        print(f"✗ 文件不存在: {args.audio_file}")
        sys.exit(1)
    
    # 语言映射到 dev_pid
    lang_map = {
        "zh": 1537,       # 中文普通话
        "en": 1737,       # 英语
        "cantonese": 1637, # 粤语
        "sichuan": 1837   # 四川话
    }
    dev_pid = lang_map.get(args.lang, 1537)
    
    print("=" * 50)
    print("百度语音识别")
    print("=" * 50)
    print(f"音频文件: {args.audio_file}")
    print(f"API 版本: {'极速版' if args.pro else '标准版'}")
    print(f"语言: {args.lang}")
    print("-" * 50)
    
    # 调用语音识别
    result = speech_to_text(args.audio_file, use_pro=args.pro, dev_pid=dev_pid)
    
    # 打印耗时统计
    _print_timings()
    
    if result:
        print("-" * 50)
        print("识别结果:")
        print(result)
        
        # 输出到文件
        if args.output:
            with open(args.output, 'w', encoding='utf-8') as f:
                f.write(result)
            print(f"\n结果已保存到: {args.output}")
    else:
        sys.exit(1)


if __name__ == "__main__":
    main()

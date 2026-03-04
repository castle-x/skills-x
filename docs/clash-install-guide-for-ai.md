# Clash 代理安装指南 (AI 执行版 - 国内镜像版)

## 1. 软件信息

| 属性 | 值 |
|------|-----|
| 软件名称 | Clash (Premium/Core) / Mihomo |
| 类型 | 代理客户端 |
| 监听端口 | HTTP: `7890`, SOCKS5: `7891` |
| 控制面板 | `http://127.0.0.1:9090` |
| 配置目录 | `~/.config/clash/` |

## 2. 国内 GitHub 镜像源

> 优先使用以下国内镜像加速下载

| 镜像名称 | 地址格式 |
|---------|---------|
| **ghproxy (推荐)** | `https://ghproxy.com/https://github.com/...` |
| **mirror.ghproxy** | `https://mirror.ghproxy.com/https://github.com/...` |
| **ghps.cc** | `https://ghps.cc/https://github.com/...` |
| **gh.api.99988866.xyz** | `https://gh.api.99988866.xyz/https://github.com/...` |
| **fastgit** | `https://hub.fastgit.xyz/...` |
| **jsdelivr CDN** | `https://cdn.jsdelivr.net/gh/user/repo@tag/file` |

## 3. 安装步骤 (国内网络环境)

### 3.1 使用 ghproxy 镜像 (推荐)

```bash
# 创建配置目录
mkdir -p ~/.config/clash

# 设置镜像前缀
GH_PROXY="https://ghproxy.com"

# 下载 Clash Premium 内核 (根据系统架构选择)
# amd64:
curl -L -o /tmp/clash.gz "${GH_PROXY}/https://github.com/Dreamacro/clash/releases/download/premium/clash-linux-amd64-v3-2023.08.17.gz"

# arm64:
# curl -L -o /tmp/clash.gz "${GH_PROXY}/https://github.com/Dreamacro/clash/releases/download/premium/clash-linux-arm64-v3-2023.08.17.gz"

# 解压并安装
gunzip -f /tmp/clash.gz
chmod +x /tmp/clash
sudo mv /tmp/clash /usr/local/bin/clash

# 验证安装
clash -v
```

### 3.2 使用 Mihomo (Clash.Meta) - 国内维护活跃

```bash
# Mihomo 是 Clash 的增强版，国内使用广泛，GitHub 镜像下载
mkdir -p ~/.config/clash
GH_PROXY="https://ghproxy.com"

# 获取最新版本号 (或使用固定版本)
VERSION="v1.18.5"

# Linux amd64
curl -L -o /tmp/mihomo.gz "${GH_PROXY}/https://github.com/MetaCubeX/mihomo/releases/download/${VERSION}/mihomo-${VERSION}-linux-amd64-compatible.gz"

# Linux arm64
# curl -L -o /tmp/mihomo.gz "${GH_PROXY}/https://github.com/MetaCubeX/mihomo/releases/download/${VERSION}/mihomo-${VERSION}-linux-arm64.gz"

# 解压安装
gunzip -f /tmp/mihomo.gz
chmod +x /tmp/mihomo-*
sudo mv /tmp/mihomo-* /usr/local/bin/clash

# 验证
clash -v
```

### 3.3 使用 jsdelivr CDN (版本较旧但稳定)

```bash
# jsdelivr 在国内有节点，但只能下载具体文件，不能下载 release
# 适合下载配置示例或配置文件

# 示例：下载配置文件模板
curl -L -o /tmp/config.yaml "https://cdn.jsdelivr.net/gh/Dreamacro/clash@master/docs/config.yaml"
```

### 3.4 手动下载安装 (如果自动下载失败)

```bash
# 如果镜像都失败，建议用户手动下载后上传
# 1. 通过浏览器/其他方式下载 clash 二进制文件
# 2. 上传到服务器
# 3. 执行安装

chmod +x /path/to/uploaded/clash
sudo mv /path/to/uploaded/clash /usr/local/bin/clash
```

## 4. 配置加载方式 (国内网络)

### 4.1 订阅 URL 使用镜像

```bash
# 如果订阅 URL 是 GitHub raw 链接，使用镜像加速

# 原始链接 (慢)
# SUB_URL="https://raw.githubusercontent.com/xxx/xxx/config.yaml"

# 使用镜像 (快)
SUB_URL="https://ghproxy.com/https://raw.githubusercontent.com/xxx/xxx/config.yaml"

# 下载配置
curl -L -o ~/.config/clash/config.yaml "${SUB_URL}"
```

### 4.2 常用订阅转换服务 (国内可用)

```bash
# 如果订阅链接需要转换，可以使用以下服务：
# - https://v1.mk/ (品云转换)
# - https://sub.xeton.dev/ 
# - https://sub.id9.cc/
# - https://sub.maoxiongnet.com/

# 转换后下载
curl -L -o ~/.config/clash/config.yaml "<转换后的URL>"
```

## 5. 环境变量设置

将以下内容添加到 `~/.bashrc` 或 `~/.zshrc`：

```bash
# Clash 代理环境变量
export http_proxy=http://127.0.0.1:7890
export https_proxy=http://127.0.0.1:7890
export HTTP_PROXY=http://127.0.0.1:7890
export HTTPS_PROXY=http://127.0.0.1:7890
export ALL_PROXY=socks5://127.0.0.1:7891
export no_proxy=localhost,127.0.0.1,*.local,*.cn,*.aliyun.com,*.tencent.com,*.baidu.com
```

使配置生效：
```bash
source ~/.bashrc  # 或 source ~/.zshrc
```

## 6. 启动 Clash

### 前台运行 (测试)
```bash
clash -d ~/.config/clash
```

### 后台运行 (生产)
```bash
# 使用 nohup
nohup clash -d ~/.config/clash > /tmp/clash.log 2>&1 &

# 或使用 systemd (推荐)
sudo tee /etc/systemd/system/clash.service > /dev/null << EOF
[Unit]
Description=Clash Daemon
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/clash -d $(eval echo ~$(whoami))/.config/clash
Restart=on-failure
User=$(whoami)

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable clash
sudo systemctl start clash
sudo systemctl status clash
```

## 7. 验证安装

```bash
# 1. 检查进程
ps aux | grep clash

# 2. 检查端口监听
netstat -tlnp | grep -E '7890|7891|9090' 2>/dev/null || \
ss -tlnp | grep -E '7890|7891|9090'

# 3. 测试代理 (需要配置已加载)
curl -x http://127.0.0.1:7890 -I https://www.google.com

# 4. 查看日志
tail -f /tmp/clash.log
```

## 8. 备用方案 (镜像失效时)

### 方案 A: 使用 GitHub 代理网站

```bash
# 使用 gh.api.99988866.xyz
curl -L -o /tmp/clash.gz "https://gh.api.99988866.xyz/https://github.com/Dreamacro/clash/releases/download/premium/clash-linux-amd64-v3-2023.08.17.gz"
```

### 方案 B: 使用 Gitee 镜像

```bash
# 部分开源项目有 Gitee 镜像
# Clash 官方没有，但可以考虑使用国内 fork
# 或者使用 sing-box (另一个代理工具，国内下载更方便)
```

### 方案 C: 使用 sing-box (推荐备选)

```bash
# sing-box 国内下载更方便，功能类似
# 官网: https://sing-box.sagernet.org/

# 下载 (国内可访问)
curl -L -o /tmp/sing-box.tar.gz "https://github.com/SagerNet/sing-box/releases/download/v1.8.14/sing-box-1.8.14-linux-amd64.tar.gz"
# 或使用镜像
curl -L -o /tmp/sing-box.tar.gz "https://ghproxy.com/https://github.com/SagerNet/sing-box/releases/download/v1.8.14/sing-box-1.8.14-linux-amd64.tar.gz"

tar -xzf /tmp/sing-box.tar.gz -C /tmp
sudo mv /tmp/sing-box-*/sing-box /usr/local/bin/
```

## 9. 常见问题

### Q: ghproxy 镜像无法访问？

```bash
# 尝试其他镜像
GH_PROXY="https://mirror.ghproxy.com"  # 备用1
GH_PROXY="https://ghps.cc"              # 备用2
GH_PROXY="https://gh.api.99988866.xyz"  # 备用3

# 或使用 wget 替代 curl
wget -O /tmp/clash.gz "${GH_PROXY}/https://github.com/Dreamacro/clash/releases/download/premium/clash-linux-amd64-v3-2023.08.17.gz"
```

### Q: 下载成功但无法运行？

```bash
# 检查架构是否匹配
uname -m  # x86_64(amd64) / aarch64(arm64)

# 检查依赖
ldd /usr/local/bin/clash  # Linux
otool -L /usr/local/bin/clash  # macOS

# 赋予执行权限
chmod +x /usr/local/bin/clash
```

### Q: 配置文件下载失败？

```bash
# 可能是订阅链接需要 UA 头
curl -L -A "clash" -o ~/.config/clash/config.yaml "<SUB_URL>"

# 或使用 wget
wget --user-agent="clash" -O ~/.config/clash/config.yaml "<SUB_URL>"
```

## 10. 待用户提供的参数

执行前需要用户确认：

- [ ] **订阅 URL** - 配置文件下载地址
- [ ] **系统架构** - `x86_64`(amd64) / `aarch64`(arm64)
- [ ] **GitHub 镜像** - 推荐 `ghproxy.com`，如失效使用备选
- [ ] **是否开机自启** - yes/no

---

**快速执行模板**:

```bash
#!/bin/bash
set -e

GH_PROXY="https://ghproxy.com"
SUB_URL="<用户提供>"
ARCH="amd64"  # 或 arm64

# 1. 安装
mkdir -p ~/.config/clash
curl -L -o /tmp/clash.gz "${GH_PROXY}/https://github.com/Dreamacro/clash/releases/download/premium/clash-linux-${ARCH}-v3-2023.08.17.gz"
gunzip -f /tmp/clash.gz
chmod +x /tmp/clash
sudo mv /tmp/clash /usr/local/bin/clash

# 2. 下载配置
curl -L -o ~/.config/clash/config.yaml "${SUB_URL}"

# 3. 启动
nohup clash -d ~/.config/clash > /tmp/clash.log 2>&1 &

# 4. 验证
sleep 2
curl -x http://127.0.0.1:7890 -s -o /dev/null -w "%{http_code}" https://www.google.com
echo "安装完成!"
```

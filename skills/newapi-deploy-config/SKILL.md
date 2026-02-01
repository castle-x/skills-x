---
name: newapi-deploy-config
description: 在本机使用 Docker 部署 New API（host 网络），并支持后续通过管理 API 自动配置模型渠道（不包含任何密钥）。
license: MIT
metadata:
  author: x
  version: "1.0"
  tags:
    - new-api
    - docker
    - deploy
    - host-network
    - proxy
---

# New API 部署与模型配置（Host 网络）

## 适用场景

- 需要在当前机器快速部署 New API
- 需要宿主机直接访问容器（host 网络模式）
- 国内网络环境无法直接访问 Docker Hub 时，需要为 Docker 守护进程配置代理

## 依赖

- Docker
- 本机可用的 HTTP 代理（示例：Clash 7890）

## 一键部署流程

### 1. 安装并启动 Docker

> 适用于 OpenCloudOS/RHEL 系列

```bash
dnf install -y docker
systemctl enable --now docker
```

### 2. 国内网络：为 Docker 守护进程配置代理（可选）

当拉取 Docker Hub 镜像失败时，需要让 Docker 守护进程走代理：

```bash
mkdir -p /etc/systemd/system/docker.service.d
cat <<'EOT' >/etc/systemd/system/docker.service.d/http-proxy.conf
[Service]
Environment="HTTP_PROXY=http://127.0.0.1:7890" "HTTPS_PROXY=http://127.0.0.1:7890" "NO_PROXY=localhost,127.0.0.1"
EOT
systemctl daemon-reload
systemctl restart docker
```

### 3. 拉取镜像

```bash
docker pull calciumion/new-api:latest
```

### 4. 使用 host 网络启动容器

```bash
mkdir -p /root/newapi/data

docker run --name new-api -d --restart always \
  --network host \
  -e TZ=Asia/Shanghai \
  -v /root/newapi/data:/data \
  calciumion/new-api:latest
```

### 5. 验证服务

```bash
curl -fsS http://127.0.0.1:3000/api/status
```

返回 `"success":true` 表示服务正常。

### 6. 初始化管理员账号

浏览器打开：

```
http://127.0.0.1:3000
```

首次访问会进入初始化页面，创建管理员账号即可使用。

## 常用运维

```bash
# 查看容器状态
docker ps --filter name=new-api

# 查看日志
docker logs -f new-api

# 重启服务
docker restart new-api

# 停止并删除
docker rm -f new-api
```

## 自动配置模型渠道（不包含密钥）

> 下面流程会通过管理 API 创建渠道。请你在控制台手动生成管理员会话或 Token，再进行调用。**不要把任何 API Key 写进文档或脚本**。

### 1. 获取会话 Cookie（示例）

在浏览器登录后台后，从开发者工具复制 `session` Cookie。命令示例：

```bash
export NEWAPI_COOKIE='session=YOUR_SESSION_COOKIE'
export NEWAPI_USER_ID='1'
```

### 2. 配置渠道（示例模板）

以下仅为模板，**请将 `key` 替换为你自己的密钥**，或使用环境变量注入。

```bash
curl -fsS -X POST http://127.0.0.1:3000/api/channel/ \
  -H 'Content-Type: application/json' \
  -H "Cookie: ${NEWAPI_COOKIE}" \
  -H "New-API-User: ${NEWAPI_USER_ID}" \
  -d '{
    "mode": "single",
    "channel": {
      "name": "OpenAI",
      "type": 1,
      "key": "<REPLACE_WITH_YOUR_KEY>",
      "base_url": "",
      "models": "gpt-5.2",
      "group": "default",
      "status": 1,
      "auto_ban": 1
    }
  }'
```

### 3. 火山方舟（Ark）模型映射模板

```bash
curl -fsS -X POST http://127.0.0.1:3000/api/channel/ \
  -H 'Content-Type: application/json' \
  -H "Cookie: ${NEWAPI_COOKIE}" \
  -H "New-API-User: ${NEWAPI_USER_ID}" \
  -d '{
    "mode": "single",
    "channel": {
      "name": "Volcengine Ark",
      "type": 45,
      "key": "<REPLACE_WITH_YOUR_KEY>",
      "base_url": "https://ark.cn-beijing.volces.com",
      "models": "doubao-seed,doubao-code,doubao-vision",
      "model_mapping": "{\\"doubao-seed\\":\\"ep-xxxxxxxx\\",\\"doubao-code\\":\\"ep-xxxxxxxx\\",\\"doubao-vision\\":\\"ep-xxxxxxxx\\"}",
      "group": "default",
      "status": 1,
      "auto_ban": 1
    }
  }'
```

### 4. 验证渠道创建

```bash
curl -fsS -H "Cookie: ${NEWAPI_COOKIE}" -H "New-API-User: ${NEWAPI_USER_ID}" \
  "http://127.0.0.1:3000/api/channel/?page=0&size=50"
```

## 安全注意事项

- **不要把 API Key 写入文档、脚本或提交到仓库**
- 建议使用环境变量或临时粘贴的方式进行配置
- 需要共享流程时，请仅分享模板

## 参考文档

- https://docs.newapi.pro/zh/docs

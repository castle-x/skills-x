# Go Embedded SPA

## 简介

该技能介绍如何将前端 SPA 静态资源（React/Vue/TSX）通过 `go:embed` 嵌入到 Go 二进制中，实现单文件部署。

## 适用场景

- 需要单一二进制的全栈交付
- SPA + Go 后端整合部署
- 跨平台（Linux/macOS/Windows）发布

## 内容要点

- `go:embed` 资源嵌入流程
- SPA 路由回退与静态文件服务
- 构建顺序与部署建议

## 目录结构

```
go-embedded-spa/
├── SKILL.md
├── assets/
│   ├── embed.go.tmpl
│   ├── Makefile.tmpl
│   └── vite.config.ts.tmpl
└── references/
    └── siteserver.md
```

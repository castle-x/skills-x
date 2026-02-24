# Skills-X 贡献指南（中文）

## 简介

该技能用于指导向 `skills-x` 收录新技能，包括外部开源技能的登记与自研技能的规范化流程。

## 适用场景

- 从外部仓库引入新技能（例如 agentskills.io、anthropics/skills）
- 创建或更新自研技能
- 校验技能目录结构与 `SKILL.md` 前置元数据
- 处理 `skills-x` 的双语 i18n 输出规范

## 目录结构

```
skills-x/
├── skills/              # 自研技能目录
├── pkg/registry/        # 外部技能索引
└── cmd/skills-x/i18n/   # 双语文案
```

## 使用提示

- 外部技能：只需更新 `pkg/registry/registry.yaml`
- 自研技能：放入 `skills/<name>/` 并补全 i18n
- 不要在单一字符串中混用中英文

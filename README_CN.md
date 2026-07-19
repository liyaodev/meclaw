# MeClaw

[English](README.md)

**产品名：** meclaw（claw）。**品类：** Agent Infra → Agent SaaS。

把微信 / 飞书 / 企微接到本地或远程 Agent（Claude、Codex 等），再用技能包做成「微信里的工作助手」——**场景 A + B 同一仓库**。

> 统一叙事：**IM Agent 基建 → 微信里能干活的助手。**

## 现状

脚手架阶段。先做场景 A（网关 + 路由）；场景 B（`skills/` + `hosted/`）同仓后做。

## 快速开始

```bash
make tidy
make run
# 打印脚手架提示 — 下一步实现 gateway
```

## 目录

```text
cmd/                 CLI 入口
internal/
  gateway/           IM 通道归一化
  agent/             ACP / CLI / HTTP
  session/           会话
  policy/            白名单 / 审计（To B）
  config/            配置
skills/              场景 B 技能包
hosted/              邀请制 / 用量（后期）
docs/                策略与架构
examples/            示例配置
```

## 文档

| 文档 | 内容 |
|------|------|
| [docs/README.md](docs/README.md) | 文档索引 |
| [docs/naming.md](docs/naming.md) | 命名：claw / Agent |
| [docs/concepts-agent-infra-saas.md](docs/concepts-agent-infra-saas.md) | Infra / SaaS 概念 |
| [docs/scenario-a-infra-im-gateway.md](docs/scenario-a-infra-im-gateway.md) | 场景 A |
| [docs/scenario-b-saas-wechat-assistant.md](docs/scenario-b-saas-wechat-assistant.md) | 场景 B |
| [docs/playbook-flywheel-90d.md](docs/playbook-flywheel-90d.md) | 飞轮与 90 天 |
| [docs/architecture.md](docs/architecture.md) | 架构 |
| [docs/repo-layout.md](docs/repo-layout.md) | 单仓约定 |

## 优先级

1. 场景 A 开源 MVP  
2. 同仓场景 B 单一技能 + 邀请制  
3. 内容 / 课 / 私有化 To B  

## License

MIT

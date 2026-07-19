# MeClaw

[English](README.md)

**产品名：** meclaw（claw）。**品类：** Agent Infra → Agent SaaS。

把微信 / 飞书 / 企微接到本地或远程 Agent（Claude、Codex 等），再用技能包做成「微信里的工作助手」——**场景 A + B 同一仓库**。

> 统一叙事：**IM Agent 基建 → 微信里能干活的助手。**

## 现状

场景 A MVP 已通：运行时管道 + stdio/HTTP 入站 + 飞书机器人适配。场景 B（`skills/` + `hosted/`）同仓后做。

## 快速开始

```bash
make tidy && make build
./bin/meclaw chat -c examples/config.example.json
# 输入消息；/agent <id> 切换 Agent

./bin/meclaw serve -c examples/config.example.json
# POST /v1/message ；可选 /v1/feishu/event
```

```bash
make test
```

详见 [docs/scenario-a-mvp.md](docs/scenario-a-mvp.md)。

## 目录

```text
cmd/                 CLI（chat / serve / version）
internal/
  gateway/           IM 通道归一化 + stdio/HTTP + feishu/
  agent/             ACP（占位）/ CLI / HTTP
  session/           会话
  policy/            白名单 / 审计（To B）
  config/            配置
  runtime/           Message → policy → session → agent
skills/              场景 B 技能包
hosted/              邀请制 / 用量（后期）
docs/                策略与架构
examples/            示例配置
```

## 文档

| 文档 | 内容 |
|------|------|
| [docs/README.md](docs/README.md) | 文档索引 |
| [docs/scenario-a-mvp.md](docs/scenario-a-mvp.md) | 场景 A MVP 用法 |
| [docs/scenario-a-infra-im-gateway.md](docs/scenario-a-infra-im-gateway.md) | 场景 A 策略 |
| [docs/scenario-b-saas-wechat-assistant.md](docs/scenario-b-saas-wechat-assistant.md) | 场景 B |
| [docs/architecture.md](docs/architecture.md) | 架构 |

## 优先级

1. 场景 A 开源 MVP（已落地核心管道 + 飞书）
2. 同仓场景 B 单一技能 + 邀请制
3. 内容 / 课 / 私有化 To B

## License

MIT

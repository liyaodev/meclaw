# MeClaw

[English](README.md)

**产品名：** meclaw（claw）。**品类：** Agent Infra → Agent SaaS。

自建 claw 全栈（A1–A6）+ 微信工作助手（场景 B）——**同一仓库**。

> 统一叙事：**如何跑通自建 MeClaw → 微信里能干活的助手。**

## 现状

- **A1–A6** 自建 claw 主干已通（见 [docs/self-hosted-claw.md](docs/self-hosted-claw.md)）
- **B**：`skills/` 示例；`hosted/` 后期

## 快速开始

```bash
make tidy && make build && make test
./bin/meclaw chat -c examples/config.example.json
./bin/meclaw serve -c examples/config.example.json
```

主路径：[docs/self-hosted-claw.md](docs/self-hosted-claw.md)  
能力矩阵：[docs/infra-capability-map.md](docs/infra-capability-map.md)

## 目录

```text
cmd/                 CLI（chat / serve / version）
internal/            A1–A5 核心（gateway/agent/session/memory/…）
deploy/              A6 Docker / K8s 骨架
skills/              场景 B 技能包
hosted/              邀请制 / 用量（后期）
docs/                策略与自建路径
examples/            示例配置
```

## 文档

| 文档 | 内容 |
|------|------|
| [docs/dev-context.md](docs/dev-context.md) | **开发上下文（后续开工先读）** |
| [docs/self-hosted-claw.md](docs/self-hosted-claw.md) | 如何跑通自建 claw |
| [docs/infra-capability-map.md](docs/infra-capability-map.md) | Infra 能力矩阵 |
| [docs/scenario-a-mvp.md](docs/scenario-a-mvp.md) | A1 MVP |
| [docs/architecture.md](docs/architecture.md) | 架构 |

## 优先级

1. 自建 claw A1→A6  
2. 同仓场景 B  
3. 内容 / 课 / 私有化 To B  

## License

MIT

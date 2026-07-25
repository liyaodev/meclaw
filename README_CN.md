# MeClaw

[English](README.md)

把飞书等 IM 接到你们现有的 Claude / Codex 一类 Agent —— **可自托管的 IM→Agent 网关**。

> **现状：** MVP 已跑通（`chat` / `serve` / `healthz`）。适合私有化试点与内训，不是完整企业编排平台。

## 30 秒跑通

```bash
make tidy && make build && make test
./bin/meclaw chat -c examples/config.example.json
./bin/meclaw serve -c examples/config.example.json
```

## 3 分钟演示

**[录制步骤与分镜 → docs/demo-recording.md](docs/demo-recording.md)**（成片上传后，把下方链接换成 B 站/网盘 URL）

<!-- 上传 meclaw-demo.mp4 后取消注释并填 URL：
[观看演示](https://www.bilibili.com/video/BVxxxxxxxx)
-->

## 需要私有化 / 内训？

私有化试点 **¥39,800** 起 · 半天内训 **¥6,800**  
微信：**coder-hs**（也可提 Issue）  
[完整报价与交付范围 → docs/pricing.md](docs/pricing.md)

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
| [docs/demo-recording.md](docs/demo-recording.md) | **3 分钟演示录制脚本** |
| [docs/pricing.md](docs/pricing.md) | **报价与交付范围（对外）** |
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

# MeClaw

[中文](README_CN.md)

Connect Feishu and other IM channels to your existing Claude / Codex-style agents — a **self-hosted IM→Agent gateway**.

> **Status:** MVP is runnable (`chat` / `serve` / `healthz`). Good for private pilots and workshops; not a full enterprise orchestration platform.

## Run in 30 seconds

```bash
make tidy && make build && make test
./bin/meclaw chat -c examples/config.example.json
./bin/meclaw serve -c examples/config.example.json
```

## 3-minute demo

**[Script & QuickTime steps → docs/demo-recording.md](docs/demo-recording.md)** (after upload, add a public link below)

<!-- After uploading meclaw-demo.mp4, uncomment and set URL:
[Watch demo](https://www.bilibili.com/video/BVxxxxxxxx)
-->

## Private deploy / training?

Private pilot from **¥39,800** · Half-day workshop **¥6,800**  
WeChat: **coder-hs** (Issues welcome too)  
[Full pricing & delivery scope → docs/pricing.md](docs/pricing.md)

## Layout

```text
cmd/                 CLI (chat / serve / version)
internal/            A1–A5 core
deploy/              A6 Docker / K8s skeletons
skills/              scenario B skill packs
hosted/              invite / metering (later)
docs/                strategy + self-host path
examples/            sample config
```

## Docs

| Doc | Topic |
|-----|--------|
| [docs/demo-recording.md](docs/demo-recording.md) | **3-min demo recording script** |
| [docs/pricing.md](docs/pricing.md) | **Pricing & delivery scope** |
| [docs/dev-context.md](docs/dev-context.md) | **Dev handoff (read first)** |
| [docs/self-hosted-claw.md](docs/self-hosted-claw.md) | How to run self-hosted claw |
| [docs/infra-capability-map.md](docs/infra-capability-map.md) | Infra capability map |
| [docs/scenario-a-mvp.md](docs/scenario-a-mvp.md) | A1 MVP |
| [docs/architecture.md](docs/architecture.md) | Architecture |

## License

MIT

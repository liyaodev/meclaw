# MeClaw

[中文](README_CN.md)

**Product name:** meclaw (claw). **Category:** Agent Infra → Agent SaaS.

Self-hosted claw stack (scenarios A1–A6) plus a WeChat work assistant (scenario B) — **one repo**.

> Narrative: **run a self-hosted MeClaw → a work assistant inside chat.**

## Status

- **A1–A6** self-hosted claw path is runnable ([docs/self-hosted-claw.md](docs/self-hosted-claw.md))
- **B**: `skills/` examples; `hosted/` later

## Quick start

```bash
make tidy && make build && make test
./bin/meclaw chat -c examples/config.example.json
./bin/meclaw serve -c examples/config.example.json
```

Path: [docs/self-hosted-claw.md](docs/self-hosted-claw.md)  
Capability map: [docs/infra-capability-map.md](docs/infra-capability-map.md)

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
| [docs/dev-context.md](docs/dev-context.md) | **Dev handoff (read first)** |
| [docs/self-hosted-claw.md](docs/self-hosted-claw.md) | How to run self-hosted claw |
| [docs/infra-capability-map.md](docs/infra-capability-map.md) | Infra capability map |
| [docs/scenario-a-mvp.md](docs/scenario-a-mvp.md) | A1 MVP |
| [docs/architecture.md](docs/architecture.md) | Architecture |

## License

MIT

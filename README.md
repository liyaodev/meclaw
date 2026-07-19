# MeClaw

[中文](README_CN.md)

**Product name:** meclaw (claw). **Category:** Agent Infra → Agent SaaS.

IM → Agent gateway: bridge WeChat / Feishu / WeCom to local or remote agents (Claude, Codex, …), then package skills into a WeChat work assistant — **one repo for scenario A + B**.

> Narrative: **IM Agent infrastructure → a work assistant inside chat.**

## Status

Scenario A MVP: runtime pipeline + stdio/HTTP ingress + Feishu bot adapter. Scenario B (`skills/` + `hosted/`) later on the same codebase.

## Quick start

```bash
make tidy && make build
./bin/meclaw chat -c examples/config.example.json
# type a message; /agent <id> to switch

./bin/meclaw serve -c examples/config.example.json
# POST /v1/message  and optional /v1/feishu/event
```

```bash
make test
```

## Layout

```text
cmd/                 CLI (chat / serve / version)
internal/
  gateway/           IM normalization + stdio/HTTP + feishu/
  agent/             ACP (stub) / CLI / HTTP runners
  session/           chat sessions
  policy/            allow-lists / audit
  config/            runtime config
  runtime/           Message → policy → session → agent
skills/              scenario B skill packs
hosted/              invite / metering (later)
docs/                strategy + architecture
examples/            sample config
```

## Docs

| Doc | Topic |
|-----|--------|
| [docs/README.md](docs/README.md) | Doc index |
| [docs/scenario-a-mvp.md](docs/scenario-a-mvp.md) | Scenario A MVP how-to |
| [docs/scenario-a-infra-im-gateway.md](docs/scenario-a-infra-im-gateway.md) | Scenario A strategy |
| [docs/scenario-b-saas-wechat-assistant.md](docs/scenario-b-saas-wechat-assistant.md) | Scenario B |
| [docs/architecture.md](docs/architecture.md) | Architecture |

## License

MIT

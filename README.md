# MeClaw

[中文](README_CN.md)

**Product name:** meclaw (claw). **Category:** Agent Infra → Agent SaaS.

IM → Agent gateway: bridge WeChat / Feishu / WeCom to local or remote agents (Claude, Codex, …), then package skills into a WeChat work assistant — **one repo for scenario A + B**.

> Narrative: **IM Agent infrastructure → a work assistant inside chat.**

## Status

Scaffold only. Scenario A (gateway + routing) first; scenario B (`skills/` + `hosted/`) later on the same codebase.

## Quick start

```bash
make tidy
make run
# prints scaffold banner — implement gateway next
```

## Layout

```text
cmd/                 CLI entry
internal/
  gateway/           IM channel normalization
  agent/             ACP / CLI / HTTP runners
  session/           chat sessions
  policy/            allow-lists / audit (To B)
  config/            runtime config
skills/              scenario B skill packs
hosted/              invite / metering (later)
docs/                strategy + architecture
examples/            sample config
```

## Docs

| Doc | Topic |
|-----|--------|
| [docs/README.md](docs/README.md) | Doc index |
| [docs/naming.md](docs/naming.md) | claw vs Agent naming |
| [docs/scenario-a-infra-im-gateway.md](docs/scenario-a-infra-im-gateway.md) | Scenario A |
| [docs/scenario-b-saas-wechat-assistant.md](docs/scenario-b-saas-wechat-assistant.md) | Scenario B |
| [docs/playbook-flywheel-90d.md](docs/playbook-flywheel-90d.md) | 90-day playbook |

## License

MIT

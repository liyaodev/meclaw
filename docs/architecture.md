# 架构（草案）

```text
                    ┌─────────────┐
   WeChat/Feishu ──►│   gateway   │
   WeCom          │  (IM 归一化) │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │   policy    │  allow-list / audit
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
       ┌──────▼──┐   ┌─────▼────┐  ┌────▼─────┐
       │ session │   │  agent   │  │  skills  │
       │  store  │   │ ACP/CLI  │  │ (场景 B) │
       └─────────┘   │  /HTTP   │  └──────────┘
                     └─────┬────┘
                           │
                     Claude/Codex/…
```

## 模块

| 包 | 职责 | 场景 |
|----|------|------|
| `internal/gateway` | IM 接入与回复 | A |
| `internal/agent` | 后端 Agent 执行 | A |
| `internal/session` | 会话连续性 | A |
| `internal/policy` | 权限与审计 | A / To B |
| `internal/config` | 配置 | A |
| `skills/` | 技能包 | B |
| `hosted/` | 邀请制与计费 | B（后期） |

## 非目标（个人阶段）

- 通用 Agent 编排框架
- 开放注册的多租户云（过早）
- 与硬件仓强耦合

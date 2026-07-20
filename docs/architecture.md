# 架构（草案）

主路径：[self-hosted-claw.md](./self-hosted-claw.md)。能力矩阵：[infra-capability-map.md](./infra-capability-map.md)。

```text
   Feishu / HTTP / stdio
            │
            ▼
      ┌───────────┐
      │  gateway  │  A1
      └─────┬─────┘
            │
      ┌─────▼─────┐
      │  policy   │  A1 / A3
      └─────┬─────┘
            │
      ┌─────▼─────┐
      │  runtime  │  A1–A5 编排
      └─────┬─────┘
            │
   ┌────────┼────────┬──────────┬──────────┐
   │        │        │          │          │
session  memory   agent     tools*    orchestrate*
  A2       A2    A1/A2       A3*         A5*
            │        │
            │   openai/cli/http/acp
            │
      sandbox*  observe*  eval*     deploy*
        A3*       A4*      A4*        A6*
```

\*骨架或下一波次。编排入口：`internal/runtime.Handle`。

## 模块

| 包 | 职责 | 场景 |
|----|------|------|
| `internal/gateway` | IM / HTTP / stdio | A1 |
| `internal/agent` | CLI / HTTP / OpenAI / ACP | A1 / A2 |
| `internal/session` | 会话持久化 | A2 |
| `internal/memory` | 简易记忆 | A2 |
| `internal/policy` | 白名单 / 审计 | A1 / A3 |
| `internal/tools` | 工具注册（骨架） | A3 |
| `internal/sandbox` | 沙箱接口（骨架） | A3 |
| `internal/observe` | Trace（骨架） | A4 |
| `internal/eval` | 评测（骨架） | A4 |
| `internal/orchestrate` | 多 Agent 绑定（骨架） | A5 |
| `internal/config` | 配置 | A1–A6 |
| `internal/runtime` | 管道编排 | A1–A5 |
| `deploy/` | Compose / K8s | A6 |
| `skills/` | 技能包 | B |
| `hosted/` | 邀请制与计费 | B（后期） |

## 非目标

- 通用 LangChain 竞品叙事
- MeClaw 自营公有多租户云 / 开放注册 Token 池
- 与硬件仓强耦合

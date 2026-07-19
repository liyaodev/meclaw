# Agent Infra 能力矩阵

MeClaw 用场景 **A1–A6** 覆盖完整 Agent Infra；**B** 是叠在同一底座上的 SaaS 层。

主路径：[self-hosted-claw.md](./self-hosted-claw.md)。

| Infra 能力 | 场景 | 包路径 | 状态 |
|------------|------|--------|------|
| IM / HTTP / stdio 通道接入与归一化 | A1 | `internal/gateway` | **done** |
| 飞书机器人 webhook | A1 | `internal/gateway/feishu` | **done** |
| 微信扫码协议 | A1 | `internal/gateway/wechat`（预留） | later |
| 企微适配器 | A1 | `internal/gateway/wecom`（预留） | later |
| CLI / HTTP Agent 执行 | A1 | `internal/agent` | **done** |
| ACP 协议 | A1/A2 | `internal/agent` | later（占位错误） |
| 用户白名单 / 审计日志 | A1/A3 | `internal/policy` | **done** |
| 配置加载 | A1–A6 | `internal/config` | **done** |
| 管道编排 | A1–A5 | `internal/runtime` | **done** |
| OpenAI-compatible 模型接入与路由 | A2 | `internal/agent` | **done** |
| 会话持久化 | A2 | `internal/session` | **done** |
| 记忆（按 chat） | A2 | `internal/memory` | **done** |
| 工具注册与调用 | A3 | `internal/tools` | **done** |
| 工具白名单执行面 | A3 | `internal/policy` + `tools` | **done** |
| 进程沙箱 | A3 | `internal/sandbox` | **done** |
| 浏览器操控 | A3 | `internal/sandbox` | later（接口） |
| 电脑操控 | A3 | `internal/sandbox` | later（接口） |
| Trace / 结构化观测 | A4 | `internal/observe` | **done** |
| 健康检查扩展 | A4 | `cmd` + `observe` | **done** |
| 评测 Eval | A4 | `internal/eval` + `evals/` | **done** |
| 多 Agent 绑定 / 编排 | A5 | `internal/orchestrate` | **done** |
| Skills 钩子 | A5 / B | `skills/` + `orchestrate` | **done** |
| Docker Compose 私有化 | A6 | `deploy/docker` | **done** |
| K8s 清单部署 | A6 | `deploy/k8s` | **done**（骨架） |
| 单集群多 workspace | A6 | — | later |
| 邀请制托管 / 计费 | B | `hosted/` | later |
| 成品技能包（周报等） | B | `skills/examples` | 示例 |

## 状态说明

| 状态 | 含义 |
|------|------|
| **done** | 可测、可演示 |
| **later** | 地图内保留，不阻塞主干 |

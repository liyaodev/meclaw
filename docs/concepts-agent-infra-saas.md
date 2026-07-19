# Agent Infra + Agent SaaS（概念）

## Agent Infra（基础设施层）

给别人「造 Agent / 跑 Agent」用的底座，例如：

- 模型接入与路由
- 工具调用、记忆、编排
- 沙箱 / 浏览器 / 电脑操控
- 可观测、评测、权限与多租户

卖的是能力与平台，客户是开发者或要自建 Agent 的公司。

## Agent SaaS（应用层）

面向最终业务场景的成品 Agent 产品，例如客服、编程、销售、运维助手。用户直接用，按席位/用量订阅，不一定自己搭底层。

## 「Agent Infra + Agent SaaS」

通常指同一家公司（或同一战略）两边都做：底层平台 + 上层垂直产品。

- **Infra** 沉淀复用能力
- **SaaS** 拿场景和收入

也有人用这句话概括整条 Agent 价值链。

## 与 MeClaw 的关系

MeClaw 是这条价值链上的**单一开源仓库**。产品承诺是 **完整 Agent Infra**，组织方式是「自建 claw」路径上的多场景，而不是只做一个 IM 插件。

| 层 | 场景 | 卖什么 |
|----|------|--------|
| Infra | **A1–A6** | 自建 claw 全栈能力（开发者 / 企业 IT） |
| SaaS | **B** | 微信工作助手（最终用户 → 小 B） |

Infra 场景拆分：

| 场景 | 内容 |
|------|------|
| **A1** Channel Gateway | IM / HTTP / stdio 接入（第一英里） |
| **A2** Model / Session / Memory | 模型、会话持久化、记忆 |
| **A3** Tools / Policy / Sandbox | 工具、权限、沙箱 |
| **A4** Observability / Eval | Trace、健康、评测 |
| **A5** Multi-Agent Orchestration | 多 Agent 绑定与 skills 钩子 |
| **A6** Cloud / Private Deploy | 用户自有云主机 / K8s 私有化部署 |

主路径文档：[self-hosted-claw.md](./self-hosted-claw.md)。能力矩阵：[infra-capability-map.md](./infra-capability-map.md)。

做事方法论见 [methodology.md](./methodology.md)。

硬件互动线（抓娃娃 / 遥控车等）独立仓库；可后接「远程设备技能」，不强绑。

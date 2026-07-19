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

MeClaw 是这条价值链上的**单一开源仓库**。

| 层 | 产品场景 | 卖什么 | 备注 |
|----|----------|--------|------|
| Infra | **A（A1–A6）** | 自建 claw 能力 | 主干；To B 私有化挂这里 |
| SaaS | **B** | 微信工作助手 | **为搞懂 Agent SaaS 业务**的落地项目；To C 挂这里 |

**To B / To C 不是第三、第四个产品场景**——它们是 [methodology.md](./methodology.md) 里的变现环节：大 To B 偏 A，邀请制 To C 偏 B。

Infra 关卡拆分见 [self-hosted-claw.md](./self-hosted-claw.md)；能力矩阵见 [infra-capability-map.md](./infra-capability-map.md)。

做事方法论见 [methodology.md](./methodology.md)。场景 B 说明见 [scenario-b-saas-wechat-assistant.md](./scenario-b-saas-wechat-assistant.md)。

硬件互动线（抓娃娃 / 遥控车等）独立仓库；可后接「远程设备技能」，不强绑。

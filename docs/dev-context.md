# 开发上下文（handoff）

> 给后续开发 / Agent 用的仓库现状快照。先读本文，再按需下钻其它 docs。  
> 更新日期：2026-07-19。  
> **侧业收款 / 赚钱 Agent：** 已迁至私有仓 `earnclaw`（与本仓并列：`../earnclaw`），勿把 KPI/名单写回本开源仓。

## 1. 产品在做什么

| 项 | 内容 |
|----|------|
| 产品名 | **meclaw**（claw） |
| 品类 | Agent Infra → Agent SaaS |
| 主叙事 | **如何跑通自建 MeClaw → 微信里能干活的助手** |
| 仓库策略 | **单仓**：Infra A1–A6 + SaaS B；不做无关赛道；不做自营公有多租户云 |

对照品类全景与 MeClaw 切片：[concepts-agent-infra-saas.md](./concepts-agent-infra-saas.md)、[methodology.md](./methodology.md)。

## 2. 当前结论（读完这一节即可开工）

**主干 0→1 已完成（MVP 级全覆盖）：** A1→A6 可测、可演示。

**产品场景只有 2 个：** A = Infra；B = SaaS 业务练习项目。To B / To C / 课 / 自媒体是方法论环节，**不要再开场景 C/D**。

**A 看起来像聊天产品 ≠ A 是 SaaS：** A 卖能力与自托管；B 才卖「扫码即用的结果」并用来搞懂 SaaS 生意。详见 [methodology.md](./methodology.md)。

**不是生产级「每一项都做完」：** 微信/企微、真 ACP、浏览器/电脑操控、向量记忆、强隔离沙箱、多 workspace 等仍是 **later**。

| 判断 | 说明 |
|------|------|
| Infra 主干 | 已通，优先**加深**而非再贴新场景标签 |
| 场景 B | 后置；为懂 Agent SaaS；`skills/` 示例，`hosted/` 空壳 |
| 加深候选 | 微信或企微通道 / 真 ACP / 沙箱或观测加深 |

能力明细：[infra-capability-map.md](./infra-capability-map.md)。  
关卡命令：[self-hosted-claw.md](./self-hosted-claw.md)。

## 3. 场景地图

```text
A1 Channel Gateway     → 消息进得来、回得出（飞书/HTTP/stdio）
A2 Model/Session/Memory → openai + 文件会话 + 简易记忆
A3 Tools/Policy/Sandbox → /tool + 白名单 + 进程沙箱
A4 Observability/Eval   → Trace JSONL + healthz + meclaw eval
A5 Orchestration        → bindings + skill 注入 + /skill
A6 Cloud Private Deploy → Compose + .env；K8s 骨架
B  WeChat Assistant     → skills/ + hosted/（后置）
```

编排入口：**`internal/runtime.Handle`**  
管道：`Message → policy → (tool|skill|agent) → session/memory → reply`（并写 audit / trace）。

## 4. 代码地图

```text
main.go                 → cmd.Execute()
cmd/                    chat | serve | eval | version
internal/
  gateway/              归一化 Message；stdio；HTTP /v1/message；feishu/
  agent/                cli | http | openai | acp(占位)
  session/              MemoryStore + FileStore（data/sessions）
  memory/               FileStore JSONL（data/memory）
  policy/               用户/工具白名单 + Auditor
  tools/                MapRegistry；内置 echo / shell
  sandbox/              LocalExecutor 命令白名单；Browser/Computer 仅接口
  observe/              JSONLTracer；Health JSON
  eval/                 Case Runner；evals/smoke.json
  orchestrate/          RuleResolver；SkillLoader
  config/               JSON Load + Validate
  runtime/              上述全部接线
deploy/docker|k8s/      A6
skills/                 场景 B 技能示例
hosted/                 邀请制壳（后期）
examples/config.example.json
evals/smoke.json
data/                   本地运行时（gitignore）
```

## 5. 关键行为约定

| 能力 | 用法 / 约定 |
|------|-------------|
| 切换 Agent | `/agent <id> [prompt]` |
| 调工具 | `/tool echo text=hi`；`/tool shell command=echo args=hi` |
| 临时技能 | `/skill <path-under-skills_dir> [prompt]` |
| 会话 key | `{channel}:{chatID}:{userID}` |
| OpenAI 密钥 | 配置 `api_key` 或环境变量 `MECLAW_OPENAI_API_KEY` |
| 落盘 | `data_dir` 默认 `./data`（sessions / memory / traces） |
| 策略空列表 | `allow_users` / `allow_tools` 为空 = **全放行**（开发默认） |
| 沙箱默认 | `allow_commands`: echo / date / uname（`Load` 时补齐） |
| 绑定优先级 | user > chat > channel（`orchestrate.RuleResolver`） |
| Agent 挂技能 | `agents.<id>.skill` 相对 `skills_dir` |

配置样例：[examples/config.example.json](../examples/config.example.json)。

## 6. 常用命令

```bash
make tidy && make build && make test

./bin/meclaw chat -c examples/config.example.json
./bin/meclaw serve -c examples/config.example.json
./bin/meclaw eval -c examples/config.example.json -f evals/smoke.json
./bin/meclaw version

# HTTP
curl -s http://127.0.0.1:8080/healthz
curl -s -X POST http://127.0.0.1:8080/v1/message \
  -H 'content-type: application/json' \
  -d '{"user_id":"u1","chat_id":"c1","text":"hello"}'

# 云主机 Compose：docs/deploy-cloud-30m.md
```

模块：`github.com/meclaw/meclaw`，Go **1.22+**，CLI 用 cobra。

## 7. 明确不做（防跑偏）

- 通用 LangChain 竞品叙事  
- MeClaw **自营**公有多租户云 / 开放注册 Token 池  
- 为「看起来像完整 Infra」而堆未接线的空包（later 项保持诚实）  
- 硬件仓强耦合（可后接技能，不强绑）

## 8. 建议的下一波工作（任选加深）

1. **通道：** 微信或企微适配器（契约同飞书：归一化 → `runtime.Handle` → 回写）  
2. **协议：** 真 ACP，替换占位错误  
3. **执行面：** 更强沙箱或真实浏览器自动化  
4. **记忆：** 摘要 / 向量检索，替代纯 JSONL 追加  
5. **场景 B：** 单一可演示技能 + 邀请制 hosted 壳  

改代码时同步更新：`infra-capability-map.md` 状态列、`self-hosted-claw.md` 对应关卡勾选。

## 9. 文档索引（按需）

| 文档 | 用途 |
|------|------|
| [self-hosted-claw.md](./self-hosted-claw.md) | 关卡验收命令 |
| [infra-capability-map.md](./infra-capability-map.md) | done / later 矩阵 |
| [architecture.md](./architecture.md) | 模块边界 |
| [deploy-cloud-30m.md](./deploy-cloud-30m.md) | 私有化部署 |
| [scenario-a-mvp.md](./scenario-a-mvp.md) | A1 飞书/HTTP 细节 |
| [scenario-b-saas-wechat-assistant.md](./scenario-b-saas-wechat-assistant.md) | B 策略 |
| [naming.md](./naming.md) | claw vs Agent 用词 |

`todo/` 目录可能含外部参考材料，**不是**当前 meclaw 运行时的一部分；以 `cmd/` + `internal/` + 本文为准。

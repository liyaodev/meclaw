# 如何跑通自建 MeClaw

> 主线：本机 → 可干活 → 可观测 → 可上云私有化（对标 openclaw 式自托管）。

完整能力归属见 [infra-capability-map.md](./infra-capability-map.md)。

## 关卡总览

| 关卡 | 场景 | 目标 | 状态 |
|------|------|------|------|
| 1 | [A1](./scenario-a-mvp.md) Channel Gateway | 消息进得来、回得出 | **已通** |
| 2 | A2 Model / Session / Memory | 换模型能聊、重启不丢上下文 | **已通** |
| 3 | A3 Tools / Policy / Sandbox | Agent 能安全执行工具 | **已通** |
| 4 | A4 Observability / Eval | 出事能查、可回归 | **已通** |
| 5 | A5 Multi-Agent Orchestration | 绑定规则 + skills 钩子 | **已通** |
| 6 | A6 Cloud / Private Deploy | 云主机一键起 | **已通（Compose）** |
| — | [B](./scenario-b-saas-wechat-assistant.md) SaaS | 微信里出活 | 后置 |

---

## 关卡 1｜A1 Channel Gateway（已通）

```bash
make tidy && make build && make test
./bin/meclaw chat -c examples/config.example.json
./bin/meclaw serve -c examples/config.example.json
```

详见 [scenario-a-mvp.md](./scenario-a-mvp.md)。

**验收：** stdio / HTTP / 飞书单测 — 已勾选。

---

## 关卡 2｜A2 Model / Session / Memory（已通）

- `data_dir`、`memory.*`、`mode: openai` + `MECLAW_OPENAI_API_KEY`

```bash
./bin/meclaw chat -c examples/config.example.json
ls data/sessions data/memory
```

**验收：** openai 单测 / 会话落盘 / 记忆注入 — 已勾选。

---

## 关卡 3｜A3 Tools / Policy / Sandbox（已通）

- `sandbox.allow_commands`：进程白名单
- `policy.allow_tools`：工具白名单（空 = 全放行）
- 消息：`/tool echo text=hi`、`/tool shell command=echo args=hi`

```bash
./bin/meclaw chat -c examples/config.example.json
# > /tool echo text=ping
# > /tool shell command=date
```

**验收勾选**

- [x] 工具注册表 List / Call
- [x] `allow_tools` 拦截未授权工具
- [x] 本地 process 沙箱白名单
- [x] 浏览器 / 电脑操控接口占位（`sandbox.BrowserController` / `ComputerController`）

---

## 关卡 4｜A4 Observability / Eval（已通）

- Trace：`data/traces/traces.jsonl`
- `/healthz`：JSON（version + data_dir）
- Eval：`meclaw eval -f evals/smoke.json`

```bash
make build
./bin/meclaw eval -c examples/config.example.json -f evals/smoke.json
./bin/meclaw serve -c examples/config.example.json
# curl localhost:8080/healthz
```

**验收勾选**

- [x] Trace JSONL 落盘
- [x] `/healthz` 含版本与 data_dir
- [x] `evals/smoke.json` 可跑分

---

## 关卡 5｜A5 Multi-Agent Orchestration（已通）

- `bindings`：channel / chat / user → agent
- `agents.*.skill` + `skills_dir`：自动注入技能说明
- `/skill <path> [prompt]`：临时加载技能包

```bash
# config 中 weekly.skill = examples/weekly-report
./bin/meclaw chat -c examples/config.example.json
# > /agent weekly
# > /skill examples/weekly-report 本周完成了 A3
```

**验收勾选**

- [x] binding 规则解析
- [x] skills 目录 README/SKILL.md 可挂接

---

## 关卡 6｜A6 Cloud / Private Deploy（已通）

见 **[deploy-cloud-30m.md](./deploy-cloud-30m.md)**。

```bash
cp examples/config.example.json deploy/docker/config.json
cp deploy/docker/.env.example deploy/docker/.env
cd deploy/docker && docker compose --env-file .env up -d --build
```

**验收勾选**

- [x] Compose + Dockerfile 可构建启动路径
- [x] `.env` 注入密钥约定
- [x] 云主机 30 分钟文档
- [x] K8s Deployment 骨架

---

## 明确不做

- MeClaw 自营公有多租户云 / 开放注册 Token 池
- 把叙事做成通用 LangChain 竞品

# A1 MVP：Channel Gateway + 飞书

> 自建 claw **关卡 1**。总路径：[self-hosted-claw.md](./self-hosted-claw.md)。A2+ 见同文档后续关卡。

本文描述已实现的 **A1 Channel Gateway** MVP。

## 管道

```text
IM / stdio / HTTP
       │
       ▼
  gateway（归一化 Message）
       │
       ▼
  policy（用户白名单 + 审计）
       │
       ▼
  session（chat 级会话 id / agent 选择）
       │
       ▼
  agent router（cli | http | acp占位）
       │
       ▼
  回复 → 通道（stdio 打印 / HTTP JSON / 飞书 reply）
```

会话 key：`{channel}:{chatID}:{userID}`。

临时切换 Agent：消息以 `/agent <id>` 开头（可带后续 prompt）。

## 本地试用

```bash
make tidy && make build
./bin/meclaw chat -c examples/config.example.json
```

示例配置默认 Agent 为 `echo`（本机 `echo` 命令），无需真实模型即可验证管道。

HTTP 入站：

```bash
./bin/meclaw serve -c examples/config.example.json
curl -s -X POST http://127.0.0.1:8080/v1/message \
  -H 'content-type: application/json' \
  -d '{"user_id":"u1","chat_id":"c1","text":"hello"}'
```

## 飞书接入

1. 在[飞书开放平台](https://open.feishu.cn/)创建企业自建应用，开通「获取与发送单聊、群组消息」等机器人权限。
2. 事件订阅选择「将事件发送至开发者服务器」，请求地址指向公网可达的：

   `https://<你的域名>/v1/feishu/event`

3. 订阅事件：`im.message.receive_v1`。
4. 把 `app_id` / `app_secret` / `verification_token` 写入配置 `gateway.channels.feishu`（可用 `encrypt_key` 字段占位；当前 MVP 处理明文事件体）。
5. 启动：

```bash
./bin/meclaw serve -c /path/to/config.json
```

未填写飞书凭证时，`serve` 仍会启动通用 `/v1/message` 与 `/healthz`，并日志提示跳过飞书。

可选：`./bin/meclaw serve --stdio` 同时开本地交互。

## Agent 模式

| mode | 行为 |
|------|------|
| `cli` | 执行 `command` + `args`，把 prompt 作为最后一个参数，stdout 为回复 |
| `http` | `POST base_url`，JSON `{prompt,session,agent_id}` → `{text}` |
| `acp` | 占位，返回明确「未实现」错误 |

## 后续 IM 约定

微信 / 企微适配器应实现与飞书相同的契约：

1. 将通道事件归一化为 `gateway.Message`
2. 调用同一 `runtime.Handle`
3. 把回复写回各自通道 API

接口骨架可放在 `internal/gateway/wechat/`、`internal/gateway/wecom/`，不改变核心管道。

## 测试

```bash
make test
```

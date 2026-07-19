# A6 — 云主机 30 分钟跑通

目标：在一台云主机上用 Docker Compose 跑起 MeClaw（HTTP 入站 + 可选飞书）。

## 前置

- 云主机：2C2G 即可（Ubuntu 22.04+）
- 安装 Docker + Compose 插件
- 安全组放行 `8080`（或前面挂 Nginx/Caddy）

## 步骤

### 1. 拉取代码

```bash
git clone <your-meclaw-repo> meclaw && cd meclaw
```

### 2. 准备配置与密钥

```bash
cp examples/config.example.json deploy/docker/config.json
cp deploy/docker/.env.example deploy/docker/.env
# 编辑 .env：MECLAW_OPENAI_API_KEY=...
# 编辑 config.json：gateway.listen 保持 :8080；按需填 feishu
```

密钥**不要**写进 Git。Compose 通过环境变量注入 `MECLAW_OPENAI_API_KEY`。

### 3. 构建并启动

```bash
cd deploy/docker
docker compose --env-file .env up -d --build
curl -s http://127.0.0.1:8080/healthz
curl -s -X POST http://127.0.0.1:8080/v1/message \
  -H 'content-type: application/json' \
  -d '{"user_id":"u1","chat_id":"c1","text":"hello"}'
```

### 4. 飞书（可选）

公网 URL：`https://<域名>/v1/feishu/event` → 填入飞书事件订阅；凭证写入 `config.json` 的 `gateway.channels.feishu`。

### 5. 数据卷

会话 / 记忆 / Trace 在命名卷 `meclaw-data` → 容器内 `/data`。

```bash
docker compose exec meclaw ls /data
```

## 验收勾选

- [ ] `healthz` 返回 JSON（含 version / data_dir）
- [ ] `/v1/message` 可回复
- [ ] `.env` 未入库
- [ ] 重启容器后 `data/sessions` 仍在

## K8s

见 [k8s/README.md](../k8s/README.md)（清单骨架；首发建议 Compose）。

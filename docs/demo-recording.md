# 3 分钟演示怎么录（Mac · QuickTime）

目标：陌生人 **2–3 分钟**内看懂 **是啥 → 能聊 → 能调 tool → HTTP 网关**。文件名 **`meclaw-demo.mp4`**。

> **时长实话：** 只录终端、不口播、狂敲命令，**凑不满 3 分钟**（往往 1 分钟内打完）。  
> **「3 分钟片」** = 下面分镜里的 **口播 + 等输出 + 多一轮 HTTP**，按表走自然到 **2:30～3:30**；赶 Day 1 可录 **90 秒精华版**（见文末）。

## 用什么录

**QuickTime Player**（系统自带）：

1. **文件 → 新建屏幕录制** → 选项可勾「显示鼠标点按」
2. **录制所选部分**，框住终端（字体 **18–20pt**）；口播多就 **终端 + README 首屏** 各录一段再剪，或一镜到底先只录终端
3. 存 **`meclaw-demo.mp4`**

备选：**⌘⇧5**；手机录屏 Day 1 可接受。

录前：关通知；`cd` 到 meclaw；**先跑通一遍**（别边录边 `make build`）。

```bash
cd /path/to/meclaw
make build
./bin/meclaw version
# 演示前清本地会话/memory（同一 key：stdio:stdio:local，否则会带上 first/second/quit 等旧记录）
rm -rf ./data
```

**段 A 若出现一大坨 `Prior context` / 嵌套 `meclaw:`：** 符合当前 MVP 行为（见下「chat 输出说明」），录片前 **`rm -rf ./data`** 会干净很多。

### chat 输出说明（是不是坏了？）

| 现象 | 是否正常 |
|------|----------|
| 输入 `quit` / `exit` 仍当普通消息、不退出 | **正常** — 只有 **Ctrl-D**（EOF）退出，提示里写的是 Ctrl-D |
| 回复里反复出现 `Prior context:`、`user: first`、`User: 你好` 等 | **正常但脏** — 示例配置 **memory 开着**，历史写在 `./data`；echo agent 会把整段 prompt（含历史）原样打出，**上一轮 assistant 全文也会进 memory**，多几轮会越滚越长 |
| 前两轮演示句被埋在旧上下文里 | **录前 `rm -rf ./data`**，或临时把 config 里 `memory.enabled` 设为 `false` |
| `/tool echo text=ping` 只输出 `ping` | **正常** — tool 不走 echo agent，也不拼那段 Prior context |

接真 LLM（OpenAI 等）时 memory 走 messages 数组，不会出现这种 echo 雪球；**演示用 echo 时以短对话 + 清 data 为准**。

---

## 分镜表（推荐 · 约 3 分钟）

按 **时间轴** 录；「口播」可后期配音，Day 1 直接对着麦念也行。

| 时间 | 画面 | 做什么 | 口播（示例） |
|------|------|--------|----------------|
| 0:00–0:25 | README 或终端 | 展示仓库首屏或 `cat README_CN.md \| head -22` | 「MeClaw：把飞书等 IM 接到你们现有的 Claude/Codex Agent，可自托管。这是 MVP，适合私有化试点。」 |
| 0:25–0:40 | 终端 | `./bin/meclaw version` | 「一条二进制，chat 本地调试，serve 对外 HTTP。」 |
| 0:40–1:35 | 终端 | **段 A：`chat`**（见下） | 「同一套 runtime：会话、策略、可选 memory。」 |
| 1:35–2:05 | 终端 | **段 B：tool + 换 Agent** | 「白名单 tool；Agent 可切换。」 |
| 2:05–2:50 | 双终端或分屏 | **段 C：`serve` + curl** | 「IM 未配时也能用 HTTP 入站，和飞书走同一条链路。」 |
| 2:50–3:10 | README / pricing | 滚到「需要私有化/内训」或 `docs/pricing.md` 表头 | 「试点 39800 起，范围在 pricing 页写死，减少扯皮。」 |

### 段 A：`chat`（约 55 秒，含打字与看输出）

```bash
./bin/meclaw chat -c examples/config.example.json
```

依次输入（**每句回车后停 2–3 秒**，让观众看清回复）：

```text
你好，我是第一次用 MeClaw
今天我们要录一个演示
/tool echo text=ping
```

- 前两轮（**清过 data 后**）：短 **`meclaw: 你好，我是第一次用 MeClaw`**，第二轮回复里能带上一句上下文；仍可能带一行 `Prior context:`，属 echo + memory 的 MVP 表现。
- 第三轮：直接 **`ping`**（tool 路径）。
- 审计 JSON 在 **stderr**；要干净画面就单窗 stdout，要 B 端感可保留 stderr。

**Ctrl-D** 退出。

### 段 B：换 Agent + shell tool（约 30 秒）

重新进 chat，或接着上面（若未退出）：

```text
/agent weekly 写一句周报标题
/tool shell command=date
```

- `weekly` 会走配置里的 skill 示例（echo 前缀 **`weekly:`**）。
- `date` 需在 sandbox 白名单内（示例配置已允许）。

**Ctrl-D** 退出。

### 段 C：`serve` + HTTP（约 45 秒）

**终端 A：**

```bash
./bin/meclaw serve -c examples/config.example.json
```

等出现 **listen :8080**（口播：「飞书凭证空时会跳过 webhook，但 HTTP 仍在。」）

**终端 B：**

```bash
curl -s --noproxy '*' http://127.0.0.1:8080/healthz | python3 -m json.tool

curl -s --noproxy '*' -X POST http://127.0.0.1:8080/v1/message \
  -H 'content-type: application/json' \
  -d '{"user_id":"u1","chat_id":"c1","text":"HTTP 入站测试"}'
```

- `healthz`：JSON 含 **`ok`**、**version**、**data_dir**。
- `/v1/message`：返回 JSON **reply**（与 IM 归一后的同一条 runtime）。

**Ctrl-C** 停 serve。

---

## 90 秒精华版（不口播、不 README）

适合 B 站「先上片」：

1. `version`（5s）
2. `chat`：`你好` → `/tool echo text=ping` → **Ctrl-D**（40s，含停顿）
3. `serve` + `healthz` + 一条 `POST /v1/message`（40s）
4. 片尾字幕：微信 **coder-hs** + README 链接（5s）

---

## 导出与 README

1. 导出 **`meclaw-demo.mp4`**（1080p，勿提交 Git，已 gitignore `*.mp4`）。
2. 上传 B 站 / 网盘 / 公众号素材库。
3. 在 [README_CN.md](../README_CN.md) / [README.md](../README.md) **「3 分钟演示」** 取消注释，填 `[观看演示](URL)`。

B 站示例：`https://www.bilibili.com/video/BVxxxxxxxx`

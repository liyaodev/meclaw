# 单仓管理 Infra A1–A6 + SaaS B

**结论：合适。** 个人开发者阶段，完整自建 claw 与场景 B 放在 **一个 GitHub 仓库 `meclaw`**。

## 目录约定

```text
meclaw/
├── README.md / README_CN.md
├── docs/                 # 自建路径 + 能力矩阵 + 架构
├── cmd/ + internal/      # Infra A1–A5
├── deploy/               # A6 Docker / K8s
├── skills/               # 场景 B 技能
├── examples/             # 示例配置
└── hosted/               # 邀请制壳（后期；密钥不入库）
```

## 何时再拆第 2 个仓库

- 开源协议要求托管侧闭源且混仓说不清
- 支付密钥不能进公开仓库（可先私有 fork）
- To B 合同要「开源核心 + 商业插件」分开授权

在此之前：**不要拆。**

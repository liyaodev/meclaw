# A6 Docker Compose

Quick path: [docs/deploy-cloud-30m.md](../../docs/deploy-cloud-30m.md)

```bash
cp ../../examples/config.example.json ./config.json
cp .env.example .env
docker compose --env-file .env up -d --build
```

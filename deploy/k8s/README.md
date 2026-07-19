# A6 Kubernetes skeleton

Apply after building/pushing `meclaw:dev` and creating ConfigMap/Secret:

```bash
kubectl create configmap meclaw-config --from-file=config.json=../../examples/config.example.json
kubectl apply -f deployment.yaml
```

Prefer Compose for first private deploy: [docs/deploy-cloud-30m.md](../../docs/deploy-cloud-30m.md).

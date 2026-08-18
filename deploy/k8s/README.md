# Kubernetes (kind) deployment

## Prerequisites (one-time on your machine)

### 1. Install kubectl
```bash
# Linux
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl && sudo mv kubectl /usr/local/bin/

# macOS
brew install kubectl
```

### 2. Install kind
```bash
# Linux
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.27.0/kind-linux-amd64
chmod +x kind && sudo mv kind /usr/local/bin/

# macOS
brew install kind
```

Docker must be running.

## Deploy

```bash
# From repo root

# 1. Create kind cluster (maps host 8080 -> NodePort 30080)
kind create cluster --config deploy/k8s/kind-config.yaml

# 2. Build app image and load into kind
docker build -t task-management-api:local .
kind load docker-image task-management-api:local --name task-management

# 3. Apply manifests (order matters)
kubectl apply -f deploy/k8s/00-namespace.yaml
kubectl apply -f deploy/k8s/01-configmap.yaml
kubectl apply -f deploy/k8s/02-secret.yaml
kubectl apply -f deploy/k8s/03b-postgres-init.yaml
kubectl apply -f deploy/k8s/03-postgres.yaml
kubectl apply -f deploy/k8s/04-redis.yaml
kubectl apply -f deploy/k8s/05-app.yaml

# 4. Wait for pods
kubectl -n task-management get pods -w
# Wait until postgres, redis, task-management-api are Running / Ready

# 5. Smoke test (NodePort mapped to localhost:8080)
curl http://localhost:8080/check_health
curl http://localhost:8080/metrics | head
```

## Useful commands

```bash
kubectl -n task-management get all
kubectl -n task-management logs -l app=task-management-api --tail=50
kubectl -n task-management describe pod -l app=task-management-api

# Tear down
kind delete cluster --name task-management
```

## What this demonstrates

- Deployment + Service + ConfigMap + Secret
- Readiness/liveness probes
- Config injection (no baked secrets in image)
- Multi-replica API Deployment
- Local kind cluster with host port mapping

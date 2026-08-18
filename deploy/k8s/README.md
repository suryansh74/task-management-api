# Kubernetes (kind) deployment

## Prerequisites (one-time)

```bash
# kubectl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl && sudo mv kubectl /usr/local/bin/

# kind
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.27.0/kind-linux-amd64
chmod +x kind && sudo mv kind /usr/local/bin/
```

Docker must be running.

## Free ports before creating the cluster

`docker compose` often holds **8000** and **50051**. kind maps host ports and will fail if they are taken.

```bash
# Stop compose stack for this project (preferred)
docker compose down

# Or free specific ports (Linux)
sudo fuser -k 8080/tcp 2>/dev/null
sudo fuser -k 18080/tcp 2>/dev/null
sudo fuser -k 50051/tcp 2>/dev/null
sudo fuser -k 15051/tcp 2>/dev/null
sudo fuser -k 6443/tcp 2>/dev/null

# See what holds a port
sudo ss -tlnp | grep -E '8080|18080|50051|15051' || true
```

If an old kind cluster exists:

```bash
kind delete cluster --name task-management
```

## Deploy (from repo root)

```bash
git pull
# already inside task-management-api — do NOT cd task-management-api again

# 0. Free ports / stop compose
docker compose down
sudo fuser -k 18080/tcp 15051/tcp 2>/dev/null || true
kind delete cluster --name task-management 2>/dev/null || true

# 1. Create kind cluster (host 18080 -> API, 15051 -> gRPC)
kind create cluster --config deploy/k8s/kind-config.yaml

# 2. Build image and load into kind
docker build -t task-management-api:local .
kind load docker-image task-management-api:local --name task-management

# 3. Apply manifests
kubectl apply -f deploy/k8s/00-namespace.yaml
kubectl apply -f deploy/k8s/01-configmap.yaml
kubectl apply -f deploy/k8s/02-secret.yaml
kubectl apply -f deploy/k8s/03b-postgres-init.yaml
kubectl apply -f deploy/k8s/03-postgres.yaml
kubectl apply -f deploy/k8s/04-redis.yaml
kubectl apply -f deploy/k8s/05-app.yaml

# 4. Wait for pods
kubectl -n task-management get pods -w

# 5. Smoke test
curl http://localhost:18080/check_health
curl http://localhost:18080/metrics | head
```

| Endpoint | URL |
|----------|-----|
| REST health | http://localhost:18080/check_health |
| Metrics | http://localhost:18080/metrics |
| gRPC | localhost:15051 |

## Useful commands

```bash
kubectl -n task-management get all
kubectl -n task-management logs -l app=task-management-api --tail=50

kind delete cluster --name task-management
```

## Port map

| Host port | Purpose |
|-----------|---------|
| 18080 | REST (NodePort 30080) |
| 15051 | gRPC (NodePort 30051) |

These avoid default docker-compose ports 8000 / 50051.

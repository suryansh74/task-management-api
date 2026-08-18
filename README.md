# Task Management API

Production-style task management backend in Go with **REST + gRPC**, Redis caching, session auth, rate limiting, Prometheus metrics, and Grafana.

## Features

- Full CRUD for tasks (PostgreSQL)
- Redis cache-aside pattern (configurable TTL)
- Session-based auth (HTTP-only cookies + Redis)
- Rate limiting (Redis token bucket)
- **Ports & Adapters** architecture (handlers, services, repositories)
- **gRPC TaskService** alongside REST (port 50051)
- **Prometheus** metrics (`/metrics`) + **Grafana** dashboards
- Structured logging (Zerolog)
- Unit tests (mocks) + integration tests (testcontainers)
- Docker Compose for full local stack

## Tech Stack

| Layer | Technology |
|-------|------------|
| Language | Go 1.25+ |
| REST | Fiber v2 |
| gRPC | google.golang.org/grpc |
| DB | PostgreSQL 16 |
| Cache / Sessions | Redis 7 |
| Metrics | Prometheus + Grafana |
| Validation | go-playground/validator |
| Logging | zerolog |

## Architecture

```
REST handlers  ──┐
                 ├──► ports.TaskService / UserService ──► service ──► repositories
gRPC TaskServer ─┘                                              ├── Postgres
                                                                └── Redis
```

## Quick Start (Docker)

```bash
git clone https://github.com/suryansh74/task-management-api.git
cd task-management-api
cp .env.example .env   # optional
docker compose up -d --build
```

| Service | URL |
|---------|-----|
| REST API | http://localhost:8000 |
| gRPC | localhost:50051 |
| Metrics | http://localhost:8000/metrics |
| Prometheus | http://localhost:9090 |
| Grafana | http://localhost:3000 (admin / admin) |

`init.sql` is mounted into Postgres so tables are created automatically on first start.

## API (REST)

### Auth
| Method | Path | Auth |
|--------|------|------|
| POST | `/register` | No |
| POST | `/login` | No |
| POST | `/logout` | Yes |

### Tasks
| Method | Path | Auth |
|--------|------|------|
| GET | `/tasks` | Yes |
| POST | `/tasks` | Yes |
| GET | `/tasks/:id` | Yes |
| PUT | `/tasks/:id` | Yes |
| DELETE | `/tasks/:id` | Yes |

### Health & Metrics
| Method | Path |
|--------|------|
| GET | `/check_health` |
| GET | `/metrics` |

### Example

```bash
curl -X POST http://localhost:8000/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Suryansh","email":"suryansh@example.com","password":"password123"}'

curl -X POST http://localhost:8000/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{"email":"suryansh@example.com","password":"password123"}'

curl -X POST http://localhost:8000/tasks \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"title":"Ship gRPC","content":"Days 3-4 done"}'
```

## gRPC

Service: `task.v1.TaskService`  
Methods: `GetTasks`, `GetTask`, `CreateTask`, `UpdateTask`, `DeleteTask`

```bash
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext -d '{"user_id":"<uuid>"}' \
  localhost:50051 task.v1.TaskService/GetTasks
```

## Observability

### Metrics

- `http_requests_total{method,path,status}`
- `http_request_duration_seconds{method,path}`
- `cache_hits_total` / `cache_misses_total`
- `tasks_created_total`
- `grpc_requests_total{method,code}`

### Grafana

1. Open http://localhost:3000 (admin / admin)
2. Prometheus datasource is auto-provisioned
3. Explore or build dashboards from the metrics above

## Tests

```bash
make test-unit
go test ./internal/... -count=1

make test-integration   # needs Docker
go test ./internal/repository/ -tags=integration -count=1 -timeout 5m -v
```

## Project layout

```
.
├── api/proto + api/gen     # gRPC
├── deploy/prometheus       # scrape config
├── deploy/grafana          # provisioning
├── internal/
│   ├── adapter/grpc/       # gRPC adapter
│   ├── handler/            # REST adapters
│   ├── service/            # core
│   ├── repository/         # Postgres/Redis + integration tests
│   ├── ports/              # interfaces
│   └── metrics/            # Prometheus
├── docker-compose.yml
├── Makefile
└── .env.example
```

## Configuration

See `.env.example`. Key vars: `SERVER_PORT`, `GRPC_PORT`, `SESSION_EXPIRATION`, `CACHE_EXPIRATION`.

## License

MIT

## Kubernetes (kind)

See [deploy/k8s/README.md](deploy/k8s/README.md).

```bash
docker compose down
sudo fuser -k 18080/tcp 15051/tcp 2>/dev/null || true
kind delete cluster --name task-management 2>/dev/null || true
kind create cluster --config deploy/k8s/kind-config.yaml
docker build -t task-management-api:local .
kind load docker-image task-management-api:local --name task-management
kubectl apply -f deploy/k8s/
curl http://localhost:18080/check_health
```

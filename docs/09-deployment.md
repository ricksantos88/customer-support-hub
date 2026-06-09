# Deployment Guide

## Infrastructure Components

- api
- postgres
- redis

## Deployment Strategy

Use Docker Compose initially.

### Local development

```bash
make migrate-up
make dev
```

This path expects PostgreSQL to be reachable on `localhost:5432`.

### Docker development

```bash
make docker-up
make migrate-up
make dev-docker
```

This path uses the services defined in `docker-compose.yml`.

## Environment Variables

```env
APP_PORT=8080
DB_URL=
REDIS_URL=
JWT_SECRET=
LOG_LEVEL=info
AUTH_ACCESS_TOKEN_TTL_MINUTES=15
AUTH_REFRESH_TOKEN_TTL_HOURS=24
AUTH_SESSION_TTL_HOURS=24
AUTH_CACHE_TTL_MINUTES=30
AUTH_RATE_LIMIT_PER_MINUTE=60
```

## Future

- Kubernetes
- ECS/EKS

## Monitoring

- logs
- metrics
- tracing

Tools:

- Prometheus
- Grafana

## Healthchecks

```http
GET /health
```

## Security

- HTTPS
- secrets management
- token rotation
- rate limiting
- session fallback through PostgreSQL when Redis is unavailable

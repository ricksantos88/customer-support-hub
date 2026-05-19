# Deployment Guide

## Infrastructure Components

* api
* postgres
* redis

## Deployment Strategy

Use Docker Compose initially.

Future:

* Kubernetes
* ECS/EKS

## Environment Variables

```env
APP_PORT=8080
DB_URL=
REDIS_URL=
JWT_SECRET=
LOG_LEVEL=info
```

## Monitoring

* logs
* metrics
* tracing

Tools:

* Prometheus
* Grafana

## Healthchecks

```http
GET /health
```

## Security

* HTTPS
* secrets management
* token rotation
* rate l

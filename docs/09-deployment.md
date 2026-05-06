# Deployment

## Infrastructure
- Docker
- Docker Compose
- PostgreSQL
- Redis
- Nginx

## Services
- api
- db
- redis

## Healthcheck

```http
GET /health
```

Response:

```json
{
  "status": "ok"
}
```

## Environment Variables

```env
APP_PORT=8080
DB_URL=
REDIS_URL=
JWT_SECRET=
```

## Observability
- logs
- metrics
- health checks
- alerting

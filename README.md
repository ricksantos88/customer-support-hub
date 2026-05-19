# Customer Support Hub

Backend API em Go para suporte ao atendimento via WhatsApp.

## Pré-requisitos

- Go 1.25+
- Docker + Docker Compose

## Setup rápido

1. Copie as variáveis de ambiente:
   ```bash
   cp .env.example .env
   ```
2. Defina ao menos `JWT_SECRET` no `.env`.
3. Instale dependências:
   ```bash
   make setup
   ```

## Executar localmente

```bash
make dev
```

API disponível em `http://localhost:8080` e health check em `GET /health`.

## Docker Compose

```bash
make docker-up
```

Serviços da fase 1:
- `api`
- `postgres`
- `redis`
Para parar:

```bash
make docker-down
```

## Comandos úteis

```bash
make lint
make test
make build
```

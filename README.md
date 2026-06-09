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
make migrate-up
make dev
```

API disponível em `http://localhost:8080` e health check em `GET /health`.

Esse fluxo usa `DB_HOST=localhost` e exige que o PostgreSQL esteja acessível na porta `5432`.

## Executar com Docker

```bash
make docker-up
make migrate-up
make dev-docker
```

Esse fluxo usa os serviços `api`, `postgres` e `redis` definidos no `docker-compose.yml`.

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

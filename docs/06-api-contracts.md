# API Contracts

## Authentication

JWT Bearer token.

Header:

```http
Authorization: Bearer <token>
```

### POST /auth/login

Request:

```json
{
  "email": "agent@example.com",
  "password": "secret"
}
```

Response:

```json
{
  "token_type": "access",
  "access_token": "eyJhbGciOi...",
  "refresh_token": "rt_...",
  "access_token_expires_in": 900,
  "refresh_token_expires_in": 86400,
  "session_id": "2c1a5d9a-...",
  "agent_id": "7a7c6d2c-..."
}
```

### POST /auth/refresh

Request:

```json
{
  "refresh_token": "rt_..."
}
```

Response: same shape as `/auth/login`, with rotated refresh token.

### POST /auth/logout

Requires `Authorization: Bearer <token>`.

Response:

```json
{
  "status": "logged_out"
}
```

## POST /messages/send

Request:

```json
{
  "conversation_id": "123",
  "text": "Olá"
}
```

Response:

```json
{
  "status": "sent",
  "message_id": "abc"
}
```

## GET /conversations

Response:

```json
[
  {
    "id": "1",
    "status": "open"
  }
]
```

## GET /conversations/:id/messages

Returns conversation history.

## POST /webhooks/whatsapp

Receives inbound webhook.

## GET /health

Response:

```json
{
  "status": "ok"
}
```

## Auth Rules

- Access tokens are short-lived JWTs.
- Refresh tokens are opaque and rotated on every refresh.
- Sessions are validated against PostgreSQL.
- Redis is used as a cache layer and may be unavailable without blocking auth flows.

---
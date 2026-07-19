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
- Endpoints under `/admin` require a JWT containing the claim `"agent_role": "admin"`. Agents with `agent` role will receive HTTP 403 Forbidden.

---

## Admin Endpoints

All admin endpoints are prefixed with `/admin` and require standard JWT authentication in the `Authorization: Bearer <token>` header.

### POST /admin/agents
Creates a new agent.
* **Request**:
```json
{
  "name": "Carlos",
  "email": "carlos@test.com",
  "password": "password123",
  "role": "agent"
}
```
* **Response (201 Created)**: Agent details.

### GET /admin/agents
Lists registered agents. Supports `limit` and `offset` query parameters.

### PUT /admin/agents/:id
Updates agent details. Cannot demote self.

### DELETE /admin/agents/:id
Inactivates (soft deletes) an agent, revoking all of their active sessions. Cannot delete self.

### GET /admin/sessions
Lists all active agent sessions in the system.

### DELETE /admin/sessions/:id
Revokes an active session immediately, logging out the target agent.

### GET /admin/status
Returns health check status for database, redis, and WhatsApp Meta Cloud configuration.

---
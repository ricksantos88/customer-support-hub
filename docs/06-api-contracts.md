# API Contracts

## Authentication

JWT Bearer token.

Header:

```http
Authorization: Bearer <token>
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

---
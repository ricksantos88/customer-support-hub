# API Contracts

## Send Message

```http
POST /messages/send
```

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
  "status": "sent"
}
```

## List Conversations

```http
GET /conversations
```

## Get Messages

```http
GET /conversations/:id/messages
```

## Webhook

```http
POST /webhooks/whatsapp
```
# Architecture

## Architectural Style

* Hexagonal Architecture
* DDD lite
* provider abstraction

## Folder Structure

```text
cmd/api
internal/
  domain/
    contact/
    conversation/
    message/
    agent/

  application/
    send_message/
    receive_message/
    assign_conversation/

  infrastructure/
    db/
    cache/
    websocket/
    whatsapp/
      meta/

  interfaces/
    http/
    ws/
```

## Layer Responsibilities

## Domain

Contains:

* entities
* value objects
* domain rules

## Application

Contains:

* use cases
* orchestration logic

## Infrastructure

Contains:

* postgres
* redis
* whatsapp provider
* websocket manager

## Interfaces

Contains:

* REST handlers
* middleware
* DTOs

## Core Interface

```go
type MessagingProvider interface {
    SendText(ctx context.Context, to string, text string) error
    HandleWebhook(payload []byte) error
}
```


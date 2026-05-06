# Architecture

## Project Structure

```text
cmd/api
internal/
  domain/
  application/
  infrastructure/
  interfaces/
```

## Layers

### Domain
Business entities.

### Application
Use cases.

### Infrastructure
External integrations.

### Interfaces
HTTP and WebSocket adapters.

## Patterns
- Hexagonal Architecture
- Repository Pattern
- Dependency Injection

## Provider Interface

```go
type MessagingProvider interface {
    SendText(ctx context.Context, to string, text string) error
    HandleWebhook(payload []byte) error
}
```
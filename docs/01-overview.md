
# WhatsApp Customer Support API

## Objective
Create a backend API in Go to integrate with WhatsApp for customer service operations.

The system must:
- Receive customer messages from WhatsApp
- Persist conversations and messages
- Expose data to frontend agents
- Allow agents to send messages back to customers
- Support real-time updates via WebSocket

## MVP Scope

### Features
- Receive inbound WhatsApp messages
- Send outbound WhatsApp messages
- Store conversations
- Store contacts
- WebSocket real-time updates
- Agent authentication

## Stack
- Go 1.24+
- Fiber
- PostgreSQL
- Redis
- Docker
- WebSocket
- WhatsApp Cloud API

## Architecture
- Hexagonal Architecture
- Domain Driven Design (lightweight)
- Provider abstraction

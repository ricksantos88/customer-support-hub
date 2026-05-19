# WhatsApp Customer Support API

## Project Objective

Build a robust backend API in Go responsible for integrating with WhatsApp and serving a customer support frontend application.

The platform must allow customer service teams to communicate with end users through WhatsApp while maintaining persistence, observability, and scalability.

## Business Goals

* Centralize customer conversations
* Enable multiple agents to manage conversations
* Persist all message history
* Support real-time customer service operations
* Allow future omnichannel expansion

## MVP Scope

### In Scope

* Receive inbound messages from WhatsApp
* Send outbound messages
* Store contacts
* Store conversations
* Store messages
* WebSocket real-time updates
* Agent authentication
* Healthcheck endpoints

### Out of Scope (future)

* Chatbot automation
* AI integrations
* Campaign management
* Message templates UI
* Analytics dashboards
* Multi-tenant SaaS

## Non-Functional Requirements

* Low latency (<300ms internal operations)
* Horizontal scalability
* High availability
* Secure secret management
* Structured logs

## Tech Stack

### Backend

* Go 1.24+
* Fiber

### Storage

* PostgreSQL
* Redis

### Infrastructure

* Docker
* Docker Compose



### Integration

* WhatsApp Cloud API

### Communication

* WebSocket

## High-Level Flow

1. Customer sends message
2. WhatsApp forwards webhook
3. API processes message
4. Persists data
5. Pushes WebSocket event to frontend
6. Agent replies
7. API sends outbound message


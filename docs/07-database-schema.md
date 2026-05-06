# Database Schema

## contacts

Stores customer records.

## conversations

Stores support sessions.

Statuses:

* open
* pending
* closed

## messages

Stores inbound/outbound messages.

Directions:

* inbound
* outbound

## agents

Stores support agents.

Suggested schema includes:

* UUID primary keys
* timestamps
* indexes
* foreign keys

Recommended indexes:

* phone
* conversation_id
* assigned_agent_id
* created_at


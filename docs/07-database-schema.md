# Database Schema

This schema is defined in `migrations/001_initial_schema.up.sql`.

## contacts

Stores customer records.

Columns:
- `id` UUID PK
- `phone` (E.164, unique when not soft-deleted)
- `name`
- `created_at`
- `updated_at`
- `deleted_at` (soft delete)

Indexes:
- `idx_contacts_phone_active_unique` on `phone` where `deleted_at IS NULL`

## agents

Stores support agents.

Columns:
- `id` UUID PK
- `name`
- `email` (unique)
- `jwt_hash` (bcrypt hash of JWT secret)
- `created_at`
- `last_active`

## conversations

Stores support sessions.

Columns:
- `id` UUID PK
- `contact_id` (FK -> `contacts.id`)
- `assigned_agent_id` (FK -> `agents.id`, nullable)
- `status` (`open` | `pending` | `closed`)
- `created_at`
- `updated_at`

Constraints:
- `status` must be `open`, `pending`, or `closed`
- FKs use `ON DELETE RESTRICT`

Indexes:
- `idx_conversations_assigned_agent_status` on (`assigned_agent_id`, `status`)

## messages

Stores inbound/outbound messages.

Columns:
- `id` UUID PK
- `conversation_id` (FK -> `conversations.id`)
- `content` (TEXT)
- `direction` (`inbound` | `outbound`)
- `sender_id` (polymorphic: contact or agent)
- `created_at`

Constraints:
- `direction` must be `inbound` or `outbound`
- FKs use `ON DELETE RESTRICT`

Indexes:
- `idx_messages_conversation_created_at` on (`conversation_id`, `created_at`)

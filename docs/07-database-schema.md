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
- `password_hash` (bcrypt hash of agent credential)
- `created_at`
- `last_active`

## auth_sessions

Stores login sessions and refresh-token rotation state.

Columns:
- `id` UUID PK
- `agent_id` (FK -> `agents.id`)
- `refresh_token_hash` (unique)
- `user_agent` nullable
- `ip_address` nullable
- `created_at`
- `last_used_at`
- `expires_at`
- `revoked_at` nullable

Constraints:
- `refresh_token_hash` must be unique
- FKs use `ON DELETE RESTRICT`

Indexes:
- `idx_auth_sessions_agent_active` on (`agent_id`, `revoked_at`, `expires_at`)
- `idx_auth_sessions_expires_at` on (`expires_at`)

Operational Notes:
- PostgreSQL is the source of truth for session validity.
- Redis may cache sessions and refresh-token lookups, but auth must continue when Redis is unavailable.

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

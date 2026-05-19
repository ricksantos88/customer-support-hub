CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE TABLE agents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    jwt_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_active TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contact_id UUID NOT NULL,
    assigned_agent_id UUID NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_conversations_contact
        FOREIGN KEY (contact_id) REFERENCES contacts(id) ON DELETE RESTRICT,
    CONSTRAINT fk_conversations_assigned_agent
        FOREIGN KEY (assigned_agent_id) REFERENCES agents(id) ON DELETE RESTRICT,
    CONSTRAINT chk_conversations_status
        CHECK (status IN ('open', 'pending', 'closed'))
);

CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL,
    content TEXT NOT NULL,
    direction VARCHAR(20) NOT NULL,
    sender_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_messages_conversation
        FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE RESTRICT,
    CONSTRAINT chk_messages_direction
        CHECK (direction IN ('inbound', 'outbound'))
);

CREATE UNIQUE INDEX idx_contacts_phone_active_unique
    ON contacts (phone)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_conversations_assigned_agent_status
    ON conversations (assigned_agent_id, status);

CREATE INDEX idx_messages_conversation_created_at
    ON messages (conversation_id, created_at);

COMMENT ON TABLE contacts IS 'Contacts that talk to support';
COMMENT ON COLUMN contacts.id IS 'Contact identifier';
COMMENT ON COLUMN contacts.phone IS 'Phone number in E.164 format';
COMMENT ON COLUMN contacts.name IS 'Display name for the contact';
COMMENT ON COLUMN contacts.created_at IS 'Creation timestamp';
COMMENT ON COLUMN contacts.updated_at IS 'Last update timestamp';
COMMENT ON COLUMN contacts.deleted_at IS 'Soft delete timestamp';

COMMENT ON TABLE agents IS 'Support agents authenticated via JWT';
COMMENT ON COLUMN agents.id IS 'Agent identifier';
COMMENT ON COLUMN agents.name IS 'Agent name';
COMMENT ON COLUMN agents.email IS 'Agent email (unique)';
COMMENT ON COLUMN agents.jwt_hash IS 'Bcrypt hash of JWT secret';
COMMENT ON COLUMN agents.created_at IS 'Creation timestamp';
COMMENT ON COLUMN agents.last_active IS 'Last activity timestamp';

COMMENT ON TABLE conversations IS 'Conversations between contacts and agents';
COMMENT ON COLUMN conversations.id IS 'Conversation identifier';
COMMENT ON COLUMN conversations.contact_id IS 'Contact that owns the conversation';
COMMENT ON COLUMN conversations.assigned_agent_id IS 'Assigned agent for the conversation';
COMMENT ON COLUMN conversations.status IS 'Conversation status: open, pending, closed';
COMMENT ON COLUMN conversations.created_at IS 'Creation timestamp';
COMMENT ON COLUMN conversations.updated_at IS 'Last update timestamp';

COMMENT ON TABLE messages IS 'Messages inside conversations';
COMMENT ON COLUMN messages.id IS 'Message identifier';
COMMENT ON COLUMN messages.conversation_id IS 'Conversation that owns the message';
COMMENT ON COLUMN messages.content IS 'Text content';
COMMENT ON COLUMN messages.direction IS 'Direction: inbound or outbound';
COMMENT ON COLUMN messages.sender_id IS 'Sender identifier (contact or agent)';
COMMENT ON COLUMN messages.created_at IS 'Creation timestamp';


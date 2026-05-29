ALTER TABLE messages
    ADD COLUMN sender_id UUID NULL;

UPDATE messages
SET sender_id = sender_contact_id
WHERE direction = 'inbound';

UPDATE messages
SET sender_id = sender_agent_id
WHERE direction = 'outbound';

ALTER TABLE messages
    ALTER COLUMN sender_id SET NOT NULL;

ALTER TABLE messages
    DROP CONSTRAINT IF EXISTS chk_messages_sender_by_direction,
    DROP CONSTRAINT IF EXISTS fk_messages_sender_contact,
    DROP CONSTRAINT IF EXISTS fk_messages_sender_agent;

ALTER TABLE messages
    DROP COLUMN sender_contact_id,
    DROP COLUMN sender_agent_id;

COMMENT ON COLUMN messages.sender_id IS 'Sender identifier (contact or agent)';

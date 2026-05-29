ALTER TABLE messages
    ADD COLUMN sender_contact_id UUID NULL,
    ADD COLUMN sender_agent_id UUID NULL;

UPDATE messages
SET sender_contact_id = sender_id
WHERE direction = 'inbound';

UPDATE messages
SET sender_agent_id = sender_id
WHERE direction = 'outbound';

ALTER TABLE messages
    ADD CONSTRAINT fk_messages_sender_contact
        FOREIGN KEY (sender_contact_id) REFERENCES contacts(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_messages_sender_agent
        FOREIGN KEY (sender_agent_id) REFERENCES agents(id) ON DELETE RESTRICT,
    ADD CONSTRAINT chk_messages_sender_by_direction
        CHECK (
            (direction = 'inbound' AND sender_contact_id IS NOT NULL AND sender_agent_id IS NULL)
            OR
            (direction = 'outbound' AND sender_agent_id IS NOT NULL AND sender_contact_id IS NULL)
        );

ALTER TABLE messages
    DROP COLUMN sender_id;

COMMENT ON COLUMN messages.sender_contact_id IS 'Sender contact identifier for inbound messages';
COMMENT ON COLUMN messages.sender_agent_id IS 'Sender agent identifier for outbound messages';

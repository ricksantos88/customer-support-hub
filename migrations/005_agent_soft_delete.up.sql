ALTER TABLE agents
    ADD COLUMN deleted_at TIMESTAMPTZ NULL;

CREATE INDEX idx_agents_deleted_at ON agents(deleted_at);

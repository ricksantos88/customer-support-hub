ALTER TABLE agents
    ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'agent'
        CONSTRAINT chk_agents_role CHECK (role IN ('admin', 'agent'));

COMMENT ON COLUMN agents.role IS 'Agent role: admin or agent';

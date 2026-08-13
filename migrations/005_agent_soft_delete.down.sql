DROP INDEX IF EXISTS idx_agents_deleted_at;

ALTER TABLE agents
    DROP COLUMN IF EXISTS deleted_at;

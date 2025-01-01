CREATE TABLE agents_meta (
	id UUID PRIMARY KEY,
	agent_id UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
	key VARCHAR(128) NOT NULL,
	value TEXT,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	UNIQUE (agent_id, key)
);
CREATE INDEX idx_agents_meta_agent_id ON agents_meta(agent_id);
CREATE INDEX idx_agents_meta_key ON agents_meta(key);

CREATE TABLE sessions (
	id UUID PRIMARY KEY,
	external_id VARCHAR(60) NOT NULL,
	workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
	agent_id UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
	title VARCHAR(255),
	status VARCHAR(50) NOT NULL DEFAULT 'active',
	labels JSONB,
	meta JSONB,
	last_activity_at TIMESTAMP NULL,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	UNIQUE (agent_id, external_id)
);
CREATE INDEX idx_sessions_workspace_id ON sessions(workspace_id);
CREATE INDEX idx_sessions_agent_id ON sessions(agent_id);
CREATE INDEX idx_sessions_status ON sessions(status);
CREATE INDEX idx_sessions_labels ON sessions USING GIN (labels jsonb_path_ops);

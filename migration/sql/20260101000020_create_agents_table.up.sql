CREATE TABLE agents (
	id UUID PRIMARY KEY,
	workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
	name VARCHAR(100) NOT NULL,
	kind VARCHAR(50) NOT NULL DEFAULT 'unmanaged',
	prompt_id UUID NULL REFERENCES prompts(id) ON DELETE SET NULL,
	kb_labels JSONB NULL,
	handle VARCHAR(100) NOT NULL UNIQUE,
	description TEXT NOT NULL DEFAULT '',
	status VARCHAR(50) NOT NULL DEFAULT 'active',
	config JSONB,
	meta JSONB,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
);
CREATE INDEX idx_agents_workspace_id ON agents(workspace_id);
CREATE INDEX idx_agents_status ON agents(status);

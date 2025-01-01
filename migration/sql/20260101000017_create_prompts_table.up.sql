CREATE TABLE prompts (
	id UUID PRIMARY KEY,
	workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
	name VARCHAR(100) NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	handle VARCHAR(100) NOT NULL,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	UNIQUE (workspace_id, handle)
);
CREATE INDEX idx_prompts_workspace_id ON prompts(workspace_id);

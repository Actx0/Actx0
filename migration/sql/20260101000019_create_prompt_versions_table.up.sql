CREATE TABLE prompt_versions (
	id UUID PRIMARY KEY,
	prompt_id UUID NOT NULL REFERENCES prompts(id) ON DELETE CASCADE,
	version INTEGER NOT NULL,
	type VARCHAR(20) NOT NULL DEFAULT 'text',
	content TEXT NOT NULL,
	config JSONB,
	labels JSONB,
	commit_message VARCHAR(255),
	commit_hash VARCHAR(40) NOT NULL,
	meta JSONB,
	status VARCHAR(50) NOT NULL DEFAULT 'active',
	production BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	UNIQUE (prompt_id, version)
);
CREATE INDEX idx_prompt_versions_prompt_id ON prompt_versions(prompt_id);
CREATE INDEX idx_prompt_versions_status ON prompt_versions(status);
CREATE INDEX idx_prompt_versions_labels ON prompt_versions USING GIN (labels);
CREATE UNIQUE INDEX idx_prompt_versions_one_production ON prompt_versions(prompt_id) WHERE production = TRUE;

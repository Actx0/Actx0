CREATE TABLE prompts_meta (
	id UUID PRIMARY KEY,
	prompt_id UUID NOT NULL REFERENCES prompts(id) ON DELETE CASCADE,
	key VARCHAR(128) NOT NULL,
	value TEXT,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	UNIQUE (prompt_id, key)
);
CREATE INDEX idx_prompts_meta_prompt_id ON prompts_meta(prompt_id);
CREATE INDEX idx_prompts_meta_key ON prompts_meta(key);

CREATE TABLE session_memories (
	id UUID PRIMARY KEY,
	internal_id UUID NOT NULL UNIQUE,
	session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
	kind VARCHAR(100) NOT NULL DEFAULT 'summary',
	content TEXT NOT NULL,
	meta JSONB,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
);
CREATE INDEX idx_session_memories_session_id ON session_memories(session_id);
CREATE INDEX idx_session_memories_kind ON session_memories(kind);

CREATE TABLE sessions_meta (
	id UUID PRIMARY KEY,
	session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
	key VARCHAR(128) NOT NULL,
	value TEXT,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	UNIQUE (session_id, key)
);
CREATE INDEX idx_sessions_meta_session_id ON sessions_meta(session_id);
CREATE INDEX idx_sessions_meta_key ON sessions_meta(key);

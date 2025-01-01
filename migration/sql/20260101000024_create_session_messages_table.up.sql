CREATE TABLE session_messages (
	id UUID PRIMARY KEY,
	internal_id UUID NOT NULL UNIQUE,
	session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
	role VARCHAR(50) NOT NULL,
	content TEXT NOT NULL,
	meta JSONB,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
);
CREATE INDEX idx_session_messages_session_id ON session_messages(session_id);
CREATE INDEX idx_session_messages_role ON session_messages(role);
CREATE INDEX idx_session_messages_created_at ON session_messages(created_at);

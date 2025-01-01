CREATE TABLE notifications (
	id UUID PRIMARY KEY,
	user_id UUID NULL REFERENCES users(id) ON DELETE CASCADE,
	workspace_id UUID NULL REFERENCES workspaces(id) ON DELETE CASCADE,
	type VARCHAR(60) NOT NULL,
	title VARCHAR(200) NOT NULL,
	body TEXT,
	read_at TIMESTAMP NULL,
	archived_at TIMESTAMP NULL,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	CHECK (num_nonnulls(user_id, workspace_id) = 1)
);
CREATE INDEX idx_notifications_user_id_created_at ON notifications (user_id, created_at DESC);
CREATE INDEX idx_notifications_workspace_id_created_at ON notifications (workspace_id, created_at DESC);

CREATE TABLE notifications_meta (
	id UUID PRIMARY KEY,
	notification_id UUID NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
	key VARCHAR(60) NOT NULL,
	value JSONB,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	UNIQUE (notification_id, key)
);
CREATE INDEX idx_notifications_meta_notification_id_key ON notifications_meta(notification_id, key);

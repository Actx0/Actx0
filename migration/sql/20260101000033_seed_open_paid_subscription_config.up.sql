INSERT INTO configs (id, key, value)
VALUES (
	'8e0c0f0a-4b3a-4d2e-9c1b-0a1b2c3d4e5f',
	'_open_paid_subscription',
	to_jsonb('0'::text)
)
ON CONFLICT (key) DO NOTHING;

// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"

	"github.com/samber/lo"
)

// Notification is a single row in the notifications table.
type Notification struct {
	Id          Id
	UserId      *Id
	WorkspaceId *Id
	Type        string
	Title       string
	Body        *string
	ReadAt      *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NotificationMeta is a single row in the notifications_meta table.
type NotificationMeta struct {
	Id             Id
	NotificationId Id
	Key            string
	Value          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NotificationRepository is the interface for notification persistence.
type NotificationRepository interface {
	Create(notification *Notification) error
	GetById(id Id) (*Notification, error)
	ListByUserId(userId Id, limit, offset int) ([]*Notification, error)
	ListByWorkspaceId(workspaceId Id, limit, offset int) ([]*Notification, error)
	ListByUserIdAndWorkspaceId(userId, workspaceId Id, limit, offset int) ([]*Notification, error)
	CountByUserId(userId Id) (int64, error)
	CountByWorkspaceId(workspaceId Id) (int64, error)
	CountByUserIdAndWorkspaceId(userId, workspaceId Id) (int64, error)
	MarkRead(id Id) error
	Delete(id Id) error
}

// NotificationMetaRepository is the interface for notification metadata CRUD.
type NotificationMetaRepository interface {
	Create(id Id, key, value string) error
	Get(id Id, key string) (*NotificationMeta, error)
	Update(id Id, key, value string) error
	Delete(id Id, key string) error
	ListByNotification(id Id) ([]*NotificationMeta, error)
}

// NotificationRepositoryPostgres implements NotificationRepository for PostgreSQL
type NotificationRepositoryPostgres struct {
	db *sql.DB
}

// NewNotificationRepository returns a notification repository.
func NewNotificationRepository(db *sql.DB) NotificationRepository {
	return &NotificationRepositoryPostgres{db: db}
}

// NotificationMetaRepositoryPostgres implements NotificationMetaRepository for PostgreSQL
type NotificationMetaRepositoryPostgres struct {
	db *sql.DB
}

// NewNotificationMetaRepository returns a notification meta repository.
func NewNotificationMetaRepository(db *sql.DB) NotificationMetaRepository {
	return &NotificationMetaRepositoryPostgres{db: db}
}

// Create inserts a notification row.
func (r *NotificationRepositoryPostgres) Create(notification *Notification) error {
	id, err := NewId()
	if err != nil {
		return err
	}
	notification.Id = id

	err = r.db.QueryRow(
		`INSERT INTO notifications (
			id, user_id, workspace_id, type, title, body
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`,
		notification.Id.String(),
		lo.FromPtr(notification.UserId),
		lo.FromPtr(notification.WorkspaceId),
		notification.Type,
		notification.Title,
		notification.Body,
	).Scan(&notification.CreatedAt, &notification.UpdatedAt)
	return err
}

// GetById returns a notification by id.
func (r *NotificationRepositoryPostgres) GetById(id Id) (*Notification, error) {
	notification := &Notification{}
	err := r.db.QueryRow(
		`SELECT
			id, user_id, workspace_id, type, title, body, read_at, created_at, updated_at
		FROM notifications
		WHERE id = $1`,
		id.String(),
	).Scan(
		&notification.Id,
		&notification.UserId,
		&notification.WorkspaceId,
		&notification.Type,
		&notification.Title,
		&notification.Body,
		&notification.ReadAt,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	)

	if isNotFound(err) {
		return nil, nil
	}
	return notification, err
}

// ListByUserId lists notification rows targeted to a user.
func (r *NotificationRepositoryPostgres) ListByUserId(userId Id, limit, offset int) ([]*Notification, error) {
	rows, err := r.db.Query(
		`SELECT
			id, user_id, workspace_id, type, title, body, read_at, created_at, updated_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		userId.String(),
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Notification
	for rows.Next() {
		notification := &Notification{}
		if err := rows.Scan(
			&notification.Id,
			&notification.UserId,
			&notification.WorkspaceId,
			&notification.Type,
			&notification.Title,
			&notification.Body,
			&notification.ReadAt,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, notification)
	}

	return list, rows.Err()
}

// ListByWorkspaceId lists notification rows targeted to a workspace.
func (r *NotificationRepositoryPostgres) ListByWorkspaceId(workspaceId Id, limit, offset int) ([]*Notification, error) {
	rows, err := r.db.Query(
		`SELECT
			id, user_id, workspace_id, type, title, body, read_at, created_at, updated_at
		FROM notifications
		WHERE workspace_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		workspaceId.String(),
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Notification
	for rows.Next() {
		notification := &Notification{}
		if err := rows.Scan(
			&notification.Id,
			&notification.UserId,
			&notification.WorkspaceId,
			&notification.Type,
			&notification.Title,
			&notification.Body,
			&notification.ReadAt,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, notification)
	}

	return list, rows.Err()
}

// ListByUserIdAndWorkspaceId lists personal and workspace notifications for an inbox.
func (r *NotificationRepositoryPostgres) ListByUserIdAndWorkspaceId(userId, workspaceId Id, limit, offset int) ([]*Notification, error) {
	rows, err := r.db.Query(
		`SELECT
			id, user_id, workspace_id, type, title, body, read_at, created_at, updated_at
		FROM notifications
		WHERE user_id = $1 OR workspace_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`,
		userId.String(),
		workspaceId.String(),
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Notification
	for rows.Next() {
		notification := &Notification{}
		if err := rows.Scan(
			&notification.Id,
			&notification.UserId,
			&notification.WorkspaceId,
			&notification.Type,
			&notification.Title,
			&notification.Body,
			&notification.ReadAt,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, notification)
	}

	return list, rows.Err()
}

// CountByUserId counts notification rows targeted to a user.
func (r *NotificationRepositoryPostgres) CountByUserId(userId Id) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1`,
		userId.String(),
	).Scan(&count)
	return count, err
}

// CountByWorkspaceId counts notification rows targeted to a workspace.
func (r *NotificationRepositoryPostgres) CountByWorkspaceId(workspaceId Id) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM notifications
		WHERE workspace_id = $1`,
		workspaceId.String(),
	).Scan(&count)
	return count, err
}

// CountByUserIdAndWorkspaceId counts personal and workspace notifications for an inbox.
func (r *NotificationRepositoryPostgres) CountByUserIdAndWorkspaceId(userId, workspaceId Id) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1 OR workspace_id = $2`,
		userId.String(),
		workspaceId.String(),
	).Scan(&count)
	return count, err
}

// MarkRead sets the read timestamp on a notification.
func (r *NotificationRepositoryPostgres) MarkRead(id Id) error {
	now := time.Now().UTC()
	_, err := r.db.Exec(
		`UPDATE notifications
		SET
			read_at = $1,
			updated_at = $2
		WHERE id = $3`,
		now,
		now,
		id.String(),
	)
	return err
}

// Delete removes a notification row.
func (r *NotificationRepositoryPostgres) Delete(id Id) error {
	_, err := r.db.Exec(
		`DELETE FROM notifications WHERE id = $1`,
		id.String(),
	)
	return err
}

// Create inserts a notification meta
func (r *NotificationMetaRepositoryPostgres) Create(id Id, key, value string) error {
	metaId, err := NewId()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		`INSERT INTO notifications_meta (id, notification_id, key, value)
		VALUES ($1, $2, $3, $4)`,
		metaId.String(),
		id.String(),
		key,
		value,
	)
	return err
}

// Get returns a notification meta
func (r *NotificationMetaRepositoryPostgres) Get(id Id, key string) (*NotificationMeta, error) {
	meta := &NotificationMeta{}
	err := r.db.QueryRow(
		`SELECT id, key, value, notification_id, created_at, updated_at
		FROM notifications_meta
		WHERE notification_id = $1 AND key = $2`,
		id.String(),
		key,
	).Scan(
		&meta.Id,
		&meta.Key,
		&meta.Value,
		&meta.NotificationId,
		&meta.CreatedAt,
		&meta.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return meta, err
}

// Update updates a notification meta
func (r *NotificationMetaRepositoryPostgres) Update(id Id, key, value string) error {
	_, err := r.db.Exec(
		`UPDATE notifications_meta
		SET
			value = $1,
			updated_at = $2
		WHERE notification_id = $3 AND key = $4`,
		value,
		time.Now().UTC(),
		id.String(),
		key,
	)
	return err
}

// Delete removes a notification meta
func (r *NotificationMetaRepositoryPostgres) Delete(id Id, key string) error {
	_, err := r.db.Exec(
		`DELETE FROM notifications_meta WHERE notification_id = $1 AND key = $2`,
		id.String(),
		key,
	)
	return err
}

// ListByNotification returns a list of notification meta
func (r *NotificationMetaRepositoryPostgres) ListByNotification(id Id) ([]*NotificationMeta, error) {
	rows, err := r.db.Query(
		`SELECT id, key, value, notification_id, created_at, updated_at
		FROM notifications_meta
		WHERE notification_id = $1
		ORDER BY key`,
		id.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*NotificationMeta
	for rows.Next() {
		meta := &NotificationMeta{}
		if err := rows.Scan(
			&meta.Id,
			&meta.Key,
			&meta.Value,
			&meta.NotificationId,
			&meta.CreatedAt,
			&meta.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, meta)
	}

	return list, rows.Err()
}

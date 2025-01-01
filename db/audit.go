// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// AuditEvent is a single row in the audit table.
type AuditEvent struct {
	Id           Id
	WorkspaceId  Id
	UserId       *Id
	Action       string
	ResourceType *string
	ResourceId   *Id
	IPAddress    *string
	UserAgent    *string
	Meta         *string
	CreatedAt    time.Time
}

// AuditEventRepository is the interface for audit event persistence.
type AuditEventRepository interface {
	Create(event *AuditEvent) error
	GetById(id Id) (*AuditEvent, error)
	ListByWorkspaceId(workspaceId Id, limit, offset int) ([]*AuditEvent, error)
	CountByWorkspaceId(workspaceId Id) (int64, error)
}

type AuditEventRepositoryPostgres struct {
	db *sql.DB
}

// NewAuditEventRepository returns an audit event repository.
func NewAuditEventRepository(db *sql.DB) AuditEventRepository {
	return &AuditEventRepositoryPostgres{db: db}
}

// Create inserts an audit event row.
func (r *AuditEventRepositoryPostgres) Create(event *AuditEvent) error {
	id, err := NewId()
	if err != nil {
		return err
	}
	event.Id = id

	err = r.db.QueryRow(
		`INSERT INTO audit (
			id, workspace_id, user_id, action, resource_type, resource_id,
			ip_address, user_agent, meta
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at`,
		event.Id.String(),
		event.WorkspaceId.String(),
		event.UserId.ValueOrNilPtr(),
		event.Action,
		event.ResourceType,
		event.ResourceId.ValueOrNilPtr(),
		event.IPAddress,
		event.UserAgent,
		event.Meta,
	).Scan(&event.CreatedAt)
	return err
}

// GetById returns an audit event by id.
func (r *AuditEventRepositoryPostgres) GetById(id Id) (*AuditEvent, error) {
	item := &AuditEvent{}
	err := r.db.QueryRow(
		`SELECT
			id, workspace_id, user_id, action, resource_type, resource_id,
			ip_address, user_agent, meta, created_at
		FROM audit
		WHERE id = $1`,
		id.String(),
	).Scan(
		&item.Id,
		&item.WorkspaceId,
		&item.UserId,
		&item.Action,
		&item.ResourceType,
		&item.ResourceId,
		&item.IPAddress,
		&item.UserAgent,
		&item.Meta,
		&item.CreatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// ListByWorkspaceId lists audit event rows by workspace id.
func (r *AuditEventRepositoryPostgres) ListByWorkspaceId(workspaceId Id, limit, offset int) ([]*AuditEvent, error) {
	rows, err := r.db.Query(
		`SELECT
			id, workspace_id, user_id, action, resource_type, resource_id,
			ip_address, user_agent, meta, created_at
		FROM audit
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

	var list []*AuditEvent
	for rows.Next() {
		item := &AuditEvent{}
		if err := rows.Scan(
			&item.Id,
			&item.WorkspaceId,
			&item.UserId,
			&item.Action,
			&item.ResourceType,
			&item.ResourceId,
			&item.IPAddress,
			&item.UserAgent,
			&item.Meta,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

// CountByWorkspaceId counts audit event rows by workspace id.
func (r *AuditEventRepositoryPostgres) CountByWorkspaceId(workspaceId Id) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM audit
		WHERE workspace_id = $1`,
		workspaceId.String(),
	).Scan(&count)
	return count, err
}

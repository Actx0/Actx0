// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// WorkspaceEventsDeliverable is one row in workspace_events_deliverables (event to deliver to a workspace integration).
type WorkspaceEventDeliverable struct {
	ID            ID
	EventID       ID
	IntegrationID ID
	Status        string
	DeliveredAt   *time.Time
	ErrorMessage  *string
	RetryCount    int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// WorkspaceEventsDeliverableRepository is the interface for workspace_events_deliverables persistence.
type WorkspaceEventDeliverableRepository interface {
	Create(eventID, integrationID ID) error
	GetByID(id ID) (*WorkspaceEventDeliverable, error)
	ListPending(limit int) ([]*WorkspaceEventDeliverable, error)
	ListByEventID(eventID ID, limit, offset int) ([]*WorkspaceEventDeliverable, error)
	ListByIntegrationID(integrationID ID, limit, offset int) ([]*WorkspaceEventDeliverable, error)
	UpdateStatus(id ID, status string, deliveredAt *time.Time, errorMessage *string, retryCount int) error
}

// NewWorkspaceEventsDeliverableRepository returns the repository for the current driver.
func NewWorkspaceEventDeliverableRepository(db *sql.DB) WorkspaceEventDeliverableRepository {
	return &workspaceEventDeliverableRepoPostgres{db: db}
}

func scanWorkspaceEventDeliverableRows(rows *sql.Rows) ([]*WorkspaceEventDeliverable, error) {
	var list []*WorkspaceEventDeliverable
	for rows.Next() {
		d := &WorkspaceEventDeliverable{}
		var deliveredAt sql.NullTime
		var errorMessage sql.NullString
		if err := rows.Scan(
			&d.ID,
			&d.EventID,
			&d.IntegrationID,
			&d.Status,
			&deliveredAt,
			&errorMessage,
			&d.RetryCount,
			&d.CreatedAt,
			&d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if deliveredAt.Valid {
			d.DeliveredAt = &deliveredAt.Time
		}
		if errorMessage.Valid {
			d.ErrorMessage = &errorMessage.String
		}
		list = append(list, d)
	}
	return list, rows.Err()
}

// --- Postgres ---

type workspaceEventDeliverableRepoPostgres struct {
	db *sql.DB
}

func (r *workspaceEventDeliverableRepoPostgres) Create(eventID, integrationID ID) error {
	id, err := NewID()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		`INSERT INTO workspace_events_deliverables (id, event_id, integration_id, status)
		VALUES ($1, $2, $3, 'pending')`,
		id.String(),
		eventID.String(),
		integrationID.String(),
	)
	return err
}

func (r *workspaceEventDeliverableRepoPostgres) GetByID(id ID) (*WorkspaceEventDeliverable, error) {
	d := &WorkspaceEventDeliverable{}
	var deliveredAt sql.NullTime
	var errorMessage sql.NullString
	err := r.db.QueryRow(
		`SELECT
			id, event_id, integration_id, status, delivered_at,
			error_message, retry_count, created_at, updated_at
		FROM workspace_events_deliverables
		WHERE id = $1`,
		id.String(),
	).Scan(
		&d.ID,
		&d.EventID,
		&d.IntegrationID,
		&d.Status,
		&deliveredAt,
		&errorMessage,
		&d.RetryCount,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if deliveredAt.Valid {
		d.DeliveredAt = &deliveredAt.Time
	}
	if errorMessage.Valid {
		d.ErrorMessage = &errorMessage.String
	}
	return d, nil
}

func (r *workspaceEventDeliverableRepoPostgres) ListPending(limit int) ([]*WorkspaceEventDeliverable, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := r.db.Query(
		`SELECT
			id, event_id, integration_id, status, delivered_at,
			error_message, retry_count, created_at, updated_at
		FROM workspace_events_deliverables
		WHERE status = 'pending'
		ORDER BY created_at ASC
		LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanWorkspaceEventDeliverableRows(rows)
}

func (r *workspaceEventDeliverableRepoPostgres) ListByEventID(eventID ID, limit, offset int) ([]*WorkspaceEventDeliverable, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.Query(
		`SELECT
			id, event_id, integration_id, status, delivered_at,
			error_message, retry_count, created_at, updated_at
		FROM workspace_events_deliverables
		WHERE event_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		eventID.String(),
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanWorkspaceEventDeliverableRows(rows)
}

func (r *workspaceEventDeliverableRepoPostgres) ListByIntegrationID(integrationID ID, limit, offset int) ([]*WorkspaceEventDeliverable, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.Query(
		`SELECT
			id, event_id, integration_id, status, delivered_at,
			error_message, retry_count, created_at, updated_at
		FROM workspace_events_deliverables
		WHERE integration_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		integrationID.String(),
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanWorkspaceEventDeliverableRows(rows)
}

func (r *workspaceEventDeliverableRepoPostgres) UpdateStatus(id ID, status string, deliveredAt *time.Time, errorMessage *string, retryCount int) error {
	now := time.Now().UTC()
	_, err := r.db.Exec(
		`UPDATE workspace_events_deliverables
		SET
			status = $1,
			delivered_at = $2,
			error_message = $3,
			retry_count = $4,
			updated_at = $5
		WHERE id = $6`,
		status,
		deliveredAt,
		nullStrPtr(errorMessage),
		retryCount,
		now,
		id.String(),
	)
	return err
}

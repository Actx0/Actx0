// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"encoding/json"
	"time"
)

// WorkspaceEvent is one row in the workspace event log.
type WorkspaceEvent struct {
	ID          ID
	WorkspaceID ID
	Type        string
	Data        json.RawMessage
	CreatedAt   time.Time
}

// Workspace event data structs (payload for each event type).
// Unmarshal event.Data into the appropriate struct based on event.Type.

// WorkspaceEventInventoryUpdated is the event data for "inventory_updated".
type WorkspaceEventInventoryUpdated struct {
	InventoryID string `json:"inventory_id"`
	UpdatedBy   string `json:"updated_by,omitempty"`
}

// WorkspaceEventMemberAdded is the event data for "member_added".
type WorkspaceEventMemberAdded struct {
	UserID  string `json:"user_id"`
	Role    string `json:"role"`
	AddedBy string `json:"added_by,omitempty"`
}

// WorkspaceEventRepository is the interface for workspace event persistence.
type WorkspaceEventRepository interface {
	Create(workspaceID ID, eventType string, data interface{}) error
	GetByID(id ID) (*WorkspaceEvent, error)
	ListByWorkspace(workspaceID ID, limit, offset int) ([]*WorkspaceEvent, error)
	ListByWorkspaceAndType(workspaceID ID, eventType string, limit, offset int) ([]*WorkspaceEvent, error)
	CountByWorkspace(workspaceID ID) (int64, error)
	CountByWorkspaceAndType(workspaceID ID, eventType string) (int64, error)
}

// NewWorkspaceEventRepository returns the repository for the current driver
func NewWorkspaceEventRepository(db *sql.DB) WorkspaceEventRepository {
	return &workspaceEventRepoPostgres{db: db}
}

func scanWorkspaceEventRows(rows *sql.Rows) ([]*WorkspaceEvent, error) {
	var list []*WorkspaceEvent
	for rows.Next() {
		e := &WorkspaceEvent{}
		var data []byte
		if err := rows.Scan(&e.ID, &e.WorkspaceID, &e.Type, &data, &e.CreatedAt); err != nil {
			return nil, err
		}
		if len(data) > 0 {
			e.Data = data
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

// --- Postgres ---

type workspaceEventRepoPostgres struct {
	db *sql.DB
}

func (r *workspaceEventRepoPostgres) Create(workspaceID ID, eventType string, data interface{}) error {
	id, err := NewID()
	if err != nil {
		return err
	}

	var dataJSON []byte
	if data != nil {
		dataJSON, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}
	_, err = r.db.Exec(
		`INSERT INTO workspace_events (id, workspace_id, type, data)
		VALUES ($1, $2, $3, $4)`,
		id.String(), workspaceID.String(), eventType, dataJSON,
	)
	return err
}

func (r *workspaceEventRepoPostgres) GetByID(id ID) (*WorkspaceEvent, error) {
	e := &WorkspaceEvent{}
	var data []byte
	err := r.db.QueryRow(
		`SELECT
			id,
			workspace_id,
			type,
			data,
			created_at
		FROM workspace_events
		WHERE id = $1`,
		id.String(),
	).Scan(&e.ID, &e.WorkspaceID, &e.Type, &data, &e.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if len(data) > 0 {
		e.Data = data
	}
	return e, nil
}

func (r *workspaceEventRepoPostgres) ListByWorkspace(workspaceID ID, limit, offset int) ([]*WorkspaceEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.Query(
		`SELECT
			id,
			workspace_id,
			type,
			data,
			created_at
		FROM workspace_events
		WHERE workspace_id = $1
		ORDER BY created_at DESC
		LIMIT $2
		OFFSET $3`,
		workspaceID.String(), limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanWorkspaceEventRows(rows)
}

func (r *workspaceEventRepoPostgres) ListByWorkspaceAndType(workspaceID ID, eventType string, limit, offset int) ([]*WorkspaceEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.Query(
		`SELECT
			id,
			workspace_id,
			type,
			data,
			created_at
		FROM workspace_events
		WHERE workspace_id = $1 AND type = $2
		ORDER BY created_at DESC
		LIMIT $3
		OFFSET $4`,
		workspaceID.String(), eventType, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanWorkspaceEventRows(rows)
}

func (r *workspaceEventRepoPostgres) CountByWorkspace(workspaceID ID) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM workspace_events
		WHERE workspace_id = $1`,
		workspaceID.String(),
	).Scan(&count)
	return count, err
}

func (r *workspaceEventRepoPostgres) CountByWorkspaceAndType(workspaceID ID, eventType string) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM workspace_events
		WHERE workspace_id = $1 AND type = $2`,
		workspaceID.String(),
		eventType,
	).Scan(&count)
	return count, err
}

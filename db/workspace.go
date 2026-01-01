// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// Workspace is a single row in the workspaces table.
type Workspace struct {
	ID        ID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// WorkspaceRepository is the interface for workspace CRUD.
type WorkspaceRepository interface {
	Create(workspace *Workspace) error
	GetByID(id ID) (*Workspace, error)
	Update(workspace *Workspace) error
	Delete(id ID) error
	List(limit, offset int, userID *ID) ([]*Workspace, error)
	Count(userID *ID) (int64, error)
	GetWorkspaceMembership(workspaceID, userID ID) (*WorkspaceUser, error)
}

// NewWorkspaceRepository returns the repository for the current driver
func NewWorkspaceRepository(db *sql.DB) WorkspaceRepository {
	return &workspaceRepoPostgres{db: db}
}

// scanWorkspaceRows builds a slice of Workspace from query rows
func scanWorkspaceRows(rows *sql.Rows) ([]*Workspace, error) {
	var list []*Workspace
	for rows.Next() {
		w := &Workspace{}
		if err := rows.Scan(&w.ID, &w.Name, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, w)
	}
	return list, rows.Err()
}

type workspaceRepoPostgres struct {
	db *sql.DB
}

// Create inserts a row and fills ID and timestamps via RETURNING.
func (r *workspaceRepoPostgres) Create(workspace *Workspace) error {
	id, err := NewID()
	if err != nil {
		return err
	}
	workspace.ID = id

	err = r.db.QueryRow(
		`INSERT INTO workspaces (id, name)
		VALUES ($1, $2)
		RETURNING created_at, updated_at`,
		workspace.ID.String(),
		workspace.Name,
	).Scan(&workspace.CreatedAt, &workspace.UpdatedAt)
	return err
}

// GetByID returns nil, nil if not found.
func (r *workspaceRepoPostgres) GetByID(id ID) (*Workspace, error) {
	w := &Workspace{}
	err := r.db.QueryRow(
		`SELECT
			id,
			name,
			created_at,
			updated_at
		FROM workspaces
		WHERE id = $1`,
		id.String(),
	).Scan(&w.ID, &w.Name, &w.CreatedAt, &w.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return w, nil
}

// Update updates a workspace.
func (r *workspaceRepoPostgres) Update(workspace *Workspace) error {
	_, err := r.db.Exec(
		`UPDATE workspaces
		SET
			name = $1,
			updated_at = $2
		WHERE id = $3`,
		workspace.Name, time.Now().UTC(), workspace.ID.String(),
	)
	return err
}

// Delete removes a workspace.
func (r *workspaceRepoPostgres) Delete(id ID) error {
	_, err := r.db.Exec(
		`DELETE FROM workspaces WHERE id = $1`,
		id.String(),
	)
	return err
}

// List returns workspaces ordered by created_at DESC. If userID is nil or empty, returns all; else only workspaces the user is in.
func (r *workspaceRepoPostgres) List(limit, offset int, userID *ID) ([]*Workspace, error) {
	if userID == nil || *userID == "" {
		rows, err := r.db.Query(
			`SELECT
				id,
				name,
				created_at,
				updated_at
			FROM workspaces
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2`,
			limit, offset,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return scanWorkspaceRows(rows)
	}
	rows, err := r.db.Query(
		`SELECT
			w.id,
			w.name,
			w.created_at,
			w.updated_at
		FROM workspaces w
		INNER JOIN workspace_users wu ON w.id = wu.workspace_id
		WHERE wu.user_id = $1
		ORDER BY w.created_at DESC
		LIMIT $2 OFFSET $3`,
		userID.String(), limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanWorkspaceRows(rows)
}

// Count returns total workspaces, or only those the user is attached to if userID is set.
func (r *workspaceRepoPostgres) Count(userID *ID) (int64, error) {
	var count int64
	if userID == nil || *userID == "" {
		err := r.db.QueryRow(
			`SELECT COUNT(*)
			FROM workspaces`,
		).Scan(&count)
		return count, err
	}
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM workspaces w
		INNER JOIN workspace_users wu ON w.id = wu.workspace_id
		WHERE wu.user_id = $1`,
		userID.String(),
	).Scan(&count)
	return count, err
}

// GetWorkspaceMembership returns the user's membership (role, etc.) for the workspace, or nil if not a member.
func (r *workspaceRepoPostgres) GetWorkspaceMembership(workspaceID, userID ID) (*WorkspaceUser, error) {
	m := &WorkspaceUser{}
	err := r.db.QueryRow(
		`SELECT
			id,
			workspace_id,
			user_id,
			role,
			created_at,
			updated_at
		FROM workspace_users
		WHERE workspace_id = $1
			AND user_id = $2`,
		workspaceID.String(), userID.String(),
	).Scan(&m.ID, &m.WorkspaceID, &m.UserID, &m.Role, &m.CreatedAt, &m.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

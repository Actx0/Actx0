// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// WorkspaceBackendLock is a single row in the workspace_backend_lock table.
type WorkspaceBackendLock struct {
	ID        ID
	BackendID ID
	LockID    string // Client lock ID (e.g. Terraform's UUID) for unlock
	Operation string
	Info      string
	Who       string
	Version   string
	Path      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// WorkspaceBackendLockRepository is the interface for workspace backend lock CRUD.
type WorkspaceBackendLockRepository interface {
	Create(l *WorkspaceBackendLock) error
	GetByID(id ID) (*WorkspaceBackendLock, error)
	GetAnyActiveByBackendID(backendID ID) (*WorkspaceBackendLock, error)
	GetActiveByBackendID(backendID ID, lockID string) (*WorkspaceBackendLock, error)
	HasActiveByBackendID(backendID ID) (bool, error)
	IsLocked(backendID ID) bool
	Update(l *WorkspaceBackendLock) error
	Delete(id ID) error
	DeleteOlderThan(days int) (int64, error)
	ListByBackendID(backendID ID, limit, offset int) ([]*WorkspaceBackendLock, error)
	CountByBackendID(backendID ID) (int64, error)
	SetInactive(id ID) error
}

// NewWorkspaceBackendLockRepository returns the repository for the current driver.
func NewWorkspaceBackendLockRepository(db *sql.DB) WorkspaceBackendLockRepository {
	return &workspaceBackendLockRepoPostgres{db: db}
}

// scanIsActive reads is_active from PostgreSQL.
func scanIsActive(raw interface{}) bool {
	if raw == nil {
		return false
	}
	switch v := raw.(type) {
	case bool:
		return v
	default:
		return false
	}
}

func scanWorkspaceBackendLockRows(rows *sql.Rows) ([]*WorkspaceBackendLock, error) {
	var list []*WorkspaceBackendLock
	for rows.Next() {
		l := &WorkspaceBackendLock{}
		var isActiveRaw interface{}
		if err := rows.Scan(
			&l.ID,
			&l.BackendID,
			&l.LockID,
			&l.Operation,
			&l.Info,
			&l.Who,
			&l.Version,
			&l.Path,
			&isActiveRaw,
			&l.CreatedAt,
			&l.UpdatedAt,
		); err != nil {
			return nil, err
		}
		l.IsActive = scanIsActive(isActiveRaw)
		list = append(list, l)
	}
	return list, rows.Err()
}

type workspaceBackendLockRepoPostgres struct {
	db *sql.DB
}

func (r *workspaceBackendLockRepoPostgres) Create(l *WorkspaceBackendLock) error {
	id, err := NewID()
	if err != nil {
		return err
	}
	l.ID = id

	return r.db.QueryRow(
		`INSERT INTO workspace_backend_lock
		(id, backend_id, lock_id, operation, info, who, version, path, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at`,
		l.ID.String(),
		l.BackendID.String(),
		l.LockID,
		l.Operation,
		l.Info,
		l.Who,
		l.Version,
		l.Path,
		l.IsActive,
	).Scan(
		&l.CreatedAt,
		&l.UpdatedAt,
	)
}

func (r *workspaceBackendLockRepoPostgres) GetByID(id ID) (*WorkspaceBackendLock, error) {
	l := &WorkspaceBackendLock{}
	var isActiveRaw interface{}
	err := r.db.QueryRow(
		`SELECT
			id, backend_id, lock_id, operation, info, who,
			version, path, is_active, created_at, updated_at
		FROM workspace_backend_lock
		WHERE id = $1`,
		id.String(),
	).Scan(
		&l.ID,
		&l.BackendID,
		&l.LockID,
		&l.Operation,
		&l.Info,
		&l.Who,
		&l.Version,
		&l.Path,
		&isActiveRaw,
		&l.CreatedAt,
		&l.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	l.IsActive = scanIsActive(isActiveRaw)
	return l, nil
}

func (r *workspaceBackendLockRepoPostgres) GetAnyActiveByBackendID(backendID ID) (*WorkspaceBackendLock, error) {
	l := &WorkspaceBackendLock{}
	var isActiveRaw interface{}
	err := r.db.QueryRow(
		`SELECT
			id, backend_id, lock_id, operation, info, who,
			version, path, is_active, created_at, updated_at
		FROM workspace_backend_lock
		WHERE backend_id = $1 AND is_active = true
		LIMIT 1`,
		backendID.String(),
	).Scan(
		&l.ID,
		&l.BackendID,
		&l.LockID,
		&l.Operation,
		&l.Info,
		&l.Who,
		&l.Version,
		&l.Path,
		&isActiveRaw,
		&l.CreatedAt,
		&l.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	l.IsActive = true
	return l, nil
}

func (r *workspaceBackendLockRepoPostgres) HasActiveByBackendID(backendID ID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM workspace_backend_lock
			WHERE backend_id = $1 AND is_active = true
		)`,
		backendID.String(),
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *workspaceBackendLockRepoPostgres) IsLocked(backendID ID) bool {
	ok, _ := r.HasActiveByBackendID(backendID)
	return ok
}

func (r *workspaceBackendLockRepoPostgres) GetActiveByBackendID(backendID ID, lockID string) (*WorkspaceBackendLock, error) {
	l := &WorkspaceBackendLock{}
	var isActiveRaw interface{}
	err := r.db.QueryRow(
		`SELECT
			id, backend_id, lock_id, operation, info, who,
			version, path, is_active, created_at, updated_at
		FROM workspace_backend_lock
		WHERE backend_id = $1 AND lock_id = $2 AND is_active = true
		LIMIT 1`,
		backendID.String(),
		lockID,
	).Scan(
		&l.ID,
		&l.BackendID,
		&l.LockID,
		&l.Operation,
		&l.Info,
		&l.Who,
		&l.Version,
		&l.Path,
		&isActiveRaw,
		&l.CreatedAt,
		&l.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	l.IsActive = true
	return l, nil
}

func (r *workspaceBackendLockRepoPostgres) Update(l *WorkspaceBackendLock) error {
	_, err := r.db.Exec(
		`UPDATE workspace_backend_lock
		SET
			operation = $1,
			info = $2,
			who = $3,
			version = $4,
			path = $5,
			is_active = $6,
			updated_at = $7
		WHERE id = $8`,
		l.Operation,
		l.Info,
		l.Who,
		l.Version,
		l.Path,
		l.IsActive,
		time.Now().UTC(),
		l.ID.String(),
	)
	return err
}

func (r *workspaceBackendLockRepoPostgres) Delete(id ID) error {
	_, err := r.db.Exec(
		`DELETE FROM workspace_backend_lock WHERE id = $1`,
		id.String(),
	)
	return err
}

func (r *workspaceBackendLockRepoPostgres) DeleteOlderThan(days int) (int64, error) {
	if days <= 0 {
		return 0, nil
	}
	cutoff := time.Now().UTC().AddDate(0, 0, -days)
	result, err := r.db.Exec(
		`DELETE FROM workspace_backend_lock
		WHERE is_active = false AND updated_at < $1`,
		cutoff,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *workspaceBackendLockRepoPostgres) SetInactive(id ID) error {
	_, err := r.db.Exec(
		`UPDATE workspace_backend_lock
		SET
			is_active = false,
			updated_at = $1
		WHERE id = $2`,
		time.Now().UTC(),
		id.String(),
	)
	return err
}

func (r *workspaceBackendLockRepoPostgres) ListByBackendID(backendID ID, limit, offset int) ([]*WorkspaceBackendLock, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.Query(
		`SELECT
			id, backend_id, lock_id, operation, info, who,
			version, path, is_active, created_at, updated_at
		FROM workspace_backend_lock
		WHERE backend_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		backendID.String(),
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanWorkspaceBackendLockRows(rows)
}

func (r *workspaceBackendLockRepoPostgres) CountByBackendID(backendID ID) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM workspace_backend_lock
		WHERE backend_id = $1`,
		backendID.String(),
	).Scan(&count)
	return count, err
}

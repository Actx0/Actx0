// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/clivern/actx0/pkg/util"
)

// PromptVersion is a prompt version joined with its parent prompt row.
type PromptVersion struct {
	Id            Id
	PromptId      Id
	WorkspaceId   Id
	Name          string
	Handle        string
	Description   string
	Version       int
	Type          string
	Content       string
	Config        *string
	Labels        *string
	CommitMessage *string
	CommitHash    string
	Meta          *string
	Status        string
	Production    bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// PromptSummary aggregates versions for one prompt.
type PromptSummary struct {
	PromptId        Id
	Name            string
	Handle          string
	Description     string
	VersionCount    int64
	Type            string
	LatestVersion   int
	LatestCreatedAt time.Time
	LatestLabels    *string
	HasProduction   bool
}

// Prompt is metadata stored on the prompts table.
type Prompt struct {
	Id          Id
	WorkspaceId Id
	Name        string
	Handle      string
	Description string
}

// PromptRepository is the interface for prompt and version CRUD.
type PromptRepository interface {
	CreatePrompt(workspaceId Id, name, handle, description string) (Id, error)
	ExistsByHandle(workspaceId Id, handle string) (bool, error)
	Create(prompt *PromptVersion) error
	GetById(id Id) (*PromptVersion, error)
	GetByHandleAndVersion(workspaceId Id, handle string, version int) (*PromptVersion, error)
	GetLatestByHandle(workspaceId Id, handle string) (*PromptVersion, error)
	GetByHandleAndLabel(workspaceId Id, handle, label string) (*PromptVersion, error)
	Update(prompt *PromptVersion) error
	DeleteVersion(id Id) error
	DeletePrompt(workspaceId Id, promptId Id) error
	ListByWorkspaceId(workspaceId Id, limit, offset int) ([]*PromptVersion, error)
	ListSummariesByWorkspace(workspaceId Id, limit, offset int) ([]*PromptSummary, error)
	GetPrompt(workspaceId Id, promptId Id) (*Prompt, error)
	GetSummaryByPromptId(workspaceId Id, promptId Id) (*PromptSummary, error)
	ListVersionsByPromptId(workspaceId Id, promptId Id, limit, offset int) ([]*PromptVersion, error)
	Count() (int64, error)
	CountByWorkspaceId(workspaceId Id) (int64, error)
	CountVersionsByPromptId(workspaceId Id, promptId Id) (int64, error)
	ClearLabelFromPromptId(promptId Id, label string, exceptId *Id) error
	ClearProductionFromPromptId(promptId Id, exceptId *Id) error
}

type PromptRepositoryPostgres struct {
	db *sql.DB
}

// NewPromptRepository returns a prompt repository.
func NewPromptRepository(db *sql.DB) PromptRepository {
	return &PromptRepositoryPostgres{db: db}
}

// CreatePrompt inserts a new prompt row.
func (r *PromptRepositoryPostgres) CreatePrompt(workspaceId Id, name, handle, description string) (Id, error) {
	id, err := NewId()
	if err != nil {
		return "", err
	}
	_, err = r.db.Exec(
		`INSERT INTO prompts (id, workspace_id, name, description, handle) VALUES ($1, $2, $3, $4, $5)`,
		id.String(),
		workspaceId.String(),
		name,
		description,
		handle,
	)
	return id, err
}

// ExistsByHandle reports whether a handle is already used in the workspace.
func (r *PromptRepositoryPostgres) ExistsByHandle(workspaceId Id, handle string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM prompts WHERE workspace_id = $1 AND handle = $2)`,
		workspaceId.String(),
		handle,
	).Scan(&exists)
	return exists, err
}

// Create inserts a prompt version row for an existing prompt.
func (r *PromptRepositoryPostgres) Create(prompt *PromptVersion) error {
	if prompt.PromptId == "" {
		return fmt.Errorf("prompt id required")
	}

	id, err := NewId()
	if err != nil {
		return err
	}
	prompt.Id = id

	if prompt.Version == 0 {
		var maxVersion int
		err := r.db.QueryRow(
			`SELECT COALESCE(MAX(version), 0) FROM prompt_versions WHERE prompt_id = $1`,
			prompt.PromptId.String(),
		).Scan(&maxVersion)
		if err != nil {
			return err
		}
		prompt.Version = maxVersion + 1
	}

	commitHash, err := util.RandomHash()
	if err != nil {
		return err
	}
	prompt.CommitHash = commitHash

	if prompt.Production {
		if err := r.ClearProductionFromPromptId(prompt.PromptId, nil); err != nil {
			return err
		}
	}

	err = r.db.QueryRow(
		`INSERT INTO prompt_versions (
			id, prompt_id, version, type, content, config, labels,
			commit_message, commit_hash, meta, status, production
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING created_at, updated_at`,
		prompt.Id.String(),
		prompt.PromptId.String(),
		prompt.Version,
		prompt.Type,
		prompt.Content,
		prompt.Config,
		prompt.Labels,
		prompt.CommitMessage,
		prompt.CommitHash,
		prompt.Meta,
		prompt.Status,
		prompt.Production,
	).Scan(&prompt.CreatedAt, &prompt.UpdatedAt)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(
		`UPDATE prompts SET updated_at = $1 WHERE id = $2`,
		time.Now().UTC(),
		prompt.PromptId.String(),
	)
	return err
}

// GetById returns a prompt by id.
func (r *PromptRepositoryPostgres) GetById(id Id) (*PromptVersion, error) {
	item := &PromptVersion{}
	err := r.db.QueryRow(`
		SELECT
			v.id, v.prompt_id, p.workspace_id, p.name, p.handle, p.description,
			v.version, v.type, v.content, v.config, v.labels,
			v.commit_message, v.commit_hash, v.meta, v.status, v.production, v.created_at, v.updated_at
		FROM prompt_versions v
		JOIN prompts p ON p.id = v.prompt_id
		WHERE v.id = $1`, id.String(),
	).Scan(
		&item.Id,
		&item.PromptId,
		&item.WorkspaceId,
		&item.Name,
		&item.Handle,
		&item.Description,
		&item.Version,
		&item.Type,
		&item.Content,
		&item.Config,
		&item.Labels,
		&item.CommitMessage,
		&item.CommitHash,
		&item.Meta,
		&item.Status,
		&item.Production,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// GetByHandleAndVersion returns a prompt by handle and version.
func (r *PromptRepositoryPostgres) GetByHandleAndVersion(workspaceId Id, handle string, version int) (*PromptVersion, error) {
	item := &PromptVersion{}
	err := r.db.QueryRow(`
		SELECT
			v.id, v.prompt_id, p.workspace_id, p.name, p.handle, p.description,
			v.version, v.type, v.content, v.config, v.labels,
			v.commit_message, v.commit_hash, v.meta, v.status, v.production, v.created_at, v.updated_at
		FROM prompt_versions v
		JOIN prompts p ON p.id = v.prompt_id
		WHERE p.workspace_id = $1 AND p.handle = $2 AND v.version = $3`,
		workspaceId.String(),
		handle,
		version,
	).Scan(
		&item.Id,
		&item.PromptId,
		&item.WorkspaceId,
		&item.Name,
		&item.Handle,
		&item.Description,
		&item.Version,
		&item.Type,
		&item.Content,
		&item.Config,
		&item.Labels,
		&item.CommitMessage,
		&item.CommitHash,
		&item.Meta,
		&item.Status,
		&item.Production,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// GetLatestByHandle returns the newest version of a prompt by handle.
func (r *PromptRepositoryPostgres) GetLatestByHandle(workspaceId Id, handle string) (*PromptVersion, error) {
	item := &PromptVersion{}
	err := r.db.QueryRow(`
		SELECT
			v.id, v.prompt_id, p.workspace_id, p.name, p.handle, p.description,
			v.version, v.type, v.content, v.config, v.labels,
			v.commit_message, v.commit_hash, v.meta, v.status, v.production, v.created_at, v.updated_at
		FROM prompt_versions v
		JOIN prompts p ON p.id = v.prompt_id
		WHERE p.workspace_id = $1 AND p.handle = $2
		ORDER BY v.version DESC
		LIMIT 1`, workspaceId.String(), handle,
	).Scan(
		&item.Id,
		&item.PromptId,
		&item.WorkspaceId,
		&item.Name,
		&item.Handle,
		&item.Description,
		&item.Version,
		&item.Type,
		&item.Content,
		&item.Config,
		&item.Labels,
		&item.CommitMessage,
		&item.CommitHash,
		&item.Meta,
		&item.Status,
		&item.Production,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// GetByHandleAndLabel returns a prompt by handle and label.
func (r *PromptRepositoryPostgres) GetByHandleAndLabel(workspaceId Id, handle, label string) (*PromptVersion, error) {
	item := &PromptVersion{}
	var err error
	if label == "production" {
		err = r.db.QueryRow(`
			SELECT
				v.id, v.prompt_id, p.workspace_id, p.name, p.handle, p.description,
				v.version, v.type, v.content, v.config, v.labels,
				v.commit_message, v.commit_hash, v.meta, v.status, v.production, v.created_at, v.updated_at
			FROM prompt_versions v
			JOIN prompts p ON p.id = v.prompt_id
			WHERE p.workspace_id = $1 AND p.handle = $2 AND v.production = TRUE
			LIMIT 1`, workspaceId.String(), handle,
		).Scan(
			&item.Id,
			&item.PromptId,
			&item.WorkspaceId,
			&item.Name,
			&item.Handle,
			&item.Description,
			&item.Version,
			&item.Type,
			&item.Content,
			&item.Config,
			&item.Labels,
			&item.CommitMessage,
			&item.CommitHash,
			&item.Meta,
			&item.Status,
			&item.Production,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
	} else {
		labelFilter, marshalErr := json.Marshal([]string{label})
		if marshalErr != nil {
			return nil, marshalErr
		}
		err = r.db.QueryRow(`
			SELECT
				v.id, v.prompt_id, p.workspace_id, p.name, p.handle, p.description,
				v.version, v.type, v.content, v.config, v.labels,
				v.commit_message, v.commit_hash, v.meta, v.status, v.production, v.created_at, v.updated_at
			FROM prompt_versions v
			JOIN prompts p ON p.id = v.prompt_id
			WHERE p.workspace_id = $1 AND p.handle = $2 AND v.labels @> $3::jsonb
			ORDER BY v.version DESC
			LIMIT 1`, workspaceId.String(), handle, string(labelFilter),
		).Scan(
			&item.Id,
			&item.PromptId,
			&item.WorkspaceId,
			&item.Name,
			&item.Handle,
			&item.Description,
			&item.Version,
			&item.Type,
			&item.Content,
			&item.Config,
			&item.Labels,
			&item.CommitMessage,
			&item.CommitHash,
			&item.Meta,
			&item.Status,
			&item.Production,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
	}
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// Update updates an existing prompt row.
func (r *PromptRepositoryPostgres) Update(prompt *PromptVersion) error {
	_, err := r.db.Exec(
		`UPDATE prompt_versions
		SET
			type = $1,
			content = $2,
			config = $3,
			labels = $4,
			commit_message = $5,
			meta = $6,
			status = $7,
			production = $8,
			updated_at = $9
		WHERE id = $10`,
		prompt.Type,
		prompt.Content,
		prompt.Config,
		prompt.Labels,
		prompt.CommitMessage,
		prompt.Meta,
		prompt.Status,
		prompt.Production,
		time.Now().UTC(),
		prompt.Id.String(),
	)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(
		`UPDATE prompts SET updated_at = $1 WHERE id = $2`,
		time.Now().UTC(),
		prompt.PromptId.String(),
	)
	return err
}

// DeleteVersion deletes a single prompt version.
func (r *PromptRepositoryPostgres) DeleteVersion(id Id) error {
	var promptId Id
	err := r.db.QueryRow(
		`SELECT prompt_id FROM prompt_versions WHERE id = $1`,
		id.String(),
	).Scan(&promptId)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`DELETE FROM prompt_versions WHERE id = $1`, id.String())
	if err != nil {
		return err
	}

	var remaining int
	err = r.db.QueryRow(
		`SELECT COUNT(*) FROM prompt_versions WHERE prompt_id = $1`,
		promptId.String(),
	).Scan(&remaining)
	if err != nil {
		return err
	}
	if remaining == 0 {
		_, err = r.db.Exec(`DELETE FROM prompts WHERE id = $1`, promptId.String())
	}
	return err
}

// DeletePrompt deletes a prompt and all of its versions.
func (r *PromptRepositoryPostgres) DeletePrompt(workspaceId Id, promptId Id) error {
	result, err := r.db.Exec(
		`DELETE FROM prompts WHERE id = $1 AND workspace_id = $2`,
		promptId.String(),
		workspaceId.String(),
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// ClearLabelFromPromptId removes a label from all versions of a prompt.
func (r *PromptRepositoryPostgres) ClearLabelFromPromptId(promptId Id, label string, exceptId *Id) error {
	rows, err := r.db.Query(
		`SELECT id, labels FROM prompt_versions WHERE prompt_id = $1`,
		promptId.String(),
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id Id
		var labels *string
		if err := rows.Scan(&id, &labels); err != nil {
			return err
		}
		if exceptId != nil && id == *exceptId {
			continue
		}
		updated, changed := util.RemoveLabelFromJSON(labels, label)
		if !changed {
			continue
		}
		_, err := r.db.Exec(
			`UPDATE prompt_versions SET labels = $1, updated_at = $2 WHERE id = $3`,
			updated,
			time.Now().UTC(),
			id.String(),
		)
		if err != nil {
			return err
		}
	}
	return rows.Err()
}

// ClearProductionFromPromptId removes the production label from a prompt.
func (r *PromptRepositoryPostgres) ClearProductionFromPromptId(promptId Id, exceptId *Id) error {
	query := `UPDATE prompt_versions SET production = FALSE, updated_at = $1 WHERE prompt_id = $2 AND production = TRUE`
	args := []interface{}{time.Now().UTC(), promptId.String()}
	if exceptId != nil {
		query += ` AND id != $3`
		args = append(args, exceptId.String())
	}
	_, err := r.db.Exec(query, args...)
	return err
}

// ListByWorkspaceId lists prompt rows by workspace id.
func (r *PromptRepositoryPostgres) ListByWorkspaceId(workspaceId Id, limit, offset int) ([]*PromptVersion, error) {
	rows, err := r.db.Query(`
		SELECT
			v.id, v.prompt_id, p.workspace_id, p.name, p.handle, p.description,
			v.version, v.type, v.content, v.config, v.labels,
			v.commit_message, v.commit_hash, v.meta, v.status, v.production, v.created_at, v.updated_at
		FROM prompt_versions v
		JOIN prompts p ON p.id = v.prompt_id
		WHERE p.workspace_id = $1
		ORDER BY p.handle ASC, v.version DESC
		LIMIT $2 OFFSET $3`, workspaceId.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*PromptVersion
	for rows.Next() {
		item := &PromptVersion{}
		if err := rows.Scan(
			&item.Id,
			&item.PromptId,
			&item.WorkspaceId,
			&item.Name,
			&item.Handle,
			&item.Description,
			&item.Version,
			&item.Type,
			&item.Content,
			&item.Config,
			&item.Labels,
			&item.CommitMessage,
			&item.CommitHash,
			&item.Meta,
			&item.Status,
			&item.Production,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

// ListVersionsByPromptId lists all versions for a prompt.
func (r *PromptRepositoryPostgres) ListVersionsByPromptId(workspaceId Id, promptId Id, limit, offset int) ([]*PromptVersion, error) {
	rows, err := r.db.Query(`
		SELECT
			v.id, v.prompt_id, p.workspace_id, p.name, p.handle, p.description,
			v.version, v.type, v.content, v.config, v.labels,
			v.commit_message, v.commit_hash, v.meta, v.status, v.production, v.created_at, v.updated_at
		FROM prompt_versions v
		JOIN prompts p ON p.id = v.prompt_id
		WHERE p.workspace_id = $1 AND p.id = $2
		ORDER BY v.version DESC
		LIMIT $3 OFFSET $4`, workspaceId.String(), promptId.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*PromptVersion
	for rows.Next() {
		item := &PromptVersion{}
		if err := rows.Scan(
			&item.Id,
			&item.PromptId,
			&item.WorkspaceId,
			&item.Name,
			&item.Handle,
			&item.Description,
			&item.Version,
			&item.Type,
			&item.Content,
			&item.Config,
			&item.Labels,
			&item.CommitMessage,
			&item.CommitHash,
			&item.Meta,
			&item.Status,
			&item.Production,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

// GetPrompt returns prompt metadata from the prompts table.
func (r *PromptRepositoryPostgres) GetPrompt(workspaceId, promptId Id) (*Prompt, error) {
	item := &Prompt{}
	err := r.db.QueryRow(
		`SELECT id, workspace_id, name, handle, description
		FROM prompts WHERE workspace_id = $1 AND id = $2`,
		workspaceId.String(),
		promptId.String(),
	).Scan(
		&item.Id,
		&item.WorkspaceId,
		&item.Name,
		&item.Handle,
		&item.Description,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// ListSummariesByWorkspace lists prompt summaries for a workspace.
func (r *PromptRepositoryPostgres) ListSummariesByWorkspace(workspaceId Id, limit, offset int) ([]*PromptSummary, error) {
	rows, err := r.db.Query(`
		SELECT
			p.id,
			p.name,
			p.handle,
			p.description,
			vc.version_count,
			lv.type,
			lv.version,
			lv.created_at,
			lv.labels,
			vc.has_production
		FROM prompts p
		JOIN (
			SELECT prompt_id, COUNT(*) AS version_count, BOOL_OR(production) AS has_production
			FROM prompt_versions
			GROUP BY prompt_id
		) vc ON vc.prompt_id = p.id
		JOIN LATERAL (
			SELECT type, version, created_at, labels
			FROM prompt_versions
			WHERE prompt_id = p.id
			ORDER BY version DESC
			LIMIT 1
		) lv ON TRUE
		WHERE p.workspace_id = $1
		ORDER BY p.handle ASC
		LIMIT $2 OFFSET $3`, workspaceId.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*PromptSummary
	for rows.Next() {
		item := &PromptSummary{}
		if err := rows.Scan(
			&item.PromptId,
			&item.Name,
			&item.Handle,
			&item.Description,
			&item.VersionCount,
			&item.Type,
			&item.LatestVersion,
			&item.LatestCreatedAt,
			&item.LatestLabels,
			&item.HasProduction,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

// GetSummaryByPromptId returns the summary for a prompt.
func (r *PromptRepositoryPostgres) GetSummaryByPromptId(workspaceId Id, promptId Id) (*PromptSummary, error) {
	item := &PromptSummary{}
	err := r.db.QueryRow(`
		SELECT
			p.id,
			p.name,
			p.handle,
			p.description,
			(SELECT COUNT(*) FROM prompt_versions v WHERE v.prompt_id = p.id),
			lv.type,
			lv.version,
			lv.created_at,
			lv.labels,
			EXISTS(SELECT 1 FROM prompt_versions v WHERE v.prompt_id = p.id AND v.production)
		FROM prompts p
		JOIN LATERAL (
			SELECT type, version, created_at, labels
			FROM prompt_versions
			WHERE prompt_id = p.id
			ORDER BY version DESC
			LIMIT 1
		) lv ON TRUE
		WHERE p.workspace_id = $1 AND p.id = $2`,
		workspaceId.String(),
		promptId.String(),
	).Scan(
		&item.PromptId,
		&item.Name,
		&item.Handle,
		&item.Description,
		&item.VersionCount,
		&item.Type,
		&item.LatestVersion,
		&item.LatestCreatedAt,
		&item.LatestLabels,
		&item.HasProduction,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return item, err
}

// CountByWorkspaceId counts prompt rows by workspace id.
func (r *PromptRepositoryPostgres) CountByWorkspaceId(workspaceId Id) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM prompts
		WHERE workspace_id = $1`,
		workspaceId.String(),
	).Scan(&count)
	return count, err
}

// CountVersionsByPromptId counts versions for a prompt.
func (r *PromptRepositoryPostgres) CountVersionsByPromptId(workspaceId Id, promptId Id) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM prompt_versions v
		JOIN prompts p ON p.id = v.prompt_id
		WHERE p.workspace_id = $1 AND p.id = $2`,
		workspaceId.String(),
		promptId.String(),
	).Scan(&count)
	return count, err
}

// Count returns the total number of prompt rows.
func (r *PromptRepositoryPostgres) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow(`SELECT COUNT(*) FROM prompts`).Scan(&count)
	return count, err
}

// PromptMeta is a single row in the prompts_meta table.
type PromptMeta struct {
	Id        Id
	PromptId  Id
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PromptMetaRepository is the interface for prompt metadata CRUD.
type PromptMetaRepository interface {
	Create(id Id, key, value string) error
	Get(id Id, key string) (*PromptMeta, error)
	Update(id Id, key, value string) error
	Delete(id Id, key string) error
	ListByPromptId(id Id) ([]*PromptMeta, error)
	Upsert(id Id, key, value string) error
}

type PromptMetaRepositoryPostgres struct {
	db *sql.DB
}

// NewPromptMetaRepository returns the repository for prompt metadata.
func NewPromptMetaRepository(db *sql.DB) PromptMetaRepository {
	return &PromptMetaRepositoryPostgres{db: db}
}

// Create inserts a prompt metadata row.
func (r *PromptMetaRepositoryPostgres) Create(id Id, key, value string) error {
	metaId, err := NewId()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(
		`INSERT INTO prompts_meta (id, prompt_id, key, value)
		VALUES ($1, $2, $3, $4)`,
		metaId.String(), id.String(), key, value,
	)
	return err
}

// Get returns prompt metadata by key.
func (r *PromptMetaRepositoryPostgres) Get(id Id, key string) (*PromptMeta, error) {
	meta := &PromptMeta{}
	err := r.db.QueryRow(
		`SELECT id, prompt_id, key, value, created_at, updated_at
		FROM prompts_meta
		WHERE prompt_id = $1 AND key = $2`,
		id.String(), key,
	).Scan(
		&meta.Id,
		&meta.PromptId,
		&meta.Key,
		&meta.Value,
		&meta.CreatedAt,
		&meta.UpdatedAt,
	)
	if isNotFound(err) {
		return nil, nil
	}
	return meta, err
}

// Update updates an existing prompt metadata row.
func (r *PromptMetaRepositoryPostgres) Update(id Id, key, value string) error {
	_, err := r.db.Exec(
		`UPDATE prompts_meta
		SET value = $1, updated_at = $2
		WHERE prompt_id = $3 AND key = $4`,
		value, time.Now().UTC(), id.String(), key,
	)
	return err
}

// Delete deletes a prompt metadata row.
func (r *PromptMetaRepositoryPostgres) Delete(id Id, key string) error {
	_, err := r.db.Exec(
		`DELETE FROM prompts_meta
		WHERE prompt_id = $1 AND key = $2`,
		id.String(), key,
	)
	return err
}

// ListByPromptId lists prompt metadata rows by prompt id.
func (r *PromptMetaRepositoryPostgres) ListByPromptId(id Id) ([]*PromptMeta, error) {
	rows, err := r.db.Query(
		`SELECT id, prompt_id, key, value, created_at, updated_at
		FROM prompts_meta
		WHERE prompt_id = $1
		ORDER BY key`,
		id.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*PromptMeta
	for rows.Next() {
		meta := &PromptMeta{}
		err := rows.Scan(
			&meta.Id,
			&meta.PromptId,
			&meta.Key,
			&meta.Value,
			&meta.CreatedAt,
			&meta.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, meta)
	}
	return list, rows.Err()
}

// Upsert creates or updates prompt metadata.
func (r *PromptMetaRepositoryPostgres) Upsert(id Id, key, value string) error {
	existing, err := r.Get(id, key)
	if err != nil {
		return err
	}
	if existing == nil {
		return r.Create(id, key, value)
	}
	return r.Update(id, key, value)
}

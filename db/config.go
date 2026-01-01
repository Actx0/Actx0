// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// ConfigEntry is a key-value configuration setting in the DB.
type ConfigEntry struct {
	ID        ID
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ConfigRepository is the interface for config CRUD.
type ConfigRepository interface {
	Create(key, value string) error
	Get(key string) (*ConfigEntry, error)
	Update(key, value string) error
	Delete(key string) error
	List() ([]*ConfigEntry, error)
}

// NewConfigRepository returns the config repository.
func NewConfigRepository(db *sql.DB) ConfigRepository {
	return &configRepoPostgres{db: db}
}

// scanConfigRows builds a slice of ConfigEntry from query rows
func scanConfigRows(rows *sql.Rows) ([]*ConfigEntry, error) {
	var list []*ConfigEntry
	for rows.Next() {
		o := &ConfigEntry{}
		if err := rows.Scan(&o.ID, &o.Key, &o.Value, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, o)
	}
	return list, rows.Err()
}

// --- Postgres ---

// configRepoPostgres implements ConfigRepository for PostgreSQL
type configRepoPostgres struct {
	db *sql.DB
}

// Create inserts a row
func (r *configRepoPostgres) Create(key, value string) error {
	id, err := NewID()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		`INSERT INTO configs (id, key, value)
		VALUES ($1, $2, $3)`,
		id.String(), key, value,
	)
	return err
}

// Get returns a config by key
func (r *configRepoPostgres) Get(key string) (*ConfigEntry, error) {
	o := &ConfigEntry{}
	err := r.db.QueryRow(
		`SELECT id, key, value, created_at, updated_at
		FROM configs
		WHERE key = $1`,
		key,
	).Scan(&o.ID, &o.Key, &o.Value, &o.CreatedAt, &o.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return o, nil
}

// Update updates a config
func (r *configRepoPostgres) Update(key, value string) error {
	_, err := r.db.Exec(
		`UPDATE configs
		SET
			value = $1,
			updated_at = $2
		WHERE key = $3`,
		value, time.Now().UTC(), key,
	)
	return err
}

// Delete removes a config
func (r *configRepoPostgres) Delete(key string) error {
	_, err := r.db.Exec(
		`DELETE FROM configs WHERE key = $1`,
		key,
	)
	return err
}

// List returns a list of configs
func (r *configRepoPostgres) List() ([]*ConfigEntry, error) {
	rows, err := r.db.Query(
		`SELECT id, key, value, created_at, updated_at
		FROM configs
		ORDER BY key`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanConfigRows(rows)
}

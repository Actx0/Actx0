// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// APIKey is the DB row for a user API key.
type APIKey struct {
	ID        ID
	UserID    ID
	Name      *string
	Key       string
	ExpiresAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// APIKeyRepository is the interface for API key CRUD.
type APIKeyRepository interface {
	Create(apiKey *APIKey) error
	GetByID(id ID) (*APIKey, error)
	GetByKey(key string) (*APIKey, error)
	GetByUserID(userID ID) ([]*APIKey, error)
	Delete(id ID) error
	DeleteByUserID(userID ID) error
	DeleteExpired() (int64, error)
	Count() (int64, error)
	CountByUserID(userID ID) (int64, error)
}

// NewAPIKeyRepository returns the repository for the current driver
func NewAPIKeyRepository(db *sql.DB) APIKeyRepository {
	return &apiKeyRepoPostgres{db: db}
}

// scanAPIKeyRows builds a slice of APIKey from query rows
func scanAPIKeyRows(rows *sql.Rows) ([]*APIKey, error) {
	var list []*APIKey
	for rows.Next() {
		k := &APIKey{}
		if err := rows.Scan(
			&k.ID,
			&k.UserID,
			&k.Name,
			&k.Key,
			&k.ExpiresAt,
			&k.CreatedAt,
			&k.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, k)
	}
	return list, rows.Err()
}

// --- Postgres ---

// apiKeyRepoPostgres implements APIKeyRepository for PostgreSQL
type apiKeyRepoPostgres struct {
	db *sql.DB
}

// Create inserts a row
func (r *apiKeyRepoPostgres) Create(apiKey *APIKey) error {
	id, err := NewID()
	if err != nil {
		return err
	}
	apiKey.ID = id

	err = r.db.QueryRow(
		`INSERT INTO user_api_keys
		(id, user_id, name, key, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at`,
		apiKey.ID.String(),
		apiKey.UserID.String(),
		apiKey.Name,
		apiKey.Key,
		apiKey.ExpiresAt,
	).Scan(
		&apiKey.CreatedAt,
		&apiKey.UpdatedAt,
	)
	return err
}

// GetByID returns a API key by ID
func (r *apiKeyRepoPostgres) GetByID(id ID) (*APIKey, error) {
	k := &APIKey{}
	err := r.db.QueryRow(
		`SELECT
			id, user_id, name, key, expires_at, created_at, updated_at
		FROM user_api_keys
		WHERE id = $1`,
		id.String(),
	).Scan(
		&k.ID,
		&k.UserID,
		&k.Name,
		&k.Key,
		&k.ExpiresAt,
		&k.CreatedAt,
		&k.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return k, nil
}

// GetByKey returns a API key by raw key
func (r *apiKeyRepoPostgres) GetByKey(key string) (*APIKey, error) {
	if key == "" {
		return nil, nil
	}
	k := &APIKey{}
	err := r.db.QueryRow(
		`SELECT
			id, user_id, name, key, expires_at, created_at, updated_at
		FROM user_api_keys
		WHERE key = $1 AND (expires_at IS NULL OR expires_at > $2)`,
		key, time.Now().UTC(),
	).Scan(
		&k.ID,
		&k.UserID,
		&k.Name,
		&k.Key,
		&k.ExpiresAt,
		&k.CreatedAt,
		&k.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return k, nil
}

// GetByUserID returns a list of API keys by user ID
func (r *apiKeyRepoPostgres) GetByUserID(userID ID) ([]*APIKey, error) {
	rows, err := r.db.Query(
		`SELECT
			id, user_id, name, key, expires_at, created_at, updated_at
		FROM user_api_keys
		WHERE user_id = $1
		ORDER BY created_at DESC`,
		userID.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAPIKeyRows(rows)
}

// Delete removes a API key
func (r *apiKeyRepoPostgres) Delete(id ID) error {
	_, err := r.db.Exec(
		`DELETE FROM user_api_keys WHERE id = $1`,
		id.String(),
	)
	return err
}

// DeleteByUserID removes a API key by user ID
func (r *apiKeyRepoPostgres) DeleteByUserID(userID ID) error {
	_, err := r.db.Exec(
		`DELETE FROM user_api_keys WHERE user_id = $1`,
		userID.String(),
	)
	return err
}

// DeleteExpired removes expired API keys
func (r *apiKeyRepoPostgres) DeleteExpired() (int64, error) {
	result, err := r.db.Exec(
		`DELETE FROM user_api_keys
		WHERE expires_at IS NOT NULL AND expires_at < $1`,
		time.Now().UTC(),
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Count returns the total number of API keys
func (r *apiKeyRepoPostgres) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM user_api_keys
		WHERE (expires_at IS NULL OR expires_at > $1)`,
		time.Now().UTC(),
	).Scan(&count)
	return count, err
}

// CountByUserID returns the total number of API keys by user ID
func (r *apiKeyRepoPostgres) CountByUserID(userID ID) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM user_api_keys
		WHERE user_id = $1 AND (expires_at IS NULL OR expires_at > $2)`,
		userID.String(), time.Now().UTC(),
	).Scan(&count)
	return count, err
}

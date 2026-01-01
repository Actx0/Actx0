// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// Session is the DB row for a user session.
type Session struct {
	ID        ID
	Token     string
	UserID    ID
	IPAddress *string
	UserAgent *string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SessionRepository is the interface for session CRUD.
type SessionRepository interface {
	Create(session *Session) error
	GetByToken(token string) (*Session, error)
	GetByID(id ID) (*Session, error)
	GetByUserID(userID ID) ([]*Session, error)
	Delete(id ID) error
	DeleteByToken(token string) error
	DeleteByUserID(userID ID) error
	DeleteExpired() (int64, error)
	IsValid(token string) (bool, error)
	UpdateExpiration(id ID, expiresAt time.Time) error
	Count() (int64, error)
	CountByUserID(userID ID) (int64, error)
}

// NewSessionRepository returns the repository for the current driver
func NewSessionRepository(db *sql.DB) SessionRepository {
	return &sessionRepoPostgres{db: db}
}

// scanSessionRows builds a slice of Session from query rows
func scanSessionRows(rows *sql.Rows) ([]*Session, error) {
	var list []*Session
	for rows.Next() {
		s := &Session{}
		if err := rows.Scan(
			&s.ID,
			&s.Token,
			&s.UserID,
			&s.IPAddress,
			&s.UserAgent,
			&s.ExpiresAt,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

// --- Postgres ---

// sessionRepoPostgres implements SessionRepository for PostgreSQL
type sessionRepoPostgres struct {
	db *sql.DB
}

// Create inserts a row
func (r *sessionRepoPostgres) Create(session *Session) error {
	id, err := NewID()
	if err != nil {
		return err
	}
	session.ID = id

	_, err = r.db.Exec(
		`INSERT INTO user_sessions
		(id, token, user_id, ip_address, user_agent, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		session.ID.String(),
		session.Token,
		session.UserID.String(),
		session.IPAddress,
		session.UserAgent,
		session.ExpiresAt,
	)
	return err
}

// GetByToken returns a session by token
func (r *sessionRepoPostgres) GetByToken(token string) (*Session, error) {
	s := &Session{}
	err := r.db.QueryRow(
		`SELECT
			id, token, user_id, ip_address, user_agent,
			expires_at, created_at, updated_at
		FROM user_sessions
		WHERE token = $1`,
		token,
	).Scan(
		&s.ID,
		&s.Token,
		&s.UserID,
		&s.IPAddress,
		&s.UserAgent,
		&s.ExpiresAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

// GetByID returns a session by ID
func (r *sessionRepoPostgres) GetByID(id ID) (*Session, error) {
	s := &Session{}
	err := r.db.QueryRow(
		`SELECT
			id, token, user_id, ip_address, user_agent,
			expires_at, created_at, updated_at
		FROM user_sessions
		WHERE id = $1`,
		id.String(),
	).Scan(
		&s.ID,
		&s.Token,
		&s.UserID,
		&s.IPAddress,
		&s.UserAgent,
		&s.ExpiresAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

// GetByUserID returns a list of sessions by user ID
func (r *sessionRepoPostgres) GetByUserID(userID ID) ([]*Session, error) {
	rows, err := r.db.Query(
		`SELECT
			id, token, user_id, ip_address, user_agent,
			expires_at, created_at, updated_at
		FROM user_sessions
		WHERE user_id = $1
		ORDER BY created_at DESC`,
		userID.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSessionRows(rows)
}

// Delete removes a session
func (r *sessionRepoPostgres) Delete(id ID) error {
	_, err := r.db.Exec(
		`DELETE FROM user_sessions WHERE id = $1`,
		id.String(),
	)
	return err
}

// DeleteByToken removes a session by token
func (r *sessionRepoPostgres) DeleteByToken(token string) error {
	_, err := r.db.Exec(
		`DELETE FROM user_sessions WHERE token = $1`,
		token,
	)
	return err
}

// DeleteByUserID removes a session by user ID
func (r *sessionRepoPostgres) DeleteByUserID(userID ID) error {
	_, err := r.db.Exec(
		`DELETE FROM user_sessions WHERE user_id = $1`,
		userID.String(),
	)
	return err
}

// DeleteExpired removes expired sessions
func (r *sessionRepoPostgres) DeleteExpired() (int64, error) {
	result, err := r.db.Exec(
		`DELETE FROM user_sessions
		WHERE expires_at < $1`,
		time.Now().UTC(),
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// IsValid checks if a session is valid
func (r *sessionRepoPostgres) IsValid(token string) (bool, error) {
	session, err := r.GetByToken(token)
	if err != nil {
		return false, err
	}
	if session == nil {
		return false, nil
	}
	return session.ExpiresAt.After(time.Now().UTC()), nil
}

// UpdateExpiration updates the expiration time of a session
func (r *sessionRepoPostgres) UpdateExpiration(id ID, expiresAt time.Time) error {
	_, err := r.db.Exec(
		`UPDATE user_sessions
		SET
			expires_at = $1,
			updated_at = $2
		WHERE id = $3`,
		expiresAt, time.Now().UTC(), id.String(),
	)
	return err
}

// Count returns the total number of sessions
func (r *sessionRepoPostgres) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM user_sessions
		WHERE expires_at > $1`,
		time.Now().UTC(),
	).Scan(&count)
	return count, err
}

// CountByUserID returns the total number of sessions by user ID
func (r *sessionRepoPostgres) CountByUserID(userID ID) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM user_sessions
		WHERE user_id = $1 AND expires_at > $2`,
		userID.String(), time.Now().UTC(),
	).Scan(&count)
	return count, err
}

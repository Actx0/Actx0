// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// UserInvite is the DB row for an invite (by email).
type UserInvite struct {
	ID            ID
	Email         string
	Role          string
	Token         string
	Status        string
	InviterUserID ID
	ExpiresAt     time.Time
	AcceptedAt    *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// UserInviteRepository is the interface for user invite CRUD.
type UserInviteRepository interface {
	Create(invite *UserInvite) error
	GetByID(id ID) (*UserInvite, error)
	GetByToken(token string) (*UserInvite, error)
	ListPaginated(limit, offset int) ([]*UserInvite, error)
	UpdateStatus(id ID, status string, acceptedAt *time.Time) error
	MarkExpiredAsExpired() (int64, error)
	Delete(id ID) error
	Count() (int64, error)
	CountPendingByEmail(email string) (int64, error)
}

// NewUserInviteRepository returns the repository for the current driver
func NewUserInviteRepository(db *sql.DB) UserInviteRepository {
	return &userInviteRepoPostgres{db: db}
}

// scanUserInviteRows builds a slice of UserInvite from query rows
func scanUserInviteRows(rows *sql.Rows) ([]*UserInvite, error) {
	var list []*UserInvite
	for rows.Next() {
		u := &UserInvite{}
		if err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.Role,
			&u.Token,
			&u.Status,
			&u.InviterUserID,
			&u.ExpiresAt,
			&u.AcceptedAt,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	return list, rows.Err()
}

// --- Postgres ---

// userInviteRepoPostgres implements UserInviteRepository for PostgreSQL
type userInviteRepoPostgres struct {
	db *sql.DB
}

// Create inserts a row
func (r *userInviteRepoPostgres) Create(invite *UserInvite) error {
	id, err := NewID()
	if err != nil {
		return err
	}
	invite.ID = id

	_, err = r.db.Exec(
		`INSERT INTO user_invites
		(id, email, role, token, status, inviter_user_id, expires_at, accepted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		invite.ID.String(),
		invite.Email,
		invite.Role,
		invite.Token,
		invite.Status,
		invite.InviterUserID.String(),
		invite.ExpiresAt,
		invite.AcceptedAt,
	)
	return err
}

// GetByID returns a user invite by ID
func (r *userInviteRepoPostgres) GetByID(id ID) (*UserInvite, error) {
	u := &UserInvite{}
	err := r.db.QueryRow(
		`SELECT
			id, email, role, token, status, inviter_user_id,
			expires_at, accepted_at, created_at, updated_at
		FROM user_invites
		WHERE id = $1`,
		id.String(),
	).Scan(
		&u.ID,
		&u.Email,
		&u.Role,
		&u.Token,
		&u.Status,
		&u.InviterUserID,
		&u.ExpiresAt,
		&u.AcceptedAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

// GetByToken returns a user invite by token
func (r *userInviteRepoPostgres) GetByToken(token string) (*UserInvite, error) {
	u := &UserInvite{}
	err := r.db.QueryRow(
		`SELECT
			id, email, role, token, status, inviter_user_id,
			expires_at, accepted_at, created_at, updated_at
		FROM user_invites
		WHERE token = $1`,
		token,
	).Scan(
		&u.ID,
		&u.Email,
		&u.Role,
		&u.Token,
		&u.Status,
		&u.InviterUserID,
		&u.ExpiresAt,
		&u.AcceptedAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

// ListPaginated returns a paginated list of all user invites (for admins).
func (r *userInviteRepoPostgres) ListPaginated(limit, offset int) ([]*UserInvite, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.Query(
		`SELECT
			id, email, role, token, status, inviter_user_id,
			expires_at, accepted_at, created_at, updated_at
		FROM user_invites
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanUserInviteRows(rows)
}

// UpdateStatus updates the status of a user invite (e.g. to "accepted" with acceptedAt set).
func (r *userInviteRepoPostgres) UpdateStatus(id ID, status string, acceptedAt *time.Time) error {
	_, err := r.db.Exec(
		`UPDATE user_invites
		SET
			status = $1,
			accepted_at = $2,
			updated_at = $3
		WHERE id = $4`,
		status,
		acceptedAt,
		time.Now().UTC(),
		id.String(),
	)
	return err
}

// MarkExpiredAsExpired sets status to 'expired' for all pending invites past expires_at. Returns the number updated.
func (r *userInviteRepoPostgres) MarkExpiredAsExpired() (int64, error) {
	result, err := r.db.Exec(
		`UPDATE user_invites
		SET
			status = 'expired',
			updated_at = $1
		WHERE status = 'pending' AND expires_at < $2`,
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete removes a user invite
func (r *userInviteRepoPostgres) Delete(id ID) error {
	_, err := r.db.Exec(
		`DELETE FROM user_invites WHERE id = $1`,
		id.String(),
	)
	return err
}

// Count returns the total number of user invites
func (r *userInviteRepoPostgres) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM user_invites`,
	).Scan(&count)
	return count, err
}

// CountPendingByEmail returns the total number of pending user invites by email
func (r *userInviteRepoPostgres) CountPendingByEmail(email string) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM user_invites
		WHERE email = $1 AND status = 'pending' AND expires_at > $2`,
		email,
		time.Now().UTC(),
	).Scan(&count)
	return count, err
}

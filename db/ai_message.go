// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// AIMessage is the DB row for one AI chat message.
type AIMessage struct {
	ID        ID
	UserID    ID
	Role      string // e.g. "user", "assistant", "system"
	Content   string
	CreatedAt time.Time
}

// AIMessageRepository is the interface for AI message CRUD.
type AIMessageRepository interface {
	Create(msg *AIMessage) error
	GetByID(id ID) (*AIMessage, error)
	GetByUserID(userID ID, limit, offset int) ([]*AIMessage, error)
	Delete(id ID) error
	DeleteByUserID(userID ID) error
	CountByUserID(userID ID) (int64, error)
}

// NewAIMessageRepository returns the repository for the current driver.
func NewAIMessageRepository(db *sql.DB) AIMessageRepository {
	return &aiMessageRepoPostgres{db: db}
}

func scanAIMessageRows(rows *sql.Rows) ([]*AIMessage, error) {
	var list []*AIMessage
	for rows.Next() {
		m := &AIMessage{}
		if err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.Role,
			&m.Content,
			&m.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

// --- Postgres ---

type aiMessageRepoPostgres struct {
	db *sql.DB
}

func (r *aiMessageRepoPostgres) Create(msg *AIMessage) error {
	id, err := NewID()
	if err != nil {
		return err
	}
	msg.ID = id

	return r.db.QueryRow(
		`INSERT INTO user_ai_messages (id, user_id, role, content)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at`,
		msg.ID.String(),
		msg.UserID.String(),
		msg.Role,
		msg.Content,
	).Scan(
		&msg.CreatedAt,
	)
}

func (r *aiMessageRepoPostgres) GetByID(id ID) (*AIMessage, error) {
	m := &AIMessage{}
	err := r.db.QueryRow(
		`SELECT
			id, user_id, role, content, created_at
		FROM user_ai_messages
		WHERE id = $1`,
		id.String(),
	).Scan(
		&m.ID,
		&m.UserID,
		&m.Role,
		&m.Content,
		&m.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *aiMessageRepoPostgres) GetByUserID(userID ID, limit, offset int) ([]*AIMessage, error) {
	rows, err := r.db.Query(
		`SELECT
			id, user_id, role, content, created_at
		FROM user_ai_messages
		WHERE user_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3`,
		userID.String(),
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAIMessageRows(rows)
}

func (r *aiMessageRepoPostgres) Delete(id ID) error {
	_, err := r.db.Exec(
		`DELETE FROM user_ai_messages WHERE id = $1`,
		id.String(),
	)
	return err
}

func (r *aiMessageRepoPostgres) DeleteByUserID(userID ID) error {
	_, err := r.db.Exec(
		`DELETE FROM user_ai_messages WHERE user_id = $1`,
		userID.String(),
	)
	return err
}

func (r *aiMessageRepoPostgres) CountByUserID(userID ID) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM user_ai_messages
		WHERE user_id = $1`,
		userID.String(),
	).Scan(&count)
	return count, err
}

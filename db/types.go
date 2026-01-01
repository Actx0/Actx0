// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

// ID is a PostgreSQL UUID primary or foreign key.
// It implements sql.Scanner and driver.Valuer so UUIDs can be scanned into DB structs.
type ID string

// NewID returns a new UUID-backed database ID.
func NewID() (ID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("generate id: %w", err)
	}
	return ID(id.String()), nil
}

// Scan implements sql.Scanner. Accepts []byte or string UUID values.
func (id *ID) Scan(value interface{}) error {
	if value == nil {
		*id = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*id = ID(string(v))
		return nil
	case string:
		*id = ID(v)
		return nil
	default:
		return fmt.Errorf("db: cannot scan type %T into ID", value)
	}
}

// Value implements driver.Valuer for INSERT/UPDATE.
func (id ID) Value() (driver.Value, error) {
	if id == "" {
		return nil, nil
	}
	return string(id), nil
}

// String returns the ID as string.
func (id ID) String() string {
	return string(id)
}

// ValueOrNil returns the ID as interface{} for use in queries; returns nil if id is empty (for nullable columns).
func (id ID) ValueOrNil() interface{} {
	if id == "" {
		return nil
	}
	return string(id)
}

// ValueOrNilPtr returns the ID as interface{} for nullable columns; returns nil if the pointer is nil or ID is empty.
func (id *ID) ValueOrNilPtr() interface{} {
	if id == nil || *id == "" {
		return nil
	}
	return string(*id)
}

func nullString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

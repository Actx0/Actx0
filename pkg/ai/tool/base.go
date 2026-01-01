// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package tool

import (
	"context"

	"github.com/clivern/actx0/db"
)

// UserIDContextKey is the context key for the current user ID.
const UserIDContextKey = "user_id"

// UserIDFromContext returns the user ID from context (set by Generate).
func UserIDFromContext(ctx context.Context) db.ID {
	if id, ok := ctx.Value(UserIDContextKey).(db.ID); ok {
		return id
	}
	return ""
}

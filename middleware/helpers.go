// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/module"
)

const (
	// ContextKeyWorkspace is the key for storing a loaded workspace in context.
	ContextKeyWorkspace contextKey = "workspace"
)

// IsAssetPath returns true if path is a static asset
func IsAssetPath(path string) bool {
	return strings.HasPrefix(path, "/assets/")
}

// GetPrincipal returns the authenticated user or workspace access key from context.
func GetPrincipal(r *http.Request) (*db.User, *db.WorkspaceAccessKey, bool) {
	p, ok := module.PrincipalFrom(r.Context())
	if !ok || p == nil {
		return nil, nil, false
	}
	return p.User, p.AccessKey, true
}

// GetWorkspaceFromContext returns the workspace attached by Protect with Workspace: true.
func GetWorkspaceFromContext(ctx context.Context) (*db.Workspace, bool) {
	workspace, ok := ctx.Value(ContextKeyWorkspace).(*db.Workspace)
	return workspace, ok && workspace != nil
}

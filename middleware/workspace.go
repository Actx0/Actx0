// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package middleware

import (
	"context"
	"net/http"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/pkg/util"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

const (
	// ContextKeyWorkspace is the key for storing a loaded workspace in context.
	ContextKeyWorkspace contextKey = "workspace"
)

// RequireWorkspace loads the workspace from the {workspaceId} route param and stores it
// on the request context. Handlers and RequireWorkspacePermission can reuse it to avoid
// a second database lookup.
func RequireWorkspace() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wid := chi.URLParam(r, "workspaceId")
			if wid == "" {
				log.Error().
					Str("path", r.URL.Path).
					Msg("RequireWorkspace missing workspace id param")
				util.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
					"errorMessage": "Invalid route configuration",
				})
				return
			}

			workspace, err := db.NewWorkspaceRepository(db.GetDB()).GetById(db.Id(wid))
			if err != nil {
				log.Error().
					Err(err).
					Str("workspaceId", wid).
					Msg("Failed to load workspace")
				util.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
					"errorMessage": "Failed to load workspace",
				})
				return
			}
			if workspace == nil {
				util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
					"errorMessage": "Workspace not found",
				})
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyWorkspace, workspace)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetWorkspaceFromContext returns the workspace attached by RequireWorkspace.
func GetWorkspaceFromContext(ctx context.Context) (*db.Workspace, bool) {
	workspace, ok := ctx.Value(ContextKeyWorkspace).(*db.Workspace)
	return workspace, ok && workspace != nil
}

// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package middleware

import (
	"errors"
	"net/http"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/module"
	"github.com/clivern/actx0/pkg/util"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// RequireUser ensures the request was authenticated as a user (session or API key).
// Access-key principals are rejected.
func RequireUser() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok || user == nil {
				log.Info().
					Str("path", r.URL.Path).
					Msg("User required: principal not in context")
				util.WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{
					"errorMessage": "Not authenticated",
				})
				return
			}

			if !user.IsActive {
				log.Info().
					Str("path", r.URL.Path).
					Str("userId", user.Id.String()).
					Msg("Inactive user blocked from user-only route")
				util.WriteJSON(w, http.StatusForbidden, map[string]interface{}{
					"errorMessage": "Account is inactive",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireWorkspacePermission checks that the principal may perform permission on the
// workspace from {workspaceId}. Reuses the workspace from RequireWorkspace when present.
// Chain RequireUser first on routes that must reject workspace access keys.
func RequireWorkspacePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			perm := module.NewPerm(db.GetDB())

			if workspace, ok := GetWorkspaceFromContext(r.Context()); ok {
				perm = perm.WithWorkspace(workspace)
			} else {
				wid := chi.URLParam(r, "workspaceId")
				if wid == "" {
					log.Error().
						Str("path", r.URL.Path).
						Str("permission", permission).
						Msg("Workspace permission check missing workspace id param")
					util.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
						"errorMessage": "Invalid route configuration",
					})
					return
				}
				perm = perm.WithWorkspaceId(db.Id(wid))
			}

			if user, ok := GetUserFromContext(r.Context()); ok && user != nil {
				if !user.IsActive {
					log.Info().
						Str("path", r.URL.Path).
						Str("userId", user.Id.String()).
						Msg("Inactive user blocked from workspace route")
					util.WriteJSON(w, http.StatusForbidden, map[string]interface{}{
						"errorMessage": "Account is inactive",
					})
					return
				}
				perm = perm.WithUser(user)
			} else if key, ok := GetAccessKeyFromContext(r.Context()); ok && key != nil {
				perm = perm.WithAccessKey(key)
			} else {
				log.Info().
					Str("path", r.URL.Path).
					Msg("Workspace permission check failed: no principal in context")
				util.WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{
					"errorMessage": "Not authenticated",
				})
				return
			}

			allowed, err := perm.Can(permission)
			if err != nil {
				if errors.Is(err, module.ErrPermissionWorkspaceMissing) {
					util.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
						"errorMessage": "Workspace not found",
					})
					return
				}
				log.Error().
					Err(err).
					Str("path", r.URL.Path).
					Str("permission", permission).
					Msg("Failed to check workspace permission")
				util.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
					"errorMessage": "Failed to check access",
				})
				return
			}
			if !allowed {
				log.Info().
					Str("path", r.URL.Path).
					Str("permission", permission).
					Msg("Workspace permission denied")
				util.WriteJSON(w, http.StatusForbidden, map[string]interface{}{
					"errorMessage": "Insufficient permissions",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetPrincipal returns the authenticated user or workspace access key from context.
func GetPrincipal(r *http.Request) (*db.User, *db.WorkspaceAccessKey, bool) {
	p, ok := module.PrincipalFrom(r.Context())
	if !ok || p == nil {
		return nil, nil, false
	}
	return p.User, p.AccessKey, true
}

// UserWorkspacePerm returns middleware for user-only workspace routes.
func UserWorkspacePerm(permission string) []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		RequireUser(),
		RequireWorkspacePermission(permission),
	}
}

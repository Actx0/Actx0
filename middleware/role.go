// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

// Package middleware provides HTTP middleware functions for the Tut application.
package middleware

import (
	"net/http"

	"github.com/clivern/actx0/pkg/util"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// RequireRole creates a middleware that checks if the authenticated user has one of the required roles
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context
			user, ok := GetUserFromContext(r.Context())
			if !ok || user == nil {
				log.Info().
					Str("path", r.URL.Path).
					Msg("Role check failed: user not in context")
				util.WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{
					"errorMessage": "Not authenticated",
				})
				return
			}

			// Check if user is active
			if !user.IsActive {
				log.Info().
					Str("path", r.URL.Path).
					Str("userId", user.Id.String()).
					Msg("Inactive user blocked from protected route")
				util.WriteJSON(w, http.StatusForbidden, map[string]interface{}{
					"errorMessage": "Account is inactive",
				})
				return
			}

			if !lo.Contains(allowedRoles, user.Role) {
				log.Info().
					Str("path", r.URL.Path).
					Str("userId", user.Id.String()).
					Str("userRole", user.Role).
					Strs("allowedRoles", allowedRoles).
					Msg("Insufficient role for route")
				util.WriteJSON(w, http.StatusForbidden, map[string]interface{}{
					"errorMessage": "Insufficient permissions",
				})
				return
			}

			// User has required role, proceed
			next.ServeHTTP(w, r)
		})
	}
}

// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/module"
	"github.com/clivern/actx0/pkg/util"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

const (
	// ContextKeyUser is the key for storing user in context
	ContextKeyUser contextKey = "user"
	// ContextKeyAccessKey is the key for storing workspace access key in context
	ContextKeyAccessKey contextKey = "access_key"
)

// SessionAuth creates a session-based authentication middleware
func SessionAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if shouldSkipAuth(r.URL.Path) {
				if !IsAssetPath(r.URL.Path) {
					log.Info().
						Str("path", r.URL.Path).
						Msg("Skipping auth for public route")
				}
				next.ServeHTTP(w, r)
				return
			}

			// API Key Check
			apiKey := r.Header.Get("X-API-Key")
			if lo.IsNotEmpty(apiKey) {
				user, err := db.NewUserRepository(db.GetDB()).GetByAPIKey(apiKey)
				if err != nil {
					log.Info().
						Err(err).
						Str("path", r.URL.Path).
						Msg("API key validation failed")
					util.WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{
						"errorMessage": "Invalid API key",
					})
					return
				}
				log.Info().
					Str("path", r.URL.Path).
					Msg("API key validated")

				ctx := WithUserContext(r.Context(), user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Access Key Check
			accessKey := r.Header.Get("X-Access-Key")
			if lo.IsNotEmpty(accessKey) {
				key, err := db.NewWorkspaceAccessKeyRepository(db.GetDB()).GetByKey(accessKey)
				if err != nil {
					log.Info().
						Err(err).
						Str("path", r.URL.Path).
						Msg("Access key validation failed")
					util.WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{
						"errorMessage": "Invalid access key",
					})
					return
				}
				if key == nil {
					log.Info().
						Str("path", r.URL.Path).
						Msg("Access key not found or expired")
					util.WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{
						"errorMessage": "Invalid access key",
					})
					return
				}
				log.Info().
					Str("path", r.URL.Path).
					Str("workspaceId", key.WorkspaceId.String()).
					Msg("Access key validated")

				ctx := WithAccessKeyContext(r.Context(), key)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Session Token Check
			sessionToken := util.GetCookie(r, "_actx0_session")
			if lo.IsEmpty(sessionToken) {
				log.Info().
					Str("path", r.URL.Path).
					Msg("No session cookie found")
				util.WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{
					"errorMessage": "Not authenticated",
				})
				return
			}

			sessionManager := module.NewSessionManager(
				db.NewSessionRepository(db.GetDB()),
				db.NewUserRepository(db.GetDB()),
			)

			user, _, err := sessionManager.ValidateSession(sessionToken)
			if err != nil {
				util.DeleteCookie(w, "_actx0_session")
				if user != nil {
					sessionManager.RevokeUserSessions(user.Id)
				}

				log.Info().
					Err(err).
					Str("path", r.URL.Path).
					Msg("Session validation failed")
				util.WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{
					"errorMessage": "Invalid or expired session",
				})
				return
			}

			log.Info().
				Str("path", r.URL.Path).
				Msg("Session validated")

			ctx := WithUserContext(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// shouldSkipAuth determines if authentication should be skipped for a given path
func shouldSkipAuth(path string) bool {
	if strings.HasPrefix(path, "/api/v1/public/") || !strings.HasPrefix(path, "/api/v1/") {
		return true
	}

	return false
}

// GetUserFromContext retrieves the user from the request context
func GetUserFromContext(ctx context.Context) (*db.User, bool) {
	return module.UserFrom(ctx)
}

// GetAccessKeyFromContext retrieves the workspace access key from the request context
func GetAccessKeyFromContext(ctx context.Context) (*db.WorkspaceAccessKey, bool) {
	return module.AccessKeyFrom(ctx)
}

// WithUserContext sets the user in the request context
func WithUserContext(ctx context.Context, user *db.User) context.Context {
	ctx = module.WithPrincipal(ctx, &module.Principal{User: user})
	return context.WithValue(ctx, ContextKeyUser, user)
}

// WithAccessKeyContext sets the workspace access key in the request context
func WithAccessKeyContext(ctx context.Context, key *db.WorkspaceAccessKey) context.Context {
	ctx = module.WithPrincipal(ctx, &module.Principal{AccessKey: key})
	return context.WithValue(ctx, ContextKeyAccessKey, key)
}

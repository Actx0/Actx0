// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package middleware

import (
	"net/http"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/locale"
	"github.com/clivern/actx0/module"
	"github.com/clivern/actx0/pkg/util"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// StatusRecorder wraps the response writer to capture the HTTP status code.
type StatusRecorder struct {
	http.ResponseWriter
	status int
}

// WriteHeader captures the HTTP status code for logging.
func (s *StatusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// IsMeteredAPICall reports whether the request is authenticated with an API or
// access key (not a browser session from the web UI).
func IsMeteredAPICall(r *http.Request) bool {
	return lo.IsNotEmpty(r.Header.Get("X-API-Key")) || lo.IsNotEmpty(r.Header.Get("X-Access-Key"))
}

// TrackWorkspaceAPICall increments workspace API usage after authenticated, authorized
// key-based requests. Browser session (web UI) calls are excluded.
func TrackWorkspaceAPICall() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			workspace, hasWorkspace := GetWorkspaceFromContext(r.Context())
			if hasWorkspace && IsMeteredAPICall(r) {
				if _, _, ok := GetPrincipal(r); ok {
					usageRepo := db.NewUsageRepository(db.GetDB(false))
					subscriptions := db.NewSubscriptionRepository(db.GetDB(false))

					reached, err := module.NewUsage().IsAPICallsReached(usageRepo, subscriptions, workspace.Id, 1)
					if err != nil {
						util.CaptureError(err)
						log.Warn().
							Err(err).
							Str("workspaceId", workspace.Id.String()).
							Str("path", r.URL.Path).
							Msg("Failed to check workspace API call capacity")
					} else if reached {
						util.WriteJSON(w, http.StatusForbidden, map[string]any{
							"errorMessage": locale.TR(r, "api_calls_capacity_reached"),
						})
						return
					}
				}
			}

			recorder := &StatusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(recorder, r)

			if recorder.status < http.StatusOK || recorder.status >= http.StatusMultipleChoices {
				return
			}

			if !IsMeteredAPICall(r) {
				return
			}

			if _, _, ok := GetPrincipal(r); !ok {
				return
			}

			if !hasWorkspace {
				return
			}

			err := module.NewUsage().IncrementAPICalls(
				db.NewUsageRepository(db.GetDB(false)),
				db.NewSubscriptionRepository(db.GetDB(false)),
				workspace.Id,
			)
			if err != nil {
				util.CaptureError(err)
				log.Warn().
					Err(err).
					Str("workspaceId", workspace.Id.String()).
					Str("path", r.URL.Path).
					Msg("Failed to record workspace API call usage")
			}
		})
	}
}

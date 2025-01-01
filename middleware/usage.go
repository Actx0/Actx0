// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package middleware

import (
	"net/http"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/module"

	"github.com/rs/zerolog/log"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rw *statusRecorder) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// TrackWorkspaceAPICall increments workspace API usage after authenticated, authorized
// requests complete successfully. Place on workspace routes after RequireWorkspace.
func TrackWorkspaceAPICall() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(recorder, r)

			if recorder.status < http.StatusOK || recorder.status >= http.StatusMultipleChoices {
				return
			}

			if _, _, ok := GetPrincipal(r); !ok {
				return
			}

			workspace, ok := GetWorkspaceFromContext(r.Context())
			if !ok {
				return
			}

			err := module.NewUsage().IncrementAPICalls(
				db.NewUsageRepository(db.GetDB(false)),
				workspace.Id,
			)
			if err != nil {
				log.Warn().
					Err(err).
					Str("workspaceId", workspace.Id.String()).
					Str("path", r.URL.Path).
					Msg("Failed to record workspace API call usage")
			}
		})
	}
}

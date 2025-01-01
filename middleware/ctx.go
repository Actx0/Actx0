// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/clivern/actx0/module"
)

// AppContext attaches request metadata as module.AppCtx on the request context.
func AppContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ac := module.NewAppCtx(
			GetRequestId(r.Context()),
			r.UserAgent(),
			clientIP(r),
			parseFeatures(r.Header.Get("X-Features")),
		)

		ctx := module.WithAppCtx(r.Context(), ac)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetAppCtx retrieves AppCtx from the request context.
func GetAppCtx(ctx context.Context) (*module.AppCtx, bool) {
	return module.AppCtxFrom(ctx)
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := strings.Index(xff, ","); i >= 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}

func parseFeatures(header string) map[string]bool {
	features := map[string]bool{}
	if header == "" {
		return features
	}

	for part := range strings.SplitSeq(header, ",") {
		name := strings.TrimSpace(part)
		if name != "" {
			features[name] = true
		}
	}

	return features
}

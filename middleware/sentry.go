// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package middleware

import (
	"net/http"

	sentryhttp "github.com/getsentry/sentry-go/http"
)

// Sentry recovers panics and reports them to Sentry with request context.
func Sentry(next http.Handler) http.Handler {
	return sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	}).Handle(next)
}

// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package core

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/viper"
)

// SetupSentry initializes the Sentry SDK from app.sentry.dsn.
func SetupSentry(release string) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              viper.GetString("app.sentry.dsn"),
		Environment:      viper.GetString("app.mode"),
		Release:          release,
		AttachStacktrace: true,
		DisableLogs:      true,
	})
	if err != nil {
		return fmt.Errorf("sentry.Init: %w", err)
	}

	return nil
}

// FlushSentry sends buffered events before shutdown.
func FlushSentry() {
	sentry.Flush(2 * time.Second)
}

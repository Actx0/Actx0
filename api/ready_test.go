// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/clivern/actx0/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationReadyEndpoint(t *testing.T) {
	t.Run("ReadyAction should return OK when database is available", func(t *testing.T) {
		db.CloseDB()

		err := db.InitDB(testPostgresConfig())
		if err != nil {
			t.Skipf("PostgreSQL test database is unavailable: %v", err)
		}
		defer db.CloseDB()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/_ready", nil)
		w := httptest.NewRecorder()

		ReadyAction(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"status":"ok"}`, strings.TrimSpace(w.Body.String()))
	})

	t.Run("ReadyAction should return ServiceUnavailable when database ping fails", func(t *testing.T) {
		db.CloseDB()

		err := db.InitDB(testPostgresConfig())
		if err != nil {
			t.Skipf("PostgreSQL test database is unavailable: %v", err)
		}

		database := db.GetDB()
		require.NotNil(t, database, "Database connection should be initialized")
		database.Close()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/_ready", nil)
		w := httptest.NewRecorder()

		ReadyAction(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		assert.Equal(t, `{"status":"not_ok"}`, strings.TrimSpace(w.Body.String()))

		db.CloseDB()
	})
}

func testPostgresConfig() db.Config {
	port, err := strconv.Atoi(envOrDefault("FINN_DATABASE_PORT", "5432"))
	if err != nil {
		port = 5432
	}

	return db.Config{
		Driver:          envOrDefault("FINN_DATABASE_DRIVER", "postgres"),
		Host:            envOrDefault("FINN_DATABASE_HOST", "localhost"),
		Port:            port,
		Username:        envOrDefault("FINN_DATABASE_USERNAME", "postgres"),
		Password:        envOrDefault("FINN_DATABASE_PASSWORD", "postgres"),
		Database:        envOrDefault("FINN_DATABASE_NAME", "finn"),
		MaxOpenConns:    5,
		MaxIdleConns:    2,
		ConnMaxLifetime: 300,
	}
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

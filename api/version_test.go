// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/clivern/actx0/build"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationVersionEndpoint(t *testing.T) {
	t.Run("VersionAction should return build info", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/public/_version", nil)
		w := httptest.NewRecorder()

		VersionAction(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var body map[string]any
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		assert.Equal(t, build.Version, body["version"])
		assert.Equal(t, build.Commit, body["commit"])
		assert.Equal(t, build.Date, body["date"])
		assert.Equal(t, build.BuiltBy, body["builtBy"])
	})
}

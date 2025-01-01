// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/clivern/actx0/pkg/util"
)

// GetMeAction returns the authenticated principal and its permissions.
func GetMeAction(w http.ResponseWriter, r *http.Request) {
	util.WriteJSON(w, http.StatusNotImplemented, map[string]interface{}{
		"errorMessage": "not implemented",
	})
}

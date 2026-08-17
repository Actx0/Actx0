// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/clivern/actx0/pkg/buildinfo"
	"github.com/clivern/actx0/pkg/util"

	"github.com/rs/zerolog/log"
)

// VersionAction returns the binary build metadata.
func VersionAction(w http.ResponseWriter, _ *http.Request) {
	log.Debug().Msg("Version check")

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"version": buildinfo.Version,
		"commit":  buildinfo.Commit,
		"date":    buildinfo.Date,
		"builtBy": buildinfo.BuiltBy,
	})
}

// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/locale"
	"github.com/clivern/actx0/middleware"
	"github.com/clivern/actx0/module"
	"github.com/clivern/actx0/pkg/util"

	"github.com/rs/zerolog/log"
)

// UpdateSettingsRequest is the body for updating app settings.
type UpdateSettingsRequest struct {
	PlatformEmail   string `json:"platformEmail" validate:"required,email,max=255" label:"Platform Email"`
	MaintenanceMode bool   `json:"maintenanceMode" label:"Maintenance Mode"`
}

// UpdateSettingsAction updates app settings (admin only).
func UpdateSettingsAction(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.GetUserFromContext(r.Context())

	log.Info().
		Str("userId", user.Id.String()).
		Msg("Updating settings")

	var req UpdateSettingsRequest
	err := util.DecodeAndValidate(r, &req)
	if err != nil {
		util.WriteValidationError(w, err)
		return
	}

	settingsModule := module.NewSettings(db.NewConfigRepository(db.GetDB()))
	err = settingsModule.Update(req.PlatformEmail, req.MaintenanceMode)
	if err != nil {
		log.Error().
			Err(err).
			Str("userId", user.Id.String()).
			Msg("Failed to update settings")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"errorMessage": locale.TR(r, "failed_update_settings"),
		})
		return
	}

	log.Info().
		Str("userId", user.Id.String()).
		Msg("Settings updated")

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"successMessage": locale.TR(r, "settings_updated_successfully"),
	})
}

// GetSettingsAction returns current app settings.
func GetSettingsAction(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.GetUserFromContext(r.Context())

	log.Info().
		Str("userId", user.Id.String()).
		Msg("Getting settings")

	settingsModule := module.NewSettings(db.NewConfigRepository(db.GetDB()))
	settings, err := settingsModule.GetSettings()

	if err != nil {
		log.Error().
			Err(err).
			Str("userId", user.Id.String()).
			Msg("Failed to get settings")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"errorMessage": locale.TR(r, "failed_get_settings"),
		})
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"settings": settings,
	})
}

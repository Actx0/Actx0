// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/locale"
	"github.com/clivern/actx0/module"
	"github.com/clivern/actx0/pkg/util"

	"github.com/rs/zerolog/log"
)

// SetupRequest is the body you send when running setup.
type SetupRequest struct {
	PlatformURL   string `json:"platformURL" validate:"required,url,min=4,max=60" label:"Platform URL"`
	PlatformEmail string `json:"platformEmail" validate:"required,email,min=4,max=60" label:"Platform Email"`
	PlatformName  string `json:"platformName" validate:"required,min=2,max=50" label:"Platform Name"`
	AdminEmail    string `json:"adminEmail" validate:"required,email,min=4,max=60" label:"Admin Email"`
	AdminPassword string `json:"adminPassword" validate:"required,strong_password,min=8,max=60" label:"Admin Password"`
}

// SetupAction runs the initial platform setup.
func SetupAction(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Setup endpoint called")

	var req SetupRequest

	if err := util.DecodeAndValidate(r, &req); err != nil {
		util.WriteValidationError(w, err)
		return
	}

	setupModule := module.NewSetup(
		db.NewConfigRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
		db.NewSystemEventRepository(db.GetDB()),
	)

	if setupModule.IsInstalled() {
		util.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": locale.TR(r, "platformAlreadyInstalled"),
		})
		return
	}

	err := setupModule.Install(&module.SetupOptions{
		PlatformURL:   req.PlatformURL,
		PlatformEmail: req.PlatformEmail,
		PlatformName:  req.PlatformName,
		AdminEmail:    req.AdminEmail,
		AdminPassword: req.AdminPassword,
	})

	if err != nil {
		log.Error().
			Err(err).
			Str("adminEmail", req.AdminEmail).
			Str("platformURL", req.PlatformURL).
			Msg("Failed to complete setup")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"errorMessage": locale.TR(r, "failedCompleteSetup"),
		})
		return
	}

	log.Info().
		Str("adminEmail", req.AdminEmail).
		Str("platformURL", req.PlatformURL).
		Msg("Platform setup completed successfully")
	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"successMessage": locale.TR(r, "setupCompletedSuccessfully"),
	})
}

// SetupStatusAction returns whether the platform is already installed.
func SetupStatusAction(w http.ResponseWriter, _ *http.Request) {
	log.Info().Msg("Setup status endpoint called")

	setupModule := module.NewSetup(
		db.NewConfigRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
		db.NewSystemEventRepository(db.GetDB()),
	)

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"installed": setupModule.IsInstalled(),
	})
}

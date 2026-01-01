// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"errors"
	"time"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/pkg/util"
)

// Setup does the initial platform install and config.
type Setup struct {
	ConfigRepository      db.ConfigRepository
	UserRepository        db.UserRepository
	SystemEventRepository db.SystemEventRepository
}

// SetupOptions is what you pass when running setup.
type SetupOptions struct {
	PlatformURL   string
	PlatformEmail string
	PlatformName  string
	AdminEmail    string
	AdminPassword string
}

// NewSetup creates a Setup with the given repositories.
func NewSetup(configRepository db.ConfigRepository, userRepository db.UserRepository, systemEventRepository db.SystemEventRepository) *Setup {
	return &Setup{
		ConfigRepository:      configRepository,
		UserRepository:        userRepository,
		SystemEventRepository: systemEventRepository,
	}
}

// IsInstalled checks whether the platform has been installed.
func (s *Setup) IsInstalled() bool {
	config, err := s.ConfigRepository.Get("is_installed")
	if err != nil {
		return false
	}
	return config != nil
}

// Install performs the initial platform installation with the provided options.
func (s *Setup) Install(options *SetupOptions) error {
	if s.IsInstalled() {
		return errors.New("platform is already installed")
	}

	hashedPassword, err := util.HashPassword(options.AdminPassword)
	if err != nil {
		return err
	}
	user := &db.User{
		Name:        "",
		Email:       options.AdminEmail,
		Password:    hashedPassword,
		Role:        db.UserRoleAdmin,
		IsActive:    true,
		LastLoginAt: time.Now().UTC(),
	}
	err = s.UserRepository.Create(user)
	if err != nil {
		return err
	}

	err = s.ConfigRepository.Create("is_installed", "1")
	if err != nil {
		return err
	}

	err = s.ConfigRepository.Create("platform_url", options.PlatformURL)
	if err != nil {
		return err
	}

	err = s.ConfigRepository.Create("platform_email", options.PlatformEmail)
	if err != nil {
		return err
	}

	err = s.ConfigRepository.Create("platform_name", options.PlatformName)
	if err != nil {
		return err
	}

	err = s.ConfigRepository.Create(
		"platform_description",
		"A Platform to Run, Schedule, and Automate Ansible and Terraform Workflows.",
	)
	if err != nil {
		return err
	}

	err = s.ConfigRepository.Create("maintenance_mode", "0")
	if err != nil {
		return err
	}

	err = s.ConfigRepository.Create("smtp_server", "smtp.gmail.com")
	if err != nil {
		return err
	}

	err = s.ConfigRepository.Create("smtp_port", "587")
	if err != nil {
		return err
	}

	err = s.ConfigRepository.Create("smtp_from_email", "noreply@ctx.io")
	if err != nil {
		return err
	}

	err = s.ConfigRepository.Create("smtp_username", "")
	if err != nil {
		return err
	}

	err = s.ConfigRepository.Create("smtp_password", "")
	if err != nil {
		return err
	}

	err = s.ConfigRepository.Create("smtp_use_tls", "1")
	if err != nil {
		return err
	}

	s.SystemEventRepository.Create("platform_installed", &db.SystemEventPlatformInstalled{
		UserID:    user.ID.String(),
		UserEmail: user.Email,
	})

	return nil
}

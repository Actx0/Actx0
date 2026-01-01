// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/clivern/actx0/db"
)

// Settings manages app settings (platform URL, SMTP, etc.).
type Settings struct {
	ConfigRepository db.ConfigRepository
}

// NewSettings creates a settings module with the given config repo.
func NewSettings(configRepository db.ConfigRepository) *Settings {
	return &Settings{ConfigRepository: configRepository}
}

// UpdateSettingsOptions is what you pass when updating settings.
type UpdateSettingsOptions struct {
	PlatformName    string
	PlatformURL     string
	PlatformEmail   string
	SMTPServer      string
	SMTPPort        int
	SMTPFromEmail   string
	SMTPUsername    string
	SMTPPassword    string
	SMTPUseTLS      bool
	MaintenanceMode bool
}

// Update persists the provided settings to the config repository.
func (s *Settings) Update(options *UpdateSettingsOptions) error {
	updates := map[string]string{
		"platform_name":    strings.TrimSpace(options.PlatformName),
		"platform_url":     strings.TrimSpace(options.PlatformURL),
		"platform_email":   strings.TrimSpace(options.PlatformEmail),
		"smtp_server":      strings.TrimSpace(options.SMTPServer),
		"smtp_port":        strconv.Itoa(options.SMTPPort),
		"smtp_from_email":  strings.TrimSpace(options.SMTPFromEmail),
		"smtp_username":    strings.TrimSpace(options.SMTPUsername),
		"smtp_password":    options.SMTPPassword,
		"smtp_use_tls":     boolToConfig(options.SMTPUseTLS),
		"maintenance_mode": boolToConfig(options.MaintenanceMode),
	}

	for key, value := range updates {
		if err := s.ConfigRepository.Update(key, value); err != nil {
			return fmt.Errorf("failed to update config %s: %w", key, err)
		}
	}

	return nil
}

func boolToConfig(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

// GetSettings retrieves the current settings from the config repository.
func (s *Settings) GetSettings() (map[string]string, error) {
	keys := []string{
		"platform_name",
		"platform_url",
		"platform_email",
		"maintenance_mode",
		"smtp_server",
		"smtp_port",
		"smtp_from_email",
		"smtp_username",
		"smtp_password",
		"smtp_use_tls",
	}

	settings := make(map[string]string, len(keys))

	for _, key := range keys {
		config, err := s.ConfigRepository.Get(key)
		if err != nil {
			return nil, err
		}
		if config != nil {
			settings[key] = config.Value
		} else {
			settings[key] = ""
		}
	}

	return settings, nil
}

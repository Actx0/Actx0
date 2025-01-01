// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"fmt"
	"strings"

	"github.com/clivern/actx0/db"
)

// Settings manages app settings.
type Settings struct {
	ConfigRepository db.ConfigRepository
}

// NewSettings creates a settings module with the given config repo.
func NewSettings(cr db.ConfigRepository) *Settings {
	return &Settings{ConfigRepository: cr}
}

// Update persists the provided settings to the config repository.
func (s *Settings) Update(platformEmail string, maintenanceMode bool) error {
	updates := map[string]string{
		"platform_email":   strings.TrimSpace(platformEmail),
		"maintenance_mode": boolToConfig(maintenanceMode),
	}

	for key, value := range updates {
		err := s.ConfigRepository.Update(key, value)
		if err != nil {
			return fmt.Errorf("failed to update config %s: %w", key, err)
		}
	}

	return nil
}

// boolToConfig converts a bool to a config string value.
func boolToConfig(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

// GetSettings retrieves the current settings from the config repository.
func (s *Settings) GetSettings() (map[string]string, error) {
	keys := []string{
		"platform_email",
		"maintenance_mode",
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

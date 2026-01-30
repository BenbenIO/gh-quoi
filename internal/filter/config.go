// Package filter provides notification filtering based on configurable rules.
package filter

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the top-level configuration for notification filtering.
// It specifies the time range for fetching notifications and the filter rules to apply.
type Config struct {
	DaysAgo int    `yaml:"since_days_ago"` // Number of days to look back for notifications
	Filters []Rule `yaml:"filters"`        // List of filter rules to apply
}

// LoadDefault loads the configuration from the default location (~/.config/gh-quoi/config.yaml).
// Returns an error if the file cannot be read or parsed.
func LoadDefault() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home: %w", err)
	}

	path := filepath.Join(home, ".config", "gh-quoi", "config.yaml")

	return LoadFile(path)
}

// LoadFile loads and parses the configuration from the specified file path.
// Returns an error if the file cannot be read or contains invalid YAML.
func LoadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing filter config: %w", err)
	}

	return &cfg, nil
}

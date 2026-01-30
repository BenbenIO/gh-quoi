package filter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFile(t *testing.T) {
	t.Run("valid config file", func(t *testing.T) {
		// Create temporary config file
		// TODO: look for better golang testing fixtures style?
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		configContent := `
since_days_ago: 7
filters:
  - name: "urgent issues"
    reasons: ["assigned", "mention"]
    types: ["Issue"]
    title_regex: "^\\[URGENT\\]"
  - name: "exclude archived"
    repos:
      not: ["archived/*"]
`

		err := os.WriteFile(configPath, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("failed to write test config: %v", err)
		}

		// Load and test
		cfg, err := LoadFile(configPath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.DaysAgo != 7 {
			t.Errorf("expected DaysAgo 7, got %d", cfg.DaysAgo)
		}

		if len(cfg.Filters) != 2 {
			t.Errorf("expected 2 filters, got %d", len(cfg.Filters))
		}

		// Test first filter
		filter1 := cfg.Filters[0]
		if filter1.Name != "urgent issues" {
			t.Errorf("expected name 'urgent issues', got '%s'", filter1.Name)
		}

		if len(filter1.Reasons.Any) != 2 {
			t.Errorf("expected 2 reasons, got %d", len(filter1.Reasons.Any))
		}

		if filter1.TitleRegex != "^\\[URGENT\\]" {
			t.Errorf("expected title regex '^\\[URGENT\\]', got '%s'", filter1.TitleRegex)
		}

		// Test second filter
		filter2 := cfg.Filters[1]
		if filter2.Name != "exclude archived" {
			t.Errorf("expected name 'exclude archived', got '%s'", filter2.Name)
		}

		if len(filter2.Repos.Not) != 1 {
			t.Errorf("expected 1 repo exclusion, got %d", len(filter2.Repos.Not))
		}

		if filter2.Repos.Not[0] != "archived/*" {
			t.Errorf("expected repo exclusion 'archived/*', got '%s'", filter2.Repos.Not[0])
		}
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := LoadFile("/nonexistent/config.yaml")
		if err == nil {
			t.Error("expected error for nonexistent file")
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "invalid.yaml")

		invalidContent := `
since_days_ago: 7
filters:
  - name: "test"
    invalid_yaml: [
`

		err := os.WriteFile(configPath, []byte(invalidContent), 0644)
		if err != nil {
			t.Fatalf("failed to write test config: %v", err)
		}

		_, err = LoadFile(configPath)
		if err == nil {
			t.Error("expected error for invalid YAML")
		}
	})

	t.Run("empty config", func(t *testing.T) {
		// Assuming an empty file is valid -- we will get all the notifications without filtering
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "empty.yaml")

		err := os.WriteFile(configPath, []byte(""), 0644)
		if err != nil {
			t.Fatalf("failed to write test config: %v", err)
		}

		cfg, err := LoadFile(configPath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.DaysAgo != 0 {
			t.Errorf("expected DaysAgo 0, got %d", cfg.DaysAgo)
		}

		if len(cfg.Filters) != 0 {
			t.Errorf("expected 0 filters, got %d", len(cfg.Filters))
		}
	})
}

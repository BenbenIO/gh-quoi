package github

import (
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Run("with valid token", func(t *testing.T) {
		// Set test token
		originalToken := os.Getenv("GITHUB_TOKEN")

		defer func() {
			if originalToken == "" {
				os.Unsetenv("GITHUB_TOKEN")
			} else {
				os.Setenv("GITHUB_TOKEN", originalToken)
			}
		}()

		os.Setenv("GITHUB_TOKEN", "test-token")

		client, err := NewClient()
		if err != nil {
			t.Fatalf("NewClient() error = %v, want nil", err)
		}

		if client == nil {
			t.Error("NewClient() = nil, want non-nil client")
		}

		if client.gh == nil {
			t.Error("client.gh = nil, want non-nil GitHub client")
		}
	})

	t.Run("without token", func(t *testing.T) {
		// Unset token
		originalToken := os.Getenv("GITHUB_TOKEN")

		defer func() {
			if originalToken == "" {
				os.Unsetenv("GITHUB_TOKEN")
			} else {
				os.Setenv("GITHUB_TOKEN", originalToken)
			}
		}()

		os.Unsetenv("GITHUB_TOKEN")

		_, err := NewClient()
		if err == nil {
			t.Error("NewClient() error = nil, want error for missing token")

			return
		}

		wantErr := "GITHUB_TOKEN environment variable not set"
		if err.Error() != wantErr {
			t.Errorf("NewClient() error = %q, want %q", err.Error(), wantErr)
		}
	})

	t.Run("with empty token", func(t *testing.T) {
		// Set empty token
		originalToken := os.Getenv("GITHUB_TOKEN")

		defer func() {
			if originalToken == "" {
				os.Unsetenv("GITHUB_TOKEN")
			} else {
				os.Setenv("GITHUB_TOKEN", originalToken)
			}
		}()

		os.Setenv("GITHUB_TOKEN", "")

		_, err := NewClient()
		if err == nil {
			t.Error("NewClient() error = nil, want error for empty token")

			return
		}

		wantErr := "GITHUB_TOKEN environment variable not set"
		if err.Error() != wantErr {
			t.Errorf("NewClient() error = %q, want %q", err.Error(), wantErr)
		}
	})
}

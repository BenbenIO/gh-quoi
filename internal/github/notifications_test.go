package github

import (
	"reflect"
	"testing"
	"time"

	"github.com/BenbenIO/gh-quoi/internal/model"
	"github.com/google/go-github/v55/github"
)

func TestDaysAgo(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name string
		days int
		want time.Duration
	}{
		{"0 days", 0, 0},
		{"1 day", 1, 24 * time.Hour},
		{"7 days", 7, 7 * 24 * time.Hour},
		{"30 days", 30, 30 * 24 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DaysAgo(tt.days)
			if result == nil {
				t.Error("DaysAgo() = nil, want non-nil time")

				return
			}

			diff := now.Sub(*result)
			// Allow for small timing differences (within 1 second)
			if diff < tt.want-time.Second || diff > tt.want+time.Second {
				t.Errorf("DaysAgo(%d) = %v ago, want %v ago", tt.days, diff, tt.want)
			}
		})
	}
}

func TestConvertNotification(t *testing.T) {
	tests := []struct {
		name  string
		input *github.Notification
		want  model.Notification
	}{
		{
			name: "complete notification",
			input: &github.Notification{
				ID: github.String("test-id"),
				Repository: &github.Repository{
					FullName: github.String("owner/repo"),
				},
				Subject: &github.NotificationSubject{
					Title: github.String("Test Issue"),
					Type:  github.String("Issue"),
					URL:   github.String("https://api.github.com/repos/owner/repo/issues/123"),
				},
				Reason: github.String("assigned"),
				UpdatedAt: &github.Timestamp{
					Time: time.Date(2023, 1, 15, 14, 30, 0, 0, time.UTC),
				},
				Unread: github.Bool(true),
			},
			want: model.Notification{
				ID:         "test-id",
				Repository: "owner/repo",
				Title:      "Test Issue",
				Reason:     "assigned",
				Type:       "Issue",
				HTMLURL:    "https://github.com/owner/repo/issues/123",
				UpdatedAt:  "2023-01-15 14h",
				Unread:     true,
			},
		},
		{
			name:  "minimal notification with nils",
			input: &github.Notification{},
			want: model.Notification{
				ID:         "",
				Repository: "",
				Title:      "",
				Reason:     "",
				Type:       "",
				HTMLURL:    "",
				UpdatedAt:  "",
				Unread:     false,
			},
		},
		{
			name: "notification with pull request",
			input: &github.Notification{
				ID: github.String("pr-id"),
				Repository: &github.Repository{
					FullName: github.String("owner/repo"),
				},
				Subject: &github.NotificationSubject{
					Title: github.String("Test PR"),
					Type:  github.String("PullRequest"),
					URL:   github.String("https://api.github.com/repos/owner/repo/pulls/456"),
				},
				Reason: github.String("mention"),
				UpdatedAt: &github.Timestamp{
					Time: time.Date(2023, 12, 1, 9, 15, 0, 0, time.UTC),
				},
				Unread: github.Bool(false),
			},
			want: model.Notification{
				ID:         "pr-id",
				Repository: "owner/repo",
				Title:      "Test PR",
				Reason:     "mention",
				Type:       "PullRequest",
				HTMLURL:    "https://github.com/owner/repo/pull/456",
				UpdatedAt:  "2023-12-01 09h",
				Unread:     false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertNotification(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertNotification() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestApiToHTMLURL(t *testing.T) {
	tests := []struct {
		name   string
		apiURL string
		want   string
	}{
		{
			name:   "issue URL",
			apiURL: "https://api.github.com/repos/owner/repo/issues/123",
			want:   "https://github.com/owner/repo/issues/123",
		},
		{
			name:   "pull request URL",
			apiURL: "https://api.github.com/repos/owner/repo/pulls/456",
			want:   "https://github.com/owner/repo/pull/456",
		},
		{
			name:   "empty URL",
			apiURL: "",
			want:   "",
		},
		{
			name:   "non-API URL",
			apiURL: "https://github.com/owner/repo",
			want:   "https://github.com/owner/repo",
		},
		{
			name:   "complex repository name",
			apiURL: "https://api.github.com/repos/org-name/repo-name.example/issues/789",
			want:   "https://github.com/org-name/repo-name.example/issues/789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := apiToHTMLURL(tt.apiURL)
			if got != tt.want {
				t.Errorf("apiToHTMLURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestListOptions(t *testing.T) {
	t.Run("struct initialization", func(t *testing.T) {
		since := time.Now().AddDate(0, 0, -7)
		opts := ListOptions{
			All:           true,
			Participating: false,
			Since:         &since,
		}

		if !opts.All {
			t.Errorf("All = %t, want true", opts.All)
		}

		if opts.Participating {
			t.Errorf("Participating = %t, want false", opts.Participating)
		}

		if opts.Since == nil {
			t.Error("Since = nil, want non-nil time")

			return
		}

		if !opts.Since.Equal(since) {
			t.Errorf("Since = %v, want %v", opts.Since, since)
		}
	})

	t.Run("zero value", func(t *testing.T) {
		var opts ListOptions

		if opts.All {
			t.Errorf("All = %t, want false", opts.All)
		}

		if opts.Participating {
			t.Errorf("Participating = %t, want false", opts.Participating)
		}

		if opts.Since != nil {
			t.Errorf("Since = %v, want nil", opts.Since)
		}
	})
}

// TestConvertNotification_EdgeCases tests edge cases and error conditions.
func TestConvertNotification_EdgeCases(t *testing.T) {
	t.Run("nil repository", func(t *testing.T) {
		input := &github.Notification{
			Repository: nil,
		}

		got := convertNotification(input)
		if got.Repository != "" {
			t.Errorf("Repository = %q, want empty string", got.Repository)
		}
	})

	t.Run("repository with nil FullName", func(t *testing.T) {
		input := &github.Notification{
			Repository: &github.Repository{
				FullName: nil,
			},
		}

		got := convertNotification(input)
		if got.Repository != "" {
			t.Errorf("Repository = %q, want empty string", got.Repository)
		}
	})

	t.Run("nil subject", func(t *testing.T) {
		input := &github.Notification{
			Subject: nil,
		}

		got := convertNotification(input)
		if got.Title != "" {
			t.Errorf("Title = %q, want empty string", got.Title)
		}

		if got.Type != "" {
			t.Errorf("Type = %q, want empty string", got.Type)
		}

		if got.HTMLURL != "" {
			t.Errorf("HTMLURL = %q, want empty string", got.HTMLURL)
		}
	})

	t.Run("subject with nil fields", func(t *testing.T) {
		input := &github.Notification{
			Subject: &github.NotificationSubject{
				Title: nil,
				Type:  nil,
				URL:   nil,
			},
		}

		got := convertNotification(input)
		if got.Title != "" {
			t.Errorf("Title = %q, want empty string", got.Title)
		}

		if got.Type != "" {
			t.Errorf("Type = %q, want empty string", got.Type)
		}

		if got.HTMLURL != "" {
			t.Errorf("HTMLURL = %q, want empty string", got.HTMLURL)
		}
	})
}

// Benchmark tests for performance.
func BenchmarkConvertNotification(b *testing.B) {
	notification := &github.Notification{
		ID: github.String("bench-id"),
		Repository: &github.Repository{
			FullName: github.String("owner/repo"),
		},
		Subject: &github.NotificationSubject{
			Title: github.String("Benchmark Issue"),
			Type:  github.String("Issue"),
			URL:   github.String("https://api.github.com/repos/owner/repo/issues/999"),
		},
		Reason: github.String("assigned"),
		UpdatedAt: &github.Timestamp{
			Time: time.Now(),
		},
		Unread: github.Bool(true),
	}

	for b.Loop() {
		_ = convertNotification(notification)
	}
}

func BenchmarkApiToHTMLURL(b *testing.B) {
	apiURL := "https://api.github.com/repos/owner/repo/issues/123"

	for b.Loop() {
		_ = apiToHTMLURL(apiURL)
	}
}

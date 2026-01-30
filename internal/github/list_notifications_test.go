package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/v55/github"
)

func TestListNotifications_Unit(t *testing.T) {
	tests := []struct {
		name           string
		opts           ListOptions
		mockResponse   string
		mockStatusCode int
		wantErr        bool
		wantCount      int
	}{
		{
			name: "successful request with notifications",
			opts: ListOptions{
				All:           false,
				Participating: false,
				Since:         DaysAgo(1),
			},
			mockResponse: `[
				{
					"id": "1",
					"repository": {
						"full_name": "owner/repo1"
					},
					"subject": {
						"title": "Test Issue",
						"type": "Issue",
						"url": "https://api.github.com/repos/owner/repo1/issues/1"
					},
					"reason": "assigned",
					"updated_at": "2023-01-15T14:30:00Z",
					"unread": true
				},
				{
					"id": "2",
					"repository": {
						"full_name": "owner/repo2"
					},
					"subject": {
						"title": "Test PR",
						"type": "PullRequest",
						"url": "https://api.github.com/repos/owner/repo2/pulls/1"
					},
					"reason": "mention",
					"updated_at": "2023-01-14T10:15:00Z",
					"unread": false
				}
			]`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			wantCount:      2,
		},
		{
			name: "empty response",
			opts: ListOptions{
				All:           true,
				Participating: true,
			},
			mockResponse:   `[]`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			wantCount:      0,
		},
		{
			name: "API error response",
			opts: ListOptions{},
			mockResponse: `{
				"message": "Bad credentials",
				"documentation_url": "https://docs.github.com/rest"
			}`,
			mockStatusCode: http.StatusUnauthorized,
			wantErr:        true,
			wantCount:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request parameters
				if tt.opts.All {
					if r.URL.Query().Get("all") != "true" {
						t.Errorf("expected all=true in query params")
					}
				}

				if tt.opts.Participating {
					if r.URL.Query().Get("participating") != "true" {
						t.Errorf("expected participating=true in query params")
					}
				}

				if tt.opts.Since != nil {
					if r.URL.Query().Get("since") == "" {
						t.Errorf("expected since parameter in query")
					}
				}

				w.WriteHeader(tt.mockStatusCode)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client := &Client{
				gh: github.NewClient(nil),
			}
			baseURL, _ := url.Parse(server.URL + "/")
			client.gh.BaseURL = baseURL

			ctx := context.Background()
			notifications, err := client.ListNotifications(ctx, tt.opts)

			if tt.wantErr {
				if err == nil {
					t.Error("ListNotifications() error = nil, want error")
				}

				return
			}

			if err != nil {
				t.Fatalf("ListNotifications() error = %v, want nil", err)
			}

			if len(notifications) != tt.wantCount {
				t.Errorf("len(notifications) = %d, want %d", len(notifications), tt.wantCount)
			}

			if len(notifications) > 0 {
				n := notifications[0]
				if n.ID == "" {
					t.Error("notification.ID = empty, want non-empty")
				}

				if n.Repository == "" {
					t.Error("notification.Repository = empty, want non-empty")
				}
			}
		})
	}
}

func TestListNotifications_OptionsMapping(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		if query.Get("all") != "true" {
			t.Errorf("expected all=true, got %s", query.Get("all"))
		}

		if query.Get("participating") != "true" {
			t.Errorf("expected participating=true, got %s", query.Get("participating"))
		}

		if query.Get("since") == "" {
			t.Error("expected since parameter to be set")
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	}))
	defer server.Close()

	client := &Client{
		gh: github.NewClient(nil),
	}
	baseURL, _ := url.Parse(server.URL + "/")
	client.gh.BaseURL = baseURL

	ctx := context.Background()
	opts := ListOptions{
		All:           true,
		Participating: true,
		Since:         DaysAgo(1),
	}

	_, err := client.ListNotifications(ctx, opts)
	if err != nil {
		t.Fatalf("ListNotifications() error = %v", err)
	}
}

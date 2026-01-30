package filter

import (
	"testing"

	"github.com/BenbenIO/gh-quoi/internal/model"
)

func TestNewMatcher(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg := &Config{
			Filters: []Rule{
				{
					Name:       "test rule",
					Reasons:    ValueFilter{Any: []string{"assigned"}},
					TitleRegex: "^\\[BUG\\]",
				},
			},
		}

		matcher, err := NewMatcher(cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := len(matcher.rules); got != 1 {
			t.Errorf("len(rules) = %d, want 1", got)
		}

		rule := matcher.rules[0]
		if got := rule.Name; got != "test rule" {
			t.Errorf("rule.Name = %q, want %q", got, "test rule")
		}

		if rule.TitleRegex == nil {
			t.Error("rule.TitleRegex = nil, want compiled regex")
		}
	})

	t.Run("invalid regex", func(t *testing.T) {
		cfg := &Config{
			Filters: []Rule{
				{
					TitleRegex: "[invalid",
				},
			},
		}

		_, err := NewMatcher(cfg)
		if err == nil {
			t.Error("NewMatcher() error = nil, want error for invalid regex")

			return
		}
	})

	t.Run("empty config", func(t *testing.T) {
		cfg := &Config{Filters: []Rule{}}

		matcher, err := NewMatcher(cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := len(matcher.rules); got != 0 {
			t.Errorf("len(rules) = %d, want 0", got)
		}
	})
}

func TestMatcher_Match(t *testing.T) {
	cfg := &Config{
		Filters: []Rule{
			{
				Name:       "urgent issues",
				Reasons:    ValueFilter{Any: []string{"assigned", "mention"}},
				Types:      ValueFilter{Any: []string{"Issue"}},
				TitleRegex: "^\\[URGENT\\]",
			},
			{
				Name:  "non-archived repos",
				Repos: ValueFilter{Not: []string{"archived/*"}},
			},
		},
	}

	matcher, err := NewMatcher(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name         string
		notification model.Notification
		want         bool
	}{
		{
			name: "matches first rule",
			notification: model.Notification{
				Repository: "owner/repo",
				Title:      "[URGENT] Fix critical bug",
				Reason:     "assigned",
				Type:       "Issue",
			},
			want: true,
		},
		{
			name: "matches second rule",
			notification: model.Notification{
				Repository: "owner/repo",
				Title:      "Regular issue",
				Reason:     "author",
				Type:       "Issue",
			},
			want: true,
		},
		{
			name: "matches second rule - wrong title for first but matches second",
			notification: model.Notification{
				Repository: "owner/repo",
				Title:      "Regular issue",
				Reason:     "assigned",
				Type:       "Issue",
			},
			want: true,
		},
		{
			name: "no match - archived repo",
			notification: model.Notification{
				Repository: "archived/old-repo",
				Title:      "Some issue",
				Reason:     "author",
				Type:       "Issue",
			},
			want: false,
		},
		{
			name: "matches second rule - wrong reason/type for first but matches second",
			notification: model.Notification{
				Repository: "owner/repo",
				Title:      "[URGENT] Fix critical bug",
				Reason:     "review_requested",
				Type:       "PullRequest",
			},
			want: true,
		},
		{
			name: "no match - archived repo excludes from second rule",
			notification: model.Notification{
				Repository: "archived/old-repo",
				Title:      "[URGENT] Fix critical bug",
				Reason:     "review_requested",
				Type:       "PullRequest",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matcher.Match(tt.notification)
			if got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatcher_Match_EmptyMatcher(t *testing.T) {
	matcher := &Matcher{rules: []CompiledRule{}}

	notification := model.Notification{
		Repository: "owner/repo",
		Title:      "Test",
		Reason:     "assigned",
		Type:       "Issue",
	}

	result := matcher.Match(notification)
	if result {
		t.Error("expected false for empty matcher, got true")
	}
}

func TestMatchValue(t *testing.T) {
	tests := []struct {
		name   string
		filter ValueFilter
		value  string
		want   bool
	}{
		{
			name:   "empty filter matches everything",
			filter: ValueFilter{},
			value:  "anything",
			want:   true,
		},
		{
			name:   "exact match in Any list",
			filter: ValueFilter{Any: []string{"assigned", "mention"}},
			value:  "assigned",
			want:   true,
		},
		{
			name:   "glob match in Any list",
			filter: ValueFilter{Any: []string{"test/*"}},
			value:  "test/repo",
			want:   true,
		},
		{
			name:   "no match in Any list",
			filter: ValueFilter{Any: []string{"assigned", "mention"}},
			value:  "review_requested",
			want:   false,
		},
		{
			name:   "excluded by Not list",
			filter: ValueFilter{Not: []string{"archived/*"}},
			value:  "archived/old-repo",
			want:   false,
		},
		{
			name:   "not excluded by Not list",
			filter: ValueFilter{Not: []string{"archived/*"}},
			value:  "active/repo",
			want:   true,
		},
		{
			name:   "matches Any but excluded by Not",
			filter: ValueFilter{Any: []string{"test/*"}, Not: []string{"test/old"}},
			value:  "test/old",
			want:   false,
		},
		{
			name:   "matches Any and not excluded",
			filter: ValueFilter{Any: []string{"test/*"}, Not: []string{"test/old"}},
			value:  "test/new",
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchValue(tt.filter, tt.value)
			if got != tt.want {
				t.Errorf("matchValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchRule(t *testing.T) {
	rule := CompiledRule{
		Name:    "test rule",
		Reasons: ValueFilter{Any: []string{"assigned"}},
		Repos:   ValueFilter{Not: []string{"archived/*"}},
		Types:   ValueFilter{Any: []string{"Issue"}},
	}

	tests := []struct {
		name         string
		notification model.Notification
		want         bool
	}{
		{
			name: "all conditions match",
			notification: model.Notification{
				Repository: "owner/repo",
				Reason:     "assigned",
				Type:       "Issue",
			},
			want: true,
		},
		{
			name: "wrong reason",
			notification: model.Notification{
				Repository: "owner/repo",
				Reason:     "mention",
				Type:       "Issue",
			},
			want: false,
		},
		{
			name: "archived repo",
			notification: model.Notification{
				Repository: "archived/old",
				Reason:     "assigned",
				Type:       "Issue",
			},
			want: false,
		},
		{
			name: "wrong type",
			notification: model.Notification{
				Repository: "owner/repo",
				Reason:     "assigned",
				Type:       "PullRequest",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchRule(rule, tt.notification)
			if got != tt.want {
				t.Errorf("matchRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

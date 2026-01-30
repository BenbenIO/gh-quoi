package filter

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestValueFilter_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		want    ValueFilter
		wantErr bool
	}{
		{
			name: "simple list",
			yaml: `["assigned", "mention"]`,
			want: ValueFilter{
				Any: []string{"assigned", "mention"},
			},
			wantErr: false,
		},
		{
			name: "empty list",
			yaml: `[]`,
			want: ValueFilter{
				Any: []string{},
			},
			wantErr: false,
		},
		{
			name: "not only",
			yaml: `{not: ["old/*", "archived/*"]}`,
			want: ValueFilter{
				Not: []string{"old/*", "archived/*"},
			},
			wantErr: false,
		},
		{
			name: "any and not",
			yaml: `{any: ["test/*"], not: ["test/old"]}`,
			want: ValueFilter{
				Any: []string{"test/*"},
				Not: []string{"test/old"},
			},
			wantErr: false,
		},
		{
			name:    "empty object",
			yaml:    `{}`,
			want:    ValueFilter{},
			wantErr: false,
		},
		{
			name:    "invalid format - string",
			yaml:    `"invalid"`,
			wantErr: true,
		},
		{
			name:    "invalid format - number",
			yaml:    `42`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var vf ValueFilter

			err := yaml.Unmarshal([]byte(tt.yaml), &vf)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(vf, tt.want) {
				t.Errorf("got %+v, want %+v", vf, tt.want)
			}
		})
	}
}

func TestRule_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		want    Rule
		wantErr bool
	}{
		{
			name: "complete rule",
			yaml: `
name: "test rule"
reasons: ["assigned", "mention"]
repos: {not: ["archived/*"]}
types: ["Issue", "PullRequest"]
title_regex: "^\\[URGENT\\]"
`,
			want: Rule{
				Name:       "test rule",
				Reasons:    ValueFilter{Any: []string{"assigned", "mention"}},
				Repos:      ValueFilter{Not: []string{"archived/*"}},
				Types:      ValueFilter{Any: []string{"Issue", "PullRequest"}},
				TitleRegex: "^\\[URGENT\\]",
			},
			wantErr: false,
		},
		{
			name: "minimal rule",
			yaml: `
name: "minimal rule"
`,
			want: Rule{
				Name: "minimal rule",
			},
			wantErr: false,
		},
		{
			name: "title_regex as list should fail",
			yaml: `
name: "regex list rule"
title_regex: ["pattern1", "pattern2"]
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rule Rule

			err := yaml.Unmarshal([]byte(tt.yaml), &rule)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(rule, tt.want) {
				t.Errorf("got %+v, want %+v", rule, tt.want)
			}
		})
	}
}

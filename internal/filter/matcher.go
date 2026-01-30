// Package filter provides notification filtering based on configurable rules.
package filter

import (
	"path/filepath"
	"regexp"

	"github.com/BenbenIO/gh-quoi/internal/model"
)

// Matcher evaluates notifications against a set of compiled filter rules.
// It determines whether a notification should be displayed based on the configured criteria.
type Matcher struct {
	rules []CompiledRule
}

// CompiledRule is a filter rule with pre-compiled regular expressions for efficient matching.
// All conditions within a rule use AND logic; all must match for the rule to apply.
type CompiledRule struct {
	Name       string         // Human-readable name of the rule
	Reasons    ValueFilter    // Filter for notification reason (assigned, mentioned, etc.)
	Repos      ValueFilter    // Filter for repository name
	Types      ValueFilter    // Filter for notification type (Issue, PullRequest, etc.)
	TitleRegex *regexp.Regexp // Optional regex pattern for notification title
}

// NewMatcher creates a new Matcher from the provided configuration.
// It compiles all regular expressions and returns an error if any regex is invalid.
func NewMatcher(cfg *Config) (*Matcher, error) {
	var compiled []CompiledRule

	for _, r := range cfg.Filters {
		cr := CompiledRule{
			Name:    r.Name,
			Reasons: r.Reasons,
			Repos:   r.Repos,
			Types:   r.Types,
		}

		if r.TitleRegex != "" {
			re, err := regexp.Compile(r.TitleRegex)
			if err != nil {
				return nil, err
			}

			cr.TitleRegex = re
		}

		compiled = append(compiled, cr)
	}

	return &Matcher{rules: compiled}, nil
}

// Match returns true if the notification matches at least one configured rule.
// Rules use OR logic: a match on any rule means the notification should be displayed.
func (m *Matcher) Match(n model.Notification) bool {
	for _, rule := range m.rules {
		if matchRule(rule, n) {
			return true
		}
	}

	return false
}

// Rule-level matching (AND inside a rule).
func matchRule(r CompiledRule, n model.Notification) bool {
	if !matchValue(r.Reasons, n.Reason) {
		return false
	}

	if !matchValue(r.Repos, n.Repository) {
		return false
	}

	if !matchValue(r.Types, n.Type) {
		return false
	}

	if r.TitleRegex != nil && !r.TitleRegex.MatchString(n.Title) {
		return false
	}

	return true
}

// matchValue matches a single string ("reason", "repo", etc)
// against include and exclude patterns.
func matchValue(vf ValueFilter, value string) bool {
	// If ANY list is present, value must match at least one pattern.
	if len(vf.Any) > 0 {
		matched := false

		for _, pat := range vf.Any {
			if ok, _ := filepath.Match(pat, value); ok {
				matched = true

				break
			}
		}

		if !matched {
			return false
		}
	}

	// Check excludes
	for _, pat := range vf.Not {
		if ok, _ := filepath.Match(pat, value); ok {
			return false
		}
	}

	return true
}

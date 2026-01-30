// Package filter provides notification filtering based on configurable rules.
package filter

import (
	"errors"

	"gopkg.in/yaml.v3"
)

// Rule defines a single filter rule that can be applied to notifications.
// All fields are optional. When multiple fields are specified, all conditions
// must match (AND logic) for the rule to apply.
type Rule struct {
	Name       string      `yaml:"name"`        // Human-readable name for the rule
	Reasons    ValueFilter `yaml:"reasons"`     // Filter by notification reason
	Repos      ValueFilter `yaml:"repos"`       // Filter by repository name
	Types      ValueFilter `yaml:"types"`       // Filter by notification type
	TitleRegex string      `yaml:"title_regex"` // Regular expression to match title
}

// ValueFilter specifies which values to include or exclude using patterns.
// Supports both simple list format and explicit include/exclude format.
//
// Examples:
//   - Simple list: reasons: ["assigned", "mentioned"]
//   - With exclusions: repos: {any: ["org/*"], not: ["org/archived-*"]}
type ValueFilter struct {
	Any []string `yaml:"-"`   // Patterns to match (glob-style). Empty means match all.
	Not []string `yaml:"not"` // Patterns to exclude (glob-style)
}

// UnmarshalYAML implements custom YAML unmarshaling to support two formats:
//  1. Simple list: reasons: ["assigned", "mention"]
//  2. Object with any/not: repos: {any: ["org/*"], not: ["org/old-*"]}
func (vf *ValueFilter) UnmarshalYAML(value *yaml.Node) error {
	// Case 1: simple list
	if value.Kind == yaml.SequenceNode {
		var list []string
		err := value.Decode(&list)
		if err != nil {
			return err
		}

		vf.Any = list

		return nil
	}

	// Case 2: object => expect "not:"
	if value.Kind == yaml.MappingNode {
		var temp map[string][]string
		err := value.Decode(&temp)
		if err != nil {
			return err
		}

		if any, ok := temp["any"]; ok {
			vf.Any = any
		}

		if notList, ok := temp["not"]; ok {
			vf.Not = notList
		}

		return nil
	}

	return errors.New("invalid ValueFilter format")
}

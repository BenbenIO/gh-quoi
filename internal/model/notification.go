// Package model defines the core data structures used throughout the application.
package model

// Notification represents a simplified GitHub notification.
// It provides a stable interface independent of the GitHub API response structure.
// See: https://docs.github.com/en/rest/activity/notifications
type Notification struct {
	Unread     bool   // Whether the notification is unread
	ID         string // Unique identifier for the notification
	Repository string // Full repository name (owner/repo)
	Title      string // Notification title/subject
	Reason     string // Why the user received this (assigned, mentioned, etc.)
	Type       string // Notification type (Issue, PullRequest, etc.)
	HTMLURL    string // Browser URL to view the notification
	UpdatedAt  string // Last update timestamp in human-readable format
}

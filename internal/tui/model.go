// Package tui implements the terminal user interface for displaying and interacting with notifications.
// It uses the Bubble Tea framework for building the interactive terminal application.
package tui

import (
	"time"

	"github.com/BenbenIO/gh-quoi/internal/model"
	"github.com/charmbracelet/bubbles/table"
)

// RefreshFunc is a function that retrieves the latest notifications.
// It is called when the user requests a refresh of the notification list.
type RefreshFunc func() ([]model.Notification, error)

// Model represents the application state for the terminal UI.
// It implements the tea.Model interface from Bubble Tea.
type Model struct {
	table         table.Model          // Table widget displaying notifications
	refreshFn     RefreshFunc          // Function to refresh notifications
	loading       bool                 // Whether a refresh is in progress
	err           error                // Last error encountered, if any
	notifications []model.Notification // Current list of notifications
	showHelp      bool                 // Whether to display the help overlay
	lastRefreshed time.Time            // Timestamp of last successful refresh
}

// NewModel creates and initializes a new TUI model with the given notifications.
// The refreshFn is called when the user requests to reload the notification list.
func NewModel(initial []model.Notification, refreshFn RefreshFunc) Model {
	m := Model{
		refreshFn: refreshFn,
	}
	m.setRows(initial)
	m.lastRefreshed = time.Now()

	return m
}

// Table helpers.
func (m *Model) setRows(notifications []model.Notification) {
	m.notifications = notifications

	rows := make([]table.Row, len(notifications))
	for i, n := range notifications {
		// Update title if viewed.
		check := ""
		if !n.Unread {
			check = "✓"
		}

		rows[i] = table.Row{
			check,
			n.Title,
			n.Reason,
			n.Repository,
			n.UpdatedAt,
		}
	}

	columns := []table.Column{
		{Title: " ", Width: 1},
		{Title: "Title", Width: 45},
		{Title: "Reason", Width: 16},
		{Title: "Repo", Width: 25},
		{Title: "Updated", Width: 15},
	}

	m.table = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	m.table.SetStyles(defaultStyles())
}

// Package tui implements the terminal user interface for displaying and interacting with notifications.
// It uses the Bubble Tea framework for building the interactive terminal application.
package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// View renders the current state of the model as a string.
// Implements the tea.Model interface.
func (m Model) View() string {
	if m.loading {
		return "Loading notifications...\n"
	}

	if m.showHelp {
		return helpView()
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress r to retry or q to quit.\n", m.err)
	}

	if len(m.notifications) == 0 {
		frogStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5DA9E9")).
			Bold(true)

		frog := frogStyle.Render(`
        _    _
       (o)--(o)
      /.______.\
      \________/
     ./        \.
    ( .        , )
     \ \_\\//_/ /
      ~~  ~~  ~~
    `)

		msg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5DA9E9")).
			Bold(true).
			Render("Nothing to see here...\nPress r to refresh.\n")

		return frog + "\n" + msg
	}

	header := "GitHub Notifications (↑/↓, enter=open, r=refresh, q=quit)\n\n"
	tableView := m.table.View()

	statusText := fmt.Sprintf(
		"%d notifications  |  Refreshed: %s",
		len(m.notifications),
		m.lastRefreshed.Format("2006-01-02 15:04"),
	)

	status := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#5DA9E9")). // accent color
		Render(statusText)

	return header + tableView + "\n" + status
}

func helpView() string {
	box := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#777777"))

	help := `
Keyboard Shortcuts
------------------
↑ / ↓     Move selection (wraps)
enter     Open in browser
r         Refresh
?         Toggle help
q         Quit
`

	return box.Render(help)
}

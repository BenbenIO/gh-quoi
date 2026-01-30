// Package tui implements the terminal user interface for displaying and interacting with notifications.
// It uses the Bubble Tea framework for building the interactive terminal application.
package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func defaultStyles() table.Styles {
	s := table.DefaultStyles()
	accent := lipgloss.Color("#5DA9E9")   // unified accent color (modern blue)
	accentBg := lipgloss.Color("#1C3A50") // darker blue background for selected row
	divider := lipgloss.Color("#666666")  // subtle divider

	// Header: accent blue text, no background, divider under
	s.Header = s.Header.
		Foreground(accent).
		Background(lipgloss.NoColor{}).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(divider).
		Bold(true)

	// Selected row: matching accent but with a darker blue background
	s.Selected = s.Selected.
		Foreground(accent).
		Background(accentBg).
		Bold(true)

	// Normal cells: clean - not sure if someone have some transparent/image terminal?
	s.Cell = s.Cell.
		Background(lipgloss.NoColor{})

	return s
}

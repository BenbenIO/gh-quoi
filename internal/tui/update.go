// Package tui implements the terminal user interface for displaying and interacting with notifications.
// It uses the Bubble Tea framework for building the interactive terminal application.
package tui

import (
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"time"

	"github.com/BenbenIO/gh-quoi/internal/model"
	tea "github.com/charmbracelet/bubbletea"
)

type refreshMsg struct {
	Notifications []model.Notification
	Err           error
}

func refreshCmd(fn RefreshFunc) tea.Cmd {
	return func() tea.Msg {
		notifications, err := fn()

		return refreshMsg{Notifications: notifications, Err: err}
	}
}

// Create a message to open a URL in the browser instead
// of doing it directly in the command, to keep UI responsive.
type openURLMsg struct{}

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	// Suppress command output to not break the TUI
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard

	return cmd.Start()
}

func openURLCmd(url string) tea.Cmd {
	return func() tea.Msg {
		_ = openBrowser(url)
		return openURLMsg{}
	}
}

// Init initializes the model and returns any initial commands to run.
// Implements the tea.Model interface.
func (m Model) Init() tea.Cmd {
	return nil // no auto-refresh on start
}

// Update handles incoming messages and updates the model state accordingly.
// Implements the tea.Model interface.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case refreshMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err

			return m, nil
		}

		m.err = nil
		m.setRows(msg.Notifications)
		m.lastRefreshed = time.Now()

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "?":
			m.showHelp = !m.showHelp

			return m, nil

		case "up", "k":
			if len(m.notifications) > 0 {
				cur := m.table.Cursor()
				if cur == 0 {
					// wrap to bottom
					m.table.SetCursor(len(m.notifications) - 1)

					return m, nil
				}
			}

		case "down", "j":
			if len(m.notifications) > 0 {
				cur := m.table.Cursor()
				if cur == len(m.notifications)-1 {
					// wrap to top
					m.table.SetCursor(0)

					return m, nil
				}
			}
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			if len(m.notifications) > 0 {
				idx := m.table.Cursor()
				url := m.notifications[idx].HTMLURL
				m.notifications[idx].Unread = false
				m.setRows(m.notifications)
				// Restore cursor position after resetting rows
				if idx < len(m.notifications) {
					m.table.SetCursor(idx)
				}
				return m, openURLCmd(url)
			}

		case "r":
			m.loading = true
			m.err = nil

			return m, refreshCmd(m.refreshFn)
		}
	case openURLMsg:
		// Redraw after opening URL
		return m, nil
	}

	var cmd tea.Cmd

	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

// Package main implements gh-quoi, a terminal UI for managing GitHub notifications.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/BenbenIO/gh-quoi/internal/filter"
	"github.com/BenbenIO/gh-quoi/internal/github"
	"github.com/BenbenIO/gh-quoi/internal/model"
	"github.com/BenbenIO/gh-quoi/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

const (
	defaultDaysAgo = 20
)

var rootCmd = &cobra.Command{
	Use:   "gh-quoi",
	Short: "GitHub notification TUI manager",
	RunE: func(cmd *cobra.Command, args []string) error {
		// This runs the TUI when no subcommand is specified
		return runTUI()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// runTUI initializes and starts the TUI application.
func runTUI() error {
	ctx := context.Background()

	apiClient, err := github.NewClient()
	if err != nil {
		return err
	}

	cfg, err := filter.LoadDefault()
	if err != nil {
		return err
	}

	if cfg.DaysAgo == 0 {
		cfg.DaysAgo = defaultDaysAgo
	}

	matcher, err := filter.NewMatcher(cfg)
	if err != nil {
		return err
	}

	refresh := func() ([]model.Notification, error) {
		opts := github.ListOptions{
			All:           false,
			Participating: false,
			Since:         github.DaysAgo(cfg.DaysAgo),
		}

		notes, err := apiClient.ListNotifications(ctx, opts)
		if err != nil {
			return nil, err
		}

		filtered := make([]model.Notification, 0, len(notes))

		for _, n := range notes {
			if matcher.Match(n) {
				filtered = append(filtered, n)
			}
		}

		return filtered, nil
	}

	initial, err := refresh()
	if err != nil {
		return err
	}

	m := tui.NewModel(initial, refresh)
	_, err = tea.NewProgram(m).Run()

	return err
}

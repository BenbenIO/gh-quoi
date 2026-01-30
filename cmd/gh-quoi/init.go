package main

import (
	"fmt"
	"os"
	"path/filepath"

	_ "embed" // Required for go:embed directive

	"github.com/spf13/cobra"
)

//go:embed defaults/config.yaml
var defaultConfig []byte
var force bool

func init() {
	initCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing config")
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration for gh-quoi",
	Long:  "Creates $HOME/.config/gh-quoi/config.yaml with default settings.",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		configDir := filepath.Join(home, ".config", "gh-quoi")
		configPath := filepath.Join(configDir, "config.yaml")

		if err := os.MkdirAll(configDir, 0o755); err != nil {
			return fmt.Errorf("create config dir: %w", err)
		}

		if _, err := os.Stat(configPath); err == nil {
			if !force {
				return fmt.Errorf(
					"config already exists at %s (use --force to overwrite)",
					configPath,
				)
			}
		}

		if err := os.WriteFile(configPath, defaultConfig, 0o644); err != nil {
			return fmt.Errorf("write config: %w", err)
		}

		fmt.Println("Config created at:", configPath)

		return nil
	},
}

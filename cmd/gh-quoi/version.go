package main

import (
	"fmt"

	"github.com/BenbenIO/gh-quoi/internal/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of gh-quoi",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gh-quoi \n  Version:", version.Version)
	},
}

package main

import (
	"github.com/spf13/cobra"
)

type versionInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print CLI version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		return output(cmd.OutOrStdout(), versionInfo{
			Version: version,
			Commit:  commit,
			Date:    date,
		})
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

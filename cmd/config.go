package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print current config file",
	Run: func(cmd *cobra.Command, args []string) {
		app.PrintConfig()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

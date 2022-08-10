package cmd

import (
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the program",
	Run: func(cmd *cobra.Command, args []string) {
		app.Uninstall()
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

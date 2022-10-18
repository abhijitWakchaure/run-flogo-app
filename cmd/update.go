package cmd

import (
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the app with latest version",
	Run: func(cmd *cobra.Command, args []string) {
		a.Update()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

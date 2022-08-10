package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete all the flogo apps in apps dir",
	Run: func(cmd *cobra.Command, args []string) {
		app.Delete()
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

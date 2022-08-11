package cmd

import (
	"github.com/abhijitWakchaure/run-flogo-app/files"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete all the flogo apps in apps dir",
	Run: func(cmd *cobra.Command, args []string) {
		files.DeleteApps(a.AppsDir, a.AppPattern)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

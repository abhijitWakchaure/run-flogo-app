package cmd

import (
	"github.com/abhijitWakchaure/run-flogo-app/app"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version info of the program",
	Run: func(cmd *cobra.Command, args []string) {
		app.PrintVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

package cmd

import (
	"github.com/abhijitWakchaure/run-flogo-app/app"
	"github.com/abhijitWakchaure/run-flogo-app/software"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version info of the program",
	Run: func(cmd *cobra.Command, args []string) {
		software.PrintUpdateInfo(a.UpdateConfig)
		app.PrintVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

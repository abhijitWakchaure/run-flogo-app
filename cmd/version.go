package cmd

import (
	"github.com/abhijitWakchaure/run-flogo-app/core"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version info of the program",
	Run: func(cmd *cobra.Command, args []string) {
		core.Version()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

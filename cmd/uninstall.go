package cmd

import (
	"github.com/abhijitWakchaure/run-flogo-app/software"
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the program",
	Run: func(cmd *cobra.Command, args []string) {
		software.Uninstall(a.InstallPath)
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

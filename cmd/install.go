package cmd

import (
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the program",
	Run: func(cmd *cobra.Command, args []string) {
		app.Install()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

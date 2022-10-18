package cmd

import (
	"github.com/abhijitWakchaure/run-flogo-app/software"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the program",
	Run: func(cmd *cobra.Command, args []string) {
		software.Install("")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

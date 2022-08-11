package cmd

import (
	"github.com/abhijitWakchaure/run-flogo-app/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print current config file",
	Run: func(cmd *cobra.Command, args []string) {
		config.Print(a.AppConfig)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

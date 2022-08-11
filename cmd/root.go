// Package cmd includes all the commands for run-flogo-app cli
package cmd

import (
	"fmt"
	"os"

	"github.com/abhijitWakchaure/run-flogo-app/app"
	"github.com/abhijitWakchaure/run-flogo-app/config"
	"github.com/abhijitWakchaure/run-flogo-app/software"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

var a *app.App
var GENDOCS bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "run-flogo-app",
	Short: "Run the most recent flogo app from your apps dir",
	Long:  `Run the most recent flogo app from your configured apps dir. If the apps dir is not configured, the default will be used`,
	Run: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		list, _ := cmd.Flags().GetBool("list")
		name, _ := cmd.Flags().GetString("name")
		go func() {
			updateConfig := software.CheckForUpdates()
			// TODO: write update config
			_ = updateConfig
		}()
		if list {
			a.RunWithList(debug, args)
		}
		if name != "" {
			a.RunNamedApp(name, debug, args)
		}
		a.RunLatestApp(debug, args)
	},
	DisableAutoGenTag: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if GENDOCS {
		fmt.Println("i> Generating docs...")
		err := doc.GenMarkdownTree(rootCmd, "./docs")
		if err != nil {
			fmt.Printf("E> Failed to generate markdown docs! Error: %s\n", err.Error())
			os.Exit(1)
		}
		header := &doc.GenManHeader{
			Title:   config.AppName,
			Section: "3",
		}
		err = doc.GenManTree(rootCmd, header, "./manpages")
		if err != nil {
			fmt.Printf("E> Failed to generate man pages! Error: %s\n", err.Error())
			os.Exit(1)
		}
	}
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("debug", "d", false, "Enable debug logs")
	rootCmd.Flags().StringP("name", "n", "", "Run app with given (partial) name")
	rootCmd.Flags().BoolP("list", "l", false, "List last 5 apps and choose a number to run")
}

func initConfig() {
	home := config.GetUserHomeDir()
	viper.AddConfigPath(home)
	viper.SetConfigType("json")
	viper.SetConfigName(".run-flogo-app")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("i> Using config file: %s\n", viper.ConfigFileUsed())
	}
	appsDir := viper.GetString("appsDir")
	appPattern := viper.GetString("appPattern")
	isUpdateAvailable := viper.GetBool("isUpdateAvailable")
	updateURL := viper.GetString("updateURL")
	releaseNotes := viper.GetString("releaseNotes")

	appConfig := &config.AppConfig{
		AppsDir:    appsDir,
		AppPattern: appPattern,
	}
	updateConfig := &software.UpdateConfig{
		IsUpdateAvailable: isUpdateAvailable,
		UpdateURL:         updateURL,
		ReleaseNotes:      releaseNotes,
	}
	a = app.NewApp(appConfig, updateConfig)
}

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// VERSION ...
var VERSION string

// AppConfig ...
type AppConfig struct {
	AppsDir    string `json:"appsDir"`
	AppPattern string `json:"appPattern"`
}

// Print prints the current config
func Print(appConfig *AppConfig) {
	b, err := json.MarshalIndent(appConfig, "", "  ")
	if err != nil {
		fmt.Printf("\n#> Failed to read app config! error: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("\n#> Current app config:\n%s\n", string(b))
}

// Write will write the config into file
func Write(appConfig *AppConfig) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("\nE> Failed to get user home directory! Error: %s\n", err.Error())
		os.Exit(1)
	}
	if appConfig.AppsDir == "" {
		appConfig.AppsDir = path.Join(userHome, "Downloads")
	}
	configJSON, _ := json.MarshalIndent(appConfig, "", "\t")
	err = ioutil.WriteFile(path.Join(userHome, ConfigFileName), configJSON, 0644)
	if err != nil {
		fmt.Printf("E> Error ERR_WRITE_CONFIG: %s\n", err.Error())
	}
}

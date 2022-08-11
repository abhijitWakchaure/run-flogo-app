package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	userHome := GetUserHomeDir()
	if appConfig.AppsDir == "" {
		appConfig.AppsDir = filepath.Join(userHome, "Downloads")
	}
	configJSON, _ := json.MarshalIndent(appConfig, "", "\t")
	err := ioutil.WriteFile(filepath.Join(userHome, ConfigFileName), configJSON, 0644)
	if err != nil {
		fmt.Printf("E> Error ERR_WRITE_CONFIG: %s\n", err.Error())
	}
}

// GetUserHomeDir ...
func GetUserHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		pwd, err := os.Getwd()
		if pwd == "" {
			fmt.Printf("\nE> Failed to get PWD! Error %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("\nE> Could not get user home dir! Using PWD [%s] as home dir\n", pwd)
		return pwd
	}
	return home
}

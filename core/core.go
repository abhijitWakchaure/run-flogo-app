// Copyright (c) 2019 abhijit wakchaure. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package core implements methods that are essential to core operations of the app.
package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/abhijitWakchaure/run-flogo-app/utils"
)

// Constants for local env
const (
	AppName                  = "run-flogo-app"
	DefaultAppPatternLinux   = `^.+-linux_amd64.*$`
	DefaultAppPatternWindows = `^.+-windows_amd64.*$`
	DefaultAppPatternDarwin  = `^.+-darwin_amd64.*$`
	ConfigFile               = "run-flogo-app.config"
	GithubLastestReleaseURL  = "https://api.github.com/repos/abhijitWakchaure/run-flogo-app/releases/latest"
	GithubDownloadBaseURL    = "https://github.com/abhijitWakchaure/run-flogo-app/releases/download/"
	GithubBaseURL            = "https://github.com/abhijitWakchaure/run-flogo-app"
	GithubIssuesURL          = "https://github.com/abhijitWakchaure/run-flogo-app/issues"
	CurrentAppVersion        = "v1.5"
	InstallPathLinux         = "/usr/local/bin"
	InstallPathDarwin        = "/usr/local/bin"
	InstallPathWindows       = `C:\Windows\system32`
)

// App holds the environmet variables for the user
type App struct {
	_TempAppName      string
	_InstallPath      string
	_AppConfigPath    string
	AppDir            string `json:"appsDir" binding:"required"`
	AppPattern        string `json:"appPattern"`
	IsUpdateAvailable bool   `json:"isUpdateAvailable"`
	UpdateURL         string `json:"updateURL"`
	ReleaseNotes      string `json:"releaseNotes"`
}

// Init ...
func (a *App) Init() {
	if runtime.GOOS == "linux" {
		a._TempAppName = AppName + "-linux-amd64"
		a.AppPattern = DefaultAppPatternLinux
		a._InstallPath = InstallPathLinux
	} else if runtime.GOOS == "windows" {
		a._TempAppName = AppName + "-windows-amd64.exe"
		a.AppPattern = DefaultAppPatternWindows
		a._InstallPath = InstallPathWindows
	} else if runtime.GOOS == "darwin" {
		a._TempAppName = AppName + "-darwin-amd64"
		a.AppPattern = DefaultAppPatternDarwin
		a._InstallPath = InstallPathDarwin
	}
	a._AppConfigPath = path.Join(utils.GetUserHomeDir(), ConfigFile)
}

// ReadConfig will read the config from configuration file
func (a *App) ReadConfig() {
	fileExists, err := utils.CheckFileExists(a._AppConfigPath)
	if err != nil {
		log.Fatalln("# Error: ERR_READ_CONFIG", err)
	}
	if !fileExists {
		fmt.Print("#> Creating config file...")
		a.WriteConfig()
		return
	}
	f, err := ioutil.ReadFile(a._AppConfigPath)
	if err != nil {
		fmt.Println("#> Unable to read config...ignoring config...using defaults")
		a.loadDefaultConfig()
		return
	}
	err = json.Unmarshal(f, &a)
	if err != nil {
		fmt.Print("#> Invalid config detected...rewriting config...")
		a.WriteConfig()
	}
}

// WriteConfig will write the config into file
func (a *App) WriteConfig() {
	if a.AppDir == "" {
		a.loadDefaultConfig()
	}
	configJSON, _ := json.MarshalIndent(a, "", "\t")
	err := ioutil.WriteFile(a._AppConfigPath, configJSON, 0600)
	if err != nil {
		log.Fatalln("# Error: ERR_WRITE_CONFIG", err)
	}
	// fmt.Println("done")
}

func (a *App) loadDefaultConfig() {
	fmt.Print("loading default config...")
	a.AppDir = path.Join(utils.GetUserHomeDir(), "Downloads")
	fmt.Println("done")
}

func (a *App) validateConfig() {
	if a.AppDir == "" {
		fmt.Print("#> Invalid config detected...")
		a.WriteConfig()
	}
}

// FindLatestApp will return the latest flogo app name
func (a *App) FindLatestApp() string {
	files, err := ioutil.ReadDir(a.AppDir)
	if err != nil {
		log.Fatal(err)
	}
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})
	validApp := regexp.MustCompile(a.AppPattern)
	for _, f := range files {
		if !f.IsDir() && validApp.MatchString(f.Name()) {
			return path.Join(a.AppDir, f.Name())
		}
	}
	log.Println("#> Info: No flogo apps found in " + a.AppDir)
	return ""
}

// CheckForUpdates will check for latest release
func (a *App) CheckForUpdates() {
	resp, err := http.Get(GithubLastestReleaseURL)
	if err != nil {
		fmt.Println()
		log.Println("# Info: Unable to reach server for updates.")
		fmt.Println()
		return
	}
	defer resp.Body.Close()
	var gitdata map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&gitdata)
	if err != nil {
		fmt.Println()
		log.Fatalln("# Error: ERR_CHKUPDATE_DECODE", err)
		log.Fatalln("# Request: Please create an issue here for this error:", GithubIssuesURL)
	}
	for _, d := range gitdata["assets"].([]interface{}) {
		durl := d.(map[string]interface{})["browser_download_url"].(string)
		if strings.Contains(durl, runtime.GOOS) && !strings.Contains(durl, CurrentAppVersion) {
			a.IsUpdateAvailable = true
			a.UpdateURL = durl
			a.ReleaseNotes = strings.Replace(strings.TrimSpace(gitdata["body"].(string)), "\n", "\n\t", -1)
			a.WriteConfig()
			return
		} else if strings.Contains(durl, runtime.GOOS) {
			fmt.Println()
			// fmt.Println("Your app is up to date 👍")
			return
		}
	}
}

// PrintUpdateInfo will print the update info
func (a *App) PrintUpdateInfo() {
	if a.IsUpdateAvailable {
		fmt.Println("#> New version of the app is available at:", a.UpdateURL)
		fmt.Println("#> Release Notes:")
		fmt.Printf("\t%s\n\n", a.ReleaseNotes)
	}
}

// Install will install the program
func (a *App) Install() {
	fmt.Print("#> Installing run-flogo-app...")
	ex, err := os.Executable()
	if err != nil {
		fmt.Println("failed")
		log.Fatalln("# Error: ERR_INSTALL_SELFPATH", err)
	}
	var src string
	var dst string
	if runtime.GOOS == "windows" {
		src = filepath.Dir(ex) + a._TempAppName
		dst = a._InstallPath + string(os.PathSeparator) + AppName + ".exe"
	} else {
		src = path.Join(filepath.Dir(ex), a._TempAppName)
		dst = path.Join(a._InstallPath, AppName)
	}
	err = utils.Copy(src, dst)
	if err != nil {
		fmt.Println("failed")
		log.Fatalln("# Error: ERR_INSTALL_COPY", err)
	}
	fmt.Println("done")
	fmt.Println("#> You can now directly execute ", AppName)
}

// Uninstall will install the program
func (a *App) Uninstall() {
	fmt.Println("#> Uninstalling run-flogo-app...")
	fmt.Print("...Deleting config file...")
	err := os.Remove(a._AppConfigPath)
	if err != nil {
		fmt.Println("failed")
		log.Println("# Error: ERR_UNINSTALL_CLRCONFIG", err)
	}
	fmt.Print("...Deleting main executable...")
	var target string
	if runtime.GOOS == "windows" {
		target = a._InstallPath + string(os.PathSeparator) + AppName + ".exe"
	} else {
		target = path.Join(a._InstallPath, AppName)
	}
	err = utils.Remove(target)
	if err != nil {
		fmt.Println("failed")
		fmt.Println("#> Unable to uninstall run-flogo-app...you can manually delete", path.Join(a._InstallPath, AppName))
		log.Fatalln("# Error: ERR_UNINSTALL_REMOVE", err)
	}
	fmt.Println()
	fmt.Println("#> Finished uninstalling run-flogo-app")
}

// Version ...
func (a *App) Version() {
	fmt.Println("## run-flogo-app")
	fmt.Println("#> Version:", CurrentAppVersion)
	fmt.Println("#> Developer: Abhijit Wakchaure")
	fmt.Println("#> Github:", GithubBaseURL)
}

// Main runs the core functions
func (a *App) Main() {
	go a.CheckForUpdates()
	a.ReadConfig()
	a.PrintUpdateInfo()
	a.validateConfig()
}

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
	ConfigFileName           = ".run-flogo-app"
	DefaultAppPatternLinux   = `^.+-linux_amd64.*$`
	DefaultAppPatternWindows = `^.+-windows_amd64.*$`
	DefaultAppPatternDarwin  = `^.+-darwin_amd64.*$`
	GithubLastestReleaseURL  = "https://api.github.com/repos/abhijitWakchaure/run-flogo-app/releases/latest"
	GithubDownloadBaseURL    = "https://github.com/abhijitWakchaure/run-flogo-app/releases/download/"
	GithubBaseURL            = "https://github.com/abhijitWakchaure/run-flogo-app"
	GithubIssuesURL          = "https://github.com/abhijitWakchaure/run-flogo-app/issues"
	CurrentAppVersion        = "v2.0.0"
	InstallPathLinux         = "/usr/local/bin"
	InstallPathDarwin        = "/usr/local/bin"
	InstallPathWindows       = `C:\Windows\system32`
)

// App holds the environment variables for the user
type App struct {
	_TempAppName      string
	_InstallPath      string
	AppDir            string `json:"appsDir" binding:"required"`
	AppPattern        string `json:"appPattern"`
	IsUpdateAvailable bool   `json:"isUpdateAvailable"`
	UpdateURL         string `json:"updateURL"`
	ReleaseNotes      string `json:"releaseNotes"`
}

// NewApp ...
func NewApp(AppDir, AppPattern string) *App {
	app := new(App)
	app.AppDir = AppDir
	app.AppPattern = AppPattern
	app.Init()
	go app.checkForUpdates()
	app.printUpdateInfo()
	app.validateConfig()
	return app
}

// Init ...
func (a *App) Init() {
	var appPattern string
	switch runtime.GOOS {
	case "linux":
		a._TempAppName = AppName + "-linux-amd64"
		appPattern = DefaultAppPatternLinux
		a._InstallPath = InstallPathLinux
	case "windows":
		a._TempAppName = AppName + "-windows-amd64.exe"
		appPattern = DefaultAppPatternWindows
		a._InstallPath = InstallPathWindows
	case "darwin":
		a._TempAppName = AppName + "-darwin-amd64"
		appPattern = DefaultAppPatternDarwin
		a._InstallPath = InstallPathDarwin
	default:
		fmt.Printf("\nError: OS %s is not yet supported, please contact developers\n", runtime.GOOS)
	}
	if a.AppPattern == "" {
		a.AppPattern = appPattern
	}
}

// RunLatestApp will run the latest app
func (a *App) RunLatestApp(debug bool, args []string) {
	latestFlogoApp := a.findLatestApp()
	if len(latestFlogoApp) == 0 {
		os.Exit(1)
	}
	fmt.Print("#> Do you want to execute this app \"", latestFlogoApp, "\" [Y/n]: ")
	choice := utils.HandleYNInput()
	if choice {
		utils.MakeAppExecutable(latestFlogoApp)
		utils.RunFlogoApp(latestFlogoApp, debug, args)
	}
}

// PrintConfig prints the current config
func (a *App) PrintConfig() {
	config := struct {
		AppDir     string `json:"appsDir"`
		AppPattern string `json:"appPattern"`
	}{
		AppDir:     a.AppDir,
		AppPattern: a.AppPattern,
	}
	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Printf("\n#> Failed to read app config! error: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("\n#> Current app config:\n%s\n", string(b))
}

// writeConfig will write the config into file
func (a *App) writeConfig() {
	if a.AppDir == "" {
		a.AppDir = path.Join(utils.GetUserHomeDir(), "Downloads")
	}
	configJSON, _ := json.MarshalIndent(a, "", "\t")
	err := ioutil.WriteFile(path.Join(utils.GetUserHomeDir(), ConfigFileName), configJSON, 0644)
	if err != nil {
		log.Fatalln("# Error: ERR_WRITE_CONFIG", err)
	}
}

func (a *App) validateConfig() {
	if a.AppDir == "" {
		fmt.Print("#> Invalid config detected...")
		a.writeConfig()
	}
}

// findLatestApp will return the latest flogo app name
func (a *App) findLatestApp() string {
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
	fmt.Println("#> No flogo apps found in " + a.AppDir)
	return ""
}

// checkForUpdates will check for latest release
func (a *App) checkForUpdates() {
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
	assets, ok := gitdata["assets"].([]interface{})
	if !ok {
		return
	}
	for _, d := range assets {
		durl := d.(map[string]interface{})["browser_download_url"].(string)
		if strings.Contains(durl, runtime.GOOS) && !strings.Contains(durl, CurrentAppVersion) {
			a.IsUpdateAvailable = true
			a.UpdateURL = durl
			a.ReleaseNotes = strings.Replace(strings.TrimSpace(gitdata["body"].(string)), "\n", "\n\t", -1)
			a.writeConfig()
			return
		} else if strings.Contains(durl, runtime.GOOS) {
			// fmt.Println()
			// fmt.Println("Your app is up to date 👍")
			return
		}
	}
}

// printUpdateInfo will print the update info
func (a *App) printUpdateInfo() {
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
	err := os.Remove(path.Join(utils.GetUserHomeDir(), ConfigFileName))
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

// Delete will delete all the flogo apps in apps dir
func (a *App) Delete() {
	fmt.Printf("#> Listing all the flogo app(s) inside [%s]...\n", a.AppDir)
	files, err := ioutil.ReadDir(a.AppDir)
	if err != nil {
		log.Fatal(err)
	}
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})
	validApp := regexp.MustCompile(a.AppPattern)
	var count int
	apps := []string{}
	for i, f := range files {
		if !f.IsDir() && validApp.MatchString(f.Name()) {
			apps = append(apps, path.Join(a.AppDir, f.Name()))
			count++
			fmt.Printf("%d. %s\n", i+1, path.Join(a.AppDir, f.Name()))
		}
	}
	if count == 0 {
		fmt.Println("#> No flogo app found inside apps dir.")
		os.Exit(0)
	}
	fmt.Printf("\nAre you sure you want to delete all %d app(s)? [y/n] ", count)
	choice := utils.HandleYNInput()
	if choice {
		fmt.Printf("\n#> Deleting %d app(s)...\n", count)
		for _, f := range apps {
			err = os.Remove(f)
			if err != nil {
				fmt.Printf("\n#> Failed to delete app [%s] error: %s", f, err.Error())
			}
		}
		if err != nil {
			os.Exit(1)
		}
		fmt.Printf("\n#> Finished deleting %d apps\n", count)
		os.Exit(0)
	}
	fmt.Println("No app(s) were deleted!")
}

// Version ...
func Version() {
	fmt.Println("#> Run FLOGO App")
	fmt.Println("#> Version:", CurrentAppVersion)
	fmt.Println("#> Developer: Abhijit Wakchaure")
	fmt.Println("#> Github:", GithubBaseURL)
}

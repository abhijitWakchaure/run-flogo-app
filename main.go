package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

// Constants for local env
const (
	DefaultAppPatternLinux   = `^.+-linux_amd64.*$`
	DefaultAppPatternWindows = `^.+-windows_amd64.*$`
	DefaultAppPatternDarwin  = `^.+-darwin_amd64.*$`
	ConfigFile               = "run-flogo-app-config.json"
	GithubLastestReleaseURL  = "https://api.github.com/repos/abhijitWakchaure/run-flogo-app/releases/latest"
	GithubDownloadBaseURL    = "https://github.com/abhijitWakchaure/run-flogo-app/releases/download/"
	CurrentAppVersion        = "v1.3"
)

// App holds the environmet variables for the user
type App struct {
	AppDir     string `json:"rfAppDir" binding:"required"`
	AppPattern string `json:"rfAppPattern"`
}

// NewApp ...
func NewApp() *App {
	app := new(App)
	app.init()
	return app
}

func (a *App) init() {
	CheckForUpdates()
	a.ReadConfig()
	a.validateConfig()
}

// ReadConfig will read the config from configuration file
func (a *App) ReadConfig() {
	fileExists, err := checkFileExists(ConfigFile)
	if err != nil {
		log.Fatalln("# Error: ERR_READ_CONFIG", err)
	}
	if !fileExists {
		fmt.Print("#> Creating config file...")
		a.WriteConfig()
		return
	}
	f, err := ioutil.ReadFile(ConfigFile)
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
	a.loadDefaultConfig()
	configJSON, _ := json.MarshalIndent(a, "", "\t")
	err := ioutil.WriteFile(ConfigFile, configJSON, 0600)
	if err != nil {
		log.Fatalln("# Error: ERR_WRITE_CONFIG", err)
	}
	fmt.Println("done")
}

func (a *App) loadDefaultConfig() {
	fmt.Print("loading default config...")
	a.AppDir = path.Join(GetUserHomeDir(), "Downloads")
	if runtime.GOOS == "linux" {
		a.AppPattern = DefaultAppPatternLinux
	} else if runtime.GOOS == "windows" {
		a.AppPattern = DefaultAppPatternWindows
	} else if runtime.GOOS == "darwin" {
		a.AppPattern = DefaultAppPatternDarwin
	}
}

func (a *App) validateConfig() {
	if a.AppDir == "" {
		fmt.Print("#> Invalid config detected...")
		a.WriteConfig()
	}
}

func checkFileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

// GetUserHomeDir returns the Home Directory of the user
func GetUserHomeDir() string {
	u, err := user.Current()
	if err != nil {
		log.Fatalln("# Error: ERR_GET_HOMEDIR", err)
	}
	return u.HomeDir
}

// EnableDebugLogs will add env variable to enable debug logs
func EnableDebugLogs(cmd *exec.Cmd) *exec.Cmd {
	debugFlag := `FLOGO_LOG_LEVEL=DEBUG`
	cmd.Env = append(os.Environ(), debugFlag)
	return cmd
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

// MakeAppExecutable will make the app executable
func MakeAppExecutable(app string) {
	fmt.Println("#> Making app executable...")
	err := os.Chmod(app, 500)
	if err != nil {
		log.Fatalln("# Error: ERR_MAKE_APP_EXEC", err)
	}
}

// RunFlogoApp will run the app
func RunFlogoApp(app string, debug *bool, tail []string) {
	cmd := exec.Command(app, tail...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	fmt.Println("#> Executing: ", strings.Join(cmd.Args, " "))
	if *debug {
		cmd = EnableDebugLogs(cmd)
	}
	err := cmd.Run()
	if err != nil {
		log.Fatalln("# Error: ERR_RUN_FA", err)
	}
}

// CheckForUpdates will check for latest release
func CheckForUpdates() {
	fmt.Print("#> Checking for updates...")
	resp, err := http.Get(GithubLastestReleaseURL)
	if err != nil {
		log.Println("Unable to check for updates...", err)
	}
	defer resp.Body.Close()
	var gitdata map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&gitdata)
	if err != nil {
		fmt.Println()
		log.Fatalln("# Error: ERR_CHKUPDATE_DECODE", err)
	}
	for _, d := range gitdata["assets"].([]interface{}) {
		durl := d.(map[string]interface{})["browser_download_url"].(string)
		if strings.Contains(durl, runtime.GOOS) && !strings.Contains(durl, CurrentAppVersion) {
			fmt.Println("New version of the app is available at", durl)
		} else if strings.Contains(durl, runtime.GOOS) {
			fmt.Println("Your app is up to date")
		}
	}
}

// HandleYNInput handles the Yes/No input
func HandleYNInput() rune {
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		log.Fatalln("# Error: ERR_READ_USRIN", err)
	}
	return char
}

func main() {
	app := NewApp()
	latestFlogoApp := app.FindLatestApp()
	flagDebug := flag.Bool("debug", false, "Set this to enable debug logs")
	flag.Parse()

	if len(latestFlogoApp) > 0 {
		fmt.Print("#> Do you want to execute this app \"", latestFlogoApp, "\" [Y/n]: ")
		choice := HandleYNInput()
		if choice == 'Y' || choice == 'y' {
			if runtime.GOOS == "windows" {
				// TODO: Handle for Windows
			}
			MakeAppExecutable(latestFlogoApp)
			RunFlogoApp(latestFlogoApp, flagDebug, flag.Args())
		} else {
			log.Println("# Info: Exiting...")
		}
	}
}

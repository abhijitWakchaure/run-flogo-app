package app

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/abhijitWakchaure/run-flogo-app/config"
	"github.com/abhijitWakchaure/run-flogo-app/files"
	"github.com/abhijitWakchaure/run-flogo-app/software"
)

// App holds the environment variables for the user
type App struct {
	*config.AppConfig
	*software.UpdateConfig
	InstallPath string
}

// NewApp ...
func NewApp(appConfig *config.AppConfig, updateConfig *software.UpdateConfig) *App {
	a := new(App)
	a.AppConfig = appConfig
	a.UpdateConfig = updateConfig
	var appPattern string
	switch runtime.GOOS {
	case "linux":
		appPattern = config.DefaultAppPatternLinux
		a.InstallPath = config.InstallPathLinux
	case "windows":
		appPattern = config.DefaultAppPatternWindows
		a.InstallPath = config.InstallPathWindows
	case "darwin":
		appPattern = config.DefaultAppPatternDarwin
		a.InstallPath = config.InstallPathDarwin
	default:
		fmt.Printf("\nError: OS %s is not yet supported, please contact developers\n", runtime.GOOS)
		os.Exit(1)
	}
	if a.AppPattern == "" {
		a.AppPattern = appPattern
	}
	software.PrintUpdateInfo(a.UpdateConfig)
	// app.validateConfig()
	return a
}

// PrintConfig will print the app config
func (a *App) PrintConfig() {
	c := &config.AppConfig{
		AppsDir:    a.AppsDir,
		AppPattern: a.AppPattern,
	}
	config.Print(c)
}

// RunLatestApp will run the latest app
func (a *App) RunLatestApp(debug bool, args []string) {
	latestFlogoApp := files.FindLatestApp(a.AppsDir, a.AppPattern)
	if len(latestFlogoApp) == 0 {
		os.Exit(1)
	}
	fmt.Print("#> Do you want to execute this app \"", latestFlogoApp, "\" [Y/n]: ")
	choice := software.HandleYNInput()
	if !choice {
		os.Exit(0)
	}
	fmt.Println("#> Making app executable...")
	err := os.Chmod(latestFlogoApp, 500)
	if err != nil {
		fmt.Printf("\nE> Error ERR_MAKE_APP_EXEC: %s\n", err.Error())
		os.Exit(1)
	}
	cmd := exec.Command(latestFlogoApp, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	fmt.Printf("\n#> Executing: %s\n\n", strings.Join(cmd.Args, " "))
	if debug {
		debugFlag := `FLOGO_LOG_LEVEL=DEBUG`
		cmd.Env = append(os.Environ(), debugFlag)
	}
	err = cmd.Run()
	if err != nil {
		fmt.Printf("\nE> Error ERR_RUN_FA: %s\n", err.Error())
		os.Exit(1)
	}
}

// PrintVersion ...
func PrintVersion() {
	fmt.Println("#> Run Flogo App")
	fmt.Println("#> Version:", config.VERSION)
	fmt.Println("#> Developer: Abhijit Wakchaure")
	fmt.Println("#> Github:", config.GithubBaseURL)
}

package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/abhijitWakchaure/run-flogo-app/config"
	"github.com/abhijitWakchaure/run-flogo-app/files"
	"github.com/abhijitWakchaure/run-flogo-app/software"
	"github.com/spf13/viper"
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
	if a.AppsDir == "" {
		a.AppsDir = filepath.Join(config.GetUserHomeDir(), "Downloads")
	}
	WriteAppConfig(a.AppConfig)
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

// WriteAppConfig will write the app config
func WriteAppConfig(appConfig *config.AppConfig) {
	viper.Set("appsDir", appConfig.AppsDir)
	viper.Set("appPattern", appConfig.AppPattern)
	viper.WriteConfig()
}

// RunLatestApp will run the latest app
func (a *App) RunLatestApp(logLevel string, args []string) {
	latestFlogoApp := files.FindLatestApp(a.AppsDir, a.AppPattern)
	if len(latestFlogoApp) == 0 {
		os.Exit(1)
	}
	fmt.Printf("#> Do you want to execute the app '%s' [y/n]: ", latestFlogoApp)
	choice := software.HandleYNInput()
	if !choice {
		os.Exit(0)
	}
	runExecutable(latestFlogoApp, logLevel, args)
}

// RunNamedApp will run the app with given (partial) name
// If there are multiple matches, it will ask for user to choose
func (a *App) RunNamedApp(name, logLevel string, args []string) {
	flogoApps := files.FindAppsWithName(a.AppsDir, a.AppPattern, name)
	if len(flogoApps) == 0 {
		fmt.Printf("\n#> No flogo apps found containing name [%s] in apps dir [%s]\n", name, a.AppsDir)
		os.Exit(1)
	}
	if len(flogoApps) == 1 {
		flogoApp := flogoApps[0]
		fmt.Printf("#> Do you want to execute the app '%s' [y/n]: ", flogoApp)
		choice := software.HandleYNInput()
		if !choice {
			os.Exit(0)
		}
		runExecutable(flogoApp, logLevel, args)
	}
	fmt.Printf("#> Got %d matches for query [%s]:\n", len(flogoApps), name)
	for i, v := range flogoApps {
		fmt.Printf("%d. %s\n", i+1, filepath.Base(v))
	}
	fmt.Printf("\n#> Choose an app that you want to execute [1-%d]: ", len(flogoApps))
	choice := software.HandleNumericInput()
	if choice < 1 || choice > len(flogoApps) {
		fmt.Printf("\nE> Invalid choice, please choose a number between 1 and %d\n", len(flogoApps))
		os.Exit(1)
	}
	flogoApp := flogoApps[choice-1]
	runExecutable(flogoApp, logLevel, args)
}

// RunWithList will list the last 5 apps and will ask user to select 1
func (a *App) RunWithList(logLevel string, args []string) {
	flogoApps := files.ListLastNApps(a.AppsDir, a.AppPattern, config.MaxAppsWithList)
	if len(flogoApps) == 0 {
		fmt.Printf("\n#> No flogo apps found in apps dir [%s]\n", a.AppsDir)
		os.Exit(1)
	}
	if len(flogoApps) == 1 {
		flogoApp := flogoApps[0]
		fmt.Printf("#> Do you want to execute the app '%s' [y/n]: ", flogoApp)
		choice := software.HandleYNInput()
		if !choice {
			os.Exit(0)
		}
		runExecutable(flogoApp, logLevel, args)
	}
	fmt.Printf("#> Here is the list of apps:\n")
	for i, v := range flogoApps {
		fmt.Printf("%d. %s\n", i+1, filepath.Base(v))
	}
	fmt.Printf("\n#> Choose an app that you want to execute [1-%d]: ", len(flogoApps))
	choice := software.HandleNumericInput()
	if choice < 1 || choice > len(flogoApps) {
		fmt.Printf("\nE> Invalid choice, please choose a number between 1 and %d\n", len(flogoApps))
		os.Exit(1)
	}
	flogoApp := flogoApps[choice-1]
	runExecutable(flogoApp, logLevel, args)
}

// Update will update the app to latest version released on Github
func (a *App) Update() {
	software.Update(a.AppConfig)
}

// PrintVersion ...
func PrintVersion() {
	fmt.Println("#> Run Flogo App")
	fmt.Println("#> Version:", config.VERSION)
	fmt.Println("#> Developer: Abhijit Wakchaure")
	fmt.Println("#> Github:", config.GithubBaseURL)
}

func runExecutable(path, logLevel string, args []string) {
	fmt.Println("\n#> Making app executable...")
	err := os.Chmod(path, 0700)
	if err != nil {
		fmt.Printf("\nE> Error ERR_MAKE_APP_EXEC: %s\n", err.Error())
		os.Exit(1)
	}
	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	fmt.Printf("#> Executing: %s\n\n", strings.Join(cmd.Args, " "))
	if logLevel != "" {
		logLevelEnv := fmt.Sprintf("FLOGO_LOG_LEVEL=%s", logLevel)
		cmd.Env = append(cmd.Env, logLevelEnv)
	}
	err = cmd.Run()
	if err != nil {
		fmt.Printf("\nE> Error ERR_RUN_FA: %s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

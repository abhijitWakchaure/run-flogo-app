package software

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/abhijitWakchaure/run-flogo-app/config"
	"github.com/spf13/viper"
)

// UpdateConfig ...
type UpdateConfig struct {
	IsUpdateAvailable bool   `json:"isUpdateAvailable"`
	UpdateURL         string `json:"updateURL"`
	ReleaseNotes      string `json:"releaseNotes"`
}

// Install will install the program
func Install(installPath string) {
	fmt.Print("#> Installing run-flogo-app...")
	ex, err := os.Executable()
	if err != nil {
		fmt.Println("failed")
		fmt.Printf("\n# Error: ERR_INSTALL_SELFPATH %s\n", err.Error())
		os.Exit(1)
	}
	var src string
	var dst string
	src, err = filepath.EvalSymlinks(ex)
	if err != nil {
		fmt.Println("failed")
		fmt.Printf("\n# Error: ERR_INSTALL_EVALSYMLNK %s\n", err.Error())
		os.Exit(1)
	}
	if runtime.GOOS == "windows" {
		dst = filepath.Join(installPath, config.AppName+".exe")
	} else {
		dst = filepath.Join(installPath, config.AppName)
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("sudo", "cp", "-fpv", src, dst)
	case "windows":
		cmd = exec.Command("cmd.exe", "/C", `copy /Y `+src+" "+dst)
	case "darwin":
		cmd = exec.Command("sudo", "cp", "-fpv", src, dst)
	default:
		fmt.Printf("\nError: OS %s is not yet supported, please contact developers\n", runtime.GOOS)
		os.Exit(1)
	}
	err = cmd.Run()
	if err != nil {
		fmt.Println("failed")
		fmt.Printf("\n# Error: ERR_INSTALL_COPY %s\n", err)
		os.Exit(1)
	}
	fmt.Println("done")
	fmt.Println("#> You can now directly execute", config.AppName)
}

// Uninstall will install the program
func Uninstall(installPath string) {
	fmt.Println("#> Uninstalling run-flogo-app...")
	fmt.Printf("   Deleting config file...")
	userHome := config.GetUserHomeDir()
	os.Remove(filepath.Join(userHome, config.ConfigFileName))
	fmt.Printf("\n   Deleting main executable...")
	var target string
	if runtime.GOOS == "windows" {
		target = installPath + string(os.PathSeparator) + config.AppName + ".exe"
	} else {
		target = filepath.Join(installPath, config.AppName)
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("sudo", "rm", target)
	case "windows":
		cmd = exec.Command("cmd.exe", "/C", `del `+target)
	case "darwin":
		cmd = exec.Command("sudo", "rm", target)
	default:
		fmt.Printf("\nError: OS %s is not yet supported, please contact developers\n", runtime.GOOS)
		os.Exit(1)
	}
	err := cmd.Run()
	if err != nil {
		fmt.Println("failed")
		fmt.Println("#> Unable to uninstall run-flogo-app! Error ERR_UNINSTALL_REMOVE...you can manually delete", filepath.Join(installPath, config.AppName))
		os.Exit(1)
	}
	fmt.Printf("\n#> Finished uninstalling run-flogo-app")
}

// CheckForUpdates will check for latest release
func CheckForUpdates() *UpdateConfig {
	resp, err := http.Get(config.GithubLastestReleaseURL)
	if err != nil {
		fmt.Printf("\n\nE> run-flogo-app Error: ERR_CHKUPDATE_HTTPGET %s\n", err)
		return nil
	}
	defer resp.Body.Close()
	var gitdata map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&gitdata)
	if err != nil {
		fmt.Printf("\n\nE> run-flogo-app Error: ERR_CHKUPDATE_DECODE %s\n", err)
		fmt.Printf("\nPlease create an issue here for this error: %s\n\n", config.GithubIssuesURL)
	}
	assets, ok := gitdata["assets"].([]interface{})
	if !ok {
		fmt.Printf("\nE> run-flogo-app Error: ERR_CHKUPDATE_DECODE %s\n", err)
		return nil
	}
	OSAndArch := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)
	for _, d := range assets {
		durl := d.(map[string]interface{})["browser_download_url"].(string)
		if !strings.Contains(durl, OSAndArch) {
			continue
		}
		if strings.Contains(durl, config.VERSION) {
			// fmt.Println()
			// fmt.Println("Your app is up to date 👍")
			return nil
		}
		return &UpdateConfig{
			IsUpdateAvailable: true,
			UpdateURL:         durl,
			ReleaseNotes:      strings.Replace(strings.TrimSpace(gitdata["body"].(string)), "\n", "\n\t", -1),
		}

	}
	return nil
}

// WriteUpdateConfig will write the update info
func WriteUpdateConfig(updateConfig *UpdateConfig) {
	viper.Set("isUpdateAvailable", updateConfig.IsUpdateAvailable)
	viper.Set("updateURL", updateConfig.UpdateURL)
	viper.Set("releaseNotes", updateConfig.ReleaseNotes)
	viper.WriteConfig()
}

// PrintUpdateInfo will print the update info
func PrintUpdateInfo(updateConfig *UpdateConfig) {
	if updateConfig.IsUpdateAvailable {
		fmt.Println("#> New version of the app is available!")
		fmt.Println("#> Release Notes:")
		fmt.Printf("\t%s\n\n", updateConfig.ReleaseNotes)
	}
}

// HandleYNInput handles the Yes/No input
func HandleYNInput() bool {
	reader := bufio.NewReader(os.Stdin)
	inputBytes, _, err := reader.ReadLine()
	if err != nil {
		fmt.Printf("\nE> Error ERR_READ_USRIN: %s\n", err.Error())
	}
	input := string(inputBytes)
	if strings.EqualFold(input, "y") || strings.EqualFold(input, "yes") {
		return true
	}
	if strings.EqualFold(input, "n") || strings.EqualFold(input, "no") {
		return false
	}
	choice, err := strconv.ParseBool(input)
	if err != nil {
		fmt.Printf("\nE> Error ERR_PARSE_BOOL: %s\n", err.Error())
	}
	return choice
}

// HandleNumericInput handles the numeric input
func HandleNumericInput() int {
	reader := bufio.NewReader(os.Stdin)
	inputBytes, _, err := reader.ReadLine()
	if err != nil {
		fmt.Printf("\nE> Error ERR_READ_USRIN: %s\n", err.Error())
		os.Exit(1)
	}
	input := string(inputBytes)
	n, err := strconv.Atoi(input)
	if err != nil {
		fmt.Printf("\nE> Error ERR_PARSE_NUMBER: %s\n", err.Error())
		os.Exit(1)
	}
	return n
}

package software

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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
func Install(src string) {
	fmt.Print("#> Installing run-flogo-app...")
	if src == "" {
		ex, err := os.Executable()
		if err != nil {
			fmt.Println("failed")
			fmt.Printf("\n# Error: ERR_INSTALL_SELFPATH %s\n", err.Error())
			os.Exit(1)
		}
		src, err = filepath.EvalSymlinks(ex)
		if err != nil {
			fmt.Println("failed")
			fmt.Printf("\n# Error: ERR_INSTALL_EVALSYMLNK %s\n", err.Error())
			os.Exit(1)
		}
	}
	err := os.Chmod(src, 0777)
	if err != nil {
		fmt.Println("failed")
		fmt.Printf("\n# Error: ERR_INSTALL_MAKEEXECUTABLE %s\n", err.Error())
		os.Exit(1)
	}
	var dst string
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		dst = filepath.Join(config.InstallPathLinux, config.AppName)
		cmd = exec.Command("sudo", "cp", "-fpv", src, dst)
	case "windows":
		dst = filepath.Join(config.InstallPathWindows, config.AppName+".exe")
		cmd = exec.Command("cmd.exe", "/C", `copy /Y `+src+" "+dst)
	case "darwin":
		dst = filepath.Join(config.InstallPathDarwin, config.AppName)
		cmd = exec.Command("sudo", "cp", "-fpv", src, dst)
	default:
		fmt.Printf("\nError: OS %s is not yet supported, please contact developer(s) to add support\n", runtime.GOOS)
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
	WriteUpdateConfig(nil)
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

// Update will update the app
func Update(appConfig *config.AppConfig) {
	updateConfig, err := CheckForUpdates()
	if err != nil {
		fmt.Printf("\nFailed to check for updates due to:\n%s", err.Error())
		os.Exit(1)
	}
	if updateConfig == nil || len(updateConfig.UpdateURL) <= 2 {
		fmt.Println("Your app is up to date 👍")
		WriteUpdateConfig(nil)
		os.Exit(0)
	}
	WriteUpdateConfig(updateConfig)
	binaryName := filepath.Base(updateConfig.UpdateURL)
	downloadPath := filepath.Join(appConfig.AppsDir, binaryName)
	fmt.Printf("Downloading latest version from: %s\n\n", updateConfig.UpdateURL)
	err = DownloadFile(downloadPath, updateConfig.UpdateURL)
	if err != nil {
		fmt.Printf("\nFailed to download updated app due to:\n%s", err.Error())
		os.Exit(1)
	}
	Install(downloadPath)
}

// CheckForUpdates will check for latest release
func CheckForUpdates() (*UpdateConfig, error) {
	resp, err := http.Get(config.GithubLastestReleaseURL)
	if err != nil {
		err = fmt.Errorf("E> run-flogo-app Error: ERR_CHKUPDATE_HTTPGET %s", err)
		fmt.Printf("\n\n%s\n", err)
		return nil, err
	}
	defer resp.Body.Close()
	var gitdata map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&gitdata)
	if err != nil {
		err = fmt.Errorf("E> run-flogo-app Error: ERR_CHKUPDATE_DECODE %s", err)
		fmt.Printf("\n\n%s\n", err)
		fmt.Printf("\nPlease create an issue here for this error: %s\n\n", config.GithubIssuesURL)
		return nil, err
	}
	assets, ok := gitdata["assets"].([]interface{})
	if !ok {
		err = fmt.Errorf("E> run-flogo-app Error: ERR_CHKUPDATE_NOASSETS %s", err)
		fmt.Printf("\n\n%s\n", err)
		return nil, err
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
			WriteUpdateConfig(nil)
			return nil, nil
		}
		return &UpdateConfig{
			IsUpdateAvailable: true,
			UpdateURL:         durl,
			ReleaseNotes:      strings.Replace(strings.TrimSpace(gitdata["body"].(string)), "\n", "\n\t", -1),
		}, nil

	}
	return nil, nil
}

// WriteUpdateConfig will write the update info
func WriteUpdateConfig(updateConfig *UpdateConfig) {
	if updateConfig == nil {
		updateConfig = &UpdateConfig{
			IsUpdateAvailable: false,
			UpdateURL:         "",
			ReleaseNotes:      "",
		}
	}
	viper.Set("isUpdateAvailable", updateConfig.IsUpdateAvailable)
	viper.Set("updateURL", updateConfig.UpdateURL)
	viper.Set("releaseNotes", updateConfig.ReleaseNotes)
	viper.WriteConfig()
}

// PrintUpdateInfo will print the update info
func PrintUpdateInfo(updateConfig *UpdateConfig) {
	if updateConfig == nil {
		return
	}
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

// DownloadFile will download the file specified by the URL and store it at the
// location specified
func DownloadFile(filepath string, url string) (err error) {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

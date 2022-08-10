package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
)

// Copy ...
func Copy(src, dst string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "linux" {
		cmd = exec.Command("sudo", "cp", "-fpv", src, dst)
	} else if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/C", `copy /Y `+src+" "+dst)
		fmt.Println(cmd.Args)
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("sudo", "cp", "-fpv", src, dst)
	}
	return cmd.Run()
}

// Remove will delete the target file
func Remove(target string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "linux" {
		cmd = exec.Command("sudo", "rm", target)
	} else if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/C", `del `+target)
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("sudo", "rm", target)
	}
	return cmd.Run()
}

// CheckFileExists ...
func CheckFileExists(path string) (bool, error) {
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

// MakeAppExecutable will make the app executable
func MakeAppExecutable(app string) {
	fmt.Println("#> Making app executable...")
	err := os.Chmod(app, 500)
	if err != nil {
		log.Fatalln("# Error: ERR_MAKE_APP_EXEC", err)
	}
}

// HandleYNInput handles the Yes/No input
func HandleYNInput() bool {
	reader := bufio.NewReader(os.Stdin)
	inputBytes, _, err := reader.ReadLine()
	if err != nil {
		log.Fatalln("# Error: ERR_READ_USRIN", err)
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
		log.Fatalln("# Error: ERR_PARSE_BOOL", err)
	}
	return choice
}

// RunFlogoApp will run the app
func RunFlogoApp(app string, debug bool, args []string) {
	cmd := exec.Command(app, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	fmt.Printf("\n#> Executing: %s\n\n", strings.Join(cmd.Args, " "))
	if debug {
		cmd = EnableDebugLogs(cmd)
	}
	err := cmd.Run()
	if err != nil {
		log.Fatalln("# Error: ERR_RUN_FA", err)
	}
}

package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"runtime"
)

// Copy ...
func Copy(src, dst string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "linux" {
		cmd = exec.Command("sudo", "cp", "-fpv", src, dst)
	} else if runtime.GOOS == "windows" {
		// TODO
	} else if runtime.GOOS == "darwin" {
		// TODO
	}
	return cmd.Run()
}

// Remove will delete the target file
func Remove(target string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "linux" {
		cmd = exec.Command("sudo", "rm", target)
	} else if runtime.GOOS == "windows" {
		// TODO
	} else if runtime.GOOS == "darwin" {
		// TODO
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
func HandleYNInput() rune {
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		log.Fatalln("# Error: ERR_READ_USRIN", err)
	}
	return char
}

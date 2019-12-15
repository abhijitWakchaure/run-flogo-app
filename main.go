package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"sort"
)

// Constants for local env
const (
	AppDir     = "/home/abhijit/Downloads"
	AppPattern = `^.+-linux_amd64.*$`
)

// EnableDebugLogs will add env variable to enable debug logs
func EnableDebugLogs(cmd *exec.Cmd) *exec.Cmd {
	debugFlag := `FLOGO_LOG_LEVEL=DEBUG`
	cmd.Env = append(os.Environ(), debugFlag)
	return cmd
}

// FindLatestApp will return the latest flogo app name
func FindLatestApp() string {
	files, err := ioutil.ReadDir(AppDir)
	if err != nil {
		log.Fatal(err)
	}
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})
	validApp := regexp.MustCompile(AppPattern)
	for _, f := range files {
		if !f.IsDir() && validApp.MatchString(f.Name()) {
			return path.Join(AppDir, f.Name())
		}
	}
	log.Println("# Info: No flogo apps found in " + AppDir)
	return ""
}

// MakeAppExecutable will make the app executable
func MakeAppExecutable(app string) {
	err := os.Chmod(app, 500)
	if err != nil {
		log.Fatalln("# Error:", err)
	}
}

// RunFlogoApp will run the app
func RunFlogoApp(app string, args []string) {
	cmd := exec.Command(app)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if len(args) > 1 {
		cmd = EnableDebugLogs(cmd)
	}
	err := cmd.Run()
	if err != nil {
		log.Fatalln("# Error:", err)
	}
}

func main() {
	app := FindLatestApp()
	if len(app) > 0 {
		fmt.Println(app)
		if runtime.GOOS == "windows" {
			// TODO: Handle for Windows
		}
		MakeAppExecutable(app)
		RunFlogoApp(app, os.Args)
	}
}

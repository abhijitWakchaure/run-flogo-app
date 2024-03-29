package files

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/abhijitWakchaure/run-flogo-app/config"
	"github.com/abhijitWakchaure/run-flogo-app/software"
)

// FindLatestApp will return the latest flogo app name
func FindLatestApp(dir, pattern string) string {
	fmt.Printf("#> Finding latest app inside apps dir [%s]...\n", dir)
	files := listAndSort(dir)
	validApp := regexp.MustCompile(pattern)
	for _, f := range files {
		if !f.IsDir() && validApp.MatchString(f.Name()) {
			return filepath.Join(dir, f.Name())
		}
	}
	fmt.Println("#> No flogo apps found in " + dir)
	return ""
}

// FindAppsWithName will return the list of matching flogo apps
func FindAppsWithName(dir, pattern, name string) []string {
	fmt.Printf("#> Searching apps with name containing '%s' inside apps dir [%s]...\n", name, dir)
	var apps []string
	name = strings.ToLower(name)
	files := listAndSort(dir)
	validApp := regexp.MustCompile(pattern)
	for _, f := range files {
		if !f.IsDir() && validApp.MatchString(f.Name()) && strings.Contains(strings.ToLower(f.Name()), name) {
			apps = append(apps, filepath.Join(dir, f.Name()))
		}
	}
	return apps
}

// ListLastNApps will return the list of last 'N' flogo apps
func ListLastNApps(dir, pattern string, n int) []string {
	fmt.Printf("#> Listing last %d apps inside apps dir [%s]...\n", n, dir)
	files := listAndSort(dir)
	var apps []string
	validApp := regexp.MustCompile(pattern)
	for _, f := range files {
		if !f.IsDir() && validApp.MatchString(f.Name()) {
			apps = append(apps, filepath.Join(dir, f.Name()))
			if len(apps) == n {
				return apps
			}
		}
	}
	return apps
}

// DeleteApps will delete all the flogo apps in apps dir
func DeleteApps(dir, pattern string) {
	fmt.Printf("#> Listing all the flogo apps inside apps dir [%s]...\n", dir)
	files := listAndSort(dir)
	validApp := regexp.MustCompile(pattern)
	var count int
	apps := []string{}
	for i, f := range files {
		if !f.IsDir() && validApp.MatchString(f.Name()) {
			apps = append(apps, filepath.Join(dir, f.Name()))
			count++
			fmt.Printf("%d. %s\n", i+1, filepath.Join(dir, f.Name()))
		}
	}
	if count == 0 {
		fmt.Println("#> No flogo app found inside apps dir.")
		os.Exit(0)
	}
	fmt.Printf("\nAre you sure you want to delete all %d app(s)? [y/n] ", count)
	choice := software.HandleYNInput()
	var err error
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

func listAndSort(dir string) []fs.FileInfo {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("\n#> Failed to read apps dir [%s]! Error %s\n", dir, err.Error())
		os.Exit(1)
	}
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})
	var nFiles []fs.FileInfo
	selfName := fmt.Sprintf("%s-%s_%s", config.AppName, runtime.GOOS, runtime.GOARCH)
	for _, f := range files {
		if !strings.Contains(f.Name(), selfName) {
			nFiles = append(nFiles, f)
		}
	}
	return nFiles
}

// Copyright (c) 2019 abhijit wakchaure. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

/*
run-flogo-app
	This program will run the latest TIBCO Flogo Enterprise app in the directory specified by you.
	By default it will pickup the latest Flogo app from your Downloads directory.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/abhijitWakchaure/run-flogo-app/core"
	"github.com/abhijitWakchaure/run-flogo-app/utils"
)

// NewApp ...
func NewApp() *core.App {
	app := new(core.App)
	app.Init()
	return app
}

// RunFlogoApp will run the app
func RunFlogoApp(app string, debug *bool, tail []string) {
	cmd := exec.Command(app, tail...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	fmt.Println("#> Executing: ", strings.Join(cmd.Args, " "))
	if *debug {
		cmd = utils.EnableDebugLogs(cmd)
	}
	err := cmd.Run()
	if err != nil {
		log.Fatalln("# Error: ERR_RUN_FA", err)
	}
}

func main() {
	app := NewApp()
	flagDebug := flag.Bool("debug", false, "Set this to enable debug logs")
	flagInstall := flag.Bool("install", false, "Set this to install the program")
	flagUninstall := flag.Bool("uninstall", false, "Set this to uninstall the program")
	flagVersion := flag.Bool("version", false, "Prints the current version of the program")
	flag.Parse()

	if *flagInstall {
		app.Install()
	} else if *flagUninstall {
		app.Uninstall()
	} else if *flagVersion {
		app.Version()
	} else {
		app.Main()
		latestFlogoApp := app.FindLatestApp()
		if len(latestFlogoApp) > 0 {
			fmt.Print("#> Do you want to execute this app \"", latestFlogoApp, "\" [Y/n]: ")
			choice := utils.HandleYNInput()
			if choice == 'Y' || choice == 'y' {
				utils.MakeAppExecutable(latestFlogoApp)
				RunFlogoApp(latestFlogoApp, flagDebug, flag.Args())
			} else {
				log.Println("# Info: Exiting...")
			}
		}
	}
}

// Copyright © 2022 ABHIJIT WAKCHAURE
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

/*
run-flogo-app
	This program will run the latest TIBCO Flogo Enterprise app in the directory specified by you.
	By default it will pickup the latest Flogo app from your Downloads directory.
*/
package main

import (
	"github.com/abhijitWakchaure/run-flogo-app/cmd"
)

func main() {
	cmd.Execute()
}

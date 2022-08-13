//go:build docs
// +build docs

package cmd

import (
	"fmt"
)

func init() {
	fmt.Println("Initializing docs...")
	GENDOCS = true
}

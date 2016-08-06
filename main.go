// Package main provides entry for the command line tool
package main

import (
	"github.com/shanzi/wu/runner"
	"os"
	"path/filepath"
)

func main() {
	var path string
	if len(os.Args) > 1 {
		path, _ = filepath.Abs(os.Args[1])
	} else {
		path, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}
	r := runner.NewRunner(path, []string{}, "")
	r.Start()
}

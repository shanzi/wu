// Package main provides entry for the command line tool
package main

import (
	"flag"
	"github.com/shanzi/wu/runner"
	"os"
	"path/filepath"
)

var path = flag.String("p", filepath.Dir(os.Args[0]), "Path to the file or directory to watch")

func main() {
	flag.Parse()

	patterns := []string{"*"}
	if flag.NArg() > 0 {
		patterns = flag.Args()
	}
	abspath, _ := filepath.Abs(*path)

	r := runner.NewRunner(abspath, patterns, "")
	r.Start()
}

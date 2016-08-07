// Package main provides entry for the command line tool
package main

import (
	"flag"
	"github.com/shanzi/wu/command"
	"github.com/shanzi/wu/runner"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
)

var path = flag.String("dir", filepath.Dir(os.Args[0]), "Path of the directory to watch")
var pattern = flag.String("pattern", "*", "Patterns as a filter of filenames")

func init() {
	log.SetFlags(0) // Turn off date and time on standard logger
}

func main() {
	flag.Parse()

	abspath, _ := filepath.Abs(*path)
	patterns := parsePattern(*pattern)
	cmd := parseCommand(flag.Args())

	r := runner.New(abspath, patterns, cmd)

	go func() {
		// Handle interrupt signal
		ch := make(chan os.Signal)
		signal.Notify(ch, os.Interrupt)

		<-ch
		log.Println()
		log.Println("Shutting down...")
		r.Exit()
	}()

	r.Start()
}

func parsePattern(pat string) []string {
	patternSep, _ := regexp.Compile("[,\\s]+")
	return patternSep.Split(pat, -1)
}

func parseCommand(flagArgs []string) command.Command {
	if len(flagArgs) == 0 {
		return command.Empty()
	}

	name := flagArgs[0]
	args := flagArgs[1:]

	return command.New(name, args...)
}

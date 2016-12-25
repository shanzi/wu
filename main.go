// Package main provides entry for the command line tool
package main

import (
	"github.com/shanzi/wu/command"
	"github.com/shanzi/wu/runner"

	"log"
	"os"
	"os/signal"
	"path/filepath"
	// "wu/command"
	// "wu/runner"
)

func init() {
	log.SetFlags(0) // Turn off date and time on standard logger
}

func main() {
	conf := getConfigs()

	abspath, _ := filepath.Abs(conf.Directory)
	patterns := conf.Patterns
	cmd := command.New(conf.Command)

	r := runner.New(abspath, patterns, cmd)

	go func() {
		// Handle interrupt signal
		ch := make(chan os.Signal)
		signal.Notify(ch, os.Interrupt)

		<-ch
		r.Exit()
	}()

	r.Start()
}

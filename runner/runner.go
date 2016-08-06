// Package manager provides manager for running watch and exec loop
package runner

import (
	"log"
	"time"
)

type Runner interface {
	Path() string
	Patterns() []string
	Command() string
	Start()
	Exit()
}

type runner struct {
	path     string
	patterns []string
	command  string

	abort chan struct{}
}

func NewRunner(path string, patterns []string, command string) Runner {
	return &runner{
		path:     path,
		patterns: patterns,
		command:  command,
	}
}

func (r *runner) Path() string {
	return r.path
}

func (r *runner) Patterns() []string {
	return r.patterns
}

func (r *runner) Command() string {
	return r.command
}

func (r *runner) Start() {
	r.abort = make(chan struct{})
	changed, err := watch(r.path, r.abort)
	if err != nil {
		log.Fatal("Failed to initialize watcher: ", err)
	}
	matched := match(changed, r.patterns)
	for fp := range matched {
		files := gather(fp, matched, 500*time.Millisecond)
		log.Println(files)
	}
}

func (r *runner) Exit() {
	r.abort <- struct{}{}
	close(r.abort)
}

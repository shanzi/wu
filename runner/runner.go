// Package manager provides manager for running watch and exec loop
package runner

import (
	"github.com/fsnotify/fsnotify"
	"log"
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

	for fp := range changed {
		log.Println(fp)
	}
}

func (r *runner) Exit() {
	r.abort <- struct{}{}
	close(r.abort)
}

func watch(path string, abort <-chan struct{}) (<-chan string, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = watcher.Add(path)
	if err != nil {
		return nil, err
	}

	out := make(chan string)
	go func() {
		defer close(out)
		defer watcher.Close()
		for {
			select {
			case <-abort:
				// Abort watching
				return
			case fp := <-watcher.Events:
				out <- fp.String()
			case err := <-watcher.Errors:
				log.Println("Watch Error:", err)
			}
		}
	}()

	return out, nil
}

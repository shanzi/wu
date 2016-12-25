package command

import (
	"time"
)

type empty string

func Empty() Command {
	return empty("Empty command")
}

func (e empty) String() string {
	return string(e)
}

func (e empty) Start(delay time.Duration, haveBuild bool) {
	// Start an empty command just do nothing but delay for given duration
	<-time.After(delay)
}

func (e empty) Terminate(wait time.Duration) {
	// Terminate empty command just return immediately without any error
}

// ProcessState contains information about an exited process,
// available after a call to Wait or Run.

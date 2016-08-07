package command

import "time"

// An empty command is a command that do nothing
type empty string

func Empty() Command {
	return empty("Empty command")
}

func (c empty) String() string {
	return string(c)
}

func (c empty) Start(delay time.Duration) {
	// Start an empty command just do nothing but delay for given duration
	<-time.After(delay)
}

func (c empty) Terminate(wait time.Duration) {
	// Terminate empty command just return immediately without any error
}

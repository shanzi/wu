// Package command provides a wrap over os/exec for easier command handling
package command

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Command interface {
	String() string
	Start(delay time.Duration)
	Terminate(wait time.Duration)
}

type command struct {
	name  string
	args  []string
	cmd   *exec.Cmd
	mutex *sync.Mutex
}

func New(name string, args ...string) Command {
	return &command{
		name,
		args,
		nil,
		&sync.Mutex{},
	}
}

func (c *command) String() string {
	return fmt.Sprintf("%s %s", c.name, strings.Join(c.args, " "))
}

func (c *command) Start(delay time.Duration) {
	<-time.After(delay) // delay for a while to avoid start too frequently

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.cmd != nil && !c.cmd.ProcessState.Exited() {
		log.Fatalln("Failed to start command: previous command hasn't exit.")
	}

	cmd := exec.Command(c.name, c.args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout // Redirect stderr of sub process to stdout of parent

	log.Println("Running command:", c.String())
	log.Println()

	err := cmd.Start()
	if err != nil {
		log.Println("Failed:", err)
	} else {
		c.cmd = cmd
		go func() {
			cmd.Wait()
			if cmd.ProcessState.Success() {
				log.Println("Command exited")
			} else {
				log.Println("Command terminated")
			}
		}()
	}
}

func (c *command) Terminate(wait time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// set c.cmd to nil after finished
	defer func() {
		c.cmd = nil
	}()

	if c.cmd == nil {
		// No command is runing, just return
		return
	}

	if c.cmd.ProcessState != nil && c.cmd.ProcessState.Exited() {
		// Command has exited, just return
		return
	}

	cmd := c.cmd
	// Try to stop the process by sending a SIGINT signal
	cmd.Process.Signal(os.Interrupt)

	ch := make(chan struct{})
	go func() {
		cmd.Wait()
		ch <- struct{}{}
	}()

	select {
	case <-ch:
		return
	case <-time.After(wait):
		log.Println("Try to kill subprocess by force")
		cmd.Process.Kill()
	}
}

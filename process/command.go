package process

import (
	"bus/buffer"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	Path   string
	Stdin  buffer.Reader
	Stdout buffer.Writer
	Stderr buffer.Writer
	State  *os.ProcessState
	Pid    int

	value   []string
	current *exec.Cmd
}

func (c *Command) Start() error {
	name := c.value[0]
	args := c.value[1:]
	if name == "sh" {
		name = os.Getenv("SHELL")
		args = []string{
			"-c", strings.Join(args, " "),
		}
	}
	c.current = exec.Command(name, args...)
	c.current.Env = os.Environ()
	c.current.Dir = c.Path
	c.current.Stdin = c.Stdin
	c.current.Stdout = c.Stdout
	c.current.Stderr = c.Stderr

	err := c.current.Start()
	if err != nil {
		return err
	}
	c.Pid = c.current.Process.Pid
	return nil
}

func (c *Command) Wait() error {
	var err error = nil
	c.State, err = c.current.Process.Wait()
	if err != nil {
		return err
	}
	for !c.State.Exited() {
		c.State, err = c.current.Process.Wait()
		if err != nil {
			return err
		}
	}
	return err
}

func (c *Command) Kill() error {
	return c.current.Process.Kill()
}

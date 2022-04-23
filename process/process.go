package process

import (
	"bus/buffer"
	"errors"
	"os"
)

type Process struct {
	Path    string
	Pid     int
	Restart bool
	Stdin   buffer.Reader
	Stdout  buffer.Writer
	Stderr  buffer.Writer
	cmd     *Command

	name          string
	commandString string
}

func NewProcess(name string, command string) *Process {
	processBuffer := buffer.NewBuffer(name)
	processErrBuffer := buffer.NewBuffer(name + ":error")
	processErrBuffer.Output = os.Stderr

	return &Process{
		Path:          "",
		Pid:           -1,
		Restart:       false,
		name:          name,
		commandString: command,

		cmd:    nil,
		Stdin:  &processBuffer,
		Stdout: &processBuffer,
		Stderr: &processErrBuffer,
	}
}

func (p *Process) UseStandardIO() {
	p.Stdin = os.Stdin
	p.Stdout = os.Stdout
	p.Stderr = os.Stderr
}

func (p *Process) Mute() {
	p.Stdout = nil
	p.Stderr = nil
}

func (p *Process) Create() {
	// fmt.Println(strings.Split(p.commandString, " "), split(p.commandString))
	p.cmd = &Command{
		Path:   p.Path,
		Stdin:  p.Stdin,
		Stdout: p.Stdout,
		Stderr: p.Stderr,

		value: split(p.commandString),
	}
}

func (p *Process) Run() error {
	if p.cmd == nil {
		p.Create()
	}
	err := p.cmd.Start()
	if err != nil {
		return err
	}
	p.Pid = p.cmd.Pid
	return nil
}

func (p *Process) Wait() error {
	err := p.cmd.Wait()
	if p.cmd.State.ExitCode() == 0 || p.Restart {
		return err
	}
	err = p.Run()
	if err != nil {
		return err
	}
	err = p.Wait()
	return err
}

func (p *Process) Kill() error {
	if p.cmd == nil {
		return errors.New("you can't kill a process if you haven't created one before")
	}
	return p.cmd.Kill()
}

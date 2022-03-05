package process

import (
	"bus/buffer"
	"errors"
	"os/exec"
)

type Process struct {
	Pid    int
	Daemon bool
	Stdin  buffer.Reader
	Stdout buffer.Writer
	Stderr buffer.Writer
	cmd    *exec.Cmd

	name string
	file string
	args []string
}

func NewProcess(name string, file string, args ...string) Process {
	processBuffer := buffer.NewBuffer(name)
	processErrBuffer := buffer.NewBuffer(name + ":error")

	return Process{
		Pid:    -1,
		Daemon: false,
		name:   name,
		file:   file,
		args:   args,

		cmd:    nil,
		Stdin:  &processBuffer,
		Stdout: &processBuffer,
		Stderr: &processErrBuffer,
	}
}

func (p *Process) Create() {
	p.cmd = exec.Command(p.file, p.args...)

	p.cmd.Stdin = p.Stdin
	p.cmd.Stdout = p.Stdout
	p.cmd.Stderr = p.Stderr
}

func (p *Process) Run() error {
	if p.cmd == nil {
		p.Create()
	}
	err := p.cmd.Start()
	if err != nil {
		return err
	}
	p.Pid = p.cmd.Process.Pid
	return nil
}

func (p *Process) Wait() error {
	state, err := p.cmd.Process.Wait()
	for !state.Exited() {
		state, err = p.cmd.Process.Wait()
	}
	if state.ExitCode() == 0 || !p.Daemon {
		return err
	}
	p.Create()
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
	return p.cmd.Process.Kill()
}

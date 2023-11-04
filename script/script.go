package script

import (
	"bus/buffer"
	"bus/config"
	"bus/helper"
	"bus/process"
	"syscall"
)

type MuteKind uint8

const (
	NoMuted MuteKind = iota
	PatialMuted
	RunMuted
	Muted
)

type Script struct {
	DryRun     bool
	process    *process.Process
	pathConfig *config.Package
	cmd        string
	muteLvl    MuteKind
}

func NewScript(pathConfig *config.Package, absPath, cmd string) *Script {
	muteLvl := NoMuted
	if cmd[0] == '~' || cmd[0] == '-' || cmd[0] == '=' {
		if cmd[0] == '~' {
			muteLvl = RunMuted
		} else if cmd[0] == '=' {
			muteLvl = PatialMuted
		} else {
			muteLvl = Muted
		}
		cmd = cmd[1:]
	}
	p := process.NewProcess(pathConfig.Name, cmd)
	if muteLvl == RunMuted || muteLvl == Muted {
		p.Mute()
	}
	p.Restart = true
	p.Path = absPath

	return &Script{
		DryRun:     false,
		process:    p,
		cmd:        cmd,
		pathConfig: pathConfig,
		muteLvl:    muteLvl,
	}
}

func (s *Script) Start(done func()) error {
	defer done()

	if s.muteLvl != PatialMuted && s.muteLvl != Muted || s.DryRun {
		buffer.Printf("%v$: %v%v\n", helper.Cyan, helper.Bold+s.cmd, helper.Reset)
	}

	if s.DryRun {
		return nil
	}

	err := s.process.Run()
	if err != nil {
		buffer.Println(err)
		syscall.Exit(1)
	}
	return s.process.Wait()
}

func (s *Script) Kill() {
	s.process.Kill()
}

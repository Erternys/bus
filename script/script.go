package script

import (
	"bus/config"
	"bus/helper"
	"bus/process"
	"fmt"
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
	p.UseStandardIO()
	if muteLvl == RunMuted || muteLvl == Muted {
		p.Mute()
	}
	p.Daemon = true
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
		fmt.Printf("%v$: %v%v\n", helper.Cyan, helper.Bold+s.cmd, helper.Reset)
	}

	if s.DryRun {
		return nil
	}

	err := s.process.Run()
	if err != nil {
		fmt.Println(err)
		syscall.Exit(1)
	}
	return s.process.Wait()
}

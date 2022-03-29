package script

import (
	"bus/middleware"
	"bus/process"
	"fmt"
	"syscall"
)

type Script struct {
	process    *process.Process
	pathConfig *middleware.Package
	cmd        string
}

func NewScript(pathConfig *middleware.Package, absPath, cmd string) *Script {
	p := process.NewProcess(pathConfig.Name, cmd)
	p.Daemon = true
	p.Path = absPath

	return &Script{
		process:    p,
		cmd:        cmd,
		pathConfig: pathConfig,
	}
}

func (s *Script) Start(done func()) error {
	defer done()

	err := s.process.Run()
	if err != nil {
		fmt.Println(err)
		syscall.Exit(1)
	}
	return s.process.Wait()
}

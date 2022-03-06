package main

import (
	"bus/process"
)

func main() {
	p := process.NewProcess(
		"echo",
		"echo test 'test'",
	)
	err := p.Run()
	if err != nil {
		panic(err)
	}
	err = p.Wait()
	if err != nil {
		panic(err)
	}
}

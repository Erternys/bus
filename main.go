package main

import (
	"bus/process"
)

func main() {
	p := process.NewProcess(
		"echo",
		"echo",
		"Hello\nWorld!",
	)
	p.Create()
	p.Run()
	p.Wait()
}

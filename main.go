package main

import (
	"bus/buffer"
	"log"
	"os"
	"os/exec"
)

func main() {
	buffer := buffer.NewBuffer("echo")
	cmd := exec.Command("echo", "Hello\nWorld!")

	cmd.Stdout = &buffer
	cmd.Stderr = os.Stderr
	err := cmd.Start()

	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	cmd.Wait()
}

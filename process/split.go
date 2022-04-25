package process

import (
	"bus/helper"
	"fmt"
	"strings"
	"syscall"
)

func split(command string) []string {
	args := []string{}
	current := ""
	str := false
	command = strings.TrimSpace(command)

	for i, c := range command {
		if c == ' ' && !str {
			args = append(args, current)
			current = ""
			continue
		}
		if c == '"' && command[i-1] != '\\' {
			if str {
				args = append(args, current)
				current = ""
			}
			str = !str
			continue
		}
		current += string(c)
	}

	if str {
		fmt.Printf("%vInvalid command, please fix it%v", helper.Red+helper.Bold, helper.Reset)
		syscall.Exit(1)
	}
	if current != "" {
		args = append(args, current)
		current = ""
	}

	return args
}

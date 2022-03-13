package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

func input(prompt string, defaultValue string) string {
	fmt.Print(prompt)

	value, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	value = strings.TrimFunc(value, unicode.IsSpace)

	if value == "" {
		return defaultValue
	}
	return value
}

func getwd() string {
	pwd, _ := os.Getwd()
	arrPwd := strings.Split(pwd, "/")

	return arrPwd[len(arrPwd)-1]
}

func getRepository() string {
	out, _ := exec.Command("git", "config", "--get", "remote.origin.url").Output()

	return strings.TrimFunc(string(out), unicode.IsSpace)
}

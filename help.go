package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func input(prompt string, defaultValue string) string {
	fmt.Print(prompt)

	var value string
	fmt.Scan(&value)

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

	return string(out)
}

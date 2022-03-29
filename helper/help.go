package helper

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

func Input(prompt string, defaultValue string) string {
	fmt.Print(prompt)

	value, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	value = strings.TrimFunc(value, unicode.IsSpace)

	if value == "" {
		return defaultValue
	}
	return value
}

func Getwd() string {
	pwd, _ := os.Getwd()
	arrPwd := strings.Split(pwd, string(os.PathSeparator))

	return arrPwd[len(arrPwd)-1]
}

func GetRepository() string {
	out, _ := exec.Command("git", "config", "--get", "remote.origin.url").Output()

	return strings.TrimFunc(string(out), unicode.IsSpace)
}
package helper

import (
	"bufio"
	"bus/buffer"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"unicode"
)

var WorkSpacePath, _ = os.Getwd()

func Input(prompt string, defaultValue string) string {
	buffer.Print(prompt)

	value, _ := bufio.NewReader(&buffer.Stdin).ReadString('\n')
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

	return strings.TrimSpace(string(out))
}

func FindArray[T any](val T, array []T) bool {
	values := reflect.ValueOf(array)

	if reflect.TypeOf(array).Kind() == reflect.Slice || values.Len() > 0 {
		for i := 0; i < values.Len(); i++ {
			if reflect.DeepEqual(val, values.Index(i).Interface()) {
				return true
			}
		}
	}

	return false
}

package cli

import (
	"reflect"
	"strings"
)

type CliApp struct {
	Name        string
	Description string
	Version     string

	args     []string
	commands []Command
}

func inArray(val interface{}, array interface{}) bool {
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

func NewApp(name string, description string, version string, commands ...Command) CliApp {
	return CliApp{
		Name:        name,
		Description: description,
		Version:     version,

		args:     []string{},
		commands: commands,
	}
}

func (c *CliApp) Run(args []string) error {
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			continue
		}
		for _, command := range c.commands {
			if command.Name == arg || inArray(arg, command.Aliases) {
				command.Cli = c
			}
		}
	}

	return nil
}

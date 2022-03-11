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
	flags    []Flag
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

func NewApp(name string, description string, version string) CliApp {
	return CliApp{
		Name:        name,
		Description: description,
		Version:     version,

		args:     []string{},
		commands: []Command{},
		flags:    []Flag{},
	}
}

func (c *CliApp) AddCommand(command Command) {
	c.commands = append(c.commands, command)
}
func (c *CliApp) AddFlag(flag Flag) {
	c.flags = append(c.flags, flag)
}

func (c *CliApp) Run(args []string) error {
	var currentCommand *Command = nil
	context := NewContext(c)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			arg = strings.TrimLeft(arg, "-")
			currentFlag := DefaultFlag()
			for _, flag := range c.flags {
				if flag.Name == arg || inArray(arg, flag.Aliases) {
					currentFlag = flag.clone()
				}
			}
			if strings.Contains(arg, "=") {
				slice := strings.Split(arg, "=")
				name, value := slice[0], strings.Join(slice[1:], "=")
				if currentFlag.Name == "" {
					currentFlag.setValueAndKind(value)
				} else {
					err := currentFlag.SetValue(value)
					if err != nil {
						return err
					}
				}
				currentFlag.Name = name
			} else if len(args) > i+2 && args[i+1] == "=" {
				if currentFlag.Name == "" {
					currentFlag.setValueAndKind(args[i+2])
				} else {
					err := currentFlag.SetValue(args[i+2])
					if err != nil {
						return err
					}
				}
				currentFlag.Name = arg
				i += 2
			} else if len(args) > i+1 {
				if currentFlag.Name == "" {
					currentFlag.setValueAndKind(args[i+1])
				} else {
					err := currentFlag.SetValue(args[i+1])
					if err != nil {
						return err
					}
				}
				currentFlag.Name = arg
				i++
			} else {
				currentFlag.Name = arg
				currentFlag.Kind = Bool
				currentFlag.Value = true
			}
			context.Flags[currentFlag.Name] = currentFlag
			continue
		}
		if currentCommand == nil {
			for _, command := range c.commands {
				if command.Name == arg || inArray(arg, command.Aliases) {
					currentCommand = &command
				}
			}
		} else {
			context.Args = append(context.Args, arg)
		}
	}

	currentCommand.Handle(context)

	return nil
}

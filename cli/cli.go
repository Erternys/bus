package cli

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type CliApp struct {
	Name        string
	Description string
	Version     string

	defaultCommand string
	args           []string
	commands       []Command
	flags          []Flag
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

		defaultCommand: "help",
		args:           []string{},
		commands:       []Command{},
		flags:          []Flag{},
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
	var err error = nil
	context := NewContext(c)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			arg = strings.TrimLeft(arg, "-")
			currentFlag := DefaultFlag()
			for _, flag := range c.flags {
				if flag.Name == arg || inArray(arg, flag.Aliases) {
					fmt.Println(flag)
					currentFlag = flag.clone()
					break
				}
			}
			if currentFlag.Name == "" {
				for _, command := range c.commands {
					if inArray(arg, command.FlagAliases) {
						if currentCommand != nil {
							context.Args = append(context.Args, currentCommand.Name)
						}
						currentCommand = &command
						break
					}
				}
				if currentFlag.Name != "" {
					continue
				}
			}
			if strings.Contains(arg, "=") {
				slice := strings.Split(arg, "=")
				name, value := slice[0], strings.Join(slice[1:], "=")
				if currentFlag.Name == "" {
					currentFlag.Name = name
					currentFlag.setValueAndKind(value)
				} else {
					err = currentFlag.SetValue(value)
				}
			} else if len(args) > i+2 && args[i+1] == "=" {
				if currentFlag.Name == "" {
					currentFlag.Name = arg
					currentFlag.setValueAndKind(args[i+2])
				} else {
					err = currentFlag.SetValue(args[i+2])
				}
				i += 2
			} else if len(args) > i+1 {
				if currentFlag.Name == "" {
					currentFlag.Name = arg
					currentFlag.setValueAndKind(args[i+1])
				} else {
					err = currentFlag.SetValue(args[i+1])
					if err != nil && strings.Contains(err.Error(), "parsing") {
						currentFlag.Value = true
						err = nil
						i--
					}
				}
				i++
			} else {
				currentFlag.Name = arg
				currentFlag.Kind = Bool
				currentFlag.Value = true
			}
			fmt.Println(currentFlag.Aliases)
			context.Flags[currentFlag.Name] = currentFlag
			continue
		}
		if currentCommand == nil {
			for _, command := range c.commands {
				if command.Name == arg || inArray(arg, command.Aliases) {
					currentCommand = &command
					break
				}
			}
			if currentCommand == nil {
				context.Args = append(context.Args, arg)
			}
		} else {
			context.Args = append(context.Args, arg)
		}
	}

	if currentCommand == nil {
		for _, command := range c.commands {
			if command.Name == c.defaultCommand {
				currentCommand = &command
				break
			}
		}
		if currentCommand == nil {
			return errors.New("no command found")
		}
	}

	if currentCommand.RequiredArgs > len(context.Args) {
		err = fmt.Errorf("expected %v arguments, but got %v.", currentCommand.RequiredArgs, len(context.Args))
	}

	currentCommand.Handle(context, err)

	return err
}

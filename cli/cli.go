package cli

import (
	"bus/helper"
	"bus/process"
	"errors"
	"fmt"
	"os"
	"path"
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
func (c *CliApp) AddCommandFromExe(command ExeCommand) {
	exe_path := path.Join(os.Args[0], "..", command.File)
	_, err := os.Stat(exe_path)
	if !os.IsNotExist(err) {
		c.commands = append(c.commands, Command{
			Name:        command.Name,
			Flags:       command.Flags,
			Aliases:     command.Aliases,
			FlagAliases: command.FlagAliases,
			Description: command.Description,
			Handle: func(c *Context, err error) {
				p := process.NewProcess(command.Name, exe_path+" "+strings.Join(c.RawArgs, " "))
				p.UseStandardIO()
				p.Run()
				p.Wait()
			},
		})
	}
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
				if flag.Name == arg || helper.InArray(arg, flag.Aliases) {
					currentFlag = flag.clone()
					break
				}
			}
			if currentCommand != nil {
				for _, flag := range currentCommand.Flags {
					if flag.Name == arg || helper.InArray(arg, flag.Aliases) {
						currentFlag = flag.clone()
						break
					}
				}
			}

			if currentFlag.Name == "" {
				for _, command := range c.commands {
					if helper.InArray(arg, command.FlagAliases) {
						if currentCommand == nil {
							context.Args = append(context.Args, currentCommand.Name)
							currentCommand = &command
						}
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
				if currentFlag.Name == "" {
					currentFlag.Name = arg
				}
				currentFlag.Kind = Bool
				currentFlag.Value = true
			}
			context.Flags[currentFlag.Name] = currentFlag
			raw := currentFlag.ToString()
			if len(raw) > 0 {
				context.RawArgs = append(context.RawArgs, raw)
			}
			continue
		}
		if currentCommand == nil {
			for _, command := range c.commands {
				if command.Name == arg || helper.InArray(arg, command.Aliases) {
					currentCommand = &command
					break
				}
			}
			if currentCommand == nil {
				context.Args = append(context.Args, arg)
				context.RawArgs = append(context.RawArgs, arg)
			}
		} else {
			context.Args = append(context.Args, arg)
			context.RawArgs = append(context.RawArgs, arg)
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

package cli

import (
	"bus/buffer"
	"bus/helper"
	"fmt"
	"sort"
	"strings"
)

type HandleAction func(c *Context, err error)

type Command struct {
	Name             string
	Flags            []Flag
	Aliases          []string
	FlagAliases      []string
	Description      string
	ShortDescription string
	Usage            string
	RequiredArgs     int
	Handle           HandleAction
}

func (c *Command) AddFlag(flag Flag) {
	c.Flags = append(c.Flags, flag)
}

func (c *CliApp) SetHelpCommand() {
	c.AddCommand(Command{
		Name:             "help",
		Aliases:          []string{"h"},
		FlagAliases:      []string{"help", "h"},
		Description:      "Print help information in the console",
		ShortDescription: "Print help information",
		Usage:            "help [Flags] [Subcommand]",
		RequiredArgs:     0,
		Handle: func(c *Context, _ error) {
			message := c.App.Name + " " + c.App.Version + "\n" + c.App.Description + "\n"
			output := &buffer.Stdout
			buffer.Printf("\n\n%v\n\n", c.Args)
			if len(c.Args) == 0 {
				// Print principal help message
				message += fmt.Sprintf("\nUsage:\n    %v\n", c.Command.Usage)

				message += "\nFlags (Options):\n"
				sort.Slice(c.App.flags, func(i, j int) bool {
					return c.App.flags[i].Name < c.App.flags[j].Name
				})
				for _, flag := range c.App.flags {
					aliases := " "
					if len(flag.Aliases) > 0 {
						aliases += "(-" + strings.Join(flag.Aliases, ", -") + ")"
					}
					message += fmt.Sprintf("    --%-14v %v\n", flag.Name+aliases, flag.Description)
				}
				message += "\nCommands:\n"
				sort.Slice(c.App.commands, func(i, j int) bool {
					return c.App.commands[i].Name < c.App.commands[j].Name
				})
				for _, command := range c.App.commands {
					aliases := " "
					if len(command.Aliases) > 0 {
						aliases += "(" + strings.Join(command.Aliases, ", ") + ")"
					}
					message += fmt.Sprintf("    %-16v %v\n", command.Name+aliases, command.ShortDescription)
				}
			} else {
				// Print the help message of a specific command
				command := c.App.GetCommand(c.Args[0])

				if command == nil {
					output = &buffer.Stderr
					message += fmt.Sprintf("\nError:\n    Unknown command: \"%v\"\n To see a list of supported commands, run:\n%v%v help%v\n", c.Args[1], helper.Bold, c.App.Name, helper.Reset)
					goto end
				}

				message += fmt.Sprintf("\nUsage:\n    %v\n", command.Usage)
				message += fmt.Sprintf("\nDescription:\n    %v\n", command.Description)
			}

		end:
			message += fmt.Sprintf("\nSee '%v help <command>' for more information on a specific command.\n", c.App.Name)

			buffer.Fprintln(output, message)
		},
	})
}

func (c *CliApp) SetVersionCommand() {
	c.AddCommand(Command{
		Name:             "version",
		Aliases:          []string{"v"},
		FlagAliases:      []string{"version", "v"},
		Description:      "Print current version of " + c.Name,
		ShortDescription: "Print current version of " + c.Name,
		Usage:            "version",
		RequiredArgs:     0,
		Handle: func(c *Context, _ error) {
			buffer.Println(c.App.Name + " " + c.App.Version)
		},
	})
}

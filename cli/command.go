package cli

import (
	"fmt"
	"sort"
	"strings"
)

type HandleAction func(c *Context, err error)

type Command struct {
	Name         string
	Aliases      []string
	FlagAliases  []string
	Description  string
	RequiredArgs int
	Handle       HandleAction
}

func (c *CliApp) SetHelpCommand() {
	c.AddCommand(Command{
		Name:         "help",
		Aliases:      []string{"h"},
		FlagAliases:  []string{"help", "h"},
		Description:  "Print help information",
		RequiredArgs: 0,
		Handle: func(c *Context, _ error) {
			output := c.App.Description + "\n"

			output += fmt.Sprintf("\nUsage:\n    %v [Flags] [Subcommand]\n", c.App.Name)

			output += "\nFlags (Options):\n"
			sort.Slice(c.App.flags, func(i, j int) bool {
				return c.App.flags[i].Name < c.App.flags[j].Name
			})
			for _, flag := range c.App.flags {
				aliases := " "
				if len(flag.Aliases) > 0 {
					aliases += "(-" + strings.Join(flag.Aliases, ", -") + ")"
				}
				output += fmt.Sprintf("    --%-14v %v\n", flag.Name+aliases, flag.Description)
			}
			output += "\nCommands:\n"
			sort.Slice(c.App.commands, func(i, j int) bool {
				return c.App.commands[i].Name < c.App.commands[j].Name
			})
			for _, command := range c.App.commands {
				aliases := " "
				if len(command.Aliases) > 0 {
					aliases += "(" + strings.Join(command.Aliases, ", ") + ")"
				}
				output += fmt.Sprintf("    %-16v %v\n", command.Name+aliases, command.Description)
			}

			output += fmt.Sprintf("\nSee '%v help <command>' for more information on a specific command.\n", c.App.Name)

			fmt.Println(output)
		},
	})
}

func (c *CliApp) SetVersionCommand() {
	c.AddCommand(Command{
		Name:         "version",
		Aliases:      []string{"v"},
		FlagAliases:  []string{"version", "v"},
		Description:  "Print current version of " + c.Name,
		RequiredArgs: 0,
		Handle: func(c *Context, _ error) {
			fmt.Println(c.App.Name + " " + c.App.Version)
		},
	})
}

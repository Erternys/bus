package new

import (
	"bus/cli"
	"bus/helper"
	"fmt"
	"os"
	"strings"
	"syscall"
)

func checkNameValidity(path string) bool {
	return !strings.ContainsAny(path, "\\?%*:|\"<>/") && path != ""
}

func NewNewCommand() cli.Command {
	return cli.Command{
		Name:             "new",
		RequiredArgs:     0,
		Description:      "Generate a new folder and init the project inside",
		ShortDescription: "Generate a new project",
		Usage:            "new [name]",
		Handle: func(c *cli.Context, _ error) {
			dir := ""
			if len(c.Args) == 0 {
				helper.Input("name of the project: ", "")
			} else {
				dir = c.Args[0]
			}

			if !checkNameValidity(dir) {
				fmt.Println("the folder does not have a correct name")
				syscall.Exit(1)
			}
			os.Mkdir(dir, 0750)
			os.Chdir(dir)

			c.App.CallCommand("init", make([]string, 0))
		},
	}
}

package new

import (
	"bus/cli"
	"fmt"
	"os"
	"strings"
	"syscall"
)

func checkNameValidity(path string) bool {
	return !strings.ContainsAny(path, "\\?%*:|\"<>/")
}

func NewNewCommand() cli.Command {
	return cli.Command{
		Name:         "new",
		RequiredArgs: 1,
		Handle: func(c *cli.Context, _ error) {
			if len(c.Args) == 0 {
				fmt.Println("...")
				syscall.Exit(1)
			}
			dir := c.Args[0]
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

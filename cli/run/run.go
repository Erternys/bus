package run

import (
	"bus/cli"
	"bus/helper"
	"bus/middleware"
	"bus/script"
	"fmt"
	"path"
	"strings"
	"sync"
	"syscall"
)

var wg sync.WaitGroup

func NewRunCommand() cli.Command {
	return cli.Command{
		Name:         "run",
		Aliases:      []string{"r"},
		RequiredArgs: 1,
		Flags: []cli.Flag{
			cli.NewFlag("background", "run command(s) in background (by default: false)", cli.Bool, "bg", "b"),
		},
		Handle: func(c *cli.Context, _ error) {
			c.Execs(middleware.ReadConfigFile)

			if len(c.Args) == 0 {
				fmt.Println("Please give a script to execute `script_name` or `process@script_name`")
				syscall.Exit(1)
			}

			background := c.GetFlag("background", false).Value.(bool)
			baseConfig := c.GetState("config", nil).(middleware.Config)

			from := "*"
			scriptName := ""
			if strings.Contains(c.Args[0], "@") {
				a := strings.Split(c.Args[0], "@")
				from = a[0]
				scriptName = a[1]
			} else {
				scriptName = c.Args[0]
			}

			for _, packagePath := range baseConfig.PackagesPath {
				if from == "*" || packagePath.Name == from {
					config := packagePath.GetExtention(c).ParseConfig()
					scripts := config["scripts"].(map[string]interface{})
					cmd, ok := scripts[scriptName].(string)
					if packagePath.Extend == "nodejs" {
						manager := baseConfig.Manager
						cmd = fmt.Sprintf("%v run %v", manager, scriptName)
					}

					if !ok {
						continue
					}

					s := script.NewScript(packagePath, path.Join(helper.WorkSpacePath, packagePath.Path), cmd)
					if !background {
						wg.Add(1)
					}
					go s.Start(wg.Done)
				}
			}

			if !background {
				wg.Wait()
			}
		},
	}
}

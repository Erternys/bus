package run

import (
	"bus/buffer"
	"bus/cli"
	"bus/config"
	"bus/helper"
	"bus/middleware"
	"bus/script"
	"fmt"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"syscall"
)

var wg sync.WaitGroup
var startedScripts []*script.Script

func CtrlC(c chan os.Signal) {
	<-c
	buffer.Print(" Killing all scripts\n")
	for _, script := range startedScripts {
		script.Kill()
	}
	syscall.Exit(0)
}

func NewRunCommand() cli.Command {
	return cli.Command{
		Name:         "run",
		Aliases:      []string{"r"},
		RequiredArgs: 1,
		Flags: []cli.Flag{
			cli.NewFlag("background", "run command(s) in background (by default: false)", cli.Bool, "bg", "b"),
		},
		Description:      "Run one or more commands in parallel", // TODO: write a correct description of the run command
		ShortDescription: "Run one or more commands in parallel",
		Usage:            "run script_name | process@script_name",
		Handle: func(c *cli.Context, _ error) {
			c.Execs(middleware.ReadConfigFile)

			if len(c.Args) == 0 {
				buffer.Eprintf("%vPlease give a script to execute `script_name` or `process@script_name`%v\n", helper.Red+helper.Bold, helper.Reset)
				syscall.Exit(1)
			}

			background := c.GetFlag("background", false).Value.(bool)
			baseConfig := c.GetState("config", nil).(config.Config)

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
					config := middleware.GetPackageExtention(packagePath, c).ParseConfig()
					scripts := config["scripts"].(map[string]any)
					cmd, ok := scripts[scriptName].(string)
					if packagePath.Extend == "nodejs" {
						manager := baseConfig.JsManager
						cmd = fmt.Sprintf("%v run %v", manager, scriptName)
					}

					if !ok {
						continue
					}

					s := script.NewScript(packagePath, path.Join(helper.WorkSpacePath, packagePath.Path), cmd)
					s.DryRun = c.GetFlag("dry-run", false).Value.(bool)
					if !background {
						wg.Add(1)
					}
					go s.Start(wg.Done)
					startedScripts = append(startedScripts, s)
				}
			}
			sig := make(chan os.Signal)
			signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

			go CtrlC(sig)

			if !background {
				wg.Wait()
			}

			buffer.Println()
		},
	}
}

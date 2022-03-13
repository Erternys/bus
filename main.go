package main

import (
	"bus/cli"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"gopkg.in/yaml.v3"
)

func main() {
	app := cli.NewApp(
		"bus",
		"p",
		"0.1.0-beta",
	)
	app.AddFlag(cli.NewFlag("config", "change the config file used (by default: bus-ws.config.yaml)", cli.String, "c"))
	app.AddCommand(cli.Command{
		Name:         "init",
		RequiredArgs: 1,
		Handle: func(c *cli.Context, _ error) {
			if len(c.Args) == 0 || c.Args[0] == "repo" {
				path := c.GetFlag("config", "bus-ws.config.yaml").Value.(string)
				_, err := os.Stat(path)

				if os.IsNotExist(err) {
					fmt.Println("Press ^C (Ctrl+C) at any time to quit.")

					currentDir := getwd()

					name := input(fmt.Sprintf("project name: (%v) ", currentDir), currentDir)
					version := input(fmt.Sprintf("version: (%v) ", "1.0.0"), "1.0.0")
					description := input("description: ", "")
					repository := ""
					if getRepository() == "" {
						repository = input("git repository: ", "")
					} else {
						repository = input(fmt.Sprintf("git repository: (%v) ", getRepository()), getRepository())
					}

					content, _ := yaml.Marshal(map[string]interface{}{
						"name":        name,
						"version":     version,
						"description": description,
						"repository":  repository,
						"packages":    make(map[string]interface{}),
					})

					file, _ := os.Create(path)
					defer file.Close()

					file.Write(content)
				}
			}
		},
	})
	app.AddCommand(cli.Command{
		Name:         "run",
		Aliases:      []string{"r"},
		RequiredArgs: 1,
		Handle: func(c *cli.Context, _ error) {
			c.Execs(readConfigFile)
		},
	})

	app.SetHelpCommand()
	app.SetVersionCommand()
	app.Run(os.Args[1:])
}

func readConfigFile(c *cli.Context, next func()) {
	configFilePath := c.GetFlag("config", "bus-ws.config.yaml").Value.(string)
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Printf("the config file does not exist, you can execute the command \"%s init\" or \"%s init repo\"\n", c.App.Name, c.App.Name)
		syscall.Exit(1)
	}

	data := make(map[interface{}]interface{})
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		fmt.Println("the config file does not corect")
		syscall.Exit(1)
	}

	c.State["config"] = data

	next()
}

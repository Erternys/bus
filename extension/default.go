package extension

import (
	"bus/cli"
	"bus/helper"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	"gopkg.in/yaml.v3"
)

type Any interface{}

type Extension struct {
	Path    string
	Context *cli.Context
}

func Default() *Extension {
	return &Extension{
		Path:    "./",
		Context: nil,
	}
}

func (e *Extension) SetContext(c *cli.Context) {
	e.Context = c
}

func (e *Extension) Init(name, dir string) {
	filename := e.Context.GetFlag("config", "bus-ws.config.yaml").Value.(string)

	version := helper.Input(fmt.Sprintf("version: (%v) ", "1.0.0"), "1.0.0")
	description := helper.Input("description: ", "")
	scripts := make(map[string]string)
	for {
		sname := helper.Input("script name: ", "")
		if sname == "" {
			break
		}
		scripts[sname] = helper.Input("script command: ", "")
	}
	license := helper.Input(fmt.Sprintf("license: (%v) ", "ISC"), "ISC")

	content, _ := yaml.Marshal(map[string]Any{
		"name":        name,
		"version":     version,
		"description": description,
		"scripts":     scripts,
		"license":     license,
	})

	ioutil.WriteFile(filename, content, 0644)
}

func (e *Extension) GetConfigPath() string {
	configFileName := e.Context.GetState("filepath", "./").(string)
	config, _ := filepath.Abs(e.Path + string(os.PathSeparator) + configFileName)
	return config
}

func (e *Extension) ParseConfig() map[string]Any {
	data := make(map[string]Any)
	err := yaml.Unmarshal([]byte(e.GetConfigPath()), &data)
	if err != nil {
		fmt.Println("the config file does not correct")
		syscall.Exit(1)
	}
	return data
}

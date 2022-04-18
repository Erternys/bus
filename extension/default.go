package extension

import (
	"bus/cli"
	"bus/helper"

	// "bus/middleware"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	"gopkg.in/yaml.v3"
)

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

func (e *Extension) SetPath(p string) {
	e.Path = p
}

func (e *Extension) Init(name, dir string) {
	filename := e.Context.GetFlag("config", ".bus.yaml").Value.(string)

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

	content, _ := yaml.Marshal(map[string]interface{}{
		"name":        name,
		"version":     version,
		"description": description,
		"scripts":     scripts,
		"license":     license,
	})

	ioutil.WriteFile(filename, content, 0644)
}

func (e *Extension) InstallDep() {
	fmt.Printf("%v%v has no dependency manager%v\n", helper.Red, e.Path, helper.Reset)
}

func (e *Extension) GetConfigPath() string {
	configFileName := e.Context.GetState("filepath", "./").(string)
	config, _ := filepath.Abs(e.Path + string(os.PathSeparator) + configFileName)
	return config
}

func (e *Extension) ParseConfig() map[string]interface{} {
	data := make(map[string]interface{})
	content, err := ioutil.ReadFile(e.GetConfigPath())
	if err != nil {
		fmt.Println("the config file was remove")
		syscall.Exit(1)
	}
	err = yaml.Unmarshal(content, &data)
	if err != nil {
		fmt.Println("the config file does not correct")
		syscall.Exit(1)
	}
	return data
}

func (e *Extension) Clone() interface{} {
	return Default()
}

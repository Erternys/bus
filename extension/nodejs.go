package extension

import (
	"bus/helper"
	"bus/process"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

type NodeJSExtension struct {
	*Extension
}

func DefaultNodeJS() *NodeJSExtension {
	return &NodeJSExtension{
		&Extension{
			Path:    "./",
			Context: nil,
		},
	}
}

func (e *NodeJSExtension) Init(name, dir string) {
	fmt.Println("")

	kit := ""
	kitFlag := e.Context.GetFlag("kit", kit)
	if kitFlag.Value == "" {
		kit = helper.Input("kit: ", kit)
	}

	npm := e.Context.GetFlag("use", "npm")
	var proc *process.Process = nil

	if kit != "" {
		proc = process.NewProcess("npm create project with kit", fmt.Sprintf("%v create %v ./", npm.Value, kit))
	} else {
		proc = process.NewProcess("npm init project", fmt.Sprintf("%v init", npm.Value))
	}

	proc.UseStandardIO()
	proc.Run()
	proc.Wait()
}

func (e *NodeJSExtension) GetConfigPath() string {
	config, _ := filepath.Abs(e.Path + string(os.PathSeparator) + "package.json")
	return config
}

func (e *NodeJSExtension) ParseConfig() map[string]Any {
	data := make(map[string]Any)
	err := json.Unmarshal([]byte(e.GetConfigPath()), &data)
	if err != nil {
		fmt.Println("the config file does not correct")
		syscall.Exit(1)
	}
	return data
}

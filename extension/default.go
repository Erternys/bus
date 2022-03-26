package extension

import (
	"bus/cli"
	"fmt"
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

func (e *Extension) GetConfigPath() string {
	configFileName := e.Context.GetFlag("config", "bus-ws.config.yaml").Value.(string)
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

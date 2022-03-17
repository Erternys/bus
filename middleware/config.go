package middleware

import (
	"bus/cli"
	"fmt"
	"io/ioutil"
	"syscall"

	"gopkg.in/yaml.v3"
)

type Package struct {
	Path   string
	Extend string
}

type Config struct {
	Name         string
	Version      string
	Description  string
	Repository   string
	PackagesPath []*Package `yaml:"packages"`
}

func ReadConfigFile(c *cli.Context, next func()) {
	configFilePath := c.GetFlag("config", "bus-ws.config.yaml").Value.(string)
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Printf("the config file does not exist, you can execute the command \"%s init\" or \"%s init repo\"\n", c.App.Name, c.App.Name)
		syscall.Exit(1)
	}

	data := Config{}
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		fmt.Println("the config file does not correct")
		syscall.Exit(1)
	}

	c.State["config"] = data
	next()
}

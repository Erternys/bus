package middleware

import (
	"bus/cli"
	"bus/config"
	"fmt"
	"io/ioutil"
	"syscall"

	"gopkg.in/yaml.v3"
)

func ReadConfigFile(c *cli.Context, next func()) {
	configFilePath := c.GetFlag("config", "bus-ws.config.yaml").Value.(string)
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Printf("the config file does not exist, you can execute the command \"%s init\" or \"%s init repo\"\n", c.App.Name, c.App.Name)
		syscall.Exit(1)
	}

	data := config.Config{}
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		fmt.Println("the config file does not correct")
		syscall.Exit(1)
	}

	c.SetState("filepath", configFilePath)
	c.SetState("config", data)
	next()
}

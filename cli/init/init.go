package init

import (
	"bus/cli"
	"bus/helper"
	"bus/middleware"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"syscall"

	"gopkg.in/yaml.v3"
)

func NewInitCommand() cli.Command {
	return cli.Command{
		Name:         "init",
		RequiredArgs: 1,
		Handle: func(c *cli.Context, _ error) {
			filename := c.GetFlag("config", "bus-ws.config.yaml").Value.(string)
			_, err := os.Stat(filename)
			if len(c.Args) == 0 || c.Args[0] == "repo" {
				if os.IsNotExist(err) {
					fmt.Println("Press ^C (Ctrl+C) at any time to quit.\n")

					currentDir := helper.Getwd()

					name := helper.Input(fmt.Sprintf("project name: (%v) ", currentDir), currentDir)
					version := helper.Input(fmt.Sprintf("version: (%v) ", "1.0.0"), "1.0.0")
					description := helper.Input("description: ", "")
					repository := helper.GetRepository()

					if repository == "" {
						repository = helper.Input("git repository: ", "")
					} else {
						repository = helper.Input(fmt.Sprintf("git repository: (%v) ", repository), repository)
					}

					content, _ := yaml.Marshal(middleware.Config{
						Name:         name,
						Version:      version,
						Description:  description,
						Repository:   repository,
						Manager:      "npm",
						PackagesPath: make([]*middleware.Package, 0),
					})

					ioutil.WriteFile(filename, content, 0644)
				}
				return
			}
			c.Execs(middleware.ReadConfigFile)
			config := c.GetState("config", nil).(middleware.Config)

			dir := path.Clean(c.Args[0])

			for _, p := range config.PackagesPath {
				if p.Path == dir {
					fmt.Printf("the package `%v` already exist\n", dir)
					syscall.Exit(0)
				}
			}

			err = os.MkdirAll(dir, 0750)
			if err != nil {
				if strings.Contains(dir, "/") {
					fmt.Println("the folders do not have a correct name")
				} else {
					fmt.Println("the folder does not have a correct name")
				}
				syscall.Exit(1)
			}
			os.Chdir(dir)

			currentDir := helper.Getwd()
			name := helper.Input(fmt.Sprintf("sub-package name: (%v) ", currentDir), currentDir)
			extend := helper.Input(fmt.Sprintf("sub-package type: (%v) ", "default"), "default")

			extension, ok := middleware.Extensions[extend]
			for ; !ok; extension, ok = middleware.Extensions[extend] {
				fmt.Printf("invalid value: `%v`\n", extend)
				extend = helper.Input(fmt.Sprintf("sub-package type: (%v) ", "default"), "default")
			}

			extension.SetContext(c)
			extension.Init(name, dir)
			config.PackagesPath = append(config.PackagesPath, &middleware.Package{
				Path:   dir,
				Name:   name,
				Extend: extend,
			})
			sort.Slice(config.PackagesPath, func(i, j int) bool {
				return config.PackagesPath[i].Path < config.PackagesPath[j].Path
			})

			dirs := strings.Split(dir, "/")
			for i := range dirs {
				dirs[i] = ".."
			}

			os.Chdir(strings.Join(dirs, "/"))

			data, _ := yaml.Marshal(config)

			ioutil.WriteFile(filename, data, 0644)
		},
	}
}

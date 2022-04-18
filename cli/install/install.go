package install

import (
	"bus/cli"
	"bus/config"
	"bus/middleware"
)

func NewInstallCommand() cli.Command {
	dev := cli.NewFlag("dev", "", cli.String, "d")

	return cli.Command{
		Name:        "install",
		Description: "",
		Aliases:     []string{"i", "add"},
		Flags:       []cli.Flag{dev},
		Handle: func(c *cli.Context, err error) {
			c.Execs(middleware.ReadConfigFile)
			if len(c.Args) == 0 {
				baseConfig := c.GetState("config", nil).(config.Config)

				for _, packagePath := range baseConfig.PackagesPath {
					middleware.GetPackageExtention(packagePath, c).InstallDep()
				}
			}
		},
	}
}

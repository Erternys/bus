package main

import (
	"bus/cli"
	init_c "bus/cli/init"
	"bus/cli/install"
	"bus/cli/run"
	"os"
)

var app = cli.NewApp(
	"bus",
	"Monorepo manager usable with several programming languages (not only JS)",
	"0.1.0-beta",
)

func main() {
	app.AddFlag(cli.NewFlag("config", "Change the config file used (by default: .bus.yaml)", cli.String, "c"))
	app.AddFlag(cli.NewFlag("dry-run", "Show which script bus execute without execute it (muter characters do not work)", cli.Bool, "dr", "d"))

	app.AddCommand(init_c.NewInitCommand())
	app.AddCommand(run.NewRunCommand())
	app.AddCommand(install.NewInstallCommand())

	app.AddCommandFromExe(cli.ExeCommand{
		File:        "proxy",
		Name:        "proxy",
		Description: "Server proxy for bus",
	})

	app.SetHelpCommand()
	app.SetVersionCommand()
	app.Run(os.Args[1:])
}

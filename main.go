package main

import (
	"bus/cli"
	init_c "bus/cli/init"
	"bus/cli/run"
	"os"
)

var app = cli.NewApp(
	"bus",
	"Monorepo manager usable with several programming languages (not only JS)",
	"0.1.0-beta",
)

func main() {
	app.AddFlag(cli.NewFlag("config", "Change the config file used (by default: bus-ws.config.yaml)", cli.String, "c"))
	app.AddFlag(cli.NewFlag("use", "Set the JavaScript package manager (by default: npm)", cli.String))
	app.AddCommand(init_c.NewInitCommand())
	app.AddCommand(run.NewRunCommand())

	app.SetHelpCommand()
	app.SetVersionCommand()
	app.Run(os.Args[1:])
}

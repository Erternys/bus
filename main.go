package main

import (
	"bus/cli"
	init_c "bus/cli/init"
	"bus/cli/install"
	new_c "bus/cli/new"
	"bus/cli/run"
	"os"
)

var app = cli.NewApp(
	"bus",
	"Monorepo manager usable with several programming languages (not only JS)",
	"0.1.0-beta",
)

func main() {
	app.SetGlobal("SHELL", os.Getenv("SHELL"))

	app.AddFlag(cli.NewFlag("dry-run", "Show which script bus execute without execute it (muter characters are disable)", cli.Bool, "dr", "d"))

	app.AddCommand(new_c.NewNewCommand())
	app.AddCommand(init_c.NewInitCommand())
	app.AddCommand(run.NewRunCommand())
	app.AddCommand(install.NewInstallCommand())

	app.AddCommandFromExe(cli.ExeCommand{
		File:             "proxy",
		Name:             "proxy",
		Description:      "Server proxy for bus",
		ShortDescription: "Server proxy for bus",
		Usage:            "proxy [Subcommand]",
	})

	app.SetHelpCommand()
	app.SetVersionCommand()
	app.Run(os.Args[1:])
}

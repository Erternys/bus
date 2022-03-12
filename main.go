package main

import (
	"bus/cli"
	"fmt"
	"os"
)

func main() {
	app := cli.NewApp(
		"bus",
		"p",
		"0.1.0-beta",
	)
	app.AddFlag(cli.NewFlag("test", "", cli.Bool))
	app.AddCommand(cli.Command{
		Name:         "test",
		Aliases:      []string{"t"},
		RequiredArgs: 1,
		Handle: func(c *cli.Context, _ error) {
			f := c.GetFlag("test", false)
			fmt.Println(c.Args, f.Value)
		},
	})

	app.SetHelpCommand()
	app.SetVersionCommand()
	app.Run(os.Args[1:])
	// p := process.NewProcess(
	// 	"echo",
	// 	"echo test 'test'",
	// )
	// err := p.Run()
	// if err != nil {
	// 	panic(err)
	// }
	// err = p.Wait()
	// if err != nil {
	// 	panic(err)
	// }
}

package main

import (
	"bus/cli"
	"os"
)

func main() {
	app := cli.NewApp(
		"bus",
		"",
		"",
		cli.Command{
			Name:    "test",
			Aliases: []string{"t"},
		},
	)

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

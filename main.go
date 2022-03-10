package main

import (
	"bus/cli"
	"fmt"
	"os"
)

func main() {
	app := cli.NewApp(
		"bus",
		"",
		"",
	)
	app.AddCommand(cli.Command{
		Name:    "test",
		Aliases: []string{"t"},
		Handle: func(c *cli.Context) {
			t := c.Flags["-test"]
			fmt.Println(t.Name)
		},
	})
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

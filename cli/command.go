package cli

type Command struct {
	Name        string
	Aliases     []string
	Description string
	// Action: func(c *cli.Context) {
	// 	pe := "peppers"
	// 	peppers := append(pizza, pe)
	// 	m := strings.Join(peppers, " ")
	// 	fmt.Println(m)
	// }
	Cli *CliApp
}

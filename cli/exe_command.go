package cli

type ExeCommand struct {
	File             string
	Name             string
	Flags            []Flag
	Aliases          []string
	FlagAliases      []string
	Description      string
	ShortDescription string
	Usage            string
}

func (c *ExeCommand) AddFlag(flag Flag) {
	c.Flags = append(c.Flags, flag)
}

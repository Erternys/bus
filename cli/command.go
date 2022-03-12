package cli

type HandleAction func(c *Context, err error)

type Command struct {
	Name         string
	Aliases      []string
	Description  string
	RequiredArgs int
	Handle       HandleAction
}

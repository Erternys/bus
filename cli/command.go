package cli

type HandleAction func(c *Context)

type Command struct {
	Name        string
	Aliases     []string
	Description string
	Handle      HandleAction
}

package cli

type Any interface{}

type Context struct {
	App *CliApp

	Flags map[string]Any
	State map[string]Any

	Continue bool
}

func NewContext(app *CliApp) *Context {
	return &Context{
		App:      app,
		Flags:    map[string]Any{},
		State:    map[string]Any{},
		Continue: true,
	}
}

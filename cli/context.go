package cli

import "fmt"

type HandleExec func(c *Context, next func())

type Context struct {
	App *CliApp

	Args    []string
	RawArgs []string
	Flags   map[string]Flag
	State   map[string]any
}

func NewContext(app *CliApp) *Context {
	return &Context{
		App:   app,
		Args:  []string{},
		Flags: map[string]Flag{},
		State: map[string]any{},
	}
}

func (c *Context) StateExist(name string) bool {
	_, ok := c.State[name]
	return ok
}

func (c *Context) GetState(name string, defaultValue any) any {
	if c.StateExist(name) {
		return c.State[name]
	}
	return defaultValue
}

func (c *Context) SetState(name string, value any) {
	c.State[name] = value
}

func (c *Context) FlagExist(name string) bool {
	_, ok := c.Flags[name]
	return ok
}

func (c *Context) GetFlag(name string, defaultValue any) Flag {
	if c.FlagExist(name) {
		return c.Flags[name]
	}

	for _, flag := range c.App.flags {
		if flag.Name == name {
			current := flag.clone()
			current.SetValue(fmt.Sprintf("%v", defaultValue))

			return current
		}
	}

	flag := DefaultFlag()
	flag.Value = defaultValue
	switch defaultValue.(type) {
	case bool:
		flag.Kind = Bool
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		flag.Kind = Int
	case float32, float64:
		flag.Kind = Float
	default:
		flag.Kind = String
	}
	return flag
}

func (c *Context) Execs(callbacks ...HandleExec) {
	next(callbacks, c, 0)()
}

func next(callbacks []HandleExec, c *Context, i int) func() {
	return func() {
		if len(callbacks) > i {
			callbacks[i](c, next(callbacks, c, i+1))
		}
	}
}

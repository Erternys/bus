package cli

import "reflect"

type Any interface{}

type Context struct {
	App *CliApp

	Flags map[string]Flag
	State map[string]Any
}

func NewContext(app *CliApp) *Context {
	return &Context{
		App:   app,
		Flags: map[string]Flag{},
		State: map[string]Any{},
	}
}

func (c *Context) FlagExist(name string) bool {
	_, ok := c.Flags[name]
	return ok
}

func (c *Context) GetFlag(name string, defaultValue Any) Flag {
	if c.FlagExist(name) {
		return c.Flags[name]
	}

	flag := DefaultFlag()
	flag.Value = defaultValue
	switch reflect.ValueOf(defaultValue).Kind() {
	case reflect.Bool:
		flag.Kind = Bool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		flag.Kind = Int
	case reflect.Float32, reflect.Float64:
		flag.Kind = Float
	default:
		flag.Kind = String
	}
	return flag
}

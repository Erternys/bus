package cli

import (
	"fmt"
	"strconv"
)

type FlagKind uint8

const (
	String FlagKind = iota
	Bool
	Int
	Float
)

type Flag struct {
	Name    string
	Aliases []string
	Value   interface{}
	Kind    FlagKind
}

func NewFlag(name string, kind FlagKind, aliases ...string) Flag {
	return Flag{
		Name:    name,
		Kind:    kind,
		Aliases: aliases,
		Value:   nil,
	}
}

func DefaultFlag() Flag {
	return Flag{}
}

func (f *Flag) SetValue(value string) error {
	switch f.Kind {
	case Bool:
		switch value {
		case "true", "yes", "y", "":
			f.Value = true
			return nil
		case "false", "no", "n":
			f.Value = false
			return nil
		default:
			return fmt.Errorf("parsing \"%s\": invalid syntax, you can only give `true`, `yes`, `y`, `false`, `no` or `n`", value)
		}
	case Int:
		v, err := strconv.ParseInt(value, 10, 64)
		f.Value = v
		return err
	case Float:
		v, err := strconv.ParseFloat(value, 64)
		f.Value = v
		return err
	case String:
		f.Value = value
		return nil
	default:
		return fmt.Errorf("typing \"%s\": invalid type", f.Name)
	}
}

func (f *Flag) clone() Flag {
	return Flag{
		Name:    f.Name,
		Kind:    f.Kind,
		Aliases: f.Aliases,
		Value:   nil,
	}
}

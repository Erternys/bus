package cli

import (
	"fmt"
	"strconv"
	"strings"
)

type FlagKind uint8

const (
	String FlagKind = iota
	Bool
	Int
	Float
)

type Flag struct {
	Name        string
	Description string
	Aliases     []string
	Value       any
	Kind        FlagKind
}

func NewFlag(name string, description string, kind FlagKind, aliases ...string) Flag {
	return Flag{
		Name:        name,
		Description: description,
		Kind:        kind,
		Aliases:     aliases,
		Value:       nil,
	}
}

func DefaultFlag() Flag {
	return Flag{}
}

func (f *Flag) setValueAndKind(value string) {
	switch value {
	case "true", "yes", "y", "":
		f.Value = true
		f.Kind = Bool
	case "false", "no", "n":
		f.Value = false
		f.Kind = Bool
	default:
		var err error = nil
		var n any = nil
		if strings.Contains(value, ".") {
			n, err = strconv.ParseFloat(value, 64)
			f.Kind = Float
		} else {
			n, err = strconv.ParseInt(value, 10, 64)
			f.Kind = Int
		}
		if err == nil {
			f.Value = n
		} else {
			f.Value = value
			f.Kind = String
		}
	}
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
		Name:        f.Name,
		Description: f.Description,
		Kind:        f.Kind,
		Aliases:     f.Aliases,
		Value:       nil,
	}
}

func (f *Flag) ToString() string {
	if f.Kind == Bool {
		if f.Value.(bool) {
			return fmt.Sprintf("--%v", f.Name)
		}
		return ""
	}
	return fmt.Sprintf("--%v=%v", f.Name, f.Value)
}

package middleware

import (
	"bus/cli"
	"bus/extension"
)

type Extension interface {
	GetConfigPath() string
	ParseConfig() map[string]extension.Any
}

var extensions = map[string]extension.Any{
	"nodejs":  extension.DefaultNodeJS(),
	"default": extension.Default(),
}

func (p *Package) GetExtention(c *cli.Context) Extension {
	// TODO
	return nil
}

package middleware

import (
	"bus/cli"
	"bus/extension"
)

type Extension interface {
	SetContext(c *cli.Context)
	GetConfigPath() string
	ParseConfig() map[string]extension.Any
}

var extensions = map[string]Extension{
	"nodejs":  extension.DefaultNodeJS(),
	"default": extension.Default(),
}

func (p *Package) GetExtention(c *cli.Context) Extension {
	for key, value := range extensions {
		if key == p.Extend {
			return value
		}
	}
	return extensions["default"]
}

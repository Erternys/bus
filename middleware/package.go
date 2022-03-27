package middleware

import (
	"bus/cli"
	"bus/extension"
)

type Extension interface {
	Init(name, dir string)
	SetContext(c *cli.Context)
	GetConfigPath() string
	ParseConfig() map[string]extension.Any
}

var Extensions = map[string]Extension{
	"nodejs":  extension.DefaultNodeJS(),
	"default": extension.Default(),
}

func (p *Package) GetExtention(c *cli.Context) Extension {
	for key, value := range Extensions {
		if key == p.Extend {
			return value
		}
	}
	return Extensions["default"]
}

package middleware

import (
	"bus/cli"
	"bus/extension"
)

type Extension interface {
	Init(name, dir string)
	InstallDep()
	SetContext(c *cli.Context)
	SetPath(p string)
	GetConfigPath() string
	ParseConfig() map[string]interface{}

	Clone() interface{}
}

var Extensions = map[string]Extension{
	"nodejs":  extension.DefaultNodeJS(),
	"default": extension.Default(),
}

func (p *Package) GetExtention(c *cli.Context) Extension {
	for key, value := range Extensions {
		if key == p.Extend {
			extension := value.Clone().(Extension)
			extension.SetContext(c)
			extension.SetPath(p.Path)
			return extension
		}
	}
	return Extensions["default"]
}

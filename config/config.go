package config

type Package struct {
	Path   string
	Name   string
	Extend string
}

type Config struct {
	Name         string
	Version      string
	Description  string
	Repository   string
	JsManager    string     `yaml:"js_manager"`
	PackagesPath []*Package `yaml:"packages"`
}

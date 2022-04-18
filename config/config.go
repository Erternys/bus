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
	Manager      string
	PackagesPath []*Package `yaml:"packages"`
}

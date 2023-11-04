package config

type Proxy struct {
	OnScript *struct {
		ListenRun      []string `yaml:"listen_run"`
		RestartOnCrash bool     `default:"false" yaml:"restart_on_crash"`
	} `yaml:"on_script"`
}

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
	Proxy        Proxy
}

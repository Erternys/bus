package extension

import (
	"os"
	"path/filepath"
)

type NodeJSExtension struct {
	*Extension
}

func DefaultNodeJS() NodeJSExtension {
	return NodeJSExtension{
		&Extension{
			Path:    "./",
			Context: nil,
		},
	}
}

func (e *NodeJSExtension) GetConfigPath() string {
	config, _ := filepath.Abs(e.Path + string(os.PathSeparator) + "package.json")
	return config
}

func (e *NodeJSExtension) ParseConfig() Any {
	return nil
}

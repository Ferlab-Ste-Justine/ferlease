package fluxcd

import (
	"path"
	yaml "gopkg.in/yaml.v2"

	"github.com/Ferlab-Ste-Justine/ferlease/tplcore"
)

type FsConventions struct {
	Naming    string
	FluxcdDir string `yaml:"fluxcd_directory"`
	AppsDir   string `yaml:"apps_directory"`
}

type Orchestration struct {
	FsConventions *FsConventions 
	FluxcdFile    string
	AppFiles      map[string]string
}

func loadFsConventions(fPath string, params *tplcore.TemplateParameters) (*FsConventions, error) {
	var conv FsConventions
	
	res, err := params.LoadFile(fPath)
	if err != nil {
		return nil, err
	}

	yamlErr := yaml.Unmarshal([]byte(res), &conv)
	if yamlErr != nil {
		return nil, yamlErr
	}

	return &conv, nil
}

func LoadTemplate(tPath string, params *tplcore.TemplateParameters) (*Orchestration, error) {
	fsConventionsPath := path.Join(tPath, "filesystem-conventions.yml")
	fluxcdPath := path.Join(tPath, "fluxcd.yml")
	appPath := path.Join(tPath, "app")
	
	var o Orchestration

	var fsCoErr error
	o.FsConventions, fsCoErr = loadFsConventions(fsConventionsPath, params)
	if fsCoErr != nil {
		return nil, fsCoErr
	}

	var flFiErr error
	o.FluxcdFile, flFiErr = params.LoadFile(fluxcdPath)
	if flFiErr != nil {
		return nil, flFiErr
	}

	var appFlsErr error
	o.AppFiles, appFlsErr = params.LoadDir(appPath)
	if appFlsErr != nil {
		return nil, appFlsErr
	}

	return &o, nil
}
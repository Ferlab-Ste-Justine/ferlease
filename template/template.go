package template

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
	
	res, err := params.LoadFile(path.Join(fPath, "filesystem-conventions.yml"))
	if err != nil {
		return nil, err
	}

	yamlErr := yaml.Unmarshal(res, &conv)
	if yamlErr != nil {
		return nil, yamlErr
	}

	return &conv, nil
}

func loadFluxcdFile(fPath string, params *tplcore.TemplateParameters) (string, error) {
	res, err := params.LoadFile(path.Join(fPath, "fluxcd.yml"))
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func LoadTemplate(tPath string, params *tplcore.TemplateParameters) (*Orchestration, error) {
	var o Orchestration

	var fsCoErr error
	o.FsConventions, fsCoErr = loadFsConventions(tPath, params)
	if fsCoErr != nil {
		return nil, fsCoErr
	}

	var flFiErr error
	o.FluxcdFile, flFiErr = loadFluxcdFile(tPath, params)
	if flFiErr != nil {
		return nil, flFiErr
	}

	var appFlsErr error
	o.AppFiles, appFlsErr = params.LoadDir(tPath)
	if appFlsErr != nil {
		return nil, appFlsErr
	}

	return &o, nil
}
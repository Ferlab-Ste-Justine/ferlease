package terraform

import (
	"os"
	"path"
	yaml "gopkg.in/yaml.v2"

	"github.com/Ferlab-Ste-Justine/ferlease/tplcore"
)

type FsConventions struct {
	Naming string
	Dir    string `yaml:"directory"`
}

type Orchestration struct {
	FsConventions  *FsConventions 
	EntrypointFile string
	ModuleFiles    map[string]string
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
	entrypointPath := path.Join(tPath, "entrypoint.tf")
	modulePath := path.Join(tPath, "module")

	var o Orchestration

	var fsCoErr error
	o.FsConventions, fsCoErr = loadFsConventions(fsConventionsPath, params)
	if fsCoErr != nil {
		return nil, fsCoErr
	}

	var entFiErr error
	o.EntrypointFile, entFiErr = params.LoadFile(entrypointPath)
	if entFiErr != nil {
		return nil, entFiErr
	}

	modulePathExits, existsErr := PathExists(modulePath)
	if existsErr != nil {
		return nil, existsErr
	}

	if modulePathExits {
		var modFlsErr error
		o.ModuleFiles, modFlsErr = params.LoadDir(modulePath)
		if modFlsErr != nil {
			return nil, modFlsErr
		}
	}

	return &o, nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return true, err
		}

		return false, nil
	}

	return true, nil
}
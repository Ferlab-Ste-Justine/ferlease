package template

import (
	yaml "gopkg.in/yaml.v2"
)

type FsConventions struct {
	Naming    string
	FluxcdDir string `yaml:"fluxcd_directory"`
	AppsDir   string `yaml:"apps_directory"`
}

type Orchestration struct {
	FsConventions *FsConventions 
	FluxcdFile    string
	AppFiles      []string
}

func loadFsConventions(path string) (*FsConventions, error) {
	return nil, nil
}

func loadFluxcdFile(path string) (string, error) {
	return nil, nil
}

func loadAppFiles(path string) ([]string, error) {
	return nil, nil
}

func LoadTemplate(path string) (*Orchestration, error) {
	return nil, nil
}
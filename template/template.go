package template

import (
	"io/ioutil"
	"path"
	"text/template"
	yaml "gopkg.in/yaml.v2"
)

type TemplateParameters struct {
	Repo    string
	Service string
	Release string
}

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

func loadFsConventions(path string, params *TemplateParameters) (*FsConventions, error) {
	var conv FsConventions
	
	f, err := ioutil.ReadFile(path.Join(path, "filesystem-conventions.yml"))
	if err != nil {
		return nil, err
	}

	tmpl, tErr := template.New("FsConventions").Parse(f)
	if tErr != nil {
		return nil, tErr
	}

	var b bytes.Buffer
	exErr = tmpl.Execute(&b, params)
	if exErr != nil {
		return nil, exErr
	}

	yamlErr = yaml.Unmarshal(b, &c)
	if yamlErr != nil {
		return nil, yamlErr
	}

	return FsConventions, nil
}

func loadFluxcdFile(path string, params *TemplateParameters) (string, error) {
	//path.Join(path, "fluxcd.yml")
	return nil, nil
}

func loadAppFiles(path string, params *TemplateParameters) ([]string, error) {
	//path.Join(path, "app")
	return nil, nil
}

func LoadTemplate(path string, params *TemplateParameters) (*Orchestration, error) {
	var o Orchestration

	o.FsConvention, fsCoErr := loadFsConventions(path, params)
	if fsCoErr != nil {
		return nil, fsCoErr
	}

	o.FluxcdFile, flFiErr := loadFluxcdFile(path, params)
	if flFiErr != nil {
		return nil, flFiErr
	}


	o.AppFiles, appFlsErr := loadAppFiles(path, params)
	if appFlsErr != nil {
		return nil, appFlsErr
	}

	return o, nil
}